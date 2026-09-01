package main

import (
	"context"
	"fmt"
)

func (s *businessStore) seed(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM records`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.seedSystemDefaults(ctx)
}

func (s *businessStore) seedSystemDefaults(ctx context.Context) error {
	jobs := []Job{
		{ID: "JOB-RECOMMEND", Name: "每日推荐更新", Kind: "recommendations", Description: "定时汇总推荐设置并调用 OpenClaw 推荐 Agent。", AgentID: "sta100-recommend-curator", Prompt: "根据当前推荐设置生成推荐内容。", Schedule: "每 60 分钟", Enabled: true, BuiltIn: true, Status: "WaitingSource", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-WEEKLY", Name: "智能体周报", Kind: "weekly_report", Description: "调用专门的周报 Agent 汇总本机 Agent 会话、Token 使用和业务审计日志，生成周报草稿。", AgentID: "sta100-weekly-report", Prompt: "汇总最近 7 天 STA-100 智能体使用情况、关键业务操作和待跟进事项。", Schedule: "每周", Enabled: true, BuiltIn: true, Status: "Ready", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-NEWS", Name: "行业新闻更新", Kind: "news", Description: "按客户确认的新闻来源和频率抓取行业新闻。", AgentID: "sta100-news-curator", Prompt: "根据新闻设置整理骑行行业新闻候选。", Schedule: "每 60 分钟", Enabled: true, BuiltIn: true, Status: "WaitingSource", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-INDEX", Name: "数据索引维护", Kind: "index", Description: "扫描本机私有文件元数据，正式正文解析和向量索引等待原始数据格式。", AgentID: "sta100-knowledge", Prompt: "检查本地私有文件是否需要解析、分类、去重和索引。", Schedule: "每天", Enabled: true, BuiltIn: true, Status: "WaitingData", UpdatedAt: "2026-08-10 08:00"},
	}
	plugins := []Plugin{
		{ID: "wechat", Name: "微信", Enabled: false, Status: "Unbound", Capabilities: []string{}, UpdatedAt: "2026-08-10 08:00"},
		{ID: "feishu", Name: "飞书", Enabled: false, Status: "Unbound", Capabilities: []string{}, UpdatedAt: "2026-08-10 08:00"},
	}
	for _, entry := range []struct {
		kind  string
		items any
	}{
		{"jobs", jobs}, {"plugins", plugins},
	} {
		if err := s.seedSlice(ctx, entry.kind, entry.items); err != nil {
			return fmt.Errorf("seed %s: %w", entry.kind, err)
		}
	}
	return s.putSetting(ctx, "preferences", UserPreferences{
		RecommendationEnabled:   true,
		RecommendationShowLimit: 5,
		DiscoveryShowLimit:      10,
		NewsFetchLimit:          5,
		NewsShowLimit:           20,
		NewsFrequency:           "1小时",
		NewsCountries:           "",
		NewsTopics:              "",
		NewsSources:             "",
		AgentAllowlists:         map[string][]string{},
		AgentModelOverrides:     map[string]string{},
	})
}

func (s *businessStore) seedSlice(ctx context.Context, kind string, input any) error {
	items, err := sliceRecords(input)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.create(ctx, kind, item.id, item.value); err != nil {
			return err
		}
	}
	return nil
}

type seedRecord struct {
	id    string
	value any
}

func sliceRecords(input any) ([]seedRecord, error) {
	result := make([]seedRecord, 0)
	switch values := input.(type) {
	case []Customer:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Product:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Quote:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Order:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Document:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Supplier:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []PrivateFile:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []NewsItem:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Recommendation:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Job:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Plugin:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	default:
		return nil, fmt.Errorf("unsupported seed type %T", input)
	}
	return result, nil
}
