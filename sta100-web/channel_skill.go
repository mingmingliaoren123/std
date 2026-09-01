package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"openclaw-orchestrator/orchestrator"
)

const (
	channelSkillSessionKind = "channel_skill_sessions"
	channelSkillRouteKind   = "channel_skill_routes"
	channelSkillEventKind   = "channel_skill_events"
	channelSkillSessionTTL  = 30 * time.Minute
)

var channelSkillIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._:-]{0,63}$`)

type ChannelSkillDefinition struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AgentID     string   `json:"agentId"`
	Required    []string `json:"required,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type ChannelSkillRoute struct {
	ID             string   `json:"id"`
	Channel        string   `json:"channel"`
	Account        string   `json:"account"`
	AgentID        string   `json:"agentId"`
	AllowedSenders []string `json:"allowedSenders,omitempty"`
	AllowedChats   []string `json:"allowedChats,omitempty"`
	Enabled        bool     `json:"enabled"`
	UpdatedAt      string   `json:"updatedAt"`
}

type ChannelSkillSession struct {
	ID              string            `json:"id"`
	Channel         string            `json:"channel"`
	Account         string            `json:"account"`
	Sender          string            `json:"sender"`
	Conversation    string            `json:"conversation"`
	AgentID         string            `json:"agentId"`
	SkillID         string            `json:"skillId"`
	CandidateIDs    []string          `json:"candidateIds,omitempty"`
	State           string            `json:"state"`
	Parameters      map[string]string `json:"parameters,omitempty"`
	Missing         []string          `json:"missing,omitempty"`
	OriginalMessage string            `json:"originalMessage,omitempty"`
	LastPrompt      string            `json:"lastPrompt,omitempty"`
	RunID           string            `json:"runId,omitempty"`
	Result          string            `json:"result,omitempty"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
	ExpiresAt       string            `json:"expiresAt"`
}

type channelSkillInboundRequest struct {
	Channel      string         `json:"channel"`
	Account      string         `json:"account,omitempty"`
	Conversation string         `json:"conversation,omitempty"`
	Sender       string         `json:"sender,omitempty"`
	From         string         `json:"from,omitempty"`
	MessageID    string         `json:"messageId,omitempty"`
	Message      string         `json:"message"`
	Timestamp    int64          `json:"timestamp,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

var channelSkillDefinitions = []ChannelSkillDefinition{
	{ID: "quote.create", Name: "生成客户报价单", Description: "根据客户和产品信息生成报价单", AgentID: "export-agent", Required: []string{"customer", "products"}, Keywords: []string{"报价单", "报价", "quote"}},
	{ID: "order.create", Name: "生成客户订单", Description: "根据客户、报价和产品信息生成订单", AgentID: "invoice-agent", Required: []string{"customer", "products"}, Keywords: []string{"订单", "order"}},
	{ID: "quote.search", Name: "查询历史报价", Description: "查询客户历史报价记录", AgentID: "export-agent", Required: []string{"query"}, Keywords: []string{"历史报价", "报价记录"}},
	{ID: "order.search", Name: "查询历史订单", Description: "查询客户历史订单记录", AgentID: "invoice-agent", Required: []string{"query"}, Keywords: []string{"历史订单", "订单记录"}},
	{ID: "inventory.search", Name: "查询公司库存", Description: "查询产品库存和可售数量", AgentID: "inventory-agent", Required: []string{"query"}, Keywords: []string{"库存", "现货"}},
	{ID: "customer.communication", Name: "查询客户沟通记录", Description: "调取客户历史沟通记录", AgentID: "email-generator", Required: []string{"customer"}, Keywords: []string{"客户沟通", "客户跟进", "联系记录"}},
	{ID: "supplier.communication", Name: "查询供应商沟通记录", Description: "调取供应商历史沟通记录", AgentID: "supplier-aggregator", Required: []string{"supplier"}, Keywords: []string{"供应商沟通", "供应商记录"}},
	{ID: "market.analysis", Name: "生成市场分析", Description: "生成指定国家或行业的市场分析", AgentID: "market-analyzer", Required: []string{"query"}, Keywords: []string{"市场分析", "市场行情", "market"}},
	{ID: "exhibition.report", Name: "生成展会报告", Description: "汇总指定展会的报道并生成报告", AgentID: "exhibition-advisor", Required: []string{"query"}, Keywords: []string{"展会报告", "展会情况", "展会"}},
	{ID: "news.daily", Name: "获取行业日报", Description: "生成自行车行业日报", AgentID: "sta100-news-curator", Required: []string{"query"}, Keywords: []string{"行业日报", "每日新闻", "日报"}},
	{ID: "market.research", Name: "行业和竞品调查", Description: "调查指定行业、产品或竞争对手", AgentID: "market-analyzer", Required: []string{"query"}, Keywords: []string{"行业调查", "竞品", "数据调查", "对手"}},
	{ID: "documents.search", Name: "查询采购和单据", Description: "查询采购资料和业务单据", AgentID: "invoice-agent", Required: []string{"query"}, Keywords: []string{"采购数据", "单据", "合同", "发票", "装箱单"}},
	{ID: "supplier.analysis", Name: "分析供应商", Description: "分析供应商能力、历史和风险", AgentID: "supplier-aggregator", Required: []string{"supplier"}, Keywords: []string{"供应商分析", "供应商情况"}},
	{ID: "email.generate", Name: "生成商务邮件", Description: "根据业务上下文生成邮件内容", AgentID: "email-generator", Required: []string{"customer", "query"}, Keywords: []string{"邮件", "email"}},
}

func channelSkills() []ChannelSkillDefinition {
	return append([]ChannelSkillDefinition(nil), channelSkillDefinitions...)
}

func findChannelSkill(id string) (ChannelSkillDefinition, bool) {
	for _, skill := range channelSkillDefinitions {
		if skill.ID == id {
			return skill, true
		}
	}
	return ChannelSkillDefinition{}, false
}

func channelSkillHash(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte{0})
		h.Write([]byte(strings.TrimSpace(part)))
	}
	return hex.EncodeToString(h.Sum(nil))[:32]
}

func channelSkillSessionID(channel, account, conversation, sender string) string {
	return "css-" + channelSkillHash(channel, account, conversation, sender)
}

func channelSkillEventID(channel, account, messageID, message string) string {
	if strings.TrimSpace(messageID) == "" {
		messageID = channelSkillHash(message)
	}
	return "cse-" + channelSkillHash(channel, account, messageID)
}

func (a *businessAPI) channelSkillRoutesHandler(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 1 && parts[0] == "skills" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"skills": channelSkills()})
		return
	}
	if len(parts) == 1 && parts[0] == "sessions" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		sessions, err := listRecords[ChannelSkillSession](r.Context(), a.store, channelSkillSessionKind)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": sessions, "total": len(sessions)})
		return
	}
	if len(parts) == 2 && parts[0] == "routes" && parts[1] != "" {
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if err := a.store.deleteRecord(r.Context(), channelSkillRouteKind, parts[1]); err != nil {
			if errors.Is(err, errRecordNotFound) {
				writeAPIError(w, http.StatusNotFound, "CHANNEL_SKILL_ROUTE_NOT_FOUND", "通道 Skill 路由不存在")
				return
			}
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": parts[1]})
		return
	}
	if len(parts) != 1 || parts[0] != "routes" {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "通道 Skill 接口不存在")
		return
	}
	if r.Method == http.MethodGet {
		routes, err := listRecords[ChannelSkillRoute](r.Context(), a.store, channelSkillRouteKind)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": routes, "total": len(routes)})
		return
	}
	if !allowMutation(w, r, http.MethodPut) {
		return
	}
	var route ChannelSkillRoute
	if err := decodeJSONBody(w, r, &route); err != nil {
		return
	}
	if err := validateChannelSkillRoute(&route); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_CHANNEL_SKILL_ROUTE", err.Error())
		return
	}
	if route.ID == "" {
		// Keep compatibility with routes created before the stable account key
		// was introduced, while preventing competing routes for one account.
		existing, listErr := listRecords[ChannelSkillRoute](r.Context(), a.store, channelSkillRouteKind)
		if listErr != nil {
			writeBusinessError(w, listErr)
			return
		}
		for _, candidate := range existing {
			if candidate.Channel == route.Channel && candidate.Account == route.Account {
				route.ID = candidate.ID
				break
			}
		}
		if route.ID == "" {
			route.ID = "csr-" + channelSkillHash(route.Channel, route.Account)
		}
	}
	route.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	err := a.store.put(r.Context(), channelSkillRouteKind, route.ID, route)
	if errors.Is(err, errRecordNotFound) {
		err = a.store.create(r.Context(), channelSkillRouteKind, route.ID, route)
	}
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": true, "item": route})
}

func validateChannelSkillRoute(route *ChannelSkillRoute) error {
	route.Channel = strings.ToLower(strings.TrimSpace(route.Channel))
	route.Account = strings.TrimSpace(route.Account)
	route.AgentID = strings.ToLower(strings.TrimSpace(route.AgentID))
	if !channelSkillIDPattern.MatchString(route.Channel) || route.Account == "" {
		return errors.New("通道和账号不能为空")
	}
	if route.AgentID == "" {
		return errors.New("Agent ID 不能为空")
	}
	for _, skill := range channelSkillDefinitions {
		if skill.AgentID == route.AgentID {
			goto validAgent
		}
	}
	return errors.New("Agent ID 不在通道 Skill 支持范围内")
validAgent:
	if len(route.AllowedSenders) == 0 && len(route.AllowedChats) == 0 {
		return errors.New("至少配置一个允许的发送人或会话；明确使用 * 才允许全部")
	}
	route.AllowedSenders = cleanChannelSkillValues(route.AllowedSenders)
	route.AllowedChats = cleanChannelSkillValues(route.AllowedChats)
	return nil
}

func cleanChannelSkillValues(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	return out
}

func channelSkillRouteAllows(route ChannelSkillRoute, sender, conversation string) bool {
	if !route.Enabled || len(route.AllowedSenders) == 0 && len(route.AllowedChats) == 0 {
		return false
	}
	for _, value := range route.AllowedSenders {
		if value == "*" || value == sender {
			return true
		}
	}
	for _, value := range route.AllowedChats {
		if value == "*" || value == conversation {
			return true
		}
	}
	return false
}

func (a *businessAPI) findChannelSkillRoute(ctx context.Context, request channelSkillInboundRequest) (ChannelSkillRoute, bool, error) {
	routes, err := listRecords[ChannelSkillRoute](ctx, a.store, channelSkillRouteKind)
	if err != nil {
		return ChannelSkillRoute{}, false, err
	}
	for _, route := range routes {
		if route.Channel == request.Channel && route.Account == request.Account && channelSkillRouteAllows(route, request.Sender, request.Conversation) {
			return route, true, nil
		}
	}
	return ChannelSkillRoute{}, false, nil
}

func (a *businessAPI) channelSkillInboundHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "request method is not supported")
		return
	}
	if !channelSkillLocalRequest(r) || !channelSkillTokenMatches(r) {
		writeAPIError(w, http.StatusForbidden, "CHANNEL_SKILL_UNAUTHORIZED", "通道入站请求未授权")
		return
	}
	channel := strings.ToLower(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/openclaw/inbound/"), "/"))
	if !channelSkillIDPattern.MatchString(channel) {
		writeAPIError(w, http.StatusBadRequest, "INVALID_CHANNEL", "通道标识无效")
		return
	}
	var request channelSkillInboundRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	request.Channel = channel
	request.Account = strings.TrimSpace(request.Account)
	request.Sender = strings.TrimSpace(request.Sender)
	request.From = strings.TrimSpace(request.From)
	request.Conversation = strings.TrimSpace(request.Conversation)
	request.Message = strings.TrimSpace(request.Message)
	if request.Sender == "" {
		request.Sender = request.From
	}
	if request.Conversation == "" {
		request.Conversation = request.Sender
	}
	if request.Account == "" || request.Sender == "" || request.Message == "" || len(request.Message) > 32<<10 {
		writeAPIError(w, http.StatusBadRequest, "INVALID_CHANNEL_MESSAGE", "账号、发送人、会话和消息不能为空")
		return
	}
	route, allowed, err := a.findChannelSkillRoute(r.Context(), request)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	if !allowed {
		writeJSON(w, http.StatusForbidden, map[string]any{"accepted": false, "message": "该通道账号尚未授权执行 STA-100 功能，请先配置允许的发送人或会话。"})
		return
	}
	eventID := channelSkillEventID(request.Channel, request.Account, request.MessageID, request.Message)
	if err := a.store.create(r.Context(), channelSkillEventKind, eventID, map[string]any{"id": eventID, "messageId": request.MessageID, "createdAt": time.Now().UTC().Format(time.RFC3339Nano)}); err != nil {
		if errors.Is(err, errRecordConflict) {
			writeJSON(w, http.StatusOK, map[string]any{"accepted": true, "duplicate": true})
			return
		}
		writeBusinessError(w, err)
		return
	}
	requestJSON := request
	go a.processChannelSkill(context.Background(), requestJSON, route)
	writeJSON(w, http.StatusAccepted, map[string]any{"accepted": true, "state": "processing", "conversation": request.Conversation})
}

func channelSkillLocalRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsLoopback() || ip.IsUnspecified())
}

func channelSkillTokenMatches(r *http.Request) bool {
	expected := channelSkillToken()
	provided := strings.TrimSpace(r.Header.Get("X-STA100-Channel-Token"))
	return expected != "" && provided != "" && len(provided) == len(expected) && subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}

func channelSkillToken() string {
	if token := strings.TrimSpace(os.Getenv("STA100_CHANNEL_SKILL_TOKEN")); token != "" {
		return token
	}
	paths := []string{}
	if root := strings.TrimSpace(os.Getenv("STA100_DATA_ROOT")); root != "" {
		paths = append(paths, filepath.Join(root, "channel-skill.token"))
	}
	if executable, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(executable), "..", "data", "channel-skill.token"))
	}
	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
			return strings.TrimSpace(string(data))
		}
	}
	return ""
}

func (a *businessAPI) processChannelSkill(ctx context.Context, request channelSkillInboundRequest, route ChannelSkillRoute) {
	sessionID := channelSkillSessionID(request.Channel, request.Account, request.Conversation, request.Sender)
	a.channelSkillMu.Lock()
	if a.channelSkillLocks == nil {
		a.channelSkillLocks = make(map[string]*sync.Mutex)
	}
	lock := a.channelSkillLocks[sessionID]
	if lock == nil {
		lock = &sync.Mutex{}
		a.channelSkillLocks[sessionID] = lock
	}
	a.channelSkillMu.Unlock()
	lock.Lock()
	defer lock.Unlock()
	processCtx, cancel := context.WithTimeout(ctx, 8*time.Minute)
	defer cancel()
	session, err := a.loadChannelSkillSession(processCtx, sessionID)
	if err != nil && !errors.Is(err, errRecordNotFound) {
		return
	}
	if err != nil {
		session = ChannelSkillSession{ID: sessionID, Channel: request.Channel, Account: request.Account, Sender: request.Sender, Conversation: request.Conversation, State: "new", Parameters: map[string]string{}, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	}
	// Keep the selected route visible in the persisted session so operators can
	// audit which Agent handled a channel conversation after a restart.
	session.AgentID = route.AgentID
	if session.ExpiresAt != "" {
		if expires, parseErr := time.Parse(time.RFC3339Nano, session.ExpiresAt); parseErr == nil && time.Now().After(expires) {
			session = ChannelSkillSession{ID: sessionID, Channel: request.Channel, Account: request.Account, Sender: request.Sender, Conversation: request.Conversation, State: "new", Parameters: map[string]string{}, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
		}
	}
	text := strings.TrimSpace(request.Message)
	if isChannelSkillCancel(text) {
		session.State = "cancelled"
		session.LastPrompt = "已取消当前功能。"
		a.saveChannelSkillSession(processCtx, session)
		a.replyChannelSkill(processCtx, request, "已取消当前功能。")
		return
	}
	if session.State == "executing" {
		a.replyChannelSkill(processCtx, request, "当前功能仍在执行，请等待结果后再发送新的指令。")
		return
	}
	if session.State == "completed" && !isChannelSkillRestart(text) {
		a.replyChannelSkill(processCtx, request, "上次功能已完成。如需再次执行，请回复“重新开始”。")
		return
	}
	if isChannelSkillRestart(text) {
		session.State, session.SkillID, session.CandidateIDs = "new", "", nil
		session.Parameters = map[string]string{}
	}
	if session.State == "selecting" {
		if choice, choiceErr := strconv.Atoi(strings.TrimSpace(text)); choiceErr == nil && choice > 0 && choice <= len(session.CandidateIDs) {
			session.SkillID = session.CandidateIDs[choice-1]
			session.CandidateIDs = nil
			session.State = "collecting"
		} else {
			a.replyChannelSkill(processCtx, request, session.LastPrompt)
			return
		}
	}
	if session.State == "awaiting_confirmation" {
		if isChannelSkillConfirm(text) {
			session.State = "executing"
			session.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
			a.saveChannelSkillSession(processCtx, session)
			a.replyChannelSkill(processCtx, request, "已确认，正在执行该功能，请稍候。")
			a.executeChannelSkill(processCtx, request, route, session)
			return
		}
		if isChannelSkillEdit(text) {
			session.State = "collecting"
			mergeChannelSkillParameters(session.Parameters, text)
		} else {
			a.replyChannelSkill(processCtx, request, "请回复“确认”执行，或回复“修改：具体修改内容”。")
			return
		}
	}
	if session.State == "new" {
		candidates := matchChannelSkills(text)
		if len(candidates) == 0 {
			session.State = "selecting"
			session.CandidateIDs = make([]string, len(channelSkillDefinitions))
			lines := []string{"请回复编号选择要执行的一个功能："}
			for i, skill := range channelSkillDefinitions {
				session.CandidateIDs[i] = skill.ID
				lines = append(lines, fmt.Sprintf("%d. %s", i+1, skill.Name))
			}
			session.LastPrompt = strings.Join(lines, "\n")
			session.OriginalMessage = text
			a.saveChannelSkillSession(processCtx, session)
			a.replyChannelSkill(processCtx, request, session.LastPrompt)
			return
		}
		session.SkillID = candidates[0]
		session.State = "collecting"
		session.OriginalMessage = text
	}
	if session.Parameters == nil {
		session.Parameters = map[string]string{}
	}
	mergeChannelSkillParameters(session.Parameters, text)
	skill, ok := findChannelSkill(session.SkillID)
	if !ok {
		session.State = "new"
		a.saveChannelSkillSession(processCtx, session)
		a.replyChannelSkill(processCtx, request, "暂不支持该功能，请重新发送指令。")
		return
	}
	if route.AgentID != skill.AgentID {
		session.State = "completed"
		session.Result = "该通道路由未授权执行此功能，请联系管理员配置对应 Agent。"
		a.saveChannelSkillSession(context.Background(), session)
		a.replyChannelSkill(context.Background(), request, session.Result)
		return
	}
	session.Missing = channelSkillMissing(skill, session.Parameters)
	session.ExpiresAt = time.Now().Add(channelSkillSessionTTL).UTC().Format(time.RFC3339Nano)
	if len(session.Missing) > 0 {
		session.State = "collecting"
		session.LastPrompt = "还需要：" + strings.Join(channelSkillLabels(session.Missing), "、") + "。请直接补充，例如：客户=xxx，产品=xxx。"
		a.saveChannelSkillSession(processCtx, session)
		a.replyChannelSkill(processCtx, request, session.LastPrompt)
		return
	}
	session.State = "awaiting_confirmation"
	session.LastPrompt = channelSkillConfirmation(skill, session.Parameters)
	a.saveChannelSkillSession(processCtx, session)
	a.replyChannelSkill(processCtx, request, session.LastPrompt)
}

func (a *businessAPI) loadChannelSkillSession(ctx context.Context, id string) (ChannelSkillSession, error) {
	var session ChannelSkillSession
	err := a.store.get(ctx, channelSkillSessionKind, id, &session)
	return session, err
}

func (a *businessAPI) saveChannelSkillSession(ctx context.Context, session ChannelSkillSession) error {
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if session.CreatedAt == "" {
		session.CreatedAt = session.UpdatedAt
	}
	if err := a.store.put(ctx, channelSkillSessionKind, session.ID, session); errors.Is(err, errRecordNotFound) {
		return a.store.create(ctx, channelSkillSessionKind, session.ID, session)
	} else {
		return err
	}
}

func (a *businessAPI) executeChannelSkill(ctx context.Context, request channelSkillInboundRequest, route ChannelSkillRoute, session ChannelSkillSession) {
	skill, ok := findChannelSkill(session.SkillID)
	if !ok {
		return
	}
	prompt := fmt.Sprintf("这是 STA-100 通道确认后执行的单一功能：%s。请严格根据已确认参数完成任务，只处理这个功能，不发起其它功能。参数：%s。请给出可直接发送给用户的简洁结果；事实不足时明确列出缺口，不要编造。", skill.Name, formatChannelSkillParameters(session.Parameters))
	result, err := a.sendAssistantAgent(ctx, skill.AgentID, "", "sta100-channel-"+channelSkillHash(request.Channel, request.Account, request.Conversation), prompt, nil, nil, nil)
	if err != nil {
		session.State = "completed"
		session.Result = "执行失败：" + err.Error()
		a.saveChannelSkillSession(context.Background(), session)
		a.replyChannelSkill(context.Background(), request, session.Result)
		return
	}
	session.State = "completed"
	session.RunID = result.RunID
	session.Result = strings.TrimSpace(result.Text)
	if session.Result == "" {
		session.Result = "功能已执行，但 Agent 未返回可展示内容。"
	}
	a.saveChannelSkillSession(context.Background(), session)
	a.replyChannelSkill(context.Background(), request, session.Result)
}

func (a *businessAPI) replyChannelSkill(ctx context.Context, request channelSkillInboundRequest, message string) {
	target := request.Conversation
	if strings.TrimSpace(target) == "" {
		target = request.Sender
	}
	if target == "" || message == "" || a.openClaw == nil || a.openClaw.service == nil {
		return
	}
	_, _ = a.openClaw.service.SendChannelMessage(ctx, orchestrator.ChannelMessageInput{Channel: request.Channel, Account: request.Account, Target: target, Message: message})
}

func matchChannelSkills(text string) []string {
	text = strings.ToLower(strings.TrimSpace(text))
	matched := []string{}
	for _, skill := range channelSkillDefinitions {
		for _, keyword := range skill.Keywords {
			if strings.Contains(text, strings.ToLower(keyword)) {
				matched = append(matched, skill.ID)
				break
			}
		}
	}
	return matched
}

func channelSkillMissing(skill ChannelSkillDefinition, params map[string]string) []string {
	missing := []string{}
	for _, key := range skill.Required {
		if strings.TrimSpace(params[key]) == "" {
			missing = append(missing, key)
		}
	}
	return missing
}

func channelSkillLabels(keys []string) []string {
	labels := map[string]string{"customer": "客户", "products": "产品和数量", "product": "产品", "supplier": "供应商", "query": "具体需求"}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if label := labels[key]; label != "" {
			out = append(out, label)
		} else {
			out = append(out, key)
		}
	}
	return out
}

func mergeChannelSkillParameters(params map[string]string, text string) {
	if params == nil {
		return
	}
	parts := strings.FieldsFunc(text, func(r rune) bool { return r == '\n' || r == ',' || r == '，' || r == ';' || r == '；' })
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.TrimSpace(strings.TrimPrefix(part, "修改"))
		part = strings.TrimSpace(strings.TrimPrefix(part, "更改"))
		part = strings.TrimSpace(strings.TrimPrefix(part, ":"))
		part = strings.TrimSpace(strings.TrimPrefix(part, "："))
		key, value, found := strings.Cut(part, "=")
		if !found {
			key, value, found = strings.Cut(part, "：")
		}
		if !found {
			continue
		}
		key = strings.TrimSpace(strings.ToLower(key))
		value = strings.TrimSpace(value)
		switch key {
		case "客户", "客户名称", "customer":
			params["customer"] = value
		case "产品", "产品和数量", "products", "product":
			params["products"] = value
		case "供应商", "supplier":
			params["supplier"] = value
		case "需求", "问题", "内容", "query", "市场", "展会":
			params["query"] = value
		}
	}
	if params["query"] == "" && len(text) > 0 && !strings.Contains(text, "确认") && !strings.Contains(text, "修改") {
		params["query"] = text
	}
}

func formatChannelSkillParameters(params map[string]string) string {
	keys := []string{"customer", "products", "supplier", "query"}
	parts := []string{}
	for _, key := range keys {
		if value := strings.TrimSpace(params[key]); value != "" {
			parts = append(parts, key+"="+value)
		}
	}
	return strings.Join(parts, "; ")
}

func channelSkillConfirmation(skill ChannelSkillDefinition, params map[string]string) string {
	return fmt.Sprintf("已识别：%s\n%s\n请回复“确认”执行，或回复“修改：字段=新值”。", skill.Name, formatChannelSkillParameters(params))
}

func isChannelSkillConfirm(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	return text == "确认" || text == "确定" || text == "confirm" || text == "yes"
}

func isChannelSkillCancel(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	return text == "取消" || text == "停止" || text == "cancel"
}

func isChannelSkillRestart(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	return text == "重新开始" || text == "重来" || text == "restart"
}

func isChannelSkillEdit(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	return strings.HasPrefix(text, "修改") || strings.HasPrefix(text, "更改") || strings.HasPrefix(text, "edit")
}
