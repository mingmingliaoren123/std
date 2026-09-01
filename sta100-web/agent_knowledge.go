package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	agentKnowledgeKind = "agent_knowledge"
	agentMemoryKind    = "agent_memories"
)

func (a *businessAPI) syncAgentKnowledgeFromLocalSources(ctx context.Context) (int, error) {
	files, err := listRecords[PrivateFile](ctx, a.store, "private_files")
	if err != nil {
		return 0, err
	}
	if _, err := a.store.softDeleteKind(ctx, agentKnowledgeKind); err != nil {
		return 0, err
	}
	now := currentText()
	count := 0
	for agentID, profile := range assistantAgentProfiles {
		content := fmt.Sprintf("Agent=%s；专业方向=%s；知识库范围=%s；工作边界=%s", profile.Name, profile.Focus, profile.KnowledgeScope, profile.Boundaries)
		item := AgentKnowledgeItem{
			ID:              agentKnowledgeID(agentID, "agent-profile", agentID),
			AgentID:         agentID,
			SourceID:        agentID,
			SourceName:      profile.Name,
			SourceKind:      "agent-profile",
			Category:        "专业设定",
			Tags:            []string{"STA-100", "Agent", "专业知识"},
			Content:         content,
			SourceUpdatedAt: now,
			SyncedAt:        now,
			Status:          "synced",
		}
		if err := a.putOrCreate(ctx, agentKnowledgeKind, item.ID, item); err != nil {
			return count, err
		}
		count++
	}
	for _, file := range files {
		for _, agentID := range assistantKnowledgeAgentsForFile(file) {
			item := AgentKnowledgeItem{
				ID:              agentKnowledgeID(agentID, file.ID, file.Updated),
				AgentID:         agentID,
				SourceID:        file.ID,
				SourceName:      file.Name,
				SourceKind:      localKnowledgeSourceKind(file),
				Category:        file.Category,
				Tags:            append([]string{}, file.Tags...),
				Content:         agentKnowledgeContentFromFile(file),
				SourceUpdatedAt: file.Updated,
				SyncedAt:        now,
				Status:          "synced_metadata",
			}
			if err := a.putOrCreate(ctx, agentKnowledgeKind, item.ID, item); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

func agentKnowledgeID(agentID, sourceID, version string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(agentID) + "\x00" + strings.TrimSpace(sourceID) + "\x00" + strings.TrimSpace(version)))
	return "AK-" + hex.EncodeToString(sum[:])[:20]
}

func localKnowledgeSourceKind(file PrivateFile) string {
	source := strings.TrimSpace(file.Source)
	if strings.Contains(source, "通用") || strings.Contains(strings.ToLower(source), "shared") || strings.Contains(strings.ToLower(source), "public") {
		return "shared_knowledge"
	}
	return "private_knowledge"
}

func agentKnowledgeContentFromFile(file PrivateFile) string {
	return fmt.Sprintf("文件=%s；分类=%s；标签=%s；来源=%s；状态=%s；大小=%s；更新时间=%s；文件正文由本机知识库解析、分块并建立向量索引后供对应 Agent 使用。",
		file.Name, file.Category, strings.Join(file.Tags, "、"), file.Source, file.Status, file.Size, file.Updated)
}

func assistantKnowledgeAgentsForFile(file PrivateFile) []string {
	text := strings.ToLower(strings.Join([]string{file.Name, file.Category, file.Source, strings.Join(file.Tags, " ")}, " "))
	agents := map[string]bool{"rag-agent": true}
	add := func(ids ...string) {
		for _, id := range ids {
			if assistantDomainAgentIDs[id] {
				agents[id] = true
			}
		}
	}
	if containsAnyKnowledgeTerm(text, "合同", "agreement", "贸易", "出口", "报关", "contract") {
		add("export-agent", "payment-advisor", "invoice-agent")
	}
	if containsAnyKnowledgeTerm(text, "法规", "regulation", "欧盟", "battery", "cbam", "合规", "ce", "un38") {
		add("cbam-calculator", "country-advisor", "market-analyzer")
	}
	if containsAnyKnowledgeTerm(text, "产品", "手册", "manual", "规格", "bom", "兼容", "compatibility", "shimano", "组件") {
		add("compatibility-agent", "repair-qa", "design-advisor", "inventory-agent")
	}
	if containsAnyKnowledgeTerm(text, "客户", "渠道", "distribution", "经销", "customer", "dealer", "distributor") {
		add("customer-measurement-agent", "b2b-marketplace-agent", "market-analyzer")
	}
	if containsAnyKnowledgeTerm(text, "供应商", "工厂", "oem", "odm", "报价", "supplier") {
		add("supplier-aggregator", "price-tracker")
	}
	if containsAnyKnowledgeTerm(text, "物流", "船", "港", "shipping", "eta") {
		add("shipping-eta")
	}
	if containsAnyKnowledgeTerm(text, "展会", "eurobike", "exhibition") {
		add("exhibition-advisor", "market-analyzer")
	}
	if containsAnyKnowledgeTerm(text, "路线", "route", "地图") {
		add("route-fetcher")
	}
	if containsAnyKnowledgeTerm(text, "赛事", "车队", "race", "team") {
		add("team-race-advisor")
	}
	out := make([]string, 0, len(agents))
	for id := range agents {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func containsAnyKnowledgeTerm(text string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(text, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func (a *businessAPI) extractAndPersistAssistantMemory(ctx context.Context, request assistantQueryRequest, sourceMessageID string) {
	agentID := assistantTargetAgent(request)
	if agentID == "" {
		return
	}
	now := currentText()
	for _, memory := range assistantMemoriesFromRequest(request, agentID, sourceMessageID, now) {
		_ = a.putOrCreate(ctx, agentMemoryKind, memory.ID, memory)
	}
}

func assistantMemoriesFromRequest(request assistantQueryRequest, agentID, sourceMessageID, now string) []AgentMemory {
	memories := make([]AgentMemory, 0, 8)
	add := func(kind, key, value string, confidence float64) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		id := agentMemoryID(agentID, request.SessionKey, kind, key)
		memories = append(memories, AgentMemory{
			ID: id, AgentID: agentID, SessionKey: request.SessionKey, Kind: kind, Key: key, Value: value,
			Confidence: confidence, SourceMessageID: sourceMessageID, UpdatedAt: now,
		})
	}
	if lang := stringContextValue(request.Context, "lang"); lang != "" {
		add("preference", "language", lang, 0.95)
	}
	if model := strings.TrimSpace(request.Model); model != "" {
		add("session", "last_model", model, 0.8)
	}
	if country := stringContextValue(request.Context, "country"); country != "" {
		add("session", "country", country, 0.9)
	}
	if cities := stringArrayContextValue(request.Context, "cities"); len(cities) > 0 {
		add("session", "cities", strings.Join(cities, "、"), 0.9)
	}
	if types := stringArrayContextValue(request.Context, "types"); len(types) > 0 {
		add("session", "customer_types", strings.Join(types, "、"), 0.9)
	}
	if category := stringContextValue(request.Context, "category"); category != "" {
		add("session", "category", category, 0.85)
	}
	for _, value := range explicitMemoryValues(request.Message) {
		add("user_note", "explicit_"+shortMemoryKey(value), value, 0.9)
	}
	return memories
}

func agentMemoryID(agentID, sessionKey, kind, key string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{agentID, sessionKey, kind, key}, "\x00")))
	return "MEM-" + hex.EncodeToString(sum[:])[:20]
}

func explicitMemoryValues(message string) []string {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}
	markers := []string{"记住", "请记住", "以后", "偏好", "默认"}
	for _, marker := range markers {
		if index := strings.Index(message, marker); index >= 0 {
			value := strings.TrimSpace(message[index:])
			if len([]rune(value)) > 120 {
				value = string([]rune(value)[:120])
			}
			return []string{value}
		}
	}
	return nil
}

func shortMemoryKey(value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(value))))
	return hex.EncodeToString(sum[:])[:10]
}
