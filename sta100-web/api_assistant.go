package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"openclaw-orchestrator/orchestrator"
)

const (
	coordinatorAgentID = "sta100-coordinator"
	knowledgeAgentID   = "sta100-knowledge"
)

var assistantSessionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,79}$`)

var assistantDomainAgentIDs = map[string]bool{
	"export-agent": true, "payment-advisor": true, "shipping-eta": true, "invoice-agent": true,
	"sinosure-advisor": true, "cbam-calculator": true, "email-generator": true, "price-tracker": true,
	"inventory-agent": true, "inventory-clearance-agent": true, "b2b-marketplace-agent": true, "used-bike-trading-agent": true,
	"repair-qa": true, "compatibility-agent": true, "market-analyzer": true, "country-advisor": true,
	"exhibition-advisor": true, "team-race-advisor": true, "supplier-aggregator": true, "brand-value-crawler": true,
	"customer-measurement-agent": true, "rag-agent": true, "design-advisor": true, "route-fetcher": true,
}

type assistantQueryRequest struct {
	Page        string                           `json:"page"`
	Feature     string                           `json:"feature"`
	Model       string                           `json:"model"`
	Message     string                           `json:"message"`
	SessionKey  string                           `json:"sessionKey"`
	Context     map[string]any                   `json:"context"`
	Attachments []orchestrator.MessageAttachment `json:"attachments"`
}

type assistantEvidence struct {
	ID        string `json:"id"`
	Entity    string `json:"entity"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Source    string `json:"source"`
}

type assistantAgentOutput struct {
	AgentID        string                  `json:"agentId"`
	Text           string                  `json:"text,omitempty"`
	RunID          string                  `json:"runId,omitempty"`
	Error          string                  `json:"error,omitempty"`
	DurationMs     int64                   `json:"durationMs,omitempty"`
	Usage          orchestrator.TokenUsage `json:"usage"`
	UsageAvailable bool                    `json:"usageAvailable"`
}

type assistantAttachmentView struct {
	Name string `json:"name"`
	Mime string `json:"mime,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type assistantPipelineStage struct {
	Stage      string `json:"stage"`
	Status     string `json:"status"`
	Detail     string `json:"detail,omitempty"`
	DurationMs int64  `json:"durationMs,omitempty"`
	Reason     string `json:"reason,omitempty"`
	Data       string `json:"data,omitempty"`
}

type assistantTiming struct {
	Stage      string `json:"stage"`
	Status     string `json:"status"`
	DurationMs int64  `json:"durationMs"`
	Reason     string `json:"reason"`
}

type assistantQueryResponse struct {
	Text            string                    `json:"text"`
	Items           []map[string]any          `json:"items"`
	UsedAgents      []string                  `json:"usedAgents"`
	Evidence        []assistantEvidence       `json:"evidence"`
	Conflicts       []map[string]any          `json:"conflicts"`
	Outputs         []assistantAgentOutput    `json:"agentOutputs,omitempty"`
	Pipeline        []assistantPipelineStage  `json:"pipeline"`
	Attachments     []assistantAttachmentView `json:"attachments,omitempty"`
	AIGenerated     bool                      `json:"aiGenerated"`
	Partial         bool                      `json:"partial"`
	Todo            []string                  `json:"todo,omitempty"`
	TokenUsage      tokenUsageSummary         `json:"tokenUsage"`
	Timings         []assistantTiming         `json:"timings,omitempty"`
	TotalDurationMs int64                     `json:"totalDurationMs"`
	SlowReason      string                    `json:"slowReason,omitempty"`
}

type assistantStructuredResult struct {
	Schema string            `json:"schema,omitempty"`
	Type   string            `json:"type"`
	Items  []json.RawMessage `json:"items,omitempty"`
}

func (a *businessAPI) assistantRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 1 && parts[0] == "attachments" {
		a.assistantAttachmentUpload(w, r)
		return
	}
	if len(parts) == 1 && parts[0] == "history" {
		a.assistantHistory(w, r)
		return
	}
	if len(parts) != 1 || parts[0] != "query" {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "统一智能查询接口不存在")
		return
	}
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request assistantQueryRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	request.Page = strings.TrimSpace(request.Page)
	request.Feature = strings.TrimSpace(request.Feature)
	request.Message = strings.TrimSpace(request.Message)
	request.SessionKey = strings.TrimSpace(request.SessionKey)
	if request.Message == "" || len(request.Message) > 32<<10 {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ASSISTANT_MESSAGE", "查询内容不能为空且不能超过 32 KiB")
		return
	}
	if request.SessionKey == "" {
		request.SessionKey = fmt.Sprintf("sta100-%s-%d", normalizedFeature(request.Feature), time.Now().Unix())
	}
	if !assistantSessionPattern.MatchString(request.SessionKey) {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ASSISTANT_SESSION", "会话标识只能包含字母、数字、点、下划线、冒号或连字符，且不能超过 80 个字符")
		return
	}
	if target, ok := request.Context["targetAgent"].(string); ok && strings.TrimSpace(target) != "" && !assistantDomainAgentIDs[strings.TrimSpace(target)] {
		writeAPIError(w, http.StatusBadRequest, "INVALID_TARGET_AGENT", "目标 Agent 必须是 STA-100 的专业业务 Agent")
		return
	}
	if request.Model != "" && !a.modelPreviouslyTestedOK(r.Context(), request.Model) {
		writeAPIError(w, http.StatusBadRequest, "MODEL_NOT_TESTED", "聊天只能选择已测试通过的模型")
		return
	}
	if err := validateAssistantAttachments(request.Attachments); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ATTACHMENT", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Minute)
	defer cancel()
	requestID := fmt.Sprintf("AI-%d", time.Now().UnixNano())
	response := a.runAssistantQuery(ctx, request, requestID)
	a.persistAssistantConversation(r.Context(), request, response)
	cleanupAssistantAttachments(request.Attachments)
	a.store.audit(r.Context(), "query", "assistant", request.Feature, requestOperator(r), map[string]any{
		"page": request.Page, "usedAgents": response.UsedAgents, "partial": response.Partial, "evidenceCount": len(response.Evidence),
	})
	writeJSON(w, http.StatusOK, response)
}

func (a *businessAPI) assistantHistory(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	query := r.URL.Query()
	agentID := strings.TrimSpace(strings.ToLower(query.Get("agentId")))
	sessionKey := strings.TrimSpace(query.Get("sessionKey"))
	if !assistantDomainAgentIDs[agentID] {
		writeAPIError(w, http.StatusBadRequest, "INVALID_TARGET_AGENT", "历史记录只能查询 STA-100 的专业业务 Agent")
		return
	}
	if sessionKey == "" {
		sessionKey = "sta100-" + agentID
	}
	if !assistantSessionPattern.MatchString(sessionKey) {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ASSISTANT_SESSION", "会话标识格式无效")
		return
	}
	limit := 100
	if value := strings.TrimSpace(query.Get("limit")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 500 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_HISTORY_LIMIT", "历史记录条数必须是 1-500 的整数")
			return
		}
		limit = parsed
	}
	messages := a.localAgentMessages(r.Context(), agentID, sessionKey, limit)
	writeJSON(w, http.StatusOK, map[string]any{
		"source":     "sta100-local",
		"agentId":    agentID,
		"sessionKey": sessionKey,
		"messages":   messages,
		"total":      len(messages),
		"hasMore":    false,
		"notice":     "历史记录来自 STA-100 本机记录，保留页面展示所需的阶段进度和耗时；OpenClaw 原始会话仍由底层维护。",
	})
}

func (a *businessAPI) persistAssistantConversation(ctx context.Context, request assistantQueryRequest, response assistantQueryResponse) {
	agentID, _ := request.Context["targetAgent"].(string)
	agentID = strings.TrimSpace(agentID)
	if agentID == "" || !assistantDomainAgentIDs[agentID] {
		return
	}
	now := time.Now().UTC()
	userCreatedAt := now.Format(time.RFC3339Nano)
	assistantCreatedAt := now.Add(1 * time.Nanosecond).Format(time.RFC3339Nano)
	userID, err := a.store.nextSequence(ctx, "agent_messages", "MSG", 8)
	if err != nil {
		return
	}
	userMessage := AgentMessageRecord{
		ID: userID, AgentID: agentID, SessionKey: request.SessionKey, Role: "user",
		Text: request.Message, Sources: []string{"本地业务数据库", "客户私有知识库", "联网检索"},
		Model: request.Model, CreatedAt: userCreatedAt,
	}
	_ = a.store.create(ctx, "agent_messages", userID, userMessage)
	assistantID, err := a.store.nextSequence(ctx, "agent_messages", "MSG", 8)
	if err != nil {
		return
	}
	assistantMessage := AgentMessageRecord{
		ID: assistantID, AgentID: agentID, SessionKey: request.SessionKey, Role: "assistant",
		Text: response.Text, Sources: response.UsedAgents, Pipeline: response.Pipeline,
		TotalDurationMs: response.TotalDurationMs, Model: request.Model,
		CreatedAt: assistantCreatedAt,
	}
	_ = a.store.create(ctx, "agent_messages", assistantID, assistantMessage)
}

func (a *businessAPI) localAgentMessages(ctx context.Context, agentID, sessionKey string, limit int) []AgentMessageRecord {
	items, err := listRecords[AgentMessageRecord](ctx, a.store, "agent_messages")
	if err != nil {
		return []AgentMessageRecord{}
	}
	filtered := make([]AgentMessageRecord, 0, len(items))
	for _, item := range items {
		if item.AgentID == agentID && item.SessionKey == sessionKey {
			filtered = append(filtered, item)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].CreatedAt == filtered[j].CreatedAt {
			if filtered[i].Role != filtered[j].Role {
				return filtered[i].Role == "user"
			}
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].CreatedAt < filtered[j].CreatedAt
	})
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}
	return filtered
}

func validateAssistantAttachments(attachments []orchestrator.MessageAttachment) error {
	if len(attachments) > 8 {
		return fmt.Errorf("单次最多上传 8 个附件")
	}
	root, err := filepath.Abs(filepath.Join(os.TempDir(), "sta100-agent-attachments"))
	if err != nil {
		return fmt.Errorf("附件目录不可用")
	}
	for _, attachment := range attachments {
		path, err := filepath.Abs(strings.TrimSpace(attachment.Path))
		if err != nil || (path != root && !strings.HasPrefix(path, root+string(os.PathSeparator))) {
			return fmt.Errorf("附件路径不在 STA-100 临时目录内")
		}
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			return fmt.Errorf("附件已不存在，请重新上传")
		}
		if info.Size() > 25<<20 {
			return fmt.Errorf("单个附件必须在 25 MB 以内")
		}
	}
	return nil
}

func cleanupAssistantAttachments(attachments []orchestrator.MessageAttachment) {
	for _, attachment := range attachments {
		_ = os.Remove(attachment.Path)
	}
}

func (a *businessAPI) assistantAttachmentUpload(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	if err := r.ParseMultipartForm(26 << 20); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ATTACHMENT", "附件格式无效或超过 25 MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "ATTACHMENT_REQUIRED", "请选择要上传的图片或文件")
		return
	}
	defer file.Close()
	if header.Size <= 0 || header.Size > 25<<20 {
		writeAPIError(w, http.StatusBadRequest, "ATTACHMENT_TOO_LARGE", "单个附件必须在 25 MB 以内")
		return
	}
	root := filepath.Join(os.TempDir(), "sta100-agent-attachments")
	if err := os.MkdirAll(root, 0o700); err != nil {
		writeBusinessError(w, err)
		return
	}
	target, err := os.CreateTemp(root, "attachment-*"+filepath.Ext(header.Filename))
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	path := target.Name()
	defer target.Close()
	if err := target.Chmod(0o600); err != nil {
		_ = os.Remove(path)
		writeBusinessError(w, err)
		return
	}
	written, err := io.Copy(target, io.LimitReader(file, 25<<20))
	if err != nil {
		_ = os.Remove(path)
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"attachment": orchestrator.MessageAttachment{
		Name: filepath.Base(header.Filename), Path: path, Mime: header.Header.Get("Content-Type"), Size: written,
	}})
}

func (a *businessAPI) runAssistantQuery(ctx context.Context, request assistantQueryRequest, requestID string) assistantQueryResponse {
	totalStarted := time.Now()
	evidence := a.collectLocalEvidence(ctx, request)
	sortAssistantEvidence(evidence)
	response := assistantQueryResponse{
		Evidence: evidence, Items: a.localAssistantItems(request, evidence), Conflicts: detectEvidenceConflicts(evidence),
		Attachments: assistantAttachmentViews(request.Attachments),
		AIGenerated: true, Partial: false,
		Todo: []string{"客户私有文件正文解析和向量索引等待原始数据格式", "联网工具接入后在工具层强制执行域名白名单"},
	}
	localDuration := time.Since(totalStarted).Milliseconds()
	if normalizedFeature(request.Feature) == "customer-discovery" {
		response.Todo = []string{"客户发现允许来源白名单、客户线索保存规则和联系方式校验规则待确认"}
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "query-fields", Status: "ok", Detail: "已接收国家、城市和客户类型筛选条件", DurationMs: localDuration, Reason: "客户本地发现不读取本机客户库，直接把页面条件交给 OpenClaw 客户发现 Agent", Data: contextSearchText(request.Context)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "query-fields", Status: "ok", DurationMs: localDuration, Reason: "客户本地发现不读取本机客户库，直接把页面条件交给 OpenClaw 客户发现 Agent"})
	} else {
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "local-retrieval", Status: "ok", Detail: fmt.Sprintf("检索到 %d 条本地结构化证据", len(evidence)), DurationMs: localDuration, Reason: "读取并筛选 STA-100 SQLite 本地业务数据", Data: assistantEvidenceHoverData(evidence)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "local-retrieval", Status: "ok", DurationMs: localDuration, Reason: "读取并筛选 STA-100 SQLite 本地业务数据"})
	}
	if len(request.Attachments) > 0 {
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "attachments", Status: "ok", Detail: fmt.Sprintf("%d 个附件已校验并提交给 OpenClaw Agent", len(request.Attachments)), Reason: "校验本地附件路径并随消息提交", Data: assistantAttachmentHoverData(request.Attachments)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "attachments", Status: "ok", DurationMs: 0, Reason: "校验本地附件路径并随消息提交"})
	}
	if targetAgent := assistantTargetAgent(request); targetAgent != "" {
		return a.runTargetAssistantAgentQuery(ctx, request, requestID, targetAgent, response, evidence, totalStarted)
	}

	evidenceJSON, _ := json.Marshal(evidence)
	knowledgePrompt := fmt.Sprintf("[页面] %s\n[功能] %s\n[用户输入]\n%s\n\n[Go 本地检索证据 JSON]\n%s\n\n请按工作区规则整理证据。不得联网；冲突值必须全部保留。", request.Page, request.Feature, request.Message, evidenceJSON)
	knowledgeStarted := time.Now()
	knowledge, err := a.sendAssistantAgent(ctx, knowledgeAgentID, request.Model, assistantStageSession(request.SessionKey, "knowledge"), knowledgePrompt, []string{"本地业务数据库", "客户私有知识库"}, nil, request.Attachments)
	knowledgeDuration := time.Since(knowledgeStarted).Milliseconds()
	if err != nil {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "failed", knowledge)
		response.Partial = true
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, assistantErrorText(err)))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "failed", Detail: assistantErrorText(err), DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 失败或等待超时", Data: assistantErrorText(err)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "knowledge-agent", Status: "failed", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 失败或等待超时"})
		knowledge.Text = "Knowledge Agent 未返回；Coordinator 将直接接收 Go 检索的原始本地证据。"
	} else {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "ok", knowledge)
		response.UsedAgents = append(response.UsedAgents, knowledgeAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "ok", Detail: "本地证据已整理", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 整理本地证据", Data: assistantTextHoverData(knowledge.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "knowledge-agent", Status: "ok", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 整理本地证据"})
	}

	domainIDs := assistantDomainRoutes(request)
	domainStarted := time.Now()
	domainOutputs := a.runDomainAgents(ctx, request, requestID, knowledge.Text, domainIDs)
	domainDuration := time.Since(domainStarted).Milliseconds()
	for _, output := range domainOutputs {
		response.Outputs = append(response.Outputs, output)
		if output.Error != "" {
			response.Partial = true
			continue
		}
		response.UsedAgents = append(response.UsedAgents, output.AgentID)
	}
	domainStatus := map[bool]string{true: "partial", false: "ok"}[hasAgentOutputError(domainOutputs)]
	response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "domain-agents", Status: domainStatus, Detail: strings.Join(domainIDs, ", "), DurationMs: domainDuration, Reason: "领域 Agent 并行真实调用，总耗时取最慢的一个", Data: assistantOutputsHoverData(domainOutputs)})
	response.Timings = append(response.Timings, assistantTiming{Stage: "domain-agents", Status: domainStatus, DurationMs: domainDuration, Reason: "领域 Agent 并行真实调用，总耗时取最慢的一个"})

	coordinatorInput := map[string]any{
		"page": request.Page, "feature": request.Feature, "userMessage": request.Message,
		"localEvidence": evidence, "knowledgeSummary": knowledge.Text, "domainOutputs": domainOutputs,
		"requiredRule": "冲突信息全部并列展示并保留来源、记录编号和更新时间；不要区分为本地结果区或联网结果区。",
	}
	coordinatorJSON, _ := json.Marshal(coordinatorInput)
	coordinatorPrompt := "请根据工作区规则整合以下 STA-100 查询上下文。只使用已提供内容，不得声称调用了未返回结果的 Agent。\n\n" + string(coordinatorJSON)
	if normalizedFeature(request.Feature) == "oem-match" {
		coordinatorPrompt += "\n\n" + assistantStructuredPrompt(assistantStructuredPromptSpec{
			ModuleName: "OEM 工厂智能匹配",
			Scope:      "用户输入的骑行类 OEM 需求、当前页面类别、排序方式和 Top 数量。",
			Fields:     []string{"title", "category", "reason", "detail", "source", "sourceUrl", "time", "score", "capacity", "moq"},
			ResultType: "oem_match",
			CountHint:  "输出前 N 条最匹配结果，不要输出概览壳子。",
			Notes: []string{
				"每条 item 必须是完整候选工厂或供应商结果；title 写工厂或供应商名称；reason 写命中的业务理由；detail 写可用于验收的完整说明。",
				"source 写来源名称；sourceUrl 写原文或来源链接，没有则留空；time 写更新时间；score 写匹配度或分数；capacity 写产能；moq 写最小起订量。",
				"没有可靠结果时 items 输出空数组，不要编造事实。",
			},
		})
	}
	coordinatorStarted := time.Now()
	coordinator, coordinatorErr := a.sendAssistantAgent(ctx, coordinatorAgentID, request.Model, assistantStageSession(request.SessionKey, "coordinator"), coordinatorPrompt, []string{"本地业务数据库", "客户私有知识库"}, nil, request.Attachments)
	coordinatorDuration := time.Since(coordinatorStarted).Milliseconds()
	if coordinatorErr != nil {
		coordinatorErrorText := assistantErrorText(coordinatorErr)
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "coordinator-agent", coordinatorAgentID, "failed", coordinator)
		response.Partial = true
		response.Text = localAssistantSummary(request, response.Items, evidence, domainOutputs, coordinatorErrorText)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(coordinatorAgentID, coordinator, coordinatorErrorText))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "coordinator-agent", Status: "failed", Detail: coordinatorErrorText, DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 失败或等待超时", Data: coordinatorErrorText})
		response.Timings = append(response.Timings, assistantTiming{Stage: "coordinator-agent", Status: "failed", DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 失败或等待超时"})
	} else {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "coordinator-agent", coordinatorAgentID, "ok", coordinator)
		response.Text = coordinator.Text
		if normalizedFeature(request.Feature) == "oem-match" {
			response.Items = parseAssistantStructuredItems("oem-match", coordinator.Text)
		}
		response.UsedAgents = append(response.UsedAgents, coordinatorAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(coordinatorAgentID, coordinator, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "coordinator-agent", Status: "ok", Detail: "完成统一汇总", DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 汇总前序结果", Data: assistantTextHoverData(coordinator.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "coordinator-agent", Status: "ok", DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 汇总前序结果"})
	}
	response.UsedAgents = uniqueStrings(response.UsedAgents)
	response.TotalDurationMs = time.Since(totalStarted).Milliseconds()
	response.SlowReason = "当前请求使用真实多 Agent 链路：先检索本地数据，再调用 Knowledge Agent、领域 Agent，最后由 Coordinator Agent 汇总；领域 Agent 虽然并行，但后续汇总必须等待前序结果。"
	if usage, err := a.tokenUsageSummary(ctx, requestID); err == nil {
		response.TokenUsage = usage
	}
	return response
}

func (a *businessAPI) runTargetAssistantAgentQuery(ctx context.Context, request assistantQueryRequest, requestID, agentID string, response assistantQueryResponse, evidence []assistantEvidence, totalStarted time.Time) assistantQueryResponse {
	evidenceJSON, _ := json.Marshal(evidence)
	prompt := fmt.Sprintf("[STA-100 当前智能体]\n%s\n\n[页面] %s\n[功能] %s\n[用户输入]\n%s\n\n[Go 本地检索证据 JSON]\n%s\n\n请基于当前智能体职责回答。必须先使用已提供的本地证据；如果证据不足，明确说明缺口和需要补充的数据。冲突信息必须并列保留来源、记录编号和更新时间。附件如已提交，请结合附件内容；不能确认附件内容时要说明。", agentID, request.Page, request.Feature, request.Message, evidenceJSON)
	preferences := defaultPreferences()
	_ = a.store.getSetting(ctx, "preferences", &preferences)
	if normalizedFeature(request.Feature) == "customer-discovery" {
		limit := discoveryResultLimit(request.Context, preferences)
		country := firstNonEmpty(stringContextValue(request.Context, "country"), "未填写")
		city := firstNonEmpty(stringContextValue(request.Context, "city"), "未填写")
		scope := fmt.Sprintf("国家=%s；城市=%s", country, city)
		if country == "全球" {
			scope = "国家和城市不限（全球范围）"
		}
		prompt = fmt.Sprintf("[STA-100 客户发现]\n当前目标 Agent：%s\n\n[固定输入]\n国家：%s\n城市：%s\n客户类型：%s\n用户输入：%s\n\n[执行边界]\n请不要检索或引用 STA-100 本机客户库。本功能只把页面固定筛选条件交给 OpenClaw 客户发现 Agent，由 Agent 按自身已配置能力获取公开来源并完成核验。\n\n[业务规则]\n- 国家、城市、客户类型都是硬性筛选条件，只返回同时满足三项条件的客户或机构线索。\n- 每条线索必须有可核验来源；没有来源、城市不确定或类型不确定的候选不要进入 items。\n- 不满足筛选条件的市场概览、供应商、产品、泛行业新闻不要作为客户线索返回。\n- 无法获取可靠客户线索时，明确说明原因，并输出空 items；不要用模型常识补客户，不要编造。\n- 如果公开信息存在冲突，在 reason 中保留冲突点，不要覆盖。\n%s", agentID, firstNonEmpty(stringContextValue(request.Context, "country"), "未填写"), firstNonEmpty(stringContextValue(request.Context, "city"), "未填写"), firstNonEmpty(stringContextValue(request.Context, "type"), "未填写"), request.Message, assistantStructuredPrompt(assistantStructuredPromptSpec{
			ModuleName: "本地客户发现",
			Scope:      fmt.Sprintf("%s；客户类型=%s；只返回符合筛选条件的客户线索。", scope, firstNonEmpty(stringContextValue(request.Context, "type"), "未填写")),
			Fields:     []string{"name", "country", "city", "type", "business", "contact", "source", "sourceUrl", "updatedAt", "score", "reason"},
			ResultType: "customer_discovery",
			CountHint:  fmt.Sprintf("最多输出 %d 条候选客户；items 数量不得超过该上限。只输出候选客户，不要输出市场概览、产品或泛行业新闻。", limit),
			Notes: []string{
				"name 写客户或机构名称；country、city、type 分别写国家、城市和客户类型；business 写业务方向；contact 写联系方式或联系方式缺口。",
				"source 写来源名称；sourceUrl 写原文或来源链接，没有则留空；updatedAt 写抓取或更新时间；score 写匹配分；reason 写匹配理由。",
				"没有可靠结果时 items 输出空数组，不要编造客户。",
			},
		}))
	}
	started := time.Now()
	sources := []string{"本地业务数据库", "客户私有知识库", "联网检索"}
	if normalizedFeature(request.Feature) == "customer-discovery" {
		sources = []string{"联网检索"}
	}
	result, err := a.sendAssistantAgent(ctx, agentID, request.Model, request.SessionKey, prompt, sources, preferences.AgentAllowlists[agentID], request.Attachments)
	duration := time.Since(started).Milliseconds()
	if err != nil {
		errorText := assistantErrorText(err)
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "target-agent", agentID, "failed", result)
		output := assistantOutputFromResult(agentID, result, errorText)
		output.DurationMs = duration
		response.Partial = true
		response.Text = localTargetAgentSummary(request, agentID, len(evidence), errorText)
		response.Outputs = append(response.Outputs, output)
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "target-agent", Status: "failed", Detail: errorText, DurationMs: duration, Reason: "真实调用当前 OpenClaw Agent 失败或等待超时", Data: errorText})
		response.Timings = append(response.Timings, assistantTiming{Stage: "target-agent", Status: "failed", DurationMs: duration, Reason: "真实调用当前 OpenClaw Agent 失败或等待超时"})
	} else {
		output := assistantOutputFromResult(agentID, result, "")
		output.DurationMs = duration
		response.Text = result.Text
		stageStatus := "ok"
		stageDetail := "当前智能体已返回"
		stageReason := "本地检索后直接调用当前 OpenClaw Agent"
		if normalizedFeature(request.Feature) == "customer-discovery" {
			response.Items = parseAssistantStructuredItems("customer_discovery", result.Text)
			response.Items = limitAssistantItems(response.Items, discoveryResultLimit(request.Context, preferences))
			if issue, issueText := customerDiscoveryAccessIssue(result.Text, len(response.Items)); issue {
				response.Partial = true
				response.Text = customerDiscoveryAccessIssueSummary(request, issueText)
				output.Error = issueText
				stageStatus = "failed"
				stageDetail = issueText
				stageReason = "OpenClaw 客户发现 Agent 已响应，但公开来源能力不可用，业务检索未完成"
			}
		}
		if output.Error == "" {
			response.UsedAgents = append(response.UsedAgents, agentID)
		}
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "target-agent", agentID, stageStatus, result)
		response.Outputs = append(response.Outputs, output)
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "target-agent", Status: stageStatus, Detail: stageDetail, DurationMs: duration, Reason: stageReason, Data: assistantTextHoverData(result.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "target-agent", Status: stageStatus, DurationMs: duration, Reason: stageReason})
	}
	response.UsedAgents = uniqueStrings(response.UsedAgents)
	response.TotalDurationMs = time.Since(totalStarted).Milliseconds()
	if usage, err := a.tokenUsageSummary(ctx, requestID); err == nil {
		response.TokenUsage = usage
	}
	return response
}

func customerDiscoveryAccessIssue(text string, itemCount int) (bool, string) {
	if itemCount > 0 {
		return false, ""
	}
	normalized := strings.ToLower(strings.TrimSpace(text))
	if normalized == "" {
		return false, ""
	}
	searchUnavailable := strings.Contains(normalized, "web_search") && (strings.Contains(normalized, "disabled") || strings.Contains(normalized, "no provider") || strings.Contains(normalized, "禁用") || strings.Contains(normalized, "无法发起"))
	fetchUnavailable := strings.Contains(normalized, "web_fetch") && (strings.Contains(normalized, "private/internal/special-use") || strings.Contains(normalized, "timeout") || strings.Contains(normalized, "timed out") || strings.Contains(normalized, "超时") || strings.Contains(normalized, "被拦截") || strings.Contains(normalized, "未取到"))
	noPublicAbility := strings.Contains(normalized, "公开来源") && (strings.Contains(normalized, "不可用") || strings.Contains(normalized, "未接通") || strings.Contains(normalized, "无法访问") || strings.Contains(normalized, "无法获取"))
	noPublicBody := strings.Contains(normalized, "未取到任何页面正文") || strings.Contains(normalized, "无法获取任何满足") || strings.Contains(normalized, "无法访问所选来源")
	noOwnDatabase := strings.Contains(normalized, "无自有客户库") || strings.Contains(normalized, "no customer") && strings.Contains(normalized, "database")
	if searchUnavailable || fetchUnavailable || noPublicAbility || (noPublicBody && noOwnDatabase) {
		return true, "OpenClaw 公开来源能力不可用，客户发现未完成"
	}
	return false, ""
}

func customerDiscoveryAccessIssueSummary(request assistantQueryRequest, issueText string) string {
	country := firstNonEmpty(stringContextValue(request.Context, "country"), "未填写")
	city := firstNonEmpty(stringContextValue(request.Context, "city"), "未填写")
	customerType := firstNonEmpty(stringContextValue(request.Context, "type"), "未填写")
	if strings.TrimSpace(issueText) == "" {
		issueText = "OpenClaw 公开来源能力不可用，客户发现未完成"
	}
	return fmt.Sprintf("%s。已提交筛选条件：国家=%s，城市=%s，客户类型=%s。当前结果不代表没有客户，而是 OpenClaw Agent 尚未获取到可核验公开来源；系统已按规则避免编造客户线索。请先确认 OpenClaw 客户发现 Agent 的公开来源能力后重试。", issueText, country, city, customerType)
}

func (a *businessAPI) sendAssistantAgent(ctx context.Context, agentID, model, sessionKey, message string, sources, allowlist []string, attachments []orchestrator.MessageAttachment) (orchestrator.AgentMessageResult, error) {
	return a.openClaw.service.SendAgentMessage(ctx, orchestrator.AgentMessageInput{AgentID: agentID, Model: model, Message: message, SessionKey: sessionKey, Sources: sources, Allowlist: allowlist, Attachments: attachments})
}

type assistantStructuredPromptSpec struct {
	ModuleName string
	Scope      string
	Fields     []string
	ResultType string
	CountHint  string
	Notes      []string
}

func assistantStructuredPrompt(spec assistantStructuredPromptSpec) string {
	lines := []string{fmt.Sprintf("\n\n请按 STA-100 结构化结果协议输出%s。Go 侧只解析固定 JSON，不再从自由文本中猜业务字段。", spec.ModuleName)}
	if strings.TrimSpace(spec.Scope) != "" {
		lines = append(lines, "输入边界："+spec.Scope)
	}
	if len(spec.Fields) > 0 {
		lines = append(lines, "输出字段："+strings.Join(spec.Fields, "、"))
	}
	if strings.TrimSpace(spec.CountHint) != "" {
		lines = append(lines, spec.CountHint)
	}
	lines = append(lines, spec.Notes...)
	lines = append(lines, fmt.Sprintf("最终结果末尾必须只输出一个机器可读结果块，格式必须为：[STA100_RESULT]{\"schema\":\"%s\",\"type\":\"%s\",\"items\":[]}[/STA100_RESULT]。", structuredBusinessResultSchema, spec.ResultType))
	return strings.Join(lines, "\n\n")
}

func parseAssistantStructuredItems(expectedType, text string) []map[string]any {
	result, ok := parseAssistantStructuredResult(text)
	if !ok || normalizedStructuredType(result.Type) != normalizedStructuredType(expectedType) {
		return nil
	}
	items := make([]map[string]any, 0, len(result.Items))
	for _, raw := range result.Items {
		var item map[string]any
		if err := json.Unmarshal(raw, &item); err != nil || len(item) == 0 {
			continue
		}
		normalized := normalizeAssistantStructuredItem(expectedType, item)
		if len(normalized) > 0 {
			items = append(items, normalized)
		}
	}
	return items
}

func parseAssistantStructuredResult(text string) (assistantStructuredResult, bool) {
	for {
		start := strings.Index(text, "[STA100_RESULT]")
		if start < 0 {
			break
		}
		afterStart := start + len("[STA100_RESULT]")
		end := strings.Index(text[afterStart:], "[/STA100_RESULT")
		if end < 0 {
			break
		}
		end += afterStart
		payload := strings.TrimSpace(text[afterStart:end])
		var result assistantStructuredResult
		if err := json.Unmarshal([]byte(payload), &result); err == nil && strings.TrimSpace(result.Type) != "" {
			return result, true
		}
		text = text[end:]
	}
	return assistantStructuredResult{}, false
}

func normalizeAssistantStructuredItem(expectedType string, item map[string]any) map[string]any {
	get := func(keys ...string) string {
		for _, key := range keys {
			if value := strings.TrimSpace(fmt.Sprint(item[key])); value != "" && value != "<nil>" {
				return value
			}
		}
		return ""
	}
	switch normalizedStructuredType(expectedType) {
	case "customer_discovery":
		name := get("name", "title", "customer", "company")
		if name == "" {
			return nil
		}
		return map[string]any{
			"name":      name,
			"country":   get("country"),
			"city":      get("city"),
			"type":      get("type", "customerType"),
			"business":  get("business", "description", "detail"),
			"contact":   get("contact", "contacts"),
			"source":    get("source"),
			"sourceUrl": get("sourceUrl", "url", "link"),
			"updatedAt": get("updatedAt", "time"),
			"score":     numericStructuredValue(item["score"], 0),
			"reason":    get("reason", "why"),
		}
	case "oem_match":
		title := get("title", "name", "company", "factory")
		if title == "" {
			return nil
		}
		return map[string]any{
			"title":     title,
			"name":      title,
			"category":  get("category"),
			"reason":    get("reason", "why"),
			"detail":    get("detail", "content", "description"),
			"source":    get("source"),
			"sourceUrl": get("sourceUrl", "url", "link"),
			"time":      get("time", "updatedAt"),
			"score":     numericStructuredValue(item["score"], 0),
			"capacity":  get("capacity"),
			"moq":       get("moq", "MOQ"),
		}
	default:
		return item
	}
}

func normalizedStructuredType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "_")
	return value
}

func numericStructuredValue(value any, fallback int) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		if parsed, err := typed.Int64(); err == nil {
			return int(parsed)
		}
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(typed)); err == nil {
			return parsed
		}
	}
	return fallback
}

func assistantAttachmentViews(attachments []orchestrator.MessageAttachment) []assistantAttachmentView {
	views := make([]assistantAttachmentView, 0, len(attachments))
	for _, attachment := range attachments {
		views = append(views, assistantAttachmentView{Name: attachment.Name, Mime: attachment.Mime, Size: attachment.Size})
	}
	return views
}

func (a *businessAPI) runDomainAgents(ctx context.Context, request assistantQueryRequest, requestID, knowledgeSummary string, agentIDs []string) []assistantAgentOutput {
	outputs := make([]assistantAgentOutput, len(agentIDs))
	preferences := defaultPreferences()
	_ = a.store.getSetting(ctx, "preferences", &preferences)
	var group sync.WaitGroup
	for index, agentID := range agentIDs {
		group.Add(1)
		go func(index int, agentID string) {
			defer group.Done()
			started := time.Now()
			prompt := fmt.Sprintf("[STA-100 页面] %s / %s\n[用户输入]\n%s\n\n[Knowledge Agent 本地证据摘要]\n%s\n\n请只从你的专业职责给出结果，保留证据和不确定性。冲突数据不要覆盖。", request.Page, request.Feature, request.Message, knowledgeSummary)
			result, err := a.sendAssistantAgent(ctx, agentID, request.Model, assistantStageSession(request.SessionKey, agentID), prompt, []string{"本地业务数据库", "客户私有知识库", "联网检索"}, preferences.AgentAllowlists[agentID], request.Attachments)
			if err != nil {
				a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "domain-agent", agentID, "failed", result)
				output := assistantOutputFromResult(agentID, result, assistantErrorText(err))
				output.DurationMs = time.Since(started).Milliseconds()
				outputs[index] = output
				return
			}
			a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "domain-agent", agentID, "ok", result)
			output := assistantOutputFromResult(agentID, result, "")
			output.DurationMs = time.Since(started).Milliseconds()
			outputs[index] = output
		}(index, agentID)
	}
	group.Wait()
	return outputs
}

func assistantOutputFromResult(agentID string, result orchestrator.AgentMessageResult, errorText string) assistantAgentOutput {
	return assistantAgentOutput{
		AgentID: agentID, Text: result.Text, RunID: result.RunID, Error: errorText,
		Usage: result.Usage, UsageAvailable: result.Usage.Available(),
	}
}

func assistantEvidenceHoverData(items []assistantEvidence) string {
	if len(items) == 0 {
		return "本地检索未命中可展示证据。"
	}
	lines := make([]string, 0, minInt(len(items), 5))
	for index, item := range items {
		if index >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("%s｜%s｜%s", item.ID, item.Title, trimHoverText(item.Content, 90)))
	}
	if len(items) > 5 {
		lines = append(lines, fmt.Sprintf("另有 %d 条证据未在悬浮摘要中展示。", len(items)-5))
	}
	return strings.Join(lines, "\n")
}

func assistantAttachmentHoverData(items []orchestrator.MessageAttachment) string {
	if len(items) == 0 {
		return "本次未提交附件。"
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("%s｜%s｜%d bytes", item.Name, item.Mime, item.Size))
	}
	return strings.Join(lines, "\n")
}

func assistantOutputsHoverData(outputs []assistantAgentOutput) string {
	if len(outputs) == 0 {
		return "没有领域 Agent 返回结果。"
	}
	lines := make([]string, 0, len(outputs))
	for _, output := range outputs {
		text := output.Text
		if output.Error != "" {
			text = output.Error
		}
		lines = append(lines, fmt.Sprintf("%s｜%s", output.AgentID, trimHoverText(text, 140)))
	}
	return strings.Join(lines, "\n")
}

func assistantTextHoverData(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "该阶段未返回可展示文本。"
	}
	return trimHoverText(value, 260)
}

func trimHoverText(value string, limit int) string {
	text := strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if limit <= 0 || len([]rune(text)) <= limit {
		return text
	}
	return string([]rune(text)[:limit]) + "..."
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func assistantDomainRoutes(request assistantQueryRequest) []string {
	if target := assistantTargetAgent(request); target != "" {
		return []string{target}
	}
	switch normalizedFeature(request.Feature) {
	case "oem-match", "oem":
		return []string{"supplier-aggregator", "brand-value-crawler"}
	case "customer-discovery":
		return []string{"customer-measurement-agent"}
	case "customer-search", "customers":
		return []string{"customer-measurement-agent", "market-analyzer"}
	case "quote", "quotes":
		return []string{"export-agent", "payment-advisor"}
	case "order", "orders", "logistics":
		return []string{"shipping-eta", "inventory-agent"}
	case "document", "documents":
		return []string{"invoice-agent", "export-agent"}
	case "product", "products":
		return []string{"compatibility-agent", "design-advisor"}
	case "news":
		return []string{"sta100-news-curator", "country-advisor"}
	default:
		return []string{"market-analyzer"}
	}
}

func assistantTargetAgent(request assistantQueryRequest) string {
	target, _ := request.Context["targetAgent"].(string)
	target = strings.TrimSpace(target)
	if assistantDomainAgentIDs[target] {
		return target
	}
	return ""
}

func (a *businessAPI) collectLocalEvidence(ctx context.Context, request assistantQueryRequest) []assistantEvidence {
	if normalizedFeature(request.Feature) == "customer-discovery" {
		return nil
	}
	query := strings.ToLower(strings.TrimSpace(request.Message + " " + contextSearchText(request.Context)))
	tokens := searchTokens(query)
	evidence := make([]assistantEvidence, 0, 30)
	appendMatch := func(item assistantEvidence, searchable string) {
		if len(evidence) >= 30 || !matchesSearchTokens(strings.ToLower(searchable), tokens) {
			return
		}
		evidence = append(evidence, item)
	}
	customers, _ := a.listCustomers(ctx)
	for _, item := range customers {
		if item.Archived {
			continue
		}
		content := fmt.Sprintf("类型=%s；国家=%s；城市=%s；联系人=%s；电话=%s；邮箱=%s；网站=%s；来源=%s；订单数=%d；累计金额=%s；描述=%s", item.Type, item.Country, item.City, item.Contact, item.Phone, item.Email, item.Website, item.Source, item.Orders, item.Total, item.Description)
		appendMatch(assistantEvidence{ID: item.ID, Entity: "customer", Title: item.Name, Content: content, UpdatedAt: item.Updated, Source: "STA-100 SQLite"}, item.ID+" "+item.Name+" "+content)
	}
	suppliers, _ := listRecords[Supplier](ctx, a.store, "suppliers")
	for _, item := range suppliers {
		if item.Archived {
			continue
		}
		content := fmt.Sprintf("产品=%s；规格=%s；报价=%s；联系人=%s；电话=%s；邮箱=%s；来源=%s；备注=%s", item.Product, item.Specification, item.Quote, item.Contact, item.Phone, item.Email, item.Source, item.Notes)
		appendMatch(assistantEvidence{ID: item.ID, Entity: "supplier", Title: item.Company, Content: content, UpdatedAt: item.Updated, Source: "STA-100 SQLite"}, item.ID+" "+item.Company+" "+content)
	}
	products, _ := listRecords[Product](ctx, a.store, "products")
	for _, item := range products {
		content := fmt.Sprintf("类别=%s；制造商=%s；价格=%s；库存=%d；HS=%s；状态=%s；描述=%s", item.Category, item.Manufacturer, item.Price, item.Stock, item.HS, item.Status, item.Description)
		appendMatch(assistantEvidence{ID: item.ID, Entity: "product", Title: item.Name, Content: content, UpdatedAt: item.Updated, Source: "STA-100 SQLite"}, item.ID+" "+item.Name+" "+content)
	}
	files, _ := listRecords[PrivateFile](ctx, a.store, "private_files")
	for _, item := range files {
		content := fmt.Sprintf("分类=%s；标签=%s；文件来源=%s；索引状态=%s；当前仅使用元数据，正文解析待实现", item.Category, strings.Join(item.Tags, "、"), item.Source, item.Status)
		appendMatch(assistantEvidence{ID: item.ID, Entity: "private-file-metadata", Title: item.Name, Content: content, UpdatedAt: item.Updated, Source: "STA-100 私有文件元数据"}, item.ID+" "+item.Name+" "+content)
	}
	if len(evidence) == 0 {
		for _, item := range customers {
			if !item.Archived && len(evidence) < 5 {
				evidence = append(evidence, assistantEvidence{ID: item.ID, Entity: "customer", Title: item.Name, Content: fmt.Sprintf("类型=%s；国家=%s；联系人=%s；电话=%s；邮箱=%s", item.Type, item.Country, item.Contact, item.Phone, item.Email), UpdatedAt: item.Updated, Source: "STA-100 SQLite"})
			}
		}
	}
	return evidence
}

func stringContextValue(context map[string]any, key string) string {
	value, _ := context[key].(string)
	return value
}

func matchesDiscoveryField(actual, expected string) bool {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(actual), expected)
}

func (a *businessAPI) localAssistantItems(request assistantQueryRequest, evidence []assistantEvidence) []map[string]any {
	items := make([]map[string]any, 0)
	hasContact, _ := request.Context["hasContact"].(bool)
	feature := normalizedFeature(request.Feature)
	if feature != "customer-search" {
		return items
	}
	customers, _ := a.listCustomers(context.Background())
	discoveryCountry, _ := request.Context["country"].(string)
	discoveryCity, _ := request.Context["city"].(string)
	discoveryType, _ := request.Context["type"].(string)
	for _, item := range customers {
		if item.Archived {
			continue
		}
		if feature == "customer-discovery" && (!matchesDiscoveryField(item.Country, discoveryCountry) || !matchesDiscoveryField(item.City, discoveryCity) || !matchesDiscoveryField(item.Type, discoveryType)) {
			continue
		}
		contact := firstNonEmpty(item.Email, item.Phone, item.Website)
		if hasContact && contact == "" {
			continue
		}
		if !evidenceContainsID(evidence, item.ID) {
			continue
		}
		items = append(items, map[string]any{"id": item.ID, "name": item.Name, "country": item.Country, "city": item.City, "type": item.Type, "contact": contact, "business": item.Description, "score": 100, "integrated": true})
	}
	return items
}

func detectEvidenceConflicts(evidence []assistantEvidence) []map[string]any {
	grouped := map[string][]assistantEvidence{}
	for _, item := range evidence {
		key := strings.ToLower(strings.TrimSpace(item.Entity + ":" + item.Title))
		grouped[key] = append(grouped[key], item)
	}
	conflicts := make([]map[string]any, 0)
	for key, items := range grouped {
		values := map[string]bool{}
		for _, item := range items {
			values[item.Content] = true
		}
		if len(values) > 1 {
			conflicts = append(conflicts, map[string]any{"key": key, "evidence": items, "rule": "并列展示，不自动覆盖"})
		}
	}
	return conflicts
}

func localAssistantSummary(request assistantQueryRequest, items []map[string]any, evidence []assistantEvidence, outputs []assistantAgentOutput, coordinatorError string) string {
	successful := 0
	failed := make([]string, 0)
	for _, output := range outputs {
		if output.Error == "" {
			successful++
		} else {
			failed = append(failed, fmt.Sprintf("%s：%s", output.AgentID, output.Error))
		}
	}
	parts := []string{fmt.Sprintf("已从本机业务数据库检索到 %d 条证据和 %d 条可展示记录，%d 个领域 Agent 返回了结果。", len(evidence), len(items), successful)}
	if len(failed) > 0 {
		parts = append(parts, fmt.Sprintf("%d 个领域 Agent 未完成，具体原因可在消息下方阶段结果中查看。", len(failed)))
	}
	if strings.TrimSpace(coordinatorError) != "" {
		parts = append(parts, "最终协调 Agent 未完成，具体原因可在消息下方阶段结果中查看。")
	} else {
		parts = append(parts, "最终协调 Agent 未完成，当前为部分结果。")
	}
	parts = append(parts, "请结合证据记录编号和更新时间人工复核。")
	return strings.Join(parts, "")
}

func localTargetAgentSummary(request assistantQueryRequest, agentID string, evidenceCount int, errorText string) string {
	if normalizedFeature(request.Feature) == "customer-discovery" {
		parts := []string{fmt.Sprintf("已将国家、城市和客户类型提交给 %s 搜索，当前未完成回复。", agentID)}
		if strings.TrimSpace(errorText) != "" {
			parts = append(parts, "具体原因可在消息下方阶段结果中查看。")
		}
		return strings.Join(parts, "")
	}
	parts := []string{fmt.Sprintf("已从本机业务数据库检索到 %d 条证据，当前智能体 %s 未完成回复。", evidenceCount, agentID)}
	if strings.TrimSpace(errorText) != "" {
		parts = append(parts, "具体原因可在消息下方阶段结果中查看。")
	}
	parts = append(parts, "请结合证据记录编号和更新时间人工复核。")
	return strings.Join(parts, "")
}

func normalizedFeature(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	if value == "" {
		return "general"
	}
	return value
}

func assistantStageSession(base, stage string) string {
	value := base + "-" + stage
	if len(value) <= 80 {
		return value
	}
	return value[:80]
}

func contextSearchText(values map[string]any) string {
	data, _ := json.Marshal(values)
	return string(data)
}

func searchTokens(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ' ' || r == ',' || r == '，' || r == '、' || r == ';' || r == '；' || r == ':' || r == '：'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len([]rune(part)) > 1 {
			result = append(result, part)
		}
	}
	return uniqueStrings(result)
}

func matchesSearchTokens(value string, tokens []string) bool {
	if len(tokens) == 0 {
		return true
	}
	for _, token := range tokens {
		if strings.Contains(value, token) {
			return true
		}
	}
	return false
}

func evidenceContainsID(evidence []assistantEvidence, id string) bool {
	for _, item := range evidence {
		if item.ID == id {
			return true
		}
	}
	return false
}

func hasAgentOutputError(outputs []assistantAgentOutput) bool {
	for _, output := range outputs {
		if output.Error != "" {
			return true
		}
	}
	return false
}

func assistantErrorText(err error) string {
	if err == nil {
		return ""
	}
	detail := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "OpenClaw 调用超时"
	case strings.Contains(detail, "agent unavailable"):
		return "Agent 尚未同步到 OpenClaw"
	case strings.Contains(detail, "selected model was not found"), strings.Contains(detail, "model was not found"), strings.Contains(detail, "model not found"), strings.Contains(detail, "no access to model"):
		return "所选模型在供应商侧不存在或当前 API Key 无权使用，请在设置中更换已测试通过的模型"
	case strings.Contains(detail, "authentication"), strings.Contains(detail, "invalid api key"), strings.Contains(detail, "unauthorized"):
		return "模型提供商鉴权失败，请检查 API Key 是否有效"
	case strings.Contains(detail, "insufficient balance"), strings.Contains(detail, "quota exceeded"), strings.Contains(detail, "余额不足"):
		return "模型提供商账户余额或额度不足"
	case strings.Contains(detail, "rate limit"), strings.Contains(detail, "too many requests"), strings.Contains(detail, "限流"):
		return "模型提供商触发限流，请稍后重试"
	case strings.Contains(detail, "failovererror"):
		return "OpenClaw 模型调用失败，请检查默认模型、候选模型和供应商配置"
	default:
		return "OpenClaw 调用失败：" + strings.TrimSpace(err.Error())
	}
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func sortAssistantEvidence(items []assistantEvidence) {
	sort.SliceStable(items, func(i, j int) bool { return items[i].UpdatedAt > items[j].UpdatedAt })
}
