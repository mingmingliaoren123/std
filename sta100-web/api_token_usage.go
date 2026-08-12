package main

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"openclaw-orchestrator/orchestrator"
)

type tokenUsageRecord struct {
	ID             int64  `json:"id"`
	RequestID      string `json:"requestId"`
	Page           string `json:"page"`
	Feature        string `json:"feature"`
	Stage          string `json:"stage"`
	AgentID        string `json:"agentId"`
	RunID          string `json:"runId,omitempty"`
	Status         string `json:"status"`
	Input          int64  `json:"input"`
	Output         int64  `json:"output"`
	CacheRead      int64  `json:"cacheRead"`
	CacheWrite     int64  `json:"cacheWrite"`
	Total          int64  `json:"total"`
	UsageAvailable bool   `json:"usageAvailable"`
	CreatedAt      string `json:"createdAt"`
}

type tokenUsageSummary struct {
	Input               int64               `json:"input"`
	Output              int64               `json:"output"`
	CacheRead           int64               `json:"cacheRead"`
	CacheWrite          int64               `json:"cacheWrite"`
	Total               int64               `json:"total"`
	Calls               int64               `json:"calls"`
	MeasuredCalls       int64               `json:"measuredCalls"`
	UnavailableCalls    int64               `json:"unavailableCalls"`
	CurrentRequestID    string              `json:"currentRequestId,omitempty"`
	CurrentRequestTotal int64               `json:"currentRequestTotal"`
	ByAgent             []tokenAgentSummary `json:"byAgent"`
	Recent              []tokenUsageRecord  `json:"recent"`
	Temporary           bool                `json:"temporary"`
	UpdatedAt           string              `json:"updatedAt"`
}

type tokenAgentSummary struct {
	AgentID          string `json:"agentId"`
	Calls            int64  `json:"calls"`
	MeasuredCalls    int64  `json:"measuredCalls"`
	UnavailableCalls int64  `json:"unavailableCalls"`
	Input            int64  `json:"input"`
	Output           int64  `json:"output"`
	CacheRead        int64  `json:"cacheRead"`
	CacheWrite       int64  `json:"cacheWrite"`
	Total            int64  `json:"total"`
}

func (a *businessAPI) tokenUsageRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 0 {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "Token 测试统计接口不存在")
		return
	}
	switch r.Method {
	case http.MethodGet:
		summary, err := a.tokenUsageSummary(r.Context(), strings.TrimSpace(r.URL.Query().Get("request_id")))
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, summary)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if _, err := a.store.db.ExecContext(r.Context(), `DELETE FROM agent_token_usage`); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "clear", "agent_token_usage", "all", requestOperator(r), map[string]any{"temporary": true})
		writeJSON(w, http.StatusOK, map[string]any{"cleared": true})
	default:
		w.Header().Set("Allow", "GET, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "只支持读取或清空 Token 测试统计")
	}
}

func (a *businessAPI) recordTokenUsage(ctx context.Context, requestID, page, feature, stage, agentID, status string, result orchestrator.AgentMessageResult) {
	usage := result.Usage
	_, _ = a.store.db.ExecContext(context.Background(), `INSERT INTO agent_token_usage(
		request_id,page,feature,stage,agent_id,run_id,status,input_tokens,output_tokens,cache_read_tokens,cache_write_tokens,total_tokens,usage_available,created_at
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		requestID, page, feature, stage, agentID, result.RunID, status,
		usage.Input, usage.Output, usage.CacheRead, usage.CacheWrite, usage.Total, boolInt(usage.Available()), timeNowUTC())
}

func (a *businessAPI) tokenUsageSummary(ctx context.Context, currentRequestID string) (tokenUsageSummary, error) {
	rows, err := a.store.db.QueryContext(ctx, `SELECT id,request_id,page,feature,stage,agent_id,run_id,status,input_tokens,output_tokens,cache_read_tokens,cache_write_tokens,total_tokens,usage_available,created_at FROM agent_token_usage ORDER BY id DESC`)
	if err != nil {
		return tokenUsageSummary{}, err
	}
	defer rows.Close()
	summary := tokenUsageSummary{CurrentRequestID: currentRequestID, ByAgent: []tokenAgentSummary{}, Recent: []tokenUsageRecord{}, Temporary: true, UpdatedAt: timeNowUTC()}
	byAgent := map[string]*tokenAgentSummary{}
	for rows.Next() {
		var record tokenUsageRecord
		var available int
		if err := rows.Scan(&record.ID, &record.RequestID, &record.Page, &record.Feature, &record.Stage, &record.AgentID, &record.RunID, &record.Status, &record.Input, &record.Output, &record.CacheRead, &record.CacheWrite, &record.Total, &available, &record.CreatedAt); err != nil {
			return tokenUsageSummary{}, err
		}
		if summary.CurrentRequestID == "" {
			summary.CurrentRequestID = record.RequestID
		}
		record.UsageAvailable = available == 1
		summary.Calls++
		summary.Input += record.Input
		summary.Output += record.Output
		summary.CacheRead += record.CacheRead
		summary.CacheWrite += record.CacheWrite
		summary.Total += record.Total
		if record.UsageAvailable {
			summary.MeasuredCalls++
		} else {
			summary.UnavailableCalls++
		}
		if record.RequestID == summary.CurrentRequestID {
			summary.CurrentRequestTotal += record.Total
		}
		agent := byAgent[record.AgentID]
		if agent == nil {
			agent = &tokenAgentSummary{AgentID: record.AgentID}
			byAgent[record.AgentID] = agent
		}
		agent.Calls++
		agent.Input += record.Input
		agent.Output += record.Output
		agent.CacheRead += record.CacheRead
		agent.CacheWrite += record.CacheWrite
		agent.Total += record.Total
		if record.UsageAvailable {
			agent.MeasuredCalls++
		} else {
			agent.UnavailableCalls++
		}
		if len(summary.Recent) < 50 {
			summary.Recent = append(summary.Recent, record)
		}
	}
	if err := rows.Err(); err != nil {
		return tokenUsageSummary{}, err
	}
	for _, agent := range byAgent {
		summary.ByAgent = append(summary.ByAgent, *agent)
	}
	sort.Slice(summary.ByAgent, func(i, j int) bool { return summary.ByAgent[i].Total > summary.ByAgent[j].Total })
	return summary, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
