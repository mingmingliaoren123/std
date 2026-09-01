package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
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
	"oem-match-agent": true, "customer-measurement-agent": true, "rag-agent": true, "design-advisor": true, "route-fetcher": true,
}

type assistantQueue struct {
	limit   int
	mu      sync.Mutex
	cond    *sync.Cond
	nextSeq int64
	running []assistantQueueEntry
	waiting []assistantQueueEntry
}

type assistantQueueEntry struct {
	seq       int64
	agentID   string
	agentName string
	feature   string
	started   time.Time
}

func newAssistantQueue(limit int) *assistantQueue {
	if limit < 1 {
		limit = 1
	}
	q := &assistantQueue{limit: limit}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *assistantQueue) status() assistantQueueInfo {
	if q == nil {
		return assistantQueueInfo{Limit: 0}
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.statusLocked(nil, false, 0)
}

func (q *assistantQueue) acquire(ctx context.Context, agentID, feature string) (assistantQueueInfo, func(), error) {
	if q == nil {
		return assistantQueueInfo{}, func() {}, nil
	}
	entry := assistantQueueEntry{agentID: agentID, agentName: assistantAgentDisplayName(agentID), feature: feature}
	queuedAt := time.Now()
	queued := false
	var blockers []assistantQueueAgent
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			q.mu.Lock()
			q.cond.Broadcast()
			q.mu.Unlock()
		case <-done:
		}
	}()
	q.mu.Lock()
	defer close(done)
	q.nextSeq++
	entry.seq = q.nextSeq
	if len(q.running) >= q.limit {
		q.waiting = append(q.waiting, entry)
		queued = true
		blockers = assistantQueueAgents(q.running)
	}
	for {
		if ctx.Err() != nil {
			q.waiting = removeAssistantQueueEntry(q.waiting, entry.seq)
			q.cond.Broadcast()
			q.mu.Unlock()
			return assistantQueueInfo{}, nil, ctx.Err()
		}
		position := q.waitingPositionLocked(entry.seq)
		canRun := len(q.running) < q.limit && position <= 1
		if canRun {
			if position == 1 {
				q.waiting = removeAssistantQueueEntry(q.waiting, entry.seq)
			}
			entry.started = time.Now()
			waitedMs := time.Since(queuedAt).Milliseconds()
			q.running = append(q.running, entry)
			info := q.statusLocked(&entry, queued, waitedMs)
			if queued {
				info.RunningNames = nil
				for _, running := range blockers {
					info.RunningNames = append(info.RunningNames, running.Name)
				}
				info.Message = assistantQueueMessage(entry, blockers, waitedMs)
			}
			release := func() {
				q.mu.Lock()
				q.running = removeAssistantQueueEntry(q.running, entry.seq)
				q.cond.Broadcast()
				q.mu.Unlock()
			}
			q.mu.Unlock()
			return info, release, nil
		}
		q.cond.Wait()
	}
}

func (q *assistantQueue) statusLocked(active *assistantQueueEntry, queued bool, waitedMs int64) assistantQueueInfo {
	info := assistantQueueInfo{
		Queued: queued, Limit: q.limit, QueueLength: len(q.waiting), WaitedMs: waitedMs,
		Running: assistantQueueAgents(q.running),
	}
	if active != nil {
		info.AgentID = active.agentID
		info.AgentName = active.agentName
		if queued {
			info.Message = assistantQueueMessage(*active, info.Running, waitedMs)
		}
	}
	for _, running := range info.Running {
		info.RunningNames = append(info.RunningNames, running.Name)
	}
	return info
}

func (q *assistantQueue) waitingPositionLocked(seq int64) int {
	for index, entry := range q.waiting {
		if entry.seq == seq {
			return index + 1
		}
	}
	return 0
}

func assistantQueueAgents(entries []assistantQueueEntry) []assistantQueueAgent {
	agents := make([]assistantQueueAgent, 0, len(entries))
	for _, entry := range entries {
		agents = append(agents, assistantQueueAgent{
			AgentID: entry.agentID,
			Name:    entry.agentName,
			Feature: entry.feature,
			Started: entry.started.Format(time.RFC3339),
		})
	}
	return agents
}

func removeAssistantQueueEntry(entries []assistantQueueEntry, seq int64) []assistantQueueEntry {
	out := entries[:0]
	for _, entry := range entries {
		if entry.seq != seq {
			out = append(out, entry)
		}
	}
	return out
}

func assistantQueueMessage(active assistantQueueEntry, running []assistantQueueAgent, waitedMs int64) string {
	names := make([]string, 0, len(running))
	for _, item := range running {
		if item.Name != "" {
			names = append(names, item.Name)
		}
	}
	if len(names) == 0 {
		return fmt.Sprintf("%s 已进入排队队列，等待可用 Agent 执行位。", active.agentName)
	}
	return fmt.Sprintf("前面已有 %s 正在运行，%s 已进入排队队列，等待约 %s 后开始执行。", strings.Join(names, "、"), active.agentName, formatAssistantQueueWait(waitedMs))
}

func formatAssistantQueueWait(ms int64) string {
	if ms <= 0 {
		return "0 秒"
	}
	seconds := int64(math.Ceil(float64(ms) / 1000))
	if seconds < 60 {
		return fmt.Sprintf("%d 秒", seconds)
	}
	return fmt.Sprintf("%d 分 %d 秒", seconds/60, seconds%60)
}

type assistantAgentProfile struct {
	Name           string
	Focus          string
	KnowledgeScope string
	Boundaries     string
}

var assistantAgentProfiles = map[string]assistantAgentProfile{
	"export-agent":               {"出口业务助手", "出口流程、贸易条款、报价检查、出口风险提醒。", "出口 SOP、贸易术语、外贸合同、报关与交付资料。", "不要替代法务或报关行作最终判断；没有目的国、贸易条款或货物信息时先列缺口。"},
	"payment-advisor":            {"支付条款助手", "收款方式、账期、信用证、赊销和付款风险。", "付款条款模板、信用证条款、历史收款规则、中信保和客户信用资料。", "不直接承诺授信或放账；缺少客户信用和订单金额时必须提示待核实。"},
	"shipping-eta":               {"物流路线助手", "国际物流路线、时效、港口、运输方式和异常提醒。", "物流报价、航线时效、港口限制、交付计划和订单物流文件。", "不要编造实时船期；没有承运商或起止地时只能给方案框架。"},
	"invoice-agent":              {"外贸单据助手", "PI、CI、PL、报关单、合同等单据字段和一致性。", "单据模板、订单、报价、产品和客户抬头资料。", "不要生成缺少关键字段的正式单据；字段缺口必须逐项列出。"},
	"sinosure-advisor":           {"中信保助手", "中信保额度、买方风险、投保资料和出运前检查。", "中信保额度、客户档案、订单金额、国家风险和历史赔付资料。", "不能替代中信保官方审批；没有买方主体和金额时只输出待确认项。"},
	"cbam-calculator":            {"合规计算助手", "欧盟 CBAM、电池法规、碳足迹和合规资料清单。", "法规文件、HS CODE、产品材料、供应商合规证书和测试报告。", "不要输出无法追溯的合规结论；法规依据不足时标记不确定。"},
	"email-generator":            {"商务邮件助手", "外贸邮件、客户跟进、催款、报价和售后沟通。", "邮件模板、客户语言偏好、历史沟通记录和业务上下文。", "不要编造承诺、价格或交期；缺少事实时使用占位待确认。"},
	"price-tracker":              {"价格监测助手", "价格趋势、竞品价格、报价对比和异常波动。", "产品报价、竞品价格来源、历史报价、汇率和成本资料。", "没有可靠价格来源时不能给精确市场价。"},
	"inventory-agent":            {"库存管理助手", "库存状态、可售数量、备货建议和订单占用。", "库存台账、订单占用、采购计划、产品主数据。", "当前不直接读取业务表时，只能基于知识库证据回答。"},
	"inventory-clearance-agent":  {"库存出清助手", "滞销库存、清仓策略、组合销售和渠道建议。", "库存龄、成本、历史销售、渠道价格和促销政策。", "不要凭空给折扣底线；没有成本和库存龄时提示缺口。"},
	"b2b-marketplace-agent":      {"B2B 平台助手", "平台发布、询盘处理、关键词、产品页和线索运营。", "平台规则、产品资料、关键词库、询盘记录和客户画像。", "不编造平台实时排名；没有平台名时给通用运营建议。"},
	"used-bike-trading-agent":    {"二手交易助手", "二手车/配件交易、估价、翻新和风险检查。", "二手交易记录、检测报告、翻新成本和渠道规则。", "没有车况和来源证明时不要给确定估价。"},
	"repair-qa":                  {"维修诊断助手", "维修问答、故障诊断、备件建议和售后排查。", "产品手册、维修 SOP、故障案例、备件兼容资料。", "涉及安全风险时必须提示人工复核和停止使用。"},
	"compatibility-agent":        {"产品兼容助手", "零部件兼容性、规格匹配、替代料和装配风险。", "产品规格、BOM、兼容表、供应商规格书和测试资料。", "没有型号、规格或接口参数时不能判断兼容。"},
	"market-analyzer":            {"市场分析助手", "市场趋势、渠道结构、竞争态势和机会判断。", "行业新闻、市场报告、客户画像、销售和渠道资料。", "区分事实和推断；没有来源时不要给确定趋势。"},
	"country-advisor":            {"国家进入助手", "目标国家进入策略、法规、渠道、付款和物流风险。", "国家法规、税费、渠道资料、展会、客户和物流知识库。", "不要默认国家；用户未指定国家时先要求确认。"},
	"exhibition-advisor":         {"展会助手", "展会筛选、展商/观众线索、参展计划和跟进动作。", "展会名录、历史参展记录、客户线索、产品资料和预算。", "不编造展商名单；缺少展会名称或年份时提示补充。"},
	"team-race-advisor":          {"赛事营销助手", "车队、赛事赞助、品牌曝光和营销合作。", "赛事日历、车队资料、赞助方案、品牌传播资料。", "没有赛事级别和预算时不做确定 ROI。"},
	"supplier-aggregator":        {"供应商聚合助手", "供应商查找、能力对比、MOQ、认证和报价整理。", "供应商资料、报价、认证、产能、交付和历史合作记录。", "没有来源或认证证据时必须标注待核实。"},
	"brand-value-crawler":        {"品牌情报助手", "品牌价值、竞品动向、渠道反馈和公开情报整理。", "品牌官网、新闻、社媒摘要、渠道反馈和竞品资料。", "不能把未经核验的舆情当事实。"},
	"oem-match-agent":            {"OEM 工厂匹配助手", "骑行行业 OEM/ODM 需求拆解、工厂能力匹配、MOQ、产能、认证和交付风险。", "整车、E-bike 电池、电机、头盔等 OEM/ODM 工厂能力资料、认证、产能、MOQ 和历史合作资料。", "没有可追溯来源、产能或认证证据时不得给出确定匹配结论，必须明确待核实项。"},
	"customer-measurement-agent": {"客户洞察助手", "客户发现、客户画像、客户质量评分和跟进建议。", "客户档案、线索、公开来源、沟通记录和交易记录。", "客户发现不得编造客户；没有联系方式或来源时标记缺口。"},
	"rag-agent":                  {"知识检索助手", "私有知识库检索、共享知识库摘要和证据定位。", "本地私有文件、共享文件、索引片段和元数据。", "只回答证据支持的内容；未命中时明确说未找到。"},
	"design-advisor":             {"产品设计助手", "产品方案、规格建议、体验设计和骑行场景适配。", "产品需求、设计稿、规格书、用户反馈和竞品资料。", "不要替代工程验证；涉及安全和法规时提示测试。"},
	"route-fetcher":              {"骑行路线助手", "骑行路线、补给点、活动路线和地图信息整理。", "路线文件、活动资料、地图链接、客户目的地和公开路线来源。", "没有真实路线来源时不要编造精确路线。"},
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
	Freshness string `json:"freshness,omitempty"`
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

type assistantQueueAgent struct {
	AgentID string `json:"agentId"`
	Name    string `json:"name"`
	Feature string `json:"feature,omitempty"`
	Started string `json:"started,omitempty"`
}

type assistantQueueInfo struct {
	Queued       bool                  `json:"queued"`
	Position     int                   `json:"position,omitempty"`
	Running      []assistantQueueAgent `json:"running,omitempty"`
	WaitedMs     int64                 `json:"waitedMs,omitempty"`
	Limit        int                   `json:"limit"`
	QueueLength  int                   `json:"queueLength"`
	Message      string                `json:"message,omitempty"`
	AgentID      string                `json:"agentId,omitempty"`
	AgentName    string                `json:"agentName,omitempty"`
	RunningNames []string              `json:"runningNames,omitempty"`
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
	Queue           assistantQueueInfo        `json:"queue,omitempty"`
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
	if len(parts) == 1 && parts[0] == "queue" {
		a.assistantQueueStatus(w, r)
		return
	}
	if len(parts) == 1 && parts[0] == "cases" {
		a.assistantCaseCreate(w, r)
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
	if request.Model == "" {
		models := a.openClaw.service.ModelSnapshot()
		defaultModel := strings.TrimSpace(firstNonEmpty(models.ResolvedDefault, models.DefaultModel))
		if defaultModel == "" {
			writeAPIError(w, http.StatusBadRequest, "MODEL_NOT_READY", "请先在设置中配置模型、测试连通性，并设置一个默认模型后再使用智能体功能")
			return
		}
		if !a.modelPreviouslyTestedOK(r.Context(), defaultModel) {
			writeAPIError(w, http.StatusBadRequest, "MODEL_NOT_READY", "当前默认模型尚未通过真实连通性测试，请先在设置中测试通过后再使用智能体功能")
			return
		}
	}
	if err := validateAssistantAttachments(request.Attachments); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ATTACHMENT", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Minute)
	defer cancel()
	requestID := fmt.Sprintf("AI-%d", time.Now().UnixNano())
	queueAgentID := assistantQueueAgentID(request)
	queueInfo, releaseQueue, err := a.assistantQueue.acquire(ctx, queueAgentID, request.Feature)
	if err != nil {
		cleanupAssistantAttachments(request.Attachments)
		writeAPIError(w, http.StatusGatewayTimeout, "ASSISTANT_QUEUE_TIMEOUT", "当前 Agent 队列等待超时，请稍后重试")
		return
	}
	defer releaseQueue()
	response := a.runAssistantQuery(ctx, request, requestID)
	response.Queue = queueInfo
	if queueInfo.Queued {
		response.Pipeline = append([]assistantPipelineStage{{
			Stage: "queue", Status: "ok", Detail: queueInfo.Message, DurationMs: queueInfo.WaitedMs,
			Reason: fmt.Sprintf("后台最多允许 %d 个 Agent 任务同时运行，超出后按请求顺序排队", queueInfo.Limit),
			Data:   strings.Join(queueInfo.RunningNames, "、"),
		}}, response.Pipeline...)
		response.Timings = append([]assistantTiming{{Stage: "queue", Status: "ok", DurationMs: queueInfo.WaitedMs, Reason: queueInfo.Message}}, response.Timings...)
		response.TotalDurationMs += queueInfo.WaitedMs
	}
	a.persistAssistantConversation(r.Context(), request, response)
	cleanupAssistantAttachments(request.Attachments)
	a.store.audit(r.Context(), "query", "assistant", request.Feature, requestOperator(r), map[string]any{
		"page": request.Page, "usedAgents": response.UsedAgents, "partial": response.Partial, "evidenceCount": len(response.Evidence), "queued": queueInfo.Queued, "queueWaitedMs": queueInfo.WaitedMs,
	})
	writeJSON(w, http.StatusOK, response)
}

type assistantCaseRequest struct {
	AgentID  string   `json:"agentId"`
	Title    string   `json:"title"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags"`
}

func (a *businessAPI) assistantCaseCreate(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request assistantCaseRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	request.AgentID = strings.TrimSpace(request.AgentID)
	request.Title = strings.TrimSpace(request.Title)
	request.Question = strings.TrimSpace(request.Question)
	request.Answer = strings.TrimSpace(request.Answer)
	if !assistantDomainAgentIDs[request.AgentID] {
		writeAPIError(w, http.StatusBadRequest, "INVALID_CASE_AGENT", "案例必须关联 STA-100 专业 Agent")
		return
	}
	if request.Question == "" || request.Answer == "" || len(request.Question) > 32<<10 || len(request.Answer) > 128<<10 {
		writeAPIError(w, http.StatusBadRequest, "INVALID_CASE_CONTENT", "案例问题和回答不能为空，且不能超过限制")
		return
	}
	if request.Title == "" {
		request.Title = truncateKnowledgeText(strings.ReplaceAll(request.Question, "\n", " "), 80)
	}
	id, err := a.store.nextSequence(r.Context(), "private_files", "CASE", 8)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	root, err := privateDataRoot()
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	path := filepath.Join(root, id+".md")
	content := fmt.Sprintf("# %s\n\n- Agent：%s\n- 确认时间：%s\n\n## 问题\n\n%s\n\n## 已确认回答\n\n%s\n", request.Title, request.AgentID, time.Now().UTC().Format(time.RFC3339), request.Question, request.Answer)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		writeBusinessError(w, err)
		return
	}
	item := PrivateFile{ID: id, Name: id + ".md", Category: request.AgentID + "确认案例", Tags: append([]string{"用户确认", request.AgentID}, request.Tags...), Size: humanBytes(int64(len(content))), Bytes: int64(len(content)), Mime: "text/markdown", Source: "用户确认案例", Status: "PendingParse", Path: path, Updated: currentText()}
	if err := a.store.create(r.Context(), "private_files", id, item); err != nil {
		_ = os.Remove(path)
		writeBusinessError(w, err)
		return
	}
	indexErr := a.indexPrivateKnowledgeFile(r.Context(), item)
	if indexErr != nil {
		item.Status = "NeedsProcessing"
	} else {
		item.Status = "Indexed"
	}
	_ = a.store.put(r.Context(), "private_files", id, item)
	_, _ = a.syncAgentKnowledgeFromLocalSources(r.Context())
	response := map[string]any{"item": item, "knowledgeIndex": a.knowledgeIndexStatus(r.Context())}
	if indexErr != nil {
		response["message"] = "案例已保存，但索引待处理：" + indexErr.Error()
	} else {
		response["message"] = "已确认案例已保存并加入对应 Agent 知识库"
	}
	writeJSON(w, http.StatusCreated, response)
}

func (a *businessAPI) assistantQueueStatus(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, a.assistantQueue.status())
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
		Text: request.Message, Sources: assistantAgentSources(agentID),
		Model: request.Model, CreatedAt: userCreatedAt,
	}
	_ = a.store.create(ctx, "agent_messages", userID, userMessage)
	a.extractAndPersistAssistantMemory(ctx, request, userID)
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
	if normalizedFeature(request.Feature) == "customer-search" {
		return a.runFastCustomerSearch(ctx, request, requestID, totalStarted)
	}
	evidence := a.collectLocalEvidence(ctx, request)
	sortAssistantEvidence(evidence)
	response := assistantQueryResponse{
		Evidence: evidence, Items: a.localAssistantItems(request, evidence), Conflicts: detectEvidenceConflicts(evidence),
		Attachments: assistantAttachmentViews(request.Attachments),
		AIGenerated: true, Partial: false,
		Todo: []string{"本机知识库按文件更新时间增量同步；无法解析的文件会保留在待处理列表", "联网工具接入后在工具层强制执行域名白名单"},
	}
	localDuration := time.Since(totalStarted).Milliseconds()
	if normalizedFeature(request.Feature) == "customer-discovery" {
		response.Todo = []string{"客户发现允许来源白名单、客户线索保存规则和联系方式校验规则待确认"}
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "query-fields", Status: "ok", Detail: "已接收国家、城市和客户类型筛选条件", DurationMs: localDuration, Reason: "客户本地发现不读取本机客户库，直接把页面条件交给 OpenClaw 客户发现 Agent", Data: contextSearchText(request.Context)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "query-fields", Status: "ok", DurationMs: localDuration, Reason: "客户本地发现不读取本机客户库，直接把页面条件交给 OpenClaw 客户发现 Agent"})
	} else {
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "agent-knowledge-context", Status: "ok", Detail: fmt.Sprintf("命中 %d 条 Agent 知识库上下文", len(evidence)), DurationMs: localDuration, Reason: "私有知识库和共享知识库只作为后台同步来源，不在聊天请求中实时扫描", Data: assistantEvidenceHoverData(evidence)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "agent-knowledge-context", Status: "ok", DurationMs: localDuration, Reason: "运行态只调用对应 OpenClaw Agent 知识库"})
	}
	if len(request.Attachments) > 0 {
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "attachments", Status: "ok", Detail: fmt.Sprintf("%d 个附件已校验并提交给 OpenClaw Agent", len(request.Attachments)), Reason: "校验本地附件路径并随消息提交", Data: assistantAttachmentHoverData(request.Attachments)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "attachments", Status: "ok", DurationMs: 0, Reason: "校验本地附件路径并随消息提交"})
	}
	if targetAgent := assistantTargetAgent(request); targetAgent != "" {
		return a.runTargetAssistantAgentQuery(ctx, request, requestID, targetAgent, response, evidence, totalStarted)
	}

	localeRule := assistantLocaleRule(request)
	evidenceJSON, _ := json.Marshal(evidence)
	knowledgePrompt := fmt.Sprintf("[页面] %s\n[功能] %s\n[用户输入]\n%s\n\n[Agent 知识库命中 JSON]\n%s\n\n[运行态知识库规则]\n请按 OpenClaw 当前 Agent 已同步知识库整理上下文。不得要求 Go 网关实时扫描私有知识库、共享知识库、客户表、供应商表或产品表。私有知识库和共享知识库只作为后台同步来源；运行态只使用已经解析并建立本地向量索引的内容，无法解析的文件不得当作事实引用。\n\n%s", request.Page, request.Feature, request.Message, evidenceJSON, localeRule)
	knowledgeStarted := time.Now()
	knowledge, err := a.sendAssistantAgent(ctx, knowledgeAgentID, request.Model, assistantStageSession(request.SessionKey, "knowledge"), knowledgePrompt, assistantAgentKnowledgeSources(), nil, request.Attachments)
	knowledgeDuration := time.Since(knowledgeStarted).Milliseconds()
	if err != nil {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "failed", knowledge)
		response.Partial = true
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, assistantErrorText(err)))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "failed", Detail: assistantErrorText(err), DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 失败或等待超时", Data: assistantErrorText(err)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "knowledge-agent", Status: "failed", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 失败或等待超时"})
		knowledge.Text = "Knowledge Agent 未返回；Coordinator 将只接收用户输入和领域 Agent 输出。"
	} else {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "ok", knowledge)
		response.UsedAgents = append(response.UsedAgents, knowledgeAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "ok", Detail: "Agent 知识库上下文已整理", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 整理已同步的 Agent 知识库上下文", Data: assistantTextHoverData(knowledge.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "knowledge-agent", Status: "ok", DurationMs: knowledgeDuration, Reason: "真实调用 Knowledge Agent 整理已同步的 Agent 知识库上下文"})
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
		"knowledgeSummary": knowledge.Text, "domainOutputs": domainOutputs,
		"requiredRule": "冲突信息全部并列展示并保留来源、记录编号和更新时间；当前聊天回答以对应 Agent 知识库为准。Go 网关不在聊天请求中实时扫描客户、供应商、产品等 SQLite 业务表，也不实时扫描私有知识库或共享知识库；私有/共享知识库通过后台同步任务进入各 Agent 专题知识库。",
	}
	coordinatorJSON, _ := json.Marshal(coordinatorInput)
	coordinatorPrompt := "请根据工作区规则整合以下 STA-100 查询上下文。只使用已提供内容，不得声称调用了未返回结果的 Agent。\n\n" + assistantLocaleRule(request) + "\n\n" + string(coordinatorJSON)
	if normalizedFeature(request.Feature) == "oem-match" {
		coordinatorPrompt += "\n\n" + assistantStructuredPrompt(assistantStructuredPromptSpec{
			ModuleName: "OEM 工厂智能匹配",
			Scope:      "用户输入的骑行类 OEM 需求和当前页面类别。",
			Fields:     []string{"title", "category", "reason", "detail", "source", "sourceUrl", "time", "score", "capacity", "moq"},
			ResultType: "oem_match",
			CountHint:  "输出所有可靠匹配结果，不要输出概览壳子。",
			Notes: []string{
				"每条 item 必须是完整候选工厂或供应商结果；title 写工厂或供应商名称；reason 写命中的业务理由；detail 写可用于验收的完整说明。",
				"source 写来源名称；sourceUrl 写原文或来源链接，没有则留空；time 写更新时间；score 写匹配度或分数；capacity 写产能；moq 写最小起订量。",
				"没有可靠结果时 items 输出空数组，不要编造事实。",
			},
		})
	}
	if normalizedFeature(request.Feature) == "customer-search" {
		coordinatorPrompt += "\n\n" + assistantStructuredPrompt(assistantStructuredPromptSpec{
			ModuleName: "客户统一搜索",
			Scope:      "用户输入的客户、市场或渠道查询，以及当前 Agent 知识库命中的客户资料。",
			Fields:     []string{"name", "country", "city", "type", "business", "contact", "phone", "email", "website", "address", "source", "sourceUrl", "updatedAt", "score", "reason"},
			ResultType: "customer_search",
			CountHint:  "只输出有明确来源和可识别客户名称的结果；历史资料必须保留来源更新时间。",
			Notes:      []string{"不得把无来源的推测当作客户记录。联系方式为空时保留为空。", "同名客户或冲突信息必须分别列出来源，不能强行合并。"},
		})
	}
	coordinatorStarted := time.Now()
	coordinator, coordinatorErr := a.sendAssistantAgent(ctx, coordinatorAgentID, request.Model, assistantStageSession(request.SessionKey, "coordinator"), coordinatorPrompt, assistantAgentKnowledgeSources(), nil, request.Attachments)
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
		if normalizedFeature(request.Feature) == "oem-match" || normalizedFeature(request.Feature) == "customer-search" {
			response.Items = parseAssistantStructuredItems(normalizedFeature(request.Feature), coordinator.Text)
		}
		response.UsedAgents = append(response.UsedAgents, coordinatorAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(coordinatorAgentID, coordinator, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "coordinator-agent", Status: "ok", Detail: "完成统一汇总", DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 汇总前序结果", Data: assistantTextHoverData(coordinator.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "coordinator-agent", Status: "ok", DurationMs: coordinatorDuration, Reason: "真实调用 Coordinator Agent 汇总前序结果"})
	}
	response.UsedAgents = uniqueStrings(response.UsedAgents)
	response.TotalDurationMs = time.Since(totalStarted).Milliseconds()
	response.SlowReason = "当前请求使用真实多 Agent 链路：先由 Knowledge Agent 整理已同步的 Agent 知识库上下文，再调用领域 Agent，最后由 Coordinator Agent 汇总；领域 Agent 虽然并行，但后续汇总必须等待前序结果。"
	if usage, err := a.tokenUsageSummary(ctx, requestID); err == nil {
		response.TokenUsage = usage
	}
	return response
}

// runFastCustomerSearch keeps the common search use case bounded to one local
// read and one customer-domain Agent call. The previous generic workflow made
// this page wait for Knowledge Agent, two domain Agents and Coordinator even
// when the user only wanted a customer lookup.
func (a *businessAPI) runFastCustomerSearch(ctx context.Context, request assistantQueryRequest, requestID string, totalStarted time.Time) assistantQueryResponse {
	localStarted := time.Now()
	localItems := a.localCustomerSearchItems(ctx, request)
	localDuration := time.Since(localStarted).Milliseconds()
	localSummary := customerSearchLocalSummary(localItems)
	response := assistantQueryResponse{
		Items: localItems, AIGenerated: true, Partial: false,
		Pipeline: []assistantPipelineStage{{Stage: "local-customer-evidence", Status: "ok", Detail: localSummary, DurationMs: localDuration, Reason: "读取本地客户、线索和已索引的客户相关知识库，不触发订单聚合", Data: customerItemsHoverData(localItems)}},
		Timings:  []assistantTiming{{Stage: "local-customer-evidence", Status: "ok", DurationMs: localDuration, Reason: "读取本地客户、线索和知识库"}},
		Todo:     []string{"互联网结果需保留公开来源、原文地址和更新时间；本地记录与互联网记录不自动覆盖"},
	}
	includeInternet, _ := request.Context["includeInternet"].(bool)
	if !includeInternet {
		response.Text = fmt.Sprintf("%s。仅展示本地结果；如需公开来源补充，请勾选“补充互联网”。", localSummary)
		response.TotalDurationMs = time.Since(totalStarted).Milliseconds()
		response.SlowReason = "当前为本地快速搜索：直接检索客户、线索和已索引知识库，不等待 OpenClaw 联网 Agent。"
		return response
	}

	localJSON, _ := json.Marshal(localItems)
	prompt := fmt.Sprintf(`[STA-100 客户统一搜索]
当前 Agent：customer-measurement-agent
[用户查询]
%s
[本地客户证据 JSON]
%s

请只补充来自互联网公开来源、且能核验的客户或机构信息。不得把本地证据改写成互联网结果；不要编造联系方式。与本地记录重复时仍可返回，但必须保留公开来源和更新时间。没有可靠互联网结果时 items 输出空数组，并说明原因。
%s`, request.Message, localJSON, assistantStructuredPrompt(assistantStructuredPromptSpec{
		ModuleName: "客户统一搜索",
		Scope:      "用户查询和本地客户证据；互联网结果必须有公开来源。",
		Fields:     []string{"name", "country", "city", "type", "business", "contact", "phone", "email", "website", "address", "source", "sourceUrl", "updatedAt", "score", "reason", "sourceType", "sourceLabel", "sourceUpdatedAt", "freshness"},
		ResultType: "customer_search",
		CountHint:  "只输出有明确客户名称和公开来源的互联网结果；没有可靠结果时输出空 items。",
		Notes:      []string{"sourceType 固定写 internet；sourceLabel 写公开来源名称；sourceUpdatedAt 写来源更新时间；freshness 写近期、较早或时间待核实。", "本地证据由 Go 侧单独展示，不能伪装成互联网来源。"},
	}))
	started := time.Now()
	result, err := a.sendAssistantAgent(ctx, "customer-measurement-agent", request.Model, assistantStageSession(request.SessionKey, "customer-search"), prompt, []string{"联网检索", "Agent 知识库"}, nil, request.Attachments)
	duration := time.Since(started).Milliseconds()
	if err != nil {
		errorText := assistantErrorText(err)
		response.Partial = true
		response.Text = fmt.Sprintf("%s；互联网搜索未完成：%s", localSummary, errorText)
		response.Outputs = []assistantAgentOutput{{AgentID: "customer-measurement-agent", Error: errorText, Text: result.Text, RunID: result.RunID, DurationMs: duration}}
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "internet-customer-agent", Status: "failed", Detail: errorText, DurationMs: duration, Reason: "客户领域 Agent 互联网检索失败", Data: errorText})
		response.Timings = append(response.Timings, assistantTiming{Stage: "internet-customer-agent", Status: "failed", DurationMs: duration, Reason: "客户领域 Agent 互联网检索失败"})
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "customer-search-agent", "customer-measurement-agent", "failed", result)
	} else {
		remoteItems := parseAssistantStructuredItems("customer-search", result.Text)
		remoteItems = filterCustomerSearchItems(remoteItems, request.Message)
		for _, item := range remoteItems {
			item["sourceType"] = "internet"
			item["sourceLabel"] = firstNonEmpty(fmt.Sprint(item["source"]), "互联网公开来源")
			item["sourceUpdatedAt"] = firstNonEmpty(fmt.Sprint(item["updatedAt"]), "")
			item["freshness"] = knowledgeFreshness(fmt.Sprint(item["updatedAt"]))
		}
		response.Items = append(localItems, remoteItems...)
		response.Text = fmt.Sprintf("%s，互联网公开来源 %d 条，已分别保留来源信息。", localSummary, len(remoteItems))
		response.UsedAgents = []string{"customer-measurement-agent"}
		response.Outputs = []assistantAgentOutput{{AgentID: "customer-measurement-agent", Text: result.Text, RunID: result.RunID, DurationMs: duration, Usage: result.Usage, UsageAvailable: result.Usage.Available()}}
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "internet-customer-agent", Status: "ok", Detail: fmt.Sprintf("互联网公开来源返回 %d 条", len(remoteItems)), DurationMs: duration, Reason: "只调用一个客户领域 Agent，减少重复编排等待", Data: assistantTextHoverData(result.Text)})
		response.Timings = append(response.Timings, assistantTiming{Stage: "internet-customer-agent", Status: "ok", DurationMs: duration, Reason: "调用客户领域 Agent 获取互联网公开来源"})
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "customer-search-agent", "customer-measurement-agent", "ok", result)
	}
	response.TotalDurationMs = time.Since(totalStarted).Milliseconds()
	response.SlowReason = "客户统一搜索已优化为本地快速读取 + 一个客户领域 Agent；不再默认调用 Knowledge Agent、多个领域 Agent 和 Coordinator。"
	if usage, usageErr := a.tokenUsageSummary(ctx, requestID); usageErr == nil {
		response.TokenUsage = usage
	}
	return response
}

func filterCustomerSearchItems(items []map[string]any, query string) []map[string]any {
	profile := newCustomerSearchProfile(query)
	if len(profile.countries) == 0 {
		return items
	}
	filtered := make([]map[string]any, 0, len(items))
	for _, item := range items {
		searchable := strings.ToLower(strings.Join([]string{
			fmt.Sprint(item["name"]), fmt.Sprint(item["country"]), fmt.Sprint(item["city"]),
			fmt.Sprint(item["type"]), fmt.Sprint(item["business"]), fmt.Sprint(item["reason"]),
		}, " "))
		if profile.matchesStructuredCountry(fmt.Sprint(item["country"]), searchable) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (a *businessAPI) localCustomerSearchItems(ctx context.Context, request assistantQueryRequest) []map[string]any {
	customers, err := listRecords[Customer](ctx, a.store, "accounts")
	if err != nil {
		return []map[string]any{}
	}
	leads, err := listRecords[Lead](ctx, a.store, "leads")
	if err != nil {
		leads = nil
	}
	profile := newCustomerSearchProfile(request.Message)
	hasContact, _ := request.Context["hasContact"].(bool)
	candidates := make([]customerSearchCandidate, 0, minInt(len(customers)+len(leads)+8, 46))
	for _, customer := range customers {
		if customer.Archived {
			continue
		}
		searchable := strings.ToLower(strings.Join([]string{customer.ID, customer.Name, customer.Country, customer.City, customer.Type, customer.Contact, customer.Phone, customer.Email, customer.Website, customer.Description}, " "))
		score := profile.score(searchable, customer.Type)
		if !profile.matchesStructuredCountry(customer.Country, searchable) || !profile.matches(searchable, score) {
			continue
		}
		if hasContact && firstNonEmpty(customer.Contact, customer.Phone, customer.Email, customer.Website) == "" {
			continue
		}
		candidates = append(candidates, customerSearchCandidate{score: score + 8, item: map[string]any{"id": customer.ID, "name": customer.Name, "country": customer.Country, "city": customer.City, "type": customer.Type, "contact": customer.Contact, "phone": customer.Phone, "email": customer.Email, "website": customer.Website, "address": "", "business": customer.Description, "source": "本地业务数据库", "sourceType": "local_business", "sourceLabel": "本地业务数据库", "sourceUpdatedAt": customer.Updated, "updatedAt": customer.Updated, "freshness": knowledgeFreshness(customer.Updated), "reason": "命中本地客户记录"}})
	}
	for _, lead := range leads {
		if lead.Archived {
			continue
		}
		searchable := strings.ToLower(strings.Join([]string{lead.ID, lead.Name, lead.Country, lead.City, lead.Type, lead.Contact, lead.Phone, lead.Email, lead.Website, lead.Address, lead.Business, lead.Source, lead.Reason}, " "))
		score := profile.score(searchable, lead.Type)
		if !profile.matchesStructuredCountry(lead.Country, searchable) || !profile.matches(searchable, score) || (hasContact && firstNonEmpty(lead.Contact, lead.Phone, lead.Email, lead.Website) == "") {
			continue
		}
		candidates = append(candidates, customerSearchCandidate{score: score + 4, item: map[string]any{"id": lead.ID, "name": lead.Name, "country": lead.Country, "city": lead.City, "type": lead.Type, "contact": lead.Contact, "phone": lead.Phone, "email": lead.Email, "website": lead.Website, "address": lead.Address, "business": lead.Business, "source": firstNonEmpty(lead.Source, "本地线索库"), "sourceUrl": lead.SourceUrl, "sourceType": "local_lead", "sourceLabel": "本地线索库", "sourceUpdatedAt": lead.Updated, "updatedAt": lead.Updated, "freshness": knowledgeFreshness(lead.Updated), "reason": firstNonEmpty(lead.Reason, "命中本地线索记录")}})
	}
	if hits, err := a.localCustomerKnowledgeHits(ctx, profile); err == nil {
		for _, hit := range hits {
			searchable := strings.ToLower(hit.Title + " " + hit.Category + " " + hit.Content)
			score := profile.score(searchable, hit.Category)
			if !profile.matches(searchable, score) {
				continue
			}
			preview := trimHoverText(hit.Content, 260)
			candidates = append(candidates, customerSearchCandidate{score: score, item: map[string]any{"id": hit.ChunkID, "name": hit.Title, "type": hit.Category, "business": preview, "content": trimHoverText(hit.Content, 1600), "sourcePath": hit.SourcePath, "source": "本地知识库", "sourceType": "local_knowledge", "sourceLabel": "本地知识库 / " + firstNonEmpty(hit.Category, "未分类"), "sourceUpdatedAt": hit.UpdatedAt, "updatedAt": hit.UpdatedAt, "freshness": knowledgeFreshness(hit.UpdatedAt), "reason": "命中已索引知识库片段"}})
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return fmt.Sprint(candidates[i].item["name"]) < fmt.Sprint(candidates[j].item["name"])
		}
		return candidates[i].score > candidates[j].score
	})
	items := make([]map[string]any, 0, minInt(len(candidates), 30))
	seen := map[string]bool{}
	for _, candidate := range candidates {
		key := strings.ToLower(strings.TrimSpace(fmt.Sprint(candidate.item["name"]))) + "\x00" + strings.ToLower(strings.TrimSpace(fmt.Sprint(candidate.item["sourceType"])))
		if key == "\x00local_knowledge" || seen[key] {
			continue
		}
		seen[key] = true
		candidate.item["score"] = minInt(candidate.score, 100)
		items = append(items, candidate.item)
		if len(items) >= 30 {
			break
		}
	}
	return items
}

// localCustomerKnowledgeHits only inspects document metadata. Customer search
// is an interactive lookup, so it must not scan every knowledge chunk body or
// vector. Full chunk retrieval remains part of the dedicated Agent workflow.
func (a *businessAPI) localCustomerKnowledgeHits(ctx context.Context, profile customerSearchProfile) ([]knowledgeHit, error) {
	terms := profile.knowledgeTerms()
	if len(terms) == 0 {
		return nil, nil
	}
	clauses := make([]string, 0, len(terms))
	args := make([]any, 0, len(terms)*3)
	for _, term := range terms {
		clauses = append(clauses, "(LOWER(d.name) LIKE ? OR LOWER(d.category) LIKE ? OR LOWER(d.source_path) LIKE ?)")
		pattern := "%" + strings.ToLower(term) + "%"
		args = append(args, pattern, pattern, pattern)
	}
	countryFilter := ""
	if countryTerms := profile.countryTerms(); len(countryTerms) > 0 {
		countryClauses := make([]string, 0, len(countryTerms))
		for _, term := range countryTerms {
			countryClauses = append(countryClauses, "LOWER(c2.content) LIKE ?")
			args = append(args, "%"+strings.ToLower(term)+"%")
		}
		countryFilter = " AND EXISTS (SELECT 1 FROM knowledge_chunks c2 WHERE c2.document_id=d.id AND (" + strings.Join(countryClauses, " OR ") + "))"
	}
	// Customer search first filters documents by metadata, channel terms and
	// requested country. Then load the indexed document text so the UI can
	// show useful detail instead of only exposing a filename. A document can
	// contain several chunks, therefore country matching must inspect all of
	// them rather than only the first chunk.
	query := `SELECT d.id,d.id,d.name,d.source_path,d.category,d.visibility,
		COALESCE((SELECT GROUP_CONCAT(c.content,' ') FROM knowledge_chunks c WHERE c.document_id=d.id),''),d.source_updated_at
		FROM knowledge_documents d
		WHERE d.status='indexed' AND (` + strings.Join(clauses, " OR ") + `)
		` + countryFilter + ` LIMIT 16`
	rows, err := a.store.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	hits := make([]knowledgeHit, 0, 8)
	for rows.Next() {
		var hit knowledgeHit
		if err := rows.Scan(&hit.ChunkID, &hit.DocumentID, &hit.Title, &hit.SourcePath, &hit.Category, &hit.Visibility, &hit.Content, &hit.UpdatedAt); err != nil {
			return nil, err
		}
		hits = append(hits, hit)
	}
	return hits, rows.Err()
}

type customerSearchCandidate struct {
	score int
	item  map[string]any
}

type customerSearchProfile struct {
	query        string
	directTokens []string
	countries    []customerSearchCountry
	wantsChannel bool
	wantsCycling bool
}

type customerSearchCountry struct {
	name    string
	aliases []string
}

var customerSearchCountries = []customerSearchCountry{
	{name: "中国", aliases: []string{"中国", "china", "prc", "中国大陆"}},
	{name: "德国", aliases: []string{"德国", "germany", "german", "deutschland"}},
	{name: "法国", aliases: []string{"法国", "france", "french"}},
	{name: "英国", aliases: []string{"英国", "uk", "united kingdom", "england", "britain"}},
	{name: "美国", aliases: []string{"美国", "usa", "us", "united states", "america"}},
	{name: "加拿大", aliases: []string{"加拿大", "canada", "canadian"}},
	{name: "日本", aliases: []string{"日本", "japan", "japanese"}},
	{name: "韩国", aliases: []string{"韩国", "south korea", "korea", "korean"}},
	{name: "澳大利亚", aliases: []string{"澳大利亚", "australia", "australian"}},
	{name: "新西兰", aliases: []string{"新西兰", "new zealand"}},
	{name: "荷兰", aliases: []string{"荷兰", "netherlands", "dutch", "holland"}},
	{name: "比利时", aliases: []string{"比利时", "belgium", "belgian"}},
	{name: "西班牙", aliases: []string{"西班牙", "spain", "spanish"}},
	{name: "意大利", aliases: []string{"意大利", "italy", "italian"}},
	{name: "瑞典", aliases: []string{"瑞典", "sweden", "swedish"}},
	{name: "波兰", aliases: []string{"波兰", "poland", "polish"}},
	{name: "瑞士", aliases: []string{"瑞士", "switzerland", "swiss"}},
	{name: "奥地利", aliases: []string{"奥地利", "austria", "austrian"}},
	{name: "新加坡", aliases: []string{"新加坡", "singapore"}},
	{name: "印度", aliases: []string{"印度", "india", "indian"}},
}

func newCustomerSearchProfile(query string) customerSearchProfile {
	profile := customerSearchProfile{query: strings.ToLower(strings.TrimSpace(query))}
	for _, country := range customerSearchCountries {
		if containsAnyCustomerSearchTerm(profile.query, country.aliases...) {
			profile.countries = append(profile.countries, country)
		}
	}
	profile.wantsChannel = containsAnyCustomerSearchTerm(profile.query, "代理", "经销", "分销", "渠道", "dealer", "distributor", "importer", "reseller", "retailer", "批发", "零售")
	profile.wantsCycling = containsAnyCustomerSearchTerm(profile.query, "骑行", "自行车", "单车", "bike", "bicycle", "cycling", "e-bike", "ebike")
	for _, token := range searchTokens(profile.query) {
		if !isCustomerSearchGenericToken(token) {
			profile.directTokens = append(profile.directTokens, token)
		}
	}
	return profile
}

func (p customerSearchProfile) knowledgeTerms() []string {
	terms := append([]string{}, p.directTokens...)
	if p.wantsChannel {
		terms = append(terms, "代理", "经销", "分销", "dealer", "distributor", "importer", "reseller")
	}
	if p.wantsCycling {
		terms = append(terms, "自行车", "骑行", "bike", "bicycle", "cycling", "e-bike")
	}
	return uniqueStrings(terms)
}

func (p customerSearchProfile) countryTerms() []string {
	terms := make([]string, 0, len(p.countries)*2)
	for _, country := range p.countries {
		for _, alias := range country.aliases {
			if len([]rune(alias)) <= 2 && isASCIIWord(alias) {
				continue
			}
			terms = append(terms, alias)
		}
	}
	return uniqueStrings(terms)
}

func (p customerSearchProfile) score(searchable, customerType string) int {
	if p.query == "" {
		return 1
	}
	score := 0
	for _, token := range p.directTokens {
		if strings.Contains(searchable, token) {
			score += 28
		}
	}
	if p.wantsChannel && (isChannelCustomerType(customerType) || containsAnyCustomerSearchTerm(searchable, "代理", "经销", "分销", "渠道", "dealer", "distributor", "importer", "reseller", "retailer", "批发", "零售")) {
		score += 48
	}
	if p.wantsCycling && containsAnyCustomerSearchTerm(searchable, "骑行", "自行车", "单车", "bike", "bicycle", "cycling", "e-bike", "ebike", "自行車") {
		score += 18
	}
	return score
}

func (p customerSearchProfile) matches(searchable string, score int) bool {
	if p.query == "" {
		return true
	}
	if len(p.countries) > 0 && !p.matchesCountry(searchable) {
		return false
	}
	if len(p.directTokens) > 0 {
		for _, token := range p.directTokens {
			if strings.Contains(searchable, token) {
				return score > 0
			}
		}
		return false
	}
	return score > 0
}

func (p customerSearchProfile) matchesCountry(searchable string) bool {
	searchable = strings.ToLower(searchable)
	for _, country := range p.countries {
		for _, alias := range country.aliases {
			if customerCountryAliasMatch(searchable, alias) {
				return true
			}
		}
	}
	return false
}

func customerCountryAliasMatch(value, alias string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	alias = strings.ToLower(strings.TrimSpace(alias))
	if value == "" || alias == "" {
		return false
	}
	// Short Latin aliases such as "us" and "uk" must be whole words;
	// substring matching would incorrectly match words like "business".
	if len([]rune(alias)) <= 2 && isASCIIWord(alias) {
		for _, token := range strings.FieldsFunc(value, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
		}) {
			if token == alias {
				return true
			}
		}
		return false
	}
	return strings.Contains(value, alias)
}

func isASCIIWord(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func (p customerSearchProfile) matchesStructuredCountry(country, searchable string) bool {
	if len(p.countries) == 0 {
		return true
	}
	if strings.TrimSpace(country) != "" {
		return p.matchesCountry(country)
	}
	return p.matchesCountry(searchable)
}

func isCustomerSearchGenericToken(token string) bool {
	return containsAnyCustomerSearchTerm(token, "代理", "经销", "分销", "渠道", "dealer", "distributor", "importer", "reseller", "retailer", "批发", "零售", "骑行", "自行车", "单车", "bike", "bicycle", "cycling", "e-bike", "ebike", "相关", "客户", "公司")
}

func isChannelCustomerType(value string) bool {
	return containsAnyCustomerSearchTerm(strings.ToLower(value), "dealer", "distributor", "importer", "reseller", "retailer", "agent", "代理", "经销", "分销", "渠道", "批发", "零售")
}

func containsAnyCustomerSearchTerm(value string, terms ...string) bool {
	value = strings.ToLower(value)
	for _, term := range terms {
		if strings.Contains(value, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func customerSearchLocalSummary(items []map[string]any) string {
	counts := map[string]int{}
	for _, item := range items {
		counts[fmt.Sprint(item["sourceType"])]++
	}
	parts := make([]string, 0, 3)
	if counts["local_business"] > 0 {
		parts = append(parts, fmt.Sprintf("本地客户 %d 条", counts["local_business"]))
	}
	if counts["local_lead"] > 0 {
		parts = append(parts, fmt.Sprintf("本地线索 %d 条", counts["local_lead"]))
	}
	if counts["local_knowledge"] > 0 {
		parts = append(parts, fmt.Sprintf("本地知识库 %d 条", counts["local_knowledge"]))
	}
	if len(parts) == 0 {
		return "未命中本地客户、线索或知识库记录"
	}
	return strings.Join(parts, "，")
}

func customerItemsHoverData(items []map[string]any) string {
	if len(items) == 0 {
		return "本地客户证据未命中。"
	}
	lines := make([]string, 0, minInt(len(items), 5))
	for _, item := range items[:minInt(len(items), 5)] {
		lines = append(lines, fmt.Sprintf("%s｜%s｜%s", item["name"], item["type"], item["sourceLabel"]))
	}
	return strings.Join(lines, "\n")
}

func (a *businessAPI) runTargetAssistantAgentQuery(ctx context.Context, request assistantQueryRequest, requestID, agentID string, response assistantQueryResponse, evidence []assistantEvidence, totalStarted time.Time) assistantQueryResponse {
	evidenceJSON, _ := json.Marshal(evidence)
	memoryText := a.assistantMemoryPrompt(ctx, agentID, request.SessionKey)
	prompt := fmt.Sprintf("[STA-100 当前智能体]\n%s\n\n%s\n\n[页面] %s\n[功能] %s\n[用户输入]\n%s\n\n[当前会话 Memory]\n%s\n\n[Agent 知识库命中 JSON]\n%s\n\n%s\n\n[运行态知识库规则]\n请基于当前智能体职责和该 Agent 已同步专题知识库回答。私有知识库和共享知识库只作为后台同步来源，不在本次请求中直接扫描。Go 不读取客户、供应商、产品等 SQLite 业务表。证据不足时明确说明缺口和需要补充的数据。冲突信息必须并列保留来源、记录编号和更新时间。附件如已提交，请结合附件内容；不能确认附件内容时要说明。", agentID, assistantAgentProfilePrompt(agentID), request.Page, request.Feature, request.Message, memoryText, evidenceJSON, assistantTurnBoundaryRule())
	preferences := defaultPreferences()
	_ = a.store.getSetting(ctx, "preferences", &preferences)
	if normalizedFeature(request.Feature) == "customer-discovery" {
		limit := discoveryResultLimit(request.Context, preferences)
		country := firstNonEmpty(stringContextValue(request.Context, "country"), "未填写")
		cities := stringArrayContextValue(request.Context, "cities")
		types := stringArrayContextValue(request.Context, "types")
		citiesStr := strings.Join(cities, "、")
		typesStr := strings.Join(types, "、")
		if citiesStr == "" {
			citiesStr = "未填写"
		}
		if typesStr == "" {
			typesStr = "未填写"
		}
		scope := fmt.Sprintf("国家=%s；城市=%s", country, citiesStr)
		prompt = fmt.Sprintf("[STA-100 客户发现]\n当前目标 Agent：%s\n\n[固定输入]\n国家：%s\n城市：%s\n客户类型：%s\n用户输入：%s\n\n[执行边界]\n请不要检索或引用 STA-100 本机客户库。本功能只把页面固定筛选条件交给 OpenClaw 客户发现 Agent，由 Agent 按自身已配置能力获取公开来源并完成核验。\n\n[业务规则]\n- 国家、城市、客户类型都是硬性筛选条件，只返回同时满足三项条件的客户或机构线索。\n- 每条线索必须有可核验来源；没有来源、城市不确定或类型不确定的候选不要进入 items。\n- 不满足筛选条件的市场概览、供应商、产品、泛行业新闻不要作为客户线索返回。\n- 无法获取可靠客户线索时，明确说明原因，并输出空 items；不要用模型常识补客户，不要编造。\n- 如果公开信息存在冲突，在 reason 中保留冲突点，不要覆盖。\n%s", agentID, country, citiesStr, typesStr, request.Message, assistantStructuredPrompt(assistantStructuredPromptSpec{
			ModuleName: "本地客户发现",
			Scope:      fmt.Sprintf("%s；客户类型=%s；只返回符合筛选条件的客户线索。", scope, typesStr),
			Fields:     []string{"name", "country", "city", "type", "business", "contact", "phone", "email", "website", "address", "source", "sourceUrl", "updatedAt", "score", "reason"},
			ResultType: "customer_discovery",
			CountHint:  fmt.Sprintf("最多输出 %d 条候选客户；items 数量不得超过该上限。只输出候选客户，不要输出市场概览、产品或泛行业新闻。", limit),
			Notes: []string{
				"name 写客户或机构名称；country、city、type 分别写国家、城市和客户类型；business 写业务方向。",
				"contact 写主要联系人姓名；phone 写电话；email 写邮箱；website 写官网；address 写地址。",
				"source 写来源名称；sourceUrl 写原文或来源链接，没有则留空；updatedAt 写抓取或更新时间；score 写匹配分；reason 写匹配理由。",
				"没有可靠结果时 items 输出空数组，不要编造客户。",
			},
		}))
	}
	started := time.Now()
	sources := assistantAgentSources(agentID)
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
		stageReason := "直接调用当前 OpenClaw Agent，运行态来源为该 Agent 已同步知识库"
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
	city := firstNonEmpty(strings.Join(stringArrayContextValue(request.Context, "cities"), "、"), stringContextValue(request.Context, "city"), "未填写")
	customerType := firstNonEmpty(strings.Join(stringArrayContextValue(request.Context, "types"), "、"), stringContextValue(request.Context, "type"), "未填写")
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
	case "customer_discovery", "customer_search":
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
			"phone":     get("phone", "tel", "telephone"),
			"email":     get("email", "mail"),
			"website":   get("website", "url", "web"),
			"address":   get("address", "location"),
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
	normalizeFloat := func(v float64) int {
		if v > 0 && v <= 1 {
			return int(math.Round(v * 100))
		}
		return int(math.Round(v))
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return normalizeFloat(typed)
	case json.Number:
		if parsed, err := typed.Int64(); err == nil {
			return int(parsed)
		}
		if parsed, err := typed.Float64(); err == nil {
			return normalizeFloat(parsed)
		}
	case string:
		raw := strings.TrimSpace(strings.TrimSuffix(typed, "%"))
		if parsed, err := strconv.Atoi(raw); err == nil {
			return parsed
		}
		if parsed, err := strconv.ParseFloat(raw, 64); err == nil {
			return normalizeFloat(parsed)
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
			memoryText := a.assistantMemoryPrompt(ctx, agentID, request.SessionKey)
			prompt := fmt.Sprintf("[STA-100 页面] %s / %s\n\n%s\n\n[用户输入]\n%s\n\n[当前会话 Memory]\n%s\n\n[Knowledge Agent 已同步知识库摘要]\n%s\n\n%s\n\n%s\n\n%s\n\n请只从你的专业职责和当前 Agent 知识库给出结果，保留证据和不确定性。冲突数据不要覆盖。", request.Page, request.Feature, assistantAgentProfilePrompt(agentID), request.Message, memoryText, knowledgeSummary, assistantLocaleRule(request), assistantTurnBoundaryRule(), assistantKnowledgeSyncTodo(agentID))
			result, err := a.sendAssistantAgent(ctx, agentID, request.Model, assistantStageSession(request.SessionKey, agentID), prompt, assistantAgentSources(agentID), preferences.AgentAllowlists[agentID], request.Attachments)
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
		return "Agent 知识库未命中可展示证据。"
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
		return []string{"oem-match-agent"}
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

func assistantQueueAgentID(request assistantQueryRequest) string {
	if target := assistantTargetAgent(request); target != "" {
		return target
	}
	switch normalizedFeature(request.Feature) {
	case "oem-match", "oem":
		return "oem-match-agent"
	case "customer-search", "customers", "customer-discovery":
		return "customer-measurement-agent"
	case "quote", "quotes":
		return "export-agent"
	case "order", "orders", "logistics":
		return "shipping-eta"
	case "document", "documents":
		return "invoice-agent"
	case "product", "products":
		return "compatibility-agent"
	case "news":
		return "market-analyzer"
	default:
		return normalizedFeature(request.Feature)
	}
}

func assistantAgentDisplayName(agentID string) string {
	agentID = strings.TrimSpace(agentID)
	if profile, ok := assistantAgentProfiles[agentID]; ok {
		return fmt.Sprintf("%s（%s）", profile.Name, agentID)
	}
	switch normalizedFeature(agentID) {
	case "oem-match", "oem":
		return "OEM 工厂智能匹配"
	case "customer-search", "customers":
		return "客户统一搜索"
	case "customer-discovery":
		return "本地客户发现"
	case "news":
		return "行业新闻"
	case "":
		return "智能查询"
	default:
		return agentID
	}
}

func assistantAgentProfilePrompt(agentID string) string {
	profile, ok := assistantAgentProfiles[agentID]
	if !ok {
		return "[专业知识] 该 Agent 使用对应专题知识库和当前输入上下文进行回答；若缺少事实，请明确说明缺口。"
	}
	return fmt.Sprintf("[专业知识]\n- Agent：%s\n- 专业方向：%s\n- 对应知识库范围：%s\n- 约束：%s", profile.Name, profile.Focus, profile.KnowledgeScope, profile.Boundaries)
}

func assistantAgentSources(agentID string) []string {
	if normalizedFeature(agentID) == "customer-discovery" {
		return []string{"联网检索"}
	}
	return assistantAgentKnowledgeSources()
}

func assistantAgentKnowledgeSources() []string {
	return []string{"Agent 知识库"}
}

func assistantKnowledgeSyncTodo(agentID string) string {
	if profile, ok := assistantAgentProfiles[agentID]; ok {
		return fmt.Sprintf("对应的私有知识库和共享知识库会定时同步到「%s知识库」，该 Agent 运行时只使用已完成索引的专题资料。", profile.Name)
	}
	return "对应的私有知识库和共享知识库会定时同步到当前 Agent 知识库。"
}

func assistantTurnBoundaryRule() string {
	return "[本轮输入边界]\n当前回答必须以本轮 [用户输入] 明确提出的国家、市场、产品和任务为最高优先级。当前会话 Memory 和 OpenClaw 历史会话只能作为偏好或背景参考；除非用户明确说“继续”“同上”“上一轮”“再发”或等价表达，不得把历史会话中的国家、市场、文件或结论当成本轮问题主线。"
}

// assistantLocaleRule 返回给 Agent 的强约束规则，避免模型默认举例用户未提及的国家或地区。
// 规则基于应用当前语言（zh/en）和当前页面的国家偏好（若 context 提供）。
// 主要目的：
//  1. 强制按应用当前语言回复；
//  2. 不在证据、用户输入和上下文之外的地区做举例；
//  3. 没有数据时明确说"未找到"，不要用常识补。
func assistantLocaleRule(request assistantQueryRequest) string {
	lang := strings.ToLower(strings.TrimSpace(stringContextValue(request.Context, "lang")))
	if lang == "" {
		lang = "zh"
	}
	preferCountries := []string{}
	if raw, ok := request.Context["preferCountries"].([]any); ok {
		for _, v := range raw {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				preferCountries = append(preferCountries, strings.TrimSpace(s))
			}
		}
	}
	if raw, ok := request.Context["preferCountries"].([]string); ok {
		preferCountries = append(preferCountries, raw...)
	}
	preferCountries = uniqueNonEmpty(preferCountries)
	langRule := "中文（zh-CN）"
	if lang == "en" {
		langRule = "English (en)"
	}
	parts := []string{
		fmt.Sprintf("[语言约束] 当前应用界面语言为 %s，最终回复必须使用 %s 撰写；不要混用其他语言。", langRule, langRule),
		"[范围约束] 优先围绕用户输入中明确提到的国家、地区、市场、产品和业务场景回答；不要主动引入用户未提及的国家或地区作为主线。",
		"[证据约束] 只能使用 [Agent 知识库命中 JSON]、[当前会话 Memory] 和 [用户输入] 中已提供的内容；证据外的内容不得作为确定事实写入最终结论。",
		"[空白处理] 没有可用证据时，必须明确写'当前证据中未找到'或'No matching record found'，禁止用模型常识补客户、补市场、补国家。",
		"[冲突处理] 冲突信息必须并列展示并保留来源、记录编号和更新时间；不要覆盖、不要平均、不要择一。",
	}
	if len(preferCountries) > 0 {
		parts = append(parts, fmt.Sprintf("[地区偏好] 当前应用倾向以下国家/地区：%s；如需举例请从该范围内选择，不要使用范围外国家。", strings.Join(preferCountries, "、")))
	}
	return strings.Join(parts, "\n")
}

func uniqueNonEmpty(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		key := strings.ToLower(strings.TrimSpace(v))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, v)
	}
	return out
}

func (a *businessAPI) collectLocalEvidence(ctx context.Context, request assistantQueryRequest) []assistantEvidence {
	if normalizedFeature(request.Feature) == "customer-discovery" {
		return nil
	}
	if a.knowledgeReady != nil {
		select {
		case <-a.knowledgeReady:
		case <-time.After(5 * time.Second):
		}
	}
	query := strings.ToLower(strings.TrimSpace(request.Message + " " + contextSearchText(request.Context)))
	tokens := searchTokens(query)
	evidence := make([]assistantEvidence, 0, 30)
	appendMatch := func(item assistantEvidence, searchable string) {
		category := strings.TrimSpace(stringContextValue(request.Context, "category"))
		if normalizedFeature(request.Feature) == "oem-match" && category != "" && category != "全部" && !strings.Contains(strings.ToLower(item.Source), strings.ToLower(category)) {
			return
		}
		if len(evidence) >= 30 || !matchesSearchTokens(strings.ToLower(searchable), tokens) {
			return
		}
		evidence = append(evidence, item)
	}
	agentIDs := assistantDomainRoutes(request)
	if target := assistantTargetAgent(request); target != "" {
		agentIDs = []string{target}
	}
	indexed := a.knowledgeEvidence(ctx, agentIDs, query)
	for _, item := range indexed {
		appendMatch(item, strings.Join([]string{item.ID, item.Title, item.Content, item.Source}, " "))
	}
	// Keep the profile/metadata fallback for an upgrade that has not finished
	// its first source import yet. Indexed chunks always take precedence.
	var indexedChunkCount int
	_ = a.store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM knowledge_chunks`).Scan(&indexedChunkCount)
	if len(evidence) == 0 && indexedChunkCount == 0 {
		agentSet := map[string]bool{}
		for _, id := range agentIDs {
			agentSet[id] = true
		}
		items, _ := listRecords[AgentKnowledgeItem](ctx, a.store, agentKnowledgeKind)
		for _, item := range items {
			if !agentSet[item.AgentID] {
				continue
			}
			searchable := strings.Join([]string{item.ID, item.AgentID, item.SourceName, item.Category, strings.Join(item.Tags, " "), item.Content}, " ")
			appendMatch(assistantEvidence{ID: item.ID, Entity: "agent-knowledge", Title: item.SourceName, Content: item.Content, UpdatedAt: item.SyncedAt, Freshness: knowledgeFreshness(item.SyncedAt), Source: item.AgentID + " 知识库"}, searchable)
		}
	}
	return evidence
}

func (a *businessAPI) assistantMemoryPrompt(ctx context.Context, agentID, sessionKey string) string {
	items, err := listRecords[AgentMemory](ctx, a.store, agentMemoryKind)
	if err != nil {
		return "暂无可用 Memory。"
	}
	lines := make([]string, 0, 8)
	for _, item := range items {
		if item.AgentID != agentID || item.SessionKey != sessionKey {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s/%s：%s（更新时间：%s）", item.Kind, item.Key, item.Value, item.UpdatedAt))
		if len(lines) >= 8 {
			break
		}
	}
	if len(lines) == 0 {
		return "暂无可用 Memory。"
	}
	return strings.Join(lines, "\n")
}

func stringContextValue(context map[string]any, key string) string {
	value, _ := context[key].(string)
	return value
}

func stringArrayContextValue(context map[string]any, key string) []string {
	if arr, ok := context[key].([]any); ok {
		result := make([]string, 0, len(arr))
		for _, v := range arr {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	if s, ok := context[key].(string); ok && s != "" {
		return []string{s}
	}
	return nil
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
	parts := []string{fmt.Sprintf("已整理 %d 条待同步知识库证据和 %d 条可展示记录，%d 个领域 Agent 返回了结果。", len(evidence), len(items), successful)}
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
	parts := []string{fmt.Sprintf("已整理 %d 条待同步知识库证据，当前智能体 %s 未完成回复。", evidenceCount, agentID)}
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

func discoveryResultLimit(context map[string]any, preferences UserPreferences) int {
	if limit := intContextValue(context, "limit"); limit > 0 {
		return limit
	}
	if preferences.DiscoveryShowLimit > 0 {
		return preferences.DiscoveryShowLimit
	}
	return 10
}

func limitAssistantItems(items []map[string]any, limit int) []map[string]any {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func intContextValue(context map[string]any, key string) int {
	if v, ok := context[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}
