package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
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
	Page       string         `json:"page"`
	Feature    string         `json:"feature"`
	Message    string         `json:"message"`
	SessionKey string         `json:"sessionKey"`
	Context    map[string]any `json:"context"`
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
	Usage          orchestrator.TokenUsage `json:"usage"`
	UsageAvailable bool                    `json:"usageAvailable"`
}

type assistantPipelineStage struct {
	Stage  string `json:"stage"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type assistantQueryResponse struct {
	Text        string                   `json:"text"`
	Items       []map[string]any         `json:"items"`
	UsedAgents  []string                 `json:"usedAgents"`
	Evidence    []assistantEvidence      `json:"evidence"`
	Conflicts   []map[string]any         `json:"conflicts"`
	Outputs     []assistantAgentOutput   `json:"agentOutputs,omitempty"`
	Pipeline    []assistantPipelineStage `json:"pipeline"`
	AIGenerated bool                     `json:"aiGenerated"`
	Partial     bool                     `json:"partial"`
	Todo        []string                 `json:"todo,omitempty"`
	TokenUsage  tokenUsageSummary        `json:"tokenUsage"`
}

func (a *businessAPI) assistantRouter(w http.ResponseWriter, r *http.Request, parts []string) {
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
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Minute)
	defer cancel()
	requestID := fmt.Sprintf("AI-%d", time.Now().UnixNano())
	response := a.runAssistantQuery(ctx, request, requestID)
	a.store.audit(r.Context(), "query", "assistant", request.Feature, requestOperator(r), map[string]any{
		"page": request.Page, "usedAgents": response.UsedAgents, "partial": response.Partial, "evidenceCount": len(response.Evidence),
	})
	writeJSON(w, http.StatusOK, response)
}

func (a *businessAPI) runAssistantQuery(ctx context.Context, request assistantQueryRequest, requestID string) assistantQueryResponse {
	evidence := a.collectLocalEvidence(ctx, request)
	sortAssistantEvidence(evidence)
	response := assistantQueryResponse{
		Evidence: evidence, Items: a.localAssistantItems(request, evidence), Conflicts: detectEvidenceConflicts(evidence),
		AIGenerated: true, Partial: false,
		Todo: []string{"客户私有文件正文解析和向量索引等待原始数据格式", "联网工具接入后在工具层强制执行域名白名单"},
	}
	response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "local-retrieval", Status: "ok", Detail: fmt.Sprintf("检索到 %d 条本地结构化证据", len(evidence))})

	evidenceJSON, _ := json.Marshal(evidence)
	knowledgePrompt := fmt.Sprintf("[页面] %s\n[功能] %s\n[用户输入]\n%s\n\n[Go 本地检索证据 JSON]\n%s\n\n请按工作区规则整理证据。不得联网；冲突值必须全部保留。", request.Page, request.Feature, request.Message, evidenceJSON)
	knowledge, err := a.sendAssistantAgent(ctx, knowledgeAgentID, assistantStageSession(request.SessionKey, "knowledge"), knowledgePrompt, []string{"本地业务数据库", "客户私有知识库"}, nil)
	if err != nil {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "failed", knowledge)
		response.Partial = true
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, assistantErrorText(err)))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "failed", Detail: assistantErrorText(err)})
		knowledge.Text = "Knowledge Agent 未返回；Coordinator 将直接接收 Go 检索的原始本地证据。"
	} else {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "knowledge-agent", knowledgeAgentID, "ok", knowledge)
		response.UsedAgents = append(response.UsedAgents, knowledgeAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(knowledgeAgentID, knowledge, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "knowledge-agent", Status: "ok", Detail: "本地证据已整理"})
	}

	domainIDs := assistantDomainRoutes(request)
	domainOutputs := a.runDomainAgents(ctx, request, requestID, knowledge.Text, domainIDs)
	for _, output := range domainOutputs {
		response.Outputs = append(response.Outputs, output)
		if output.Error != "" {
			response.Partial = true
			continue
		}
		response.UsedAgents = append(response.UsedAgents, output.AgentID)
	}
	response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "domain-agents", Status: map[bool]string{true: "partial", false: "ok"}[hasAgentOutputError(domainOutputs)], Detail: strings.Join(domainIDs, ", ")})

	coordinatorInput := map[string]any{
		"page": request.Page, "feature": request.Feature, "userMessage": request.Message,
		"localEvidence": evidence, "knowledgeSummary": knowledge.Text, "domainOutputs": domainOutputs,
		"requiredRule": "冲突信息全部并列展示并保留来源、记录编号和更新时间；不要区分为本地结果区或联网结果区。",
	}
	coordinatorJSON, _ := json.Marshal(coordinatorInput)
	coordinatorPrompt := "请根据工作区规则整合以下 STA-100 查询上下文。只使用已提供内容，不得声称调用了未返回结果的 Agent。\n\n" + string(coordinatorJSON)
	coordinator, coordinatorErr := a.sendAssistantAgent(ctx, coordinatorAgentID, assistantStageSession(request.SessionKey, "coordinator"), coordinatorPrompt, []string{"本地业务数据库", "客户私有知识库"}, nil)
	if coordinatorErr != nil {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "coordinator-agent", coordinatorAgentID, "failed", coordinator)
		response.Partial = true
		response.Text = localAssistantSummary(request, response.Items, evidence, domainOutputs)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(coordinatorAgentID, coordinator, assistantErrorText(coordinatorErr)))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "coordinator-agent", Status: "failed", Detail: assistantErrorText(coordinatorErr)})
	} else {
		a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "coordinator-agent", coordinatorAgentID, "ok", coordinator)
		response.Text = coordinator.Text
		response.UsedAgents = append(response.UsedAgents, coordinatorAgentID)
		response.Outputs = append(response.Outputs, assistantOutputFromResult(coordinatorAgentID, coordinator, ""))
		response.Pipeline = append(response.Pipeline, assistantPipelineStage{Stage: "coordinator-agent", Status: "ok", Detail: "完成统一汇总"})
	}
	response.UsedAgents = uniqueStrings(response.UsedAgents)
	if usage, err := a.tokenUsageSummary(ctx, requestID); err == nil {
		response.TokenUsage = usage
	}
	return response
}

func (a *businessAPI) sendAssistantAgent(ctx context.Context, agentID, sessionKey, message string, sources, allowlist []string) (orchestrator.AgentMessageResult, error) {
	return a.openClaw.service.SendAgentMessage(ctx, orchestrator.AgentMessageInput{AgentID: agentID, Message: message, SessionKey: sessionKey, Sources: sources, Allowlist: allowlist})
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
			prompt := fmt.Sprintf("[STA-100 页面] %s / %s\n[用户输入]\n%s\n\n[Knowledge Agent 本地证据摘要]\n%s\n\n请只从你的专业职责给出结果，保留证据和不确定性。冲突数据不要覆盖。", request.Page, request.Feature, request.Message, knowledgeSummary)
			result, err := a.sendAssistantAgent(ctx, agentID, assistantStageSession(request.SessionKey, agentID), prompt, []string{"本地业务数据库", "客户私有知识库", "联网检索"}, preferences.AgentAllowlists[agentID])
			if err != nil {
				a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "domain-agent", agentID, "failed", result)
				outputs[index] = assistantOutputFromResult(agentID, result, assistantErrorText(err))
				return
			}
			a.recordTokenUsage(ctx, requestID, request.Page, request.Feature, "domain-agent", agentID, "ok", result)
			outputs[index] = assistantOutputFromResult(agentID, result, "")
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

func assistantDomainRoutes(request assistantQueryRequest) []string {
	if target, ok := request.Context["targetAgent"].(string); ok && strings.TrimSpace(target) != "" {
		return []string{strings.TrimSpace(target)}
	}
	switch normalizedFeature(request.Feature) {
	case "oem-match", "oem":
		return []string{"supplier-aggregator", "brand-value-crawler"}
	case "customer-search", "customer-discovery", "customers":
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
		return []string{"market-analyzer", "country-advisor"}
	default:
		return []string{"market-analyzer"}
	}
}

func (a *businessAPI) collectLocalEvidence(ctx context.Context, request assistantQueryRequest) []assistantEvidence {
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

func (a *businessAPI) localAssistantItems(request assistantQueryRequest, evidence []assistantEvidence) []map[string]any {
	items := make([]map[string]any, 0)
	hasContact, _ := request.Context["hasContact"].(bool)
	feature := normalizedFeature(request.Feature)
	if feature != "customer-search" && feature != "customer-discovery" {
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
		if feature == "customer-discovery" && (item.Country != discoveryCountry || item.City != discoveryCity || item.Type != discoveryType) {
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

func localAssistantSummary(request assistantQueryRequest, items []map[string]any, evidence []assistantEvidence, outputs []assistantAgentOutput) string {
	successful := 0
	for _, output := range outputs {
		if output.Error == "" {
			successful++
		}
	}
	return fmt.Sprintf("已从本机业务数据库检索到 %d 条证据和 %d 条可展示记录，%d 个领域 Agent 返回了结果。最终协调 Agent 暂不可用，因此当前为部分结果；请结合证据记录编号和更新时间人工复核。", len(evidence), len(items), successful)
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
	switch {
	case err == context.DeadlineExceeded:
		return "OpenClaw 调用超时"
	case strings.Contains(err.Error(), "agent unavailable"):
		return "Agent 尚未同步到 OpenClaw"
	default:
		return "OpenClaw 调用失败"
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
