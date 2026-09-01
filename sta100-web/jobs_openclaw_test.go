package main

import (
	"encoding/json"
	"strings"
	"testing"

	"openclaw-orchestrator/orchestrator"
)

func TestParseCronBusinessResultAcceptsUnmarkedBusinessJSON(t *testing.T) {
	content := "欧盟电池规则更新要求企业重新核查产品标签、追溯资料和供应链尽调文件。对骑行行业而言，电助力自行车、电池包和售后备件都可能受到影响。出口企业需要把法规节点、库存批次和客户交付计划放在一起复核，并保留原文来源以便人工确认。若企业同时面向德国、法国和荷兰市场，还应同步检查不同国家的本地申报口径和客户沟通话术，避免后续交付节奏被打乱。"
	result, ok := parseCronBusinessResult(`整理完成：
{"schema":"sta100.business.v1","type":"news","items":[{"category":"法规","title":"欧盟电池规则更新","summary":"欧盟电池规则更新，出口企业需要复核标签、追溯资料和交付节点。","content":"` + content + `","source":"EUR-Lex","sourceUrl":"https://example.test/news","time":"2026-08-16","relevance":"高"}]}`)
	if !ok {
		t.Fatal("expected unmarked business JSON to be accepted")
	}
	if result.Type != "news" || len(result.Items) != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestParseCronBusinessResultRejectsPromptExample(t *testing.T) {
	if _, ok := parseCronBusinessResult(`请按以下格式输出：
{"type":"news","items":[{"category":"分类","title":"标题","summary":"摘要","source":"来源","sourceUrl":"原文地址","time":"发布时间","relevance":"相关度"}]}`); ok {
		t.Fatal("prompt example must not be treated as business data")
	}
}

func TestParseCronMarkdownBusinessResultRejectsThinNewsDigest(t *testing.T) {
	if result, ok := parseCronMarkdownBusinessResult("news", `## 一、行业动态（骑行行业）

**1. Accell Group 推出欧洲经销商服务计划**
该计划面向欧洲本地渠道，重点改善交付和售后。
来源：Bike Europe（2026-08-11）· https://example.test/accell

## 二、法规趋势（欧盟）

**2. 欧盟电池法规执行节点更新**
出口企业需要复核产品标签和追溯资料。
来源：EUR-Lex（2026-08-12）· https://example.test/eu-battery`); ok {
		t.Fatalf("thin markdown digest must not be treated as complete news: %+v", result)
	}
}

func TestParseCronMarkdownBusinessResultRejectsPlainBulletDigest(t *testing.T) {
	if result, ok := parseCronMarkdownBusinessResult("news", `# 行业新闻快照

## 行业动态 / 财务
- Fox Factory 上调 2026 展望，但 Q2 营收仍承压 (Bike Europe 2026-08-14)
- Giant 7 月营收复苏加速 (Bike Europe 2026-08-11)

## 法规
- 英国将自行车反倾销税延长 3 年 (Bike Europe 2026-08-12)`); ok {
		t.Fatalf("plain bullet digest must not be treated as complete news: %+v", result)
	}
}

func TestStripSTA100ResultBlocksHandlesMalformedClosingMarker(t *testing.T) {
	text := stripSTA100ResultBlocks(`以上内容仅为展会层面信息，不构成可独立发布的新闻条目。
[STA100_RESULT]{"type":"news","items":[]}[/STA100_RESULT`)
	if text != "以上内容仅为展会层面信息，不构成可独立发布的新闻条目。" {
		t.Fatalf("unexpected text: %q", text)
	}
}

func TestDisplayableRecommendationRejectsMetadataTitle(t *testing.T) {
	item := Recommendation{
		ID:         "REC-AI-META",
		Title:      "类型: 客户机会信号",
		Desc:       "这是一条看起来较长但标题实际是元数据的推荐。",
		Content:    "这是一条看起来较长但标题实际是元数据的推荐，不能作为用户可读推荐进入概览列表，否则会把格式字段误展示成业务内容。",
		Source:     "OpenClaw Agent",
		Type:       "行业推荐",
		DataStatus: "generated",
	}
	if displayableRecommendation(item) {
		t.Fatal("metadata title recommendation must be hidden")
	}
}

func TestDisplayableNewsRejectsPlaceholderSummary(t *testing.T) {
	item := NewsItem{
		ID:         "NEWS-AI-THIN",
		Category:   "行业资讯",
		Title:      "欧洲自行车市场趋于稳定",
		Summary:    "该资讯由 OpenClaw Agent 整理，详情请结合来源复核。",
		Source:     "OpenClaw Agent",
		DataStatus: "generated",
	}
	if displayableNews(item) {
		t.Fatal("placeholder news summary must be hidden")
	}
}

func TestNewsResultRejectsURLOnlySummaryAndContent(t *testing.T) {
	item := newsResultItem{
		Category:  "行业资讯",
		Title:     "Bike Europe 新闻链接",
		Summary:   "链接：https://www.bike-eu.com/article",
		Content:   "https://www.bike-eu.com/article",
		Source:    "Bike Europe",
		SourceURL: "https://www.bike-eu.com/article",
		Time:      "2026-08-18",
		Relevance: "待复核",
	}
	if validNewsResultItem(item) {
		t.Fatal("url-only news fields must be rejected")
	}
	if displayableNews(NewsItem{ID: "NEWS-AI-URL", Category: item.Category, Title: item.Title, Summary: item.Summary, Content: item.Content, Source: item.Source, SourceURL: item.SourceURL, DataStatus: "generated"}) {
		t.Fatal("url-only stored news must not be visible")
	}
}

func TestNewsNarrativeProducesSectionItemsAndDropsIncompleteTail(t *testing.T) {
	item := NewsItem{
		ID:         "NEWS-AI-NARRATIVE",
		Category:   "行业资讯",
		Title:      "行业新闻摘要",
		Summary:    "摘要",
		Source:     "OpenClaw 行业新闻 Agent",
		Time:       "2026-08-18 12:00",
		Relevance:  "高",
		DataStatus: "generated_narrative",
		Content: `# 新闻周报

## 一、Accell 破产冲击波（法国 经销商）

摘要：Accell 旗下法国 Lapierre 进入司法重整，经销商需要复核订单和售后责任。
背景/影响/建议：Lapierre 是法国本土老牌山地车和电助力品牌，渠道网络覆盖欧洲多国。母公司破产会影响库存、质保、零部件供应和信用账期，建议经销商立即确认在途订单、售后条款和替代品牌组合。

## 二、Eurobike 2026 收官：2027 移师法兰克福（德国 经销商）

摘要：Eurobike 官网…`,
	}
	items := filterDisplayableNews([]NewsItem{item})
	if len(items) != 1 {
		t.Fatalf("expected one complete section; narrative is the separate overview, got %d: %+v", len(items), items)
	}
	if items[0].Title != "Accell 破产冲击波" || strings.Contains(items[0].Content, "Eurobike 官网…") {
		t.Fatalf("unexpected derived news items: %+v", items)
	}
	for _, item := range items {
		if item.Title == "本次行业新闻概览" || item.Title == "行业新闻摘要" {
			t.Fatalf("legacy overview must never be exposed as a news item: %+v", item)
		}
	}
}

func TestNewsCategoryFilterUsesContentKeywords(t *testing.T) {
	item := NewsItem{
		Category: "行业资讯",
		Title:    "欧洲经销商调整 E-bike 采购",
		Summary:  "北欧渠道关注轻量化电助力车型。",
		Content:  "德国和瑞典市场的经销商正在调整补货节奏。",
		Source:   "Bike Europe",
	}
	for _, category := range []string{"欧洲市场", "智能骑行", "渠道"} {
		if !newsMatchesCategory(item, category) {
			t.Fatalf("expected category %q to match item", category)
		}
	}
	if newsMatchesCategory(item, "法规") {
		t.Fatal("item should not match unrelated regulation category")
	}
}

func TestNormalizeNewsTextRemovesRepeatedLabelsAndLines(t *testing.T) {
	text := normalizeNewsText("摘要：市场出现变化。\n完整内容：市场出现变化。\n完整内容：市场出现变化。\n后续影响正在评估。")
	if strings.Contains(text, "完整内容：") || strings.Count(text, "市场出现变化。") != 1 {
		t.Fatalf("unexpected normalized news text: %q", text)
	}
}

func TestNormalizeNewsRecordFieldsSwapsBodyLikeTitle(t *testing.T) {
	category, title, summary, content := normalizeNewsRecordFields(
		"Pro-Days：法国下行市场中仍显渠道韧性",
		"完整内容：Bike Europe 在 Pro-Days 展会后撰文指出，法国自行车市场仍处于下行通道。",
		"Bike Europe 在 Pro-Days 展会后撰文指出，法国自行车市场仍处于下行通道。",
		"Bike Europe 在 Pro-Days 展会后撰文指出，法国自行车市场仍处于下行通道。建议经销商重新评估渠道备货节奏与备件供应。",
	)
	if title == "" || title == "行业新闻" || title == "完整内容：Bike Europe 在 Pro-Days 展会后撰文指出，法国自行车市场仍处于下行通道。" {
		t.Fatalf("title was not normalized: category=%q title=%q", category, title)
	}
	if category != "渠道" {
		t.Fatalf("category = %q, want 渠道", category)
	}
	if summary == "" || content == "" {
		t.Fatalf("summary/content should remain available: %q %q", summary, content)
	}
}

func TestNewsStructuredResultPromptIncludesFetchLimitConstraint(t *testing.T) {
	prompt := newsStructuredResultPrompt()
	if strings.Contains(prompt, "maxInt(") {
		t.Fatal("news prompt must not reference missing helper")
	}
	if !strings.Contains(prompt, "每次最多输出 N 条有效新闻") {
		t.Fatalf("news prompt must describe the fetch limit constraint, got: %s", prompt)
	}
}

func TestRecommendationNarrativeProducesReadableItems(t *testing.T) {
	item := Recommendation{
		ID:         "REC-AI-NARRATIVE",
		Title:      "Klever Mobility MBO",
		Source:     "OpenClaw 推荐 Agent",
		Type:       "客户机会信号",
		Time:       "2026-08-18 12:00",
		DataStatus: "generated_narrative",
		Content: `Klever Mobility MBO
为什么推荐：管理层从 Kymco 手中买回，德国 45 km/h 高速电助力车出现新一轮选型窗口。
详情：本轮交易会影响德国高速电助力车品牌与渠道合作结构，建议优先复核供货、账期和替代品牌组合。
来源：Bike Europe
来源链接：https://example.test/klever
类型：客户机会信号
时间：2026-08-18`,
	}
	items := filterDisplayableRecommendations([]Recommendation{item})
	if len(items) != 0 {
		t.Fatalf("expected narrative shell to stay hidden, got %d: %+v", len(items), items)
	}
}

func TestRecommendationFieldBlockNarrativeParsing(t *testing.T) {
	item := Recommendation{
		ID:         "REC-AI-FIELDS",
		Title:      "Accell 破产后德法子公司各自重组",
		Source:     "OpenClaw 推荐 Agent",
		Type:       "行业重大风险",
		DataStatus: "generated_narrative",
		Content: `## 推荐 1
标题：Accell 破产后德法子公司各自重组（行业重大风险）
为什么推荐：Accell 破产影响德法整车品牌归属，属于高优先级风险。
详情：Lapierre 司法重整 + Quanta Capital 拟收购 + DuTech 监管通过，德法整车品牌归属未定，Q4/2027 备货节奏被打乱。对中国出口商而言，这会影响未交付订单、信用账期、样品排期和替代品牌开发优先级，需要在本周内复核 Accell 体系客户清单、未回款风险、法国与德国独立品牌候选名单，并准备重组成功和重组失败两套报价与交付方案。
来源：Bike Europe
来源链接：https://example.test/accell
类型：行业重大风险
时间：2026-08-18`,
	}
	items := filterDisplayableRecommendations([]Recommendation{item})
	if len(items) != 1 {
		t.Fatalf("expected one parsed recommendation item, got %d: %+v", len(items), items)
	}
	if items[0].Title != "Accell 破产后德法子公司各自重组" || items[0].Why == "" || items[0].Detail == "" {
		t.Fatalf("unexpected parsed recommendation fields: %+v", items[0])
	}
	if items[0].Why == items[0].Detail || items[0].Content != items[0].Detail {
		t.Fatalf("recommendation fields were not separated correctly: %+v", items[0])
	}
}

func TestRecommendationMarkdownFallbackSkipsMetadataFieldLines(t *testing.T) {
	result, ok := parseCronMarkdownBusinessResult("recommendations", `## 推荐
- 标题：Klever 独立运营打开德国客户窗口
- 为什么推荐：直接命中德国高速电助力车渠道变化。
- 详情：Klever 管理层回购后，渠道策略和供应链合作可能调整，需要复核德国经销网络与替代供货机会。
- 来源：Bike Europe
- 来源链接：https://example.test/klever`)
	if ok {
		t.Fatalf("metadata field lines must not become standalone recommendations: %+v", result)
	}
}

func TestRecommendationPromptUsesStructuredOutput(t *testing.T) {
	prompt := cronPrompt(Job{
		OutputTarget: "recommendations",
		Prompt:       "根据关注国家和主题生成推荐。",
	})
	if !strings.Contains(prompt, `"type":"recommendations"`) {
		t.Fatalf("recommendation prompt must contain structured result block, got: %s", prompt)
	}
}

func TestRuntimeJobStatusPrioritizesLatestRunFailure(t *testing.T) {
	status := runtimeJobStatus(orchestrator.CronJob{
		Enabled: true,
		Status:  "success",
		State: orchestrator.CronState{
			LastRunStatus: "error",
			LastRunError:  "network connection error",
		},
	})
	if status != "failed" {
		t.Fatalf("status = %q, want failed", status)
	}
}

func TestBuiltInRecommendationJobUsesHiddenSystemAgent(t *testing.T) {
	job := applyBuiltInJobDefaults(Job{ID: "JOB-RECOMMEND", BuiltIn: true, AgentID: "sta100-coordinator"})
	if job.AgentID != "sta100-recommend-curator" {
		t.Fatalf("AgentID = %q, want sta100-recommend-curator", job.AgentID)
	}
}

func TestBuiltInNewsJobUsesHiddenSystemAgent(t *testing.T) {
	job := applyBuiltInJobDefaults(Job{ID: "JOB-NEWS", BuiltIn: true, AgentID: "market-analyzer"})
	if job.AgentID != "sta100-news-curator" {
		t.Fatalf("AgentID = %q, want sta100-news-curator", job.AgentID)
	}
}

func TestBuiltInWeeklyJobUsesHiddenSystemAgent(t *testing.T) {
	job := applyBuiltInJobDefaults(Job{ID: "JOB-WEEKLY", BuiltIn: true, AgentID: "sta100-coordinator"})
	if job.AgentID != "sta100-weekly-report" {
		t.Fatalf("AgentID = %q, want sta100-weekly-report", job.AgentID)
	}
}

// TestValidCronBusinessResultAcceptsEmptyRecommendations 锁定 recommendations 与 news 行为一致：
// 当 Agent 因为没有可靠证据返回 items=[] 时，结果仍应被判定为合法，
// 否则 go 侧会拒收整个结果块，概览页永远看不到任何推荐。
func TestValidCronBusinessResultAcceptsEmptyRecommendations(t *testing.T) {
	if !validCronBusinessResult(cronBusinessResult{Type: "recommendations", Items: nil}) {
		t.Fatal("empty recommendations must be accepted, matching news policy")
	}
	if !validCronBusinessResult(cronBusinessResult{Type: "recommendations", Items: []json.RawMessage{}}) {
		t.Fatal("zero-length recommendations items must be accepted, matching news policy")
	}
}

// TestValidCronBusinessResultRejectsInvalidRecommendations 保留既有严格性：
// 当 items 非空但没有任何一条合规条目时，仍然应该被拒收，避免把示例数据/格式字段
// 当成业务结果。
func TestValidCronBusinessResultRejectsInvalidRecommendations(t *testing.T) {
	raw, err := json.Marshal(map[string]any{"title": "标题", "detail": "占位"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	result := cronBusinessResult{Type: "recommendations", Items: []json.RawMessage{raw}}
	if validCronBusinessResult(result) {
		t.Fatal("recommendations with only metadata-shaped items must be rejected")
	}
}

// TestValidCronBusinessResultAcceptsPluralDetailsField 锁定推荐 Agent
// 使用 "details"（复数）字段时仍能被 Go 解析为有效推荐条目。
func TestValidCronBusinessResultAcceptsPluralDetailsField(t *testing.T) {
	body := strings.Repeat("细节文本，描述业务影响和下一步建议。", 12)
	raw, err := json.Marshal(map[string]any{
		"title":   "Accell 破产影响德法 E-bike 经销商",
		"why":     "德国子公司和法国 Lapierre 同步进入司法程序。",
		"details": body,
		"source":  "Bike Europe",
		"type":    "经销商风险预警",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	result := cronBusinessResult{Type: "recommendations", Items: []json.RawMessage{raw}}
	if !validCronBusinessResult(result) {
		t.Fatal("recommendations with plural details field must be accepted")
	}
}

// TestParseCronBusinessResultPicksLastBlockWhenFirstEmpty 锁定 Agent
// 在多个 STA100_RESULT 块中先交占位空块、后交真实数据块时，Go 端应当
// 拿到最后一个非空块，而不是永远停在第一个空块上。
func TestParseCronBusinessResultPicksLastBlockWhenFirstEmpty(t *testing.T) {
	body := strings.Repeat("细节文本，分析业务影响和下一步建议。", 12)
	final := `[STA100_RESULT]{"type":"recommendations","items":[{"title":"Accell 破产影响德法 E-bike 经销商","detail":"` + body + `","source":"Bike Europe","type":"经销商风险预警"}]}[/STA100_RESULT]`
	summary := "中间思考占位块：\n" +
		`[STA100_RESULT]{"type":"recommendations","items":[]}[/STA100_RESULT]` + "\n" +
		"经过检索，下面给出最终结果：\n" + final
	result, ok := parseCronBusinessResult(summary)
	if !ok {
		t.Fatalf("expected last STA100_RESULT block to win, got %+v", result)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item from last block, got %d", len(result.Items))
	}
}
func TestParseCronMarkdownRecommendationsCarriesDetailField(t *testing.T) {
	body := strings.Repeat("详情段：分析业务影响与下一步建议。", 8)
	summary := "**1. Accell 破产影响德法 E-bike 经销商网络**\n\n" +
		"**为什么推荐：** Accell 旗下品牌深度绑定欧洲零售门店。\n\n" +
		"**详情：** " + body + "\n\n" +
		"**来源：** Bike Europe ｜ **链接：** https://www.bike-eu.com/news\n"
	result, ok := parseCronMarkdownBusinessResult("recommendations", summary)
	if !ok {
		t.Fatalf("markdown recommendations should be parseable, got %+v", result)
	}
	if len(result.Items) == 0 {
		t.Fatal("markdown parser produced no items")
	}
	var item recommendationResultItem
	if err := json.Unmarshal(result.Items[0], &item); err != nil {
		t.Fatalf("unmarshal item: %v", err)
	}
	if item.Detail == "" {
		t.Fatal("markdown item must carry detail field for downstream validation")
	}
	if len([]rune(item.Detail)) < 120 {
		t.Fatalf("markdown item detail too short: %d runes", len([]rune(item.Detail)))
	}
}
