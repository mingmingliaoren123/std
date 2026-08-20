package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"openclaw-orchestrator/orchestrator"
)

var jobEveryMinutesPattern = regexp.MustCompile(`每\s*(\d+)\s*分钟`)
var jobEveryHoursPattern = regexp.MustCompile(`每\s*(\d+)\s*小时`)

const businessResultVersion = 8
const structuredBusinessResultSchema = "sta100.business.v1"
const newsWorkspaceOutputDir = "openclaw/workspaces/sta100-news-curator/output"
const recommendationWorkspaceOutputDir = "openclaw/workspaces/sta100-recommend-curator/output"

func (a *businessAPI) jobsRuntimeHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	statusCh := make(chan struct {
		status map[string]any
		err    error
	}, 1)
	go func() {
		status, err := a.openClaw.service.CronStatus(ctx)
		statusCh <- struct {
			status map[string]any
			err    error
		}{status: status, err: err}
	}()
	items, err := a.syncBuiltInJobs(ctx)
	statusResult := <-statusCh
	if statusResult.err != nil && err == nil {
		err = statusResult.err
	}
	if items == nil {
		items, _ = listRecords[Job](ctx, a.store, "jobs")
	}
	response := map[string]any{
		"jobs":           items,
		"openclawStatus": statusResult.status,
		"automation":     a.overviewAutomationData(ctx),
		"synced":         err == nil,
	}
	if err != nil {
		response["syncError"] = userJobSyncError(err)
	}
	writeJSON(w, http.StatusOK, response)
}

func (a *businessAPI) jobRuntimeAction(w http.ResponseWriter, r *http.Request, id, action string) {
	if action == "runs" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		var item Job
		if err := a.store.get(r.Context(), "jobs", id, &item); err != nil {
			writeBusinessError(w, err)
			return
		}
		if item.OpenClawID == "" {
			writeJSON(w, http.StatusOK, map[string]any{"entries": []orchestrator.CronRun{}, "message": "任务尚未同步到 OpenClaw，暂无运行记录"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		result, err := a.openClaw.service.CronRuns(ctx, item.OpenClawID, 50)
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, sanitizeCronRuns(result))
		return
	}
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var item Job
	if err := a.store.get(r.Context(), "jobs", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	if action == "run" {
		if item.OpenClawID == "" {
			var err error
			item, err = a.syncOneJob(ctx, item)
			if err != nil {
				writeOpenClawError(w, err)
				return
			}
			if item.OpenClawID == "" {
				writeAPIError(w, http.StatusConflict, "JOB_NOT_SYNCED", "任务尚未同步到 OpenClaw，暂时不能立即执行")
				return
			}
		}
		result, err := a.openClaw.service.CronRun(ctx, item.OpenClawID)
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		item.Status = "running"
		item.SyncStatus = "synced"
		item.SyncMessage = "已提交 OpenClaw，等待运行结果"
		item.UpdatedAt = currentText()
		_ = a.store.put(r.Context(), "jobs", item.ID, item)
		writeJSON(w, http.StatusOK, map[string]any{"job": item, "run": sanitizeCronRunSubmission(result)})
		return
	}
	if action != "enable" && action != "disable" {
		writeAPIError(w, http.StatusNotFound, "JOB_ACTION_NOT_FOUND", "任务操作不存在")
		return
	}
	updated, err := a.applyJobRuntimeAction(ctx, item, action)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	updated.UpdatedAt = currentText()
	if err := a.store.put(r.Context(), "jobs", updated.ID, updated); err != nil {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *businessAPI) syncBuiltInJobs(ctx context.Context) ([]Job, error) {
	items, err := listRecords[Job](ctx, a.store, "jobs")
	if err != nil {
		return nil, err
	}
	for index := range items {
		if items[index].BuiltIn {
			items[index] = applyBuiltInJobDefaults(items[index])
		}
	}
	runtime, err := a.openClaw.service.CronList(ctx)
	if err != nil {
		for index := range items {
			if items[index].BuiltIn {
				items[index].SyncStatus = "unavailable"
				items[index].SyncMessage = userJobSyncError(err)
				if items[index].Enabled {
					items[index].Status = "unsynced"
				}
				_ = a.store.put(ctx, "jobs", items[index].ID, items[index])
			}
		}
		return items, err
	}
	items = mergeCronRuntime(items, runtime.Jobs)
	byID := make(map[string]orchestrator.CronJob, len(runtime.Jobs))
	byDeclaration := make(map[string]orchestrator.CronJob, len(runtime.Jobs))
	for _, job := range runtime.Jobs {
		byID[job.ID] = job
		if job.DeclarationKey != "" {
			byDeclaration[job.DeclarationKey] = job
		}
	}
	var firstErr error
	for index := range items {
		item := items[index]
		if !item.BuiltIn {
			continue
		}
		runtimeJob, found := byID[item.OpenClawID]
		if !found {
			runtimeJob, found = byDeclaration["sta100:"+item.ID]
		}
		if !found {
			runtimeJob, err = a.createOpenClawJob(ctx, item)
		} else if !cronJobMatches(item, runtimeJob) {
			runtimeJob, err = a.updateOpenClawJob(ctx, item, runtimeJob.ID)
		}
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			items[index].SyncStatus = "error"
			items[index].SyncMessage = userJobSyncError(err)
			items[index].Status = "unsynced"
			_ = a.store.put(ctx, "jobs", item.ID, items[index])
			continue
		}
		items[index] = applyCronRuntime(item, runtimeJob)
		items[index] = a.syncJobBusinessResult(ctx, items[index], runtimeJob)
		_ = a.store.put(ctx, "jobs", item.ID, items[index])
	}
	return items, firstErr
}

func (a *businessAPI) createOpenClawJob(ctx context.Context, item Job) (orchestrator.CronJob, error) {
	return a.openClaw.service.CronAdd(ctx, cronInputFromJob(item))
}

func (a *businessAPI) updateOpenClawJob(ctx context.Context, item Job, openClawID string) (orchestrator.CronJob, error) {
	input := cronInputFromJob(item)
	input.ID = openClawID
	return a.openClaw.service.CronEdit(ctx, input)
}

func (a *businessAPI) syncOneJob(ctx context.Context, item Job) (Job, error) {
	var (
		runtimeJob orchestrator.CronJob
		err        error
	)
	if item.OpenClawID == "" {
		runtimeJob, err = a.createOpenClawJob(ctx, item)
	} else {
		runtimeJob, err = a.updateOpenClawJob(ctx, item, item.OpenClawID)
	}
	if err != nil {
		if errors.Is(err, orchestrator.ErrUnavailable) || errors.Is(err, context.DeadlineExceeded) {
			item.SyncStatus = "unavailable"
			item.SyncMessage = userJobSyncError(err)
			item.Status = "unsynced"
			return item, nil
		}
		return item, err
	}
	item = applyCronRuntime(item, runtimeJob)
	return a.syncJobBusinessResult(ctx, item, runtimeJob), nil
}

func (a *businessAPI) applyJobRuntimeAction(ctx context.Context, item Job, action string) (Job, error) {
	if action == "disable" && item.OpenClawID == "" {
		item.Enabled = false
		item.Status = "disabled"
		item.SyncStatus = firstNonEmpty(item.SyncStatus, "pending")
		item.SyncMessage = "任务已在本机关闭；尚未同步到 OpenClaw，重新开启时会再写入 OpenClaw Cron。"
		item.NextRun = ""
		return item, nil
	}
	if item.OpenClawID == "" {
		item, err := a.syncOneJob(ctx, item)
		if err != nil {
			return item, err
		}
	}
	var (
		runtimeJob orchestrator.CronJob
		err        error
	)
	switch action {
	case "enable":
		runtimeJob, err = a.openClaw.service.CronEnable(ctx, item.OpenClawID, true)
	case "disable":
		runtimeJob, err = a.openClaw.service.CronEnable(ctx, item.OpenClawID, false)
	default:
		return item, fmt.Errorf("unsupported job action %q", action)
	}
	if err != nil {
		return item, err
	}
	item = applyCronRuntime(item, runtimeJob)
	return a.syncJobBusinessResult(ctx, item, runtimeJob), nil
}

func applyBuiltInJobDefaults(item Job) Job {
	defaults := map[string]Job{
		"JOB-RECOMMEND": {
			Name:           "每日推荐更新",
			Kind:           "recommendations",
			Description:    "定时汇总本地推荐缓存，后续接入客户确认的数据源。",
			AgentID:        "sta100-recommend-curator",
			Prompt:         "根据当前关注国家、主题和本地业务数据生成推荐摘要。",
			Schedule:       "每 60 分钟",
			Enabled:        true,
			OutputTarget:   "recommendations",
			BusinessStatus: "waiting",
		},
		"JOB-WEEKLY": {
			Name:           "智能体周报",
			Kind:           "weekly_report",
			Description:    "调用专门的周报 Agent 汇总本机 Agent 会话、Token 使用和业务审计日志，生成周报草稿。",
			AgentID:        "sta100-weekly-report",
			Prompt:         "汇总最近 7 天 STA-100 智能体使用情况、关键业务操作和待跟进事项。",
			Schedule:       "每周",
			Enabled:        true,
			OutputTarget:   "weekly_report",
			BusinessStatus: "waiting",
		},
		"JOB-NEWS": {
			Name:           "行业新闻更新",
			Kind:           "news",
			Description:    "按客户确认的新闻来源和频率抓取行业新闻，当前来源规则待确认。",
			AgentID:        "sta100-news-curator",
			Prompt:         "围绕骑行行业、欧洲渠道、法规和产品趋势整理新闻候选。",
			Schedule:       "每 60 分钟",
			Enabled:        true,
			OutputTarget:   "news",
			BusinessStatus: "waiting",
		},
		"JOB-INDEX": {
			Name:           "数据索引维护",
			Kind:           "index",
			Description:    "扫描本机私有文件元数据，正式正文解析和向量索引等待原始数据格式。",
			AgentID:        "sta100-knowledge",
			Prompt:         "检查本地私有文件是否需要解析、分类、去重和索引。",
			Schedule:       "每天",
			Enabled:        true,
			OutputTarget:   "private_files",
			BusinessStatus: "waiting",
		},
	}
	defaultValue, ok := defaults[item.ID]
	if !ok {
		return item
	}
	if strings.TrimSpace(item.Name) == "" {
		item.Name = defaultValue.Name
	}
	if strings.TrimSpace(item.Kind) == "" {
		item.Kind = defaultValue.Kind
	}
	if strings.TrimSpace(item.Description) == "" {
		item.Description = defaultValue.Description
	}
	if strings.TrimSpace(item.AgentID) == "" {
		item.AgentID = defaultValue.AgentID
	}
	if (item.ID == "JOB-RECOMMEND" || item.ID == "JOB-NEWS" || item.ID == "JOB-WEEKLY") && item.AgentID != defaultValue.AgentID {
		item.AgentID = defaultValue.AgentID
	}
	if strings.TrimSpace(item.Prompt) == "" {
		item.Prompt = defaultValue.Prompt
	}
	if strings.TrimSpace(item.Schedule) == "" {
		item.Schedule = defaultValue.Schedule
	}
	if strings.TrimSpace(item.Status) == "" {
		item.Status = "unsynced"
	}
	if strings.TrimSpace(item.OutputTarget) == "" {
		item.OutputTarget = defaultValue.OutputTarget
	}
	if strings.TrimSpace(item.BusinessStatus) == "" {
		item.BusinessStatus = defaultValue.BusinessStatus
	}
	return item
}

type cronBusinessResult struct {
	Schema    string            `json:"schema,omitempty"`
	Type      string            `json:"type"`
	Items     []json.RawMessage `json:"items,omitempty"`
	Summary   string            `json:"summary,omitempty"`
	Generated string            `json:"generatedAt,omitempty"`
}

type recommendationResultItem struct {
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Why       string `json:"why"`
	Detail    string `json:"detail"`
	Content   string `json:"content"`
	Summary   string `json:"summary"`
	Body      string `json:"body"`
	Reason    string `json:"reason"`
	Source    string `json:"source"`
	SourceURL string `json:"sourceUrl"`
	Type      string `json:"type"`
	Time      string `json:"time"`
}

type newsResultItem struct {
	Category  string `json:"category"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content,omitempty"`
	Source    string `json:"source"`
	SourceURL string `json:"sourceUrl"`
	Time      string `json:"time"`
	Relevance string `json:"relevance"`
}

type newsWorkspaceOutput struct {
	FetchedAt string           `json:"fetched_at,omitempty"`
	Source    string           `json:"source,omitempty"`
	Filters   map[string]any   `json:"filters,omitempty"`
	Items     []newsResultItem `json:"items,omitempty"`
	Notes     string           `json:"notes,omitempty"`
}

type recommendationWorkspaceOutput struct {
	GeneratedAt   string                     `json:"generated_at,omitempty"`
	Trigger       string                     `json:"trigger,omitempty"`
	Filters       map[string]any             `json:"filters,omitempty"`
	SourceRefresh string                     `json:"source_refresh,omitempty"`
	Items         []recommendationResultItem `json:"items,omitempty"`
	Notes         string                     `json:"notes,omitempty"`
}

func (a *businessAPI) syncJobBusinessResult(ctx context.Context, item Job, runtime orchestrator.CronJob) Job {
	if !item.BuiltIn || item.OutputTarget == "" {
		return item
	}
	if runtime.State.LastRunAtMs <= 0 {
		if item.BusinessStatus == "" {
			item.BusinessStatus = "waiting"
		}
		if item.BusinessMessage == "" {
			item.BusinessMessage = "等待 OpenClaw 首次执行，执行完成后自动更新对应业务数据。"
		}
		return item
	}
	reparseOldResult := (item.BusinessStatus == "needs_review" || item.BusinessStatus == "failed") &&
		item.BusinessResultVersion < businessResultVersion
	if runtime.State.LastRunAtMs <= item.LastProcessedRunAtMs && !reparseOldResult {
		return item
	}
	item.BusinessStatus = "syncing"
	item.BusinessMessage = "任务已执行，正在整理结果并更新业务数据。"

	runs, err := a.openClaw.service.CronRuns(ctx, runtime.ID, 20)
	if err != nil {
		item.BusinessStatus = "failed"
		item.BusinessMessage = "任务已执行，但暂时无法读取 OpenClaw 运行结果：" + userJobSyncError(err)
		return item
	}
	item.LastProcessedRunAtMs = runtime.State.LastRunAtMs
	run := latestCronRun(runs.Entries)
	if run == nil {
		if runtimeStatus := cronRuntimeStatus(runtime); isCronFailureStatus(runtimeStatus) {
			item.BusinessStatus = "failed"
			item.BusinessMessage = cronFailureMessage(orchestrator.CronRun{
				Status:      runtimeStatus,
				Error:       firstNonEmpty(runtime.State.LastRunError, runtime.State.LastError),
				ErrorReason: runtime.State.LastErrorReason,
			})
			item.BusinessResultVersion = businessResultVersion
			return item
		}
		item.BusinessStatus = "needs_review"
		item.BusinessMessage = "任务执行状态已更新，但暂无可读取的结果摘要，请查看运行记录。"
		item.BusinessResultVersion = businessResultVersion
		return item
	}
	item.LastRunID = run.RunID
	runStatus := firstNonEmpty(run.Status, cronRuntimeStatus(runtime))
	if isCronFailureStatus(runStatus) || strings.TrimSpace(run.Error) != "" {
		item.BusinessStatus = "failed"
		item.BusinessMessage = cronFailureMessage(*run)
		item.BusinessResultVersion = businessResultVersion
		return item
	}

	switch item.OutputTarget {
	case "recommendations", "news", "weekly_report":
		var result cronBusinessResult
		var ok bool
		if item.OutputTarget == "news" {
			result, ok = parseNewsBusinessResult(run.Summary, runtime.State.LastRunAtMs)
		} else if item.OutputTarget == "recommendations" {
			result, ok = parseRecommendationBusinessResult(run.Summary, runtime.State.LastRunAtMs)
		} else {
			result, ok = parseCronBusinessResult(run.Summary)
			if !ok {
				result, ok = parseCronMarkdownBusinessResult(item.OutputTarget, run.Summary)
			}
		}
		if !ok {
			item.BusinessStatus = "needs_review"
			item.BusinessMessage = "任务已完成，但结果没有包含可自动入库的数据块；原始摘要保留在运行记录中。"
			item.BusinessResultVersion = businessResultVersion
			return item
		}
		switch item.OutputTarget {
		case "recommendations":
			if result.Type != "recommendations" {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = "任务已完成，但结果类型与推荐数据不一致，未自动入库。"
				item.BusinessResultVersion = businessResultVersion
				return item
			}
			count, err := a.persistRecommendations(ctx, item, result.Items)
			if err != nil {
				item.BusinessStatus = "failed"
				item.BusinessMessage = "推荐结果已返回，但写入本地推荐库失败：" + err.Error()
				item.BusinessResultVersion = businessResultVersion
				return item
			}
			if count == 0 {
				item.BusinessStatus = "needs_review"
				item.BusinessUpdatedAt = currentText()
				item.BusinessMessage = "任务已完成，但未返回可自动展示的有效推荐；页面继续展示本机缓存推荐。"
				item.BusinessResultVersion = businessResultVersion
				return item
			}
			if visibleCount, err := a.displayableRecommendationCount(ctx); err == nil {
				count = visibleCount
			}
			if count == 0 {
				item.BusinessStatus = "needs_review"
				item.BusinessUpdatedAt = currentText()
				item.BusinessMessage = "推荐结果已写入，但没有符合页面展示规则的有效推荐；请复核 Agent 输出格式。"
				item.BusinessResultVersion = businessResultVersion
				return item
			}
			item.BusinessStatus = "updated"
			item.BusinessUpdatedAt = currentText()
			item.BusinessMessage = fmt.Sprintf("已更新推荐数据 %d 条，概览页面可直接查看。", count)
			item.BusinessResultVersion = businessResultVersion
		case "news":
			if result.Type == "news" {
				count, err := a.persistNews(ctx, item, result.Items)
				if err != nil {
					item.BusinessStatus = "failed"
					item.BusinessMessage = "新闻结果已返回，但写入本地新闻库失败：" + err.Error()
					item.BusinessResultVersion = businessResultVersion
					return item
				}
				if count > 0 {
					item.BusinessStatus = "updated"
					item.BusinessUpdatedAt = currentText()
					item.BusinessMessage = fmt.Sprintf("已更新行业新闻 %d 条，概览页面可直接查看。", count)
					item.BusinessResultVersion = businessResultVersion
					return item
				}
			}
			item.BusinessStatus = "needs_review"
			item.BusinessUpdatedAt = currentText()
			item.BusinessMessage = "任务已完成，但未返回可自动展示的有效新闻；页面继续展示本机缓存资讯。"
			item.BusinessResultVersion = businessResultVersion
		case "weekly_report":
			if result.Type != "weekly_report" || strings.TrimSpace(result.Summary) == "" {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = "周报已完成，但没有返回可展示的摘要。"
				item.BusinessResultVersion = businessResultVersion
				return item
			}
			item.BusinessStatus = "updated"
			item.BusinessUpdatedAt = currentText()
			item.BusinessMessage = "周报已生成，完整内容请通过运行记录查看。"
			item.LastResult = result.Summary
			item.BusinessResultVersion = businessResultVersion
		}
	case "private_files":
		count, err := a.refreshPrivateFileMetadata(ctx)
		if err != nil {
			item.BusinessStatus = "failed"
			item.BusinessMessage = "文件元数据扫描失败：" + err.Error()
			item.BusinessResultVersion = businessResultVersion
			return item
		}
		item.BusinessStatus = "updated"
		item.BusinessUpdatedAt = currentText()
		item.BusinessMessage = fmt.Sprintf("已完成 %d 个本地文件的元数据复核；正文解析和向量索引仍按数据格式待补充。", count)
		item.BusinessResultVersion = businessResultVersion
	default:
		item.BusinessStatus = "needs_review"
		item.BusinessMessage = "任务已完成，但未配置业务结果接收目标。"
		item.BusinessResultVersion = businessResultVersion
	}
	return item
}

func latestCronRun(entries []orchestrator.CronRun) *orchestrator.CronRun {
	if len(entries) == 0 {
		return nil
	}
	latest := entries[0]
	score := func(run orchestrator.CronRun) int64 {
		if run.RunAtMs > 0 {
			return run.RunAtMs
		}
		if run.EndedAtMs > 0 {
			return run.EndedAtMs
		}
		return run.StartedAtMs
	}
	for index := 1; index < len(entries); index++ {
		if score(entries[index]) > score(latest) {
			latest = entries[index]
		}
	}
	return &latest
}

func sanitizeCronRuns(result orchestrator.CronRunsResult) map[string]any {
	entries := make([]map[string]any, 0, len(result.Entries))
	for _, entry := range result.Entries {
		entries = append(entries, sanitizeCronRun(entry))
	}
	return map[string]any{"entries": entries, "total": result.Total, "hasMore": result.HasMore, "nextOffset": result.NextOffset}
}

func sanitizeCronRunSubmission(result map[string]any) map[string]any {
	message := "请稍后查看运行记录。"
	if value, ok := result["message"].(string); ok && strings.TrimSpace(value) != "" {
		message = strings.TrimSpace(value)
	}
	return map[string]any{
		"submitted": true,
		"message":   message,
	}
}

func sanitizeCronRun(run orchestrator.CronRun) map[string]any {
	status := strings.TrimSpace(run.Status)
	if status == "" {
		status = "unknown"
	}
	message := strings.TrimSpace(run.Error)
	if message == "" {
		message = cronRunBusinessSummary(run.Summary)
	}
	return map[string]any{
		"status":         status,
		"statusLabel":    cronRunStatusLabel(status),
		"message":        message,
		"runAtMs":        firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs),
		"durationMs":     cronRunDurationMs(run),
		"deliveryStatus": cronDeliveryStatusLabel(run.DeliveryStatus),
	}
}

func cronRunBusinessSummary(summary string) string {
	text := strings.TrimSpace(summary)
	if text == "" {
		return "执行已结束，OpenClaw 未返回可展示摘要。"
	}
	if result, ok := parseCronBusinessResult(text); ok {
		switch strings.ToLower(strings.TrimSpace(result.Type)) {
		case "news":
			return fmt.Sprintf("行业新闻结果已返回 %d 条，详情已按规则写入新闻缓存或等待业务复核。", len(result.Items))
		case "recommendations":
			return fmt.Sprintf("为你推荐结果已返回 %d 条，详情已按规则写入推荐缓存或等待业务复核。", len(result.Items))
		case "weekly_report":
			return "智能体周报已返回摘要，完整内容保存在 OpenClaw 运行记录中。"
		}
	}
	if len([]rune(text)) > 120 {
		return string([]rune(text)[:120]) + "..."
	}
	return text
}

func cronRunStatusLabel(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "success", "ok", "completed", "done":
		return "执行成功"
	case "running", "active":
		return "执行中"
	case "failed", "error", "failure":
		return "执行失败"
	case "queued", "pending", "waiting":
		return "等待执行"
	default:
		return "状态待确认"
	}
}

func cronDeliveryStatusLabel(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "none":
		return "未请求推送"
	case "success", "ok", "sent":
		return "已推送"
	case "failed", "error":
		return "推送失败"
	default:
		return status
	}
}

func cronRunDurationMs(run orchestrator.CronRun) int64 {
	if run.StartedAtMs > 0 && run.EndedAtMs > run.StartedAtMs {
		return run.EndedAtMs - run.StartedAtMs
	}
	return 0
}

func firstPositiveInt64(values ...int64) int64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func parseCronBusinessResult(summary string) (cronBusinessResult, bool) {
	startMarker, endMarker := "[STA100_RESULT]", "[/STA100_RESULT]"
	start := strings.Index(summary, startMarker)
	if start >= 0 {
		payload := strings.TrimSpace(summary[start+len(startMarker):])
		if end := strings.Index(payload, endMarker); end >= 0 {
			payload = strings.TrimSpace(payload[:end])
		}
		var result cronBusinessResult
		decoder := json.NewDecoder(strings.NewReader(payload))
		if err := decoder.Decode(&result); err == nil && validCronBusinessResult(result) {
			return result, true
		}
	}
	return parseUnmarkedBusinessResult(summary)
}

func validCronBusinessResult(result cronBusinessResult) bool {
	if !validCronBusinessResultType(result.Type) {
		return false
	}
	if schema := strings.TrimSpace(result.Schema); schema != "" && schema != structuredBusinessResultSchema {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(result.Type)) {
	case "news":
		// An explicit empty result is valid: it means the Agent found no
		// reliable article. A non-empty result must contain at least one
		// complete item; URL-only or truncated items are rejected here.
		if len(result.Items) == 0 {
			return true
		}
		for _, raw := range result.Items {
			var item newsResultItem
			if json.Unmarshal(raw, &item) == nil && validNewsResultItem(item) {
				return true
			}
		}
		return false
	case "recommendations":
		return meaningfulCronBusinessResult(result)
	case "weekly_report":
		return strings.TrimSpace(result.Summary) != ""
	default:
		return false
	}
}

func validCronBusinessResultType(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "recommendations", "news", "weekly_report":
		return true
	default:
		return false
	}
}

// parseUnmarkedBusinessResult accepts a valid business JSON object when the
// model omitted the STA100 markers but still returned a usable payload. It
// deliberately rejects empty or template-like examples so the prompt itself
// cannot become business data.
func parseUnmarkedBusinessResult(summary string) (cronBusinessResult, bool) {
	text := strings.TrimSpace(summary)
	if text == "" {
		return cronBusinessResult{}, false
	}
	candidates := []string{text}
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 3 {
			candidates = append(candidates, strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n")))
		}
	}
	for _, candidate := range candidates {
		var result cronBusinessResult
		if json.Unmarshal([]byte(candidate), &result) == nil && validCronBusinessResult(result) {
			return result, true
		}
	}
	for offset := 0; offset < len(text); offset++ {
		if text[offset] != '{' {
			continue
		}
		var result cronBusinessResult
		decoder := json.NewDecoder(strings.NewReader(text[offset:]))
		if decoder.Decode(&result) == nil && validCronBusinessResult(result) {
			return result, true
		}
	}
	return cronBusinessResult{}, false
}

func meaningfulCronBusinessResult(result cronBusinessResult) bool {
	switch strings.ToLower(strings.TrimSpace(result.Type)) {
	case "weekly_report":
		return strings.TrimSpace(result.Summary) != ""
	case "recommendations", "news":
		for _, raw := range result.Items {
			var item map[string]any
			if json.Unmarshal(raw, &item) != nil {
				continue
			}
			title, _ := item["title"].(string)
			title = strings.TrimSpace(title)
			if title == "" || title == "标题" || title == "新闻标题" {
				continue
			}
			return true
		}
	}
	return false
}

var markdownBusinessItemPattern = regexp.MustCompile(`^\s*(?:[-*]\s*)?\*\*(?:\d+\s*[.、]\s*)?(.+?)\*\*(?:\s+.*)?$`)
var markdownBusinessBulletPattern = regexp.MustCompile(`^\s*[-*]\s+(.+)$`)
var markdownBusinessSourcePattern = regexp.MustCompile(`^\s*来源\s*[:：]\s*(.+)$`)
var markdownBusinessURLPattern = regexp.MustCompile(`https?://[^\s)）]+`)
var markdownBusinessDatePattern = regexp.MustCompile(`[（(]([^）)]+)[）)]`)

type markdownBusinessItem struct {
	title     string
	summary   []string
	category  string
	source    string
	sourceURL string
	time      string
}

func parseCronMarkdownBusinessResult(target, summary string) (cronBusinessResult, bool) {
	if target != "recommendations" && target != "news" {
		return cronBusinessResult{}, false
	}
	summary = stripSTA100ResultBlocks(summary)
	currentCategory := map[string]string{
		"recommendations": "行业推荐",
		"news":            "行业资讯",
	}[target]
	items := make([]markdownBusinessItem, 0)
	hasExplicitSource := false
	var current *markdownBusinessItem
	flush := func() {
		if current == nil || strings.TrimSpace(current.title) == "" {
			return
		}
		current.summary = compactMarkdownSummary(current.summary)
		if current.source == "" {
			current.source = "OpenClaw Agent"
		}
		if current.time == "" {
			current.time = currentText()
		}
		if current.category == "" {
			current.category = currentCategory
		}
		items = append(items, *current)
		current = nil
	}

	for _, rawLine := range strings.Split(summary, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || line == "---" {
			continue
		}
		if strings.Contains(line, "[STA100_RESULT]") || strings.Contains(line, "\"items\"") {
			continue
		}
		plainLine := cleanBusinessMetadataLine(line)
		if strings.HasPrefix(line, "##") {
			currentCategory = markdownSectionCategory(strings.TrimSpace(strings.TrimLeft(line, "#")))
			if currentCategory == "" {
				currentCategory = map[string]string{"recommendations": "行业推荐", "news": "行业资讯"}[target]
			}
			continue
		}
		if current != nil {
			if match := markdownBusinessSourcePattern.FindStringSubmatch(plainLine); len(match) == 2 {
				current.source, current.sourceURL, current.time = parseMarkdownSource(match[1])
				hasExplicitSource = hasExplicitSource || current.source != "" || current.sourceURL != ""
				continue
			}
			if value, ok := metadataFieldValue(plainLine, "链接：", "链接:", "URL：", "URL:", "url:", "link:"); ok {
				if url := firstURL(value); url != "" {
					current.sourceURL = url
					hasExplicitSource = true
				}
				continue
			}
			if value, ok := metadataFieldValue(plainLine, "时间：", "时间:", "发布时间：", "发布时间:", "time:"); ok {
				current.time = firstNonEmpty(cleanMarkdownText(value), current.time)
				continue
			}
			if _, ok := metadataFieldValue(plainLine, "相关度：", "相关度:", "relevance:"); ok {
				current.summary = append(current.summary, cleanMarkdownText(plainLine))
				continue
			}
		} else if businessTitleLooksLikeMetadata(plainLine) {
			continue
		}
		if match := markdownBusinessItemPattern.FindStringSubmatch(line); len(match) == 2 {
			flush()
			title := cleanMarkdownText(match[1])
			if businessTitleLooksLikeMetadata(title) {
				continue
			}
			current = &markdownBusinessItem{title: title, category: currentCategory}
			hasExplicitSource = hasExplicitSource || strings.Contains(title, "http://") || strings.Contains(title, "https://")
			continue
		}
		if match := markdownBusinessBulletPattern.FindStringSubmatch(line); len(match) == 2 {
			flush()
			title := cleanMarkdownText(match[1])
			if businessTitleLooksLikeMetadata(title) {
				continue
			}
			current = &markdownBusinessItem{title: title, category: currentCategory, summary: []string{title}}
			hasExplicitSource = hasExplicitSource || strings.Contains(title, "http://") || strings.Contains(title, "https://") || (strings.Contains(title, "(") && strings.Contains(title, ")"))
			continue
		}
		if current == nil {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		current.summary = append(current.summary, cleanMarkdownText(line))
	}
	flush()

	result := cronBusinessResult{Type: target, Items: make([]json.RawMessage, 0, len(items))}
	for _, item := range items {
		if item.title == "" || item.title == "标题" || item.title == "新闻标题" {
			continue
		}
		var value any
		if target == "news" {
			newsItem := newsResultItem{
				Category: item.category, Title: item.title, Summary: strings.Join(item.summary, " "), Content: strings.Join(item.summary, "\n"),
				Source: item.source, SourceURL: item.sourceURL, Time: item.time, Relevance: "待复核",
			}
			if !validLegacyMarkdownNewsItem(newsItem) {
				continue
			}
			value = newsItem
		} else {
			value = recommendationResultItem{
				Title: item.title, Desc: strings.Join(item.summary, " "),
				Source: item.source, SourceURL: item.sourceURL, Type: item.category, Time: item.time,
			}
		}
		raw, err := json.Marshal(value)
		if err == nil {
			result.Items = append(result.Items, raw)
		}
	}
	return result, len(result.Items) > 0 && hasExplicitSource && validCronBusinessResult(result)
}

func parseNewsBusinessResult(summary string, afterMs int64) (cronBusinessResult, bool) {
	if result, ok := parseCronBusinessResult(summary); ok && result.Type == "news" && len(result.Items) > 0 {
		return result, true
	}
	if result, ok := loadLatestNewsWorkspaceResultSince(afterMs); ok {
		return result, true
	}
	if result, ok := parseCronMarkdownBusinessResult("news", summary); ok {
		return result, true
	}
	if result, ok := parseCronBusinessResult(summary); ok && result.Type == "news" {
		return result, true
	}
	return cronBusinessResult{}, false
}

func parseRecommendationBusinessResult(summary string, afterMs int64) (cronBusinessResult, bool) {
	if result, ok := parseCronBusinessResult(summary); ok && result.Type == "recommendations" && recommendationResultHasValidItems(result) {
		return result, true
	}
	if result, ok := loadLatestRecommendationWorkspaceResultSince(afterMs); ok {
		return result, true
	}
	if result, ok := parseCronMarkdownBusinessResult("recommendations", summary); ok && recommendationResultHasValidItems(result) {
		return result, true
	}
	if result, ok := parseCronBusinessResult(summary); ok && result.Type == "recommendations" {
		return result, true
	}
	return cronBusinessResult{}, false
}

func recommendationResultHasValidItems(result cronBusinessResult) bool {
	if result.Type != "recommendations" {
		return false
	}
	for _, raw := range result.Items {
		var item recommendationResultItem
		if json.Unmarshal(raw, &item) == nil && validRecommendationResultItem(item) {
			return true
		}
	}
	return false
}

func loadLatestRecommendationWorkspaceResultSince(afterMs int64) (cronBusinessResult, bool) {
	path, err := latestWorkspaceOutputFile(recommendationWorkspaceOutputDir)
	if err != nil {
		return cronBusinessResult{}, false
	}
	if afterMs > 0 {
		if info, err := os.Stat(path); err == nil && info.ModTime().UnixMilli() < afterMs {
			return cronBusinessResult{}, false
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cronBusinessResult{}, false
	}
	var output recommendationWorkspaceOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return cronBusinessResult{}, false
	}
	items := make([]json.RawMessage, 0, len(output.Items))
	for _, item := range output.Items {
		if !validRecommendationResultItem(item) {
			continue
		}
		raw, err := json.Marshal(item)
		if err != nil {
			continue
		}
		items = append(items, raw)
	}
	if len(items) == 0 {
		return cronBusinessResult{}, false
	}
	return cronBusinessResult{
		Schema: structuredBusinessResultSchema,
		Type:   "recommendations",
		Items:  items,
	}, true
}

func loadLatestNewsWorkspaceResultSince(afterMs int64) (cronBusinessResult, bool) {
	path, err := latestWorkspaceOutputFile(newsWorkspaceOutputDir)
	if err != nil {
		return cronBusinessResult{}, false
	}
	if afterMs > 0 {
		if info, err := os.Stat(path); err == nil && info.ModTime().UnixMilli() < afterMs {
			return cronBusinessResult{}, false
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cronBusinessResult{}, false
	}
	var output newsWorkspaceOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return cronBusinessResult{}, false
	}
	items := make([]json.RawMessage, 0, len(output.Items))
	for _, item := range output.Items {
		if !validNewsResultItem(item) {
			continue
		}
		raw, err := json.Marshal(item)
		if err != nil {
			continue
		}
		items = append(items, raw)
	}
	if len(items) == 0 {
		return cronBusinessResult{}, false
	}
	return cronBusinessResult{
		Schema: structuredBusinessResultSchema,
		Type:   "news",
		Items:  items,
	}, true
}

func latestWorkspaceOutputFile(dir string) (string, error) {
	candidates := []string{dir, filepath.Clean(filepath.Join("..", dir)), filepath.Clean(filepath.Join("..", "..", dir))}
	seen := map[string]bool{}
	type fileInfo struct {
		path string
		mod  time.Time
	}
	files := make([]fileInfo, 0)
	for _, candidate := range candidates {
		if seen[candidate] {
			continue
		}
		seen[candidate] = true
		entries, err := os.ReadDir(candidate)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files = append(files, fileInfo{path: filepath.Join(candidate, entry.Name()), mod: info.ModTime()})
		}
	}
	if len(files) == 0 {
		return "", fmt.Errorf("news workspace output not found")
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].mod.Equal(files[j].mod) {
			return files[i].path > files[j].path
		}
		return files[i].mod.After(files[j].mod)
	})
	return files[0].path, nil
}

func markdownSectionCategory(section string) string {
	section = strings.TrimSpace(section)
	runes := []rune(section)
	for index, value := range runes {
		if value == '、' || value == '.' {
			section = string(runes[index+1:])
			break
		}
	}
	runes = []rune(section)
	for start, value := range runes {
		if value != '（' && value != '(' {
			continue
		}
		for end := start + 1; end < len(runes); end++ {
			if runes[end] == '）' || runes[end] == ')' {
				section = string(runes[start+1 : end])
				return cleanMarkdownText(section)
			}
		}
		break
	}
	return cleanMarkdownText(section)
}

func parseMarkdownSource(value string) (source, sourceURL, occurredAt string) {
	value = strings.TrimSpace(value)
	if match := markdownBusinessURLPattern.FindString(value); match != "" {
		sourceURL = strings.TrimRight(match, "，。.,")
		value = strings.TrimSpace(strings.TrimSuffix(value, match))
		value = strings.TrimSpace(strings.TrimRight(value, "·|"))
	}
	if match := markdownBusinessDatePattern.FindStringSubmatch(value); len(match) == 2 {
		occurredAt = strings.TrimSpace(match[1])
		value = strings.TrimSpace(value[:strings.Index(value, match[0])])
	}
	source = cleanMarkdownText(strings.TrimSpace(strings.TrimRight(value, "·|")))
	return source, sourceURL, occurredAt
}

func compactMarkdownSummary(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, line)
	}
	return result
}

func cleanMarkdownText(value string) string {
	value = strings.TrimSpace(value)
	value = stripSTA100ResultBlocks(value)
	value = strings.Trim(value, "*_` ")
	return value
}

func cleanBusinessGeneratedTitle(value string) string {
	value = cleanMarkdownText(value)
	value = strings.ReplaceAll(value, "**：", "：")
	value = strings.ReplaceAll(value, "**", "")
	value = strings.TrimSpace(value)
	return value
}

func cleanBusinessMetadataLine(value string) string {
	value = cleanMarkdownText(value)
	return strings.TrimSpace(strings.TrimLeft(value, "-*• "))
}

func metadataFieldValue(value string, prefixes ...string) (string, bool) {
	value = strings.TrimSpace(value)
	lower := strings.ToLower(value)
	for _, prefix := range prefixes {
		prefixLower := strings.ToLower(prefix)
		if strings.HasPrefix(lower, prefixLower) {
			return strings.TrimSpace(value[len(prefix):]), true
		}
	}
	return "", false
}

func businessTitleLooksLikeMetadata(value string) bool {
	value = cleanBusinessMetadataLine(value)
	if value == "" {
		return true
	}
	if _, ok := metadataFieldValue(
		value,
		"标题：", "标题:", "title:",
		"为什么推荐：", "为什么推荐:", "推荐理由：", "推荐理由:", "reason:", "why:",
		"详情：", "详情:", "内容：", "内容:", "完整内容：", "完整内容:", "detail:", "content:",
		"摘要：", "摘要:", "summary:",
		"链接：", "链接:", "URL：", "URL:", "url:", "link:", "来源链接：", "来源链接:", "原文链接：", "原文链接:",
		"来源：", "来源:", "source:",
		"时间：", "时间:", "发布时间：", "发布时间:", "time:",
		"相关度：", "相关度:", "relevance:",
		"类型：", "类型:", "分类：", "分类:", "type:", "category:",
	); ok {
		return true
	}
	return markdownBusinessURLPattern.FindString(value) == value
}

func cleanBusinessGeneratedBody(value string) string {
	value = cleanMarkdownText(value)
	value = strings.TrimSpace(strings.ReplaceAll(value, "**", ""))
	return value
}

func recommendationBody(input recommendationResultItem) string {
	return firstNonEmpty(input.Desc, input.Content, input.Summary, input.Body, input.Detail, input.Reason)
}

func shouldSkipRecommendationResult(input recommendationResultItem) bool {
	if businessTitleLooksLikeMetadata(input.Title) {
		return true
	}
	text := strings.ToLower(strings.TrimSpace(input.Title + " " + recommendationBody(input) + " " + input.Type))
	skipMarkers := []string{
		"ai 生成",
		"证据缺失声明",
		"未随请求附带",
		"本次未检索",
		"没有业务数据文件",
		"没有可靠结果",
		"无可展示结果",
		"source summaries",
		"检索源摘要",
		"运营信号",
		"备注",
		"不确定项",
		"关注主题",
		"关注国家",
		"模板占位",
	}
	for _, marker := range skipMarkers {
		if strings.Contains(text, strings.ToLower(marker)) {
			return true
		}
	}
	return false
}

func shouldSkipNewsResult(input newsResultItem) bool {
	if businessTitleLooksLikeMetadata(input.Title) {
		return true
	}
	if strings.TrimSpace(input.Content) == "" && strings.Contains(input.Summary, "详情请结合来源复核") {
		return true
	}
	text := strings.ToLower(strings.TrimSpace(input.Title + " " + input.Summary + " " + input.Content + " " + input.Category))
	skipMarkers := []string{
		"ai 生成",
		"不构成可独立发布",
		"未纳入 items",
		"没有可靠结果",
		"源摘要",
		"source summaries",
		"备注",
	}
	for _, marker := range skipMarkers {
		if strings.Contains(text, strings.ToLower(marker)) {
			return true
		}
	}
	return false
}

const (
	minNewsSummaryRunes          = 20
	minNewsContentRunes          = 120
	minRecommendationDetailRunes = 120
)

func validRecommendationResultItem(input recommendationResultItem) bool {
	title := cleanBusinessGeneratedTitle(input.Title)
	detail := cleanBusinessGeneratedBody(firstNonEmpty(input.Detail, input.Content, input.Body))
	if title == "" || businessTitleLooksLikeMetadata(title) || shouldSkipRecommendationResult(input) {
		return false
	}
	if detail == "" || len([]rune(detail)) < minRecommendationDetailRunes {
		return false
	}
	if isURLOnlyText(detail) || sameBusinessText(title, detail) {
		return false
	}
	return true
}

func validNewsResultItem(input newsResultItem) bool {
	title := cleanBusinessGeneratedTitle(input.Title)
	summary := normalizeNewsText(input.Summary)
	content := normalizeNewsText(input.Content)
	source := cleanBusinessGeneratedTitle(input.Source)
	sourceURL := strings.TrimSpace(input.SourceURL)

	if title == "" || businessTitleLooksLikeMetadata(title) || shouldSkipNewsResult(input) {
		return false
	}
	if summary == "" || content == "" || source == "" {
		return false
	}
	if len([]rune(summary)) < minNewsSummaryRunes || len([]rune(content)) < minNewsContentRunes {
		return false
	}
	// The source URL has its own field. A URL in summary/content is almost
	// always a truncated or mis-mapped Agent response, so reject it instead
	// of showing a link as the article body.
	if isURLOnlyText(summary) || isURLOnlyText(content) || firstURL(summary) != "" || firstURL(content) != "" {
		return false
	}
	compactSummary := strings.Join(strings.Fields(summary), " ")
	compactContent := strings.Join(strings.Fields(content), " ")
	if strings.EqualFold(compactSummary, compactContent) {
		return false
	}
	if looksIncompleteGeneratedSection(content) {
		return false
	}
	if sourceURL != "" && firstURL(sourceURL) != sourceURL {
		return false
	}
	return true
}

func validLegacyMarkdownNewsItem(input newsResultItem) bool {
	// Markdown is retained only as a recovery path for historical runs. It
	// must pass the same completeness rules as structured JSON.
	return validNewsResultItem(input)
}

func isURLOnlyText(value string) bool {
	text := strings.TrimSpace(strings.Trim(value, "`*_"))
	if text == "" {
		return false
	}
	url := firstURL(text)
	if url == "" {
		return false
	}
	remaining := strings.TrimSpace(strings.Replace(text, url, "", 1))
	remaining = regexp.MustCompile(`(?i)^(?:链接|原文链接|来源链接|url|link)\s*[:：]?\s*$`).ReplaceAllString(remaining, "")
	return strings.TrimSpace(remaining) == ""
}

func displayableRecommendation(item Recommendation) bool {
	if isLegacyDemoRecommendation(item) {
		return false
	}
	content := cleanBusinessGeneratedBody(item.Content)
	status := strings.ToLower(strings.TrimSpace(item.DataStatus))
	if status == "generated_narrative" {
		return false
	}
	if status != "generated_narrative" && status != "generated_section" && len([]rune(content)) < 120 {
		return false
	}
	if status != "generated_narrative" && strings.EqualFold(strings.TrimSpace(item.Title), strings.TrimSpace(content)) {
		return false
	}
	input := recommendationResultItem{
		Title: item.Title, Desc: item.Desc, Why: item.Why, Detail: item.Detail, Content: item.Content, Source: item.Source, Type: item.Type,
	}
	return validRecommendationResultItem(input) && !strings.Contains(item.Title, "**") && !strings.Contains(item.Desc, "[STA100_RESULT]")
}

func displayableNews(item NewsItem) bool {
	if isLegacyDemoNews(item) {
		return false
	}
	input := newsResultItem{
		Category: item.Category, Title: item.Title, Summary: item.Summary, Content: item.Content,
		Source: item.Source, SourceURL: item.SourceURL, Time: item.Time, Relevance: item.Relevance,
	}
	if shouldSkipNewsResult(input) || strings.Contains(item.Summary, "[STA100_RESULT]") {
		return false
	}
	return validNewsResultItem(input)
}

func isLegacyDemoRecommendation(item Recommendation) bool {
	switch strings.TrimSpace(item.ID) {
	case "REC-0001", "REC-0002", "REC-0003":
		return true
	default:
		return false
	}
}

func isLegacyDemoNews(item NewsItem) bool {
	switch strings.TrimSpace(item.ID) {
	case "NEWS-0001", "NEWS-0002", "NEWS-0003", "NEWS-0004", "NEWS-0005":
		return true
	default:
		return false
	}
}

func filterDisplayableRecommendations(items []Recommendation) []Recommendation {
	filtered := make([]Recommendation, 0, len(items))
	seen := map[string]bool{}
	for _, item := range items {
		normalized := normalizeDisplayRecommendation(item)
		if strings.ToLower(strings.TrimSpace(normalized.DataStatus)) != "generated_narrative" &&
			displayableRecommendation(normalized) && !seenDisplayKey(seen, normalized.Title, normalized.Source) {
			filtered = append(filtered, normalized)
		}
		for _, section := range recommendationSectionsFromNarrative(normalized) {
			if displayableRecommendation(section) && !seenDisplayKey(seen, section.Title, section.Source) {
				filtered = append(filtered, section)
			}
		}
	}
	return filtered
}

func filterDisplayableNews(items []NewsItem) []NewsItem {
	filtered := make([]NewsItem, 0, len(items))
	seen := map[string]bool{}
	for _, item := range items {
		normalized := normalizeDisplayNewsItem(item)
		// generated_narrative is a legacy batch response. It is never exposed
		// as an item; complete markdown sections are the only recoverable items.
		if strings.ToLower(strings.TrimSpace(normalized.DataStatus)) != "generated_narrative" &&
			displayableNews(normalized) && !seenDisplayKey(seen, normalized.Title, normalized.Source) {
			filtered = append(filtered, normalized)
		}
		// A legacy long-form Agent response may contain explicit, complete
		// headings. Recover those headings without exposing the parent narrative.
		for _, section := range newsSectionsFromNarrative(normalized) {
			if displayableNews(section) && !seenDisplayKey(seen, section.Title, section.Source) {
				filtered = append(filtered, section)
			}
		}
	}
	return filtered
}

func newsMatchesCategory(item NewsItem, category string) bool {
	label := strings.TrimSpace(category)
	if label == "" || label == "全部" {
		return true
	}
	text := strings.ToLower(strings.Join([]string{
		item.Category,
		item.Title,
		item.Summary,
		item.Content,
		item.Source,
		item.Relevance,
	}, " "))
	terms := map[string][]string{
		"欧洲市场": {"europe", "欧洲", "德国", "法国", "意大利", "西班牙", "荷兰", "波兰", "比利时", "瑞典", "奥地利", "挪威", "丹麦", "芬兰", "葡萄牙", "捷克", "匈牙利", "北欧", "欧盟"},
		"法规":   {"法规", "规则", "监管", "合规", "政策", "directive", "regulation", "eur-lex", "电池法", "反倾销", "认证", "尽职调查"},
		"智能骑行": {"智能骑行", "e-bike", "ebike", "电助力", "电动自行车", "智能设备", "码表", "功率计", "gps", "传感器", "电子变速"},
		"渠道":   {"渠道", "经销商", "分销", "dealer", "distribution", "retail", "零售", "代理", "门店", "售后"},
		"产品":   {"产品", "新品", "launch", "发布", "组件", "电池", "电机", "整车", "配件", "车架", "轮胎", "传动", "头盔"},
	}[label]
	if len(terms) == 0 {
		terms = []string{strings.ToLower(label)}
	}
	for _, term := range terms {
		if strings.Contains(text, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func seenDisplayKey(seen map[string]bool, title, source string) bool {
	key := strings.ToLower(strings.TrimSpace(title + "\x00" + source))
	if key == "\x00" {
		return true
	}
	if seen[key] {
		return true
	}
	seen[key] = true
	return false
}

func normalizeDisplayNewsItem(item NewsItem) NewsItem {
	item.Category, item.Title, item.Summary, item.Content = normalizeNewsRecordFields(item.Category, item.Title, item.Summary, item.Content)
	if strings.ToLower(strings.TrimSpace(item.DataStatus)) == "generated_narrative" {
		item.Content = sanitizeGeneratedNarrative(item.Content)
		item.Summary = ""
	}
	return item
}

func normalizeNewsRecordFields(category, title, summary, content string) (string, string, string, string) {
	category = cleanBusinessGeneratedTitle(category)
	title = cleanBusinessGeneratedTitle(title)
	summary = normalizeNewsText(summary)
	content = normalizeNewsText(content)

	canonical := canonicalNewsCategory(category, title, summary, content)
	if newsTitleLooksLikeBody(title) && newsCategoryLooksLikeTitle(category) {
		titleCandidate := category
		title = titleCandidate
		canonical = canonicalNewsCategory(titleCandidate, titleCandidate, summary, content)
	}
	if canonical == "" {
		canonical = "行业资讯"
	}
	if title == "" || newsTitleLooksLikeBody(title) {
		title = newsFallbackTitle(summary, content, category)
	}
	return canonical, title, summary, content
}

func newsCategoryLooksLikeTitle(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if _, ok := map[string]struct{}{"全部": {}, "欧洲市场": {}, "法规": {}, "智能骑行": {}, "渠道": {}, "产品": {}}[value]; ok {
		return true
	}
	if len([]rune(value)) < 8 || len([]rune(value)) > 70 {
		return false
	}
	return strings.Contains(value, "：") || strings.Contains(value, ":") || strings.Count(value, "、") <= 1
}

func newsTitleLooksLikeBody(value string) bool {
	text := strings.TrimSpace(cleanBusinessGeneratedTitle(value))
	if text == "" {
		return true
	}
	if strings.HasPrefix(text, "完整内容") || strings.HasPrefix(text, "摘要") || strings.HasPrefix(text, "详细内容") || strings.HasPrefix(text, "内容") {
		return true
	}
	if len([]rune(text)) > 70 {
		return true
	}
	if strings.Count(text, "。") >= 2 || strings.Count(text, "；") >= 2 {
		return true
	}
	return false
}

func newsFallbackTitle(summary, content, category string) string {
	for _, value := range []string{summary, content, category} {
		text := strings.TrimSpace(cleanBusinessGeneratedTitle(value))
		if text == "" {
			continue
		}
		if len([]rune(text)) <= 80 && !newsTitleLooksLikeBody(text) {
			return text
		}
		if idx := strings.IndexAny(text, "。.!！？"); idx > 0 {
			candidate := strings.TrimSpace(text[:idx])
			if len([]rune(candidate)) >= 4 && len([]rune(candidate)) <= 80 {
				return candidate
			}
		}
	}
	return "行业新闻"
}

func canonicalNewsCategory(category, title, summary, content string) string {
	if value := strings.TrimSpace(category); value != "" {
		switch value {
		case "欧洲市场", "法规", "智能骑行", "渠道", "产品":
			return value
		}
	}
	text := strings.ToLower(strings.Join([]string{title, summary, content}, " "))
	rules := []struct {
		label string
		terms []string
	}{
		{"法规", []string{"法规", "规则", "监管", "合规", "policy", "directive", "regulation", "eur-lex", "电池法", "反倾销", "认证", "尽职调查"}},
		{"智能骑行", []string{"智能骑行", "e-bike", "ebike", "电助力", "电动自行车", "智能设备", "码表", "功率计", "gps", "传感器", "电子变速"}},
		{"渠道", []string{"渠道", "经销商", "分销", "dealer", "distribution", "retail", "零售", "代理", "门店", "售后"}},
		{"产品", []string{"产品", "新品", "launch", "发布", "组件", "电池", "电机", "整车", "配件", "车架", "轮胎", "传动", "头盔"}},
		{"欧洲市场", []string{"欧洲", "europe", "德国", "法国", "意大利", "西班牙", "荷兰", "波兰", "比利时", "瑞典", "奥地利", "挪威", "丹麦", "芬兰", "葡萄牙", "捷克", "匈牙利", "北欧", "欧盟"}},
	}
	for _, rule := range rules {
		for _, term := range rule.terms {
			if strings.Contains(text, strings.ToLower(term)) {
				return rule.label
			}
		}
	}
	return "行业资讯"
}

func normalizeDisplayRecommendation(item Recommendation) Recommendation {
	if strings.ToLower(strings.TrimSpace(item.DataStatus)) == "generated_narrative" {
		item.Content = sanitizeGeneratedNarrative(item.Content)
		item.Desc = newsBrief(item.Content, 120)
	}
	return item
}

func sanitizeGeneratedNarrative(value string) string {
	text := strings.TrimSpace(value)
	if text == "" {
		return ""
	}
	sections := splitMarkdownBusinessSections(text)
	if len(sections) == 0 {
		return text
	}
	prefix := strings.TrimSpace(text[:sections[0].start])
	parts := make([]string, 0, len(sections)+1)
	if prefix != "" {
		parts = append(parts, prefix)
	}
	for _, section := range sections {
		if looksIncompleteGeneratedSection(section.body) {
			continue
		}
		parts = append(parts, strings.TrimSpace(section.raw))
	}
	return strings.TrimSpace(strings.Join(parts, "\n\n"))
}

type markdownNarrativeSection struct {
	index int
	start int
	title string
	body  string
	raw   string
}

var markdownNarrativeHeadingPattern = regexp.MustCompile(`(?m)^##\s+(.+)$`)
var numberedRecommendationPattern = regexp.MustCompile(`(?m)^\s*(\d+)[.、]\s+(.+?)(?:\s+-\s+(.+))?$`)

func splitMarkdownBusinessSections(text string) []markdownNarrativeSection {
	matches := markdownNarrativeHeadingPattern.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return nil
	}
	sections := make([]markdownNarrativeSection, 0, len(matches))
	for index, match := range matches {
		start := match[0]
		bodyStart := match[1]
		end := len(text)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		title := strings.TrimSpace(text[match[2]:match[3]])
		raw := strings.TrimSpace(text[start:end])
		body := strings.TrimSpace(text[bodyStart:end])
		sections = append(sections, markdownNarrativeSection{index: index + 1, start: start, title: cleanMarkdownText(title), body: body, raw: raw})
	}
	return sections
}

func looksIncompleteGeneratedSection(value string) bool {
	text := strings.TrimSpace(value)
	if text == "" {
		return true
	}
	if strings.Contains(text, "背景/影响/建议") || strings.Contains(text, "为什么推荐") || strings.Contains(text, "影响:") || strings.Contains(text, "下一步:") {
		return false
	}
	runes := []rune(text)
	if len(runes) < 80 {
		return true
	}
	return strings.HasSuffix(text, "…") || strings.HasSuffix(text, "...")
}

func newsSectionsFromNarrative(item NewsItem) []NewsItem {
	if strings.ToLower(strings.TrimSpace(item.DataStatus)) != "generated_narrative" {
		return nil
	}
	sections := splitMarkdownBusinessSections(item.Content)
	result := make([]NewsItem, 0, len(sections))
	for _, section := range sections {
		if looksIncompleteGeneratedSection(section.body) {
			continue
		}
		summary := extractMarkdownSummaryLine(section.body)
		result = append(result, NewsItem{
			ID:          fmt.Sprintf("%s-S%02d", item.ID, section.index),
			Category:    firstNonEmpty(extractCategoryFromTitle(section.title), item.Category, "行业资讯"),
			Title:       stripParentheticalType(stripNarrativeOrdinal(section.title)),
			Summary:     firstNonEmpty(summary, newsBrief(section.body, 120)),
			Content:     strings.TrimSpace(section.body),
			Source:      item.Source,
			SourceURL:   item.SourceURL,
			Time:        item.Time,
			Relevance:   firstNonEmpty(item.Relevance, "待复核"),
			UpdatedAt:   item.UpdatedAt,
			GeneratedBy: item.GeneratedBy,
			DataStatus:  "generated_section",
		})
	}
	return result
}

func recommendationSectionsFromNarrative(item Recommendation) []Recommendation {
	if strings.ToLower(strings.TrimSpace(item.DataStatus)) != "generated_narrative" {
		return nil
	}
	if sections := recommendationFieldSectionsFromNarrative(item); len(sections) > 0 {
		return sections
	}
	matches := numberedRecommendationPattern.FindAllStringSubmatch(item.Content, -1)
	result := make([]Recommendation, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		title := cleanMarkdownText(match[2])
		body := cleanMarkdownText(match[3])
		if body == "" || len([]rune(body)) < 25 {
			continue
		}
		content := strings.TrimSpace(title + "\n\n" + body)
		result = append(result, Recommendation{
			ID:          fmt.Sprintf("%s-R%02s", item.ID, match[1]),
			Title:       stripParentheticalType(title),
			Desc:        newsBrief(body, 120),
			Content:     content,
			Source:      item.Source,
			SourceURL:   item.SourceURL,
			Type:        firstNonEmpty(extractParentheticalType(title), item.Type, "行业推荐"),
			Time:        item.Time,
			UpdatedAt:   item.UpdatedAt,
			GeneratedBy: item.GeneratedBy,
			DataStatus:  "generated_section",
		})
	}
	return result
}

func recommendationFieldSectionsFromNarrative(item Recommendation) []Recommendation {
	lines := strings.Split(item.Content, "\n")
	blocks := make([][]string, 0)
	current := make([]string, 0)
	hasTitle := false
	flush := func() {
		if len(current) == 0 || !hasTitle {
			current = current[:0]
			hasTitle = false
			return
		}
		block := make([]string, len(current))
		copy(block, current)
		blocks = append(blocks, block)
		current = current[:0]
		hasTitle = false
	}
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			if len(current) > 0 {
				current = append(current, "")
			}
			continue
		}
		if line == "##" || strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") || strings.HasPrefix(line, "#### ") {
			continue
		}
		if strings.HasPrefix(line, "标题：") || strings.HasPrefix(line, "标题:") {
			flush()
			current = append(current, line)
			hasTitle = true
			continue
		}
		if len(current) == 0 {
			continue
		}
		current = append(current, line)
	}
	flush()
	if len(blocks) == 0 {
		return nil
	}
	result := make([]Recommendation, 0, len(blocks))
	for index, block := range blocks {
		rec := recommendationFromFieldBlock(item, block, index+1)
		if rec.Title == "" {
			continue
		}
		result = append(result, rec)
	}
	return result
}

func recommendationFromFieldBlock(item Recommendation, block []string, index int) Recommendation {
	rec := Recommendation{
		ID:          fmt.Sprintf("%s-F%02d", item.ID, index),
		Source:      item.Source,
		SourceURL:   item.SourceURL,
		Type:        firstNonEmpty(item.Type, "行业推荐"),
		Time:        item.Time,
		UpdatedAt:   item.UpdatedAt,
		GeneratedBy: item.GeneratedBy,
		DataStatus:  "generated_section",
	}
	for _, rawLine := range block {
		line := strings.TrimSpace(rawLine)
		if value, ok := fieldLineValue(line, "标题：", "标题:"); ok {
			rec.Title = stripParentheticalType(cleanMarkdownText(value))
			continue
		}
		if value, ok := fieldLineValue(line, "为什么推荐：", "为什么推荐:", "推荐理由：", "推荐理由:", "why:", "reason:"); ok {
			rec.Why = cleanBusinessGeneratedBody(value)
			continue
		}
		if value, ok := fieldLineValue(line, "详情：", "详情:", "内容：", "内容:", "完整内容：", "完整内容:", "detail:", "content:"); ok {
			rec.Detail = cleanBusinessGeneratedBody(value)
			continue
		}
		if value, ok := fieldLineValue(line, "来源链接：", "来源链接:", "原文链接：", "原文链接:", "链接：", "链接:"); ok {
			rec.SourceURL = strings.TrimSpace(value)
			continue
		}
		if value, ok := fieldLineValue(line, "来源：", "来源:"); ok {
			rec.Source = cleanBusinessGeneratedTitle(value)
			continue
		}
		if value, ok := fieldLineValue(line, "类型：", "类型:"); ok {
			rec.Type = cleanBusinessGeneratedTitle(value)
			continue
		}
		if value, ok := fieldLineValue(line, "时间：", "时间:"); ok {
			rec.Time = cleanBusinessGeneratedTitle(value)
			continue
		}
		if rec.Detail != "" {
			rec.Detail = strings.TrimSpace(rec.Detail + "\n" + cleanMarkdownText(line))
		}
	}
	if sameBusinessText(rec.Why, rec.Detail) {
		rec.Why = ""
	}
	rec.Desc = newsBrief(firstNonEmpty(rec.Why, rec.Detail), 100)
	rec.Content = strings.TrimSpace(firstNonEmpty(rec.Detail, rec.Why))
	return rec
}

func sameBusinessText(left, right string) bool {
	left = strings.Join(strings.Fields(cleanBusinessGeneratedBody(left)), " ")
	right = strings.Join(strings.Fields(cleanBusinessGeneratedBody(right)), " ")
	return left != "" && strings.EqualFold(left, right)
}

func fieldLineValue(line string, prefixes ...string) (string, bool) {
	line = strings.TrimSpace(strings.TrimLeft(line, "-*• "))
	lower := strings.ToLower(line)
	for _, prefix := range prefixes {
		prefixLower := strings.ToLower(prefix)
		if strings.HasPrefix(lower, prefixLower) {
			return strings.TrimSpace(line[len(prefix):]), true
		}
	}
	return "", false
}

func extractMarkdownSummaryLine(value string) string {
	for _, line := range strings.Split(value, "\n") {
		line = cleanMarkdownText(line)
		if summary, ok := metadataFieldValue(line, "摘要：", "摘要:", "summary:"); ok {
			return cleanBusinessGeneratedBody(summary)
		}
	}
	return ""
}

func stripNarrativeOrdinal(value string) string {
	text := cleanMarkdownText(value)
	text = regexp.MustCompile(`^[一二三四五六七八九十]+[、.]\s*`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`^\d+[、.]\s*`).ReplaceAllString(text, "")
	return text
}

func extractCategoryFromTitle(value string) string {
	inner, _ := parentheticalTail(cleanMarkdownText(value))
	return inner
}

func extractParentheticalType(value string) string {
	inner, _ := parentheticalTail(cleanMarkdownText(value))
	return inner
}

func stripParentheticalType(value string) string {
	text := cleanMarkdownText(value)
	_, stripped := parentheticalTail(text)
	if stripped != "" {
		return stripped
	}
	return text
}

func parentheticalTail(text string) (inner, stripped string) {
	text = strings.TrimSpace(text)
	start := strings.LastIndexAny(text, "（(")
	end := strings.LastIndexAny(text, "）)")
	if start < 0 || end <= start {
		return "", ""
	}
	closeLen := 1
	if strings.HasPrefix(text[end:], "）") {
		closeLen = len("）")
	}
	if strings.TrimSpace(text[end+closeLen:]) != "" {
		return "", ""
	}
	openLen := 1
	if strings.HasPrefix(text[start:], "（") {
		openLen = len("（")
	}
	return strings.TrimSpace(text[start+openLen : end]), strings.TrimSpace(text[:start])
}

func stripSTA100ResultBlocks(value string) string {
	text := strings.TrimSpace(value)
	for {
		start := strings.Index(text, "[STA100_RESULT]")
		if start < 0 {
			break
		}
		afterStart := start + len("[STA100_RESULT]")
		end := strings.Index(text[afterStart:], "[/STA100_RESULT")
		if end < 0 {
			text = strings.TrimSpace(text[:start])
			break
		}
		end += afterStart
		closeEnd := strings.Index(text[end:], "]")
		if closeEnd < 0 {
			text = strings.TrimSpace(text[:start])
			break
		}
		text = strings.TrimSpace(text[:start] + text[end+closeEnd+1:])
	}
	return text
}

func (a *businessAPI) persistRecommendations(ctx context.Context, job Job, rawItems []json.RawMessage) (int, error) {
	existing, err := listRecords[Recommendation](ctx, a.store, "recommendations")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, raw := range rawItems {
		var input recommendationResultItem
		if err := json.Unmarshal(raw, &input); err != nil || strings.TrimSpace(input.Title) == "" {
			continue
		}
		input.Title = cleanBusinessGeneratedTitle(input.Title)
		input.Why = cleanBusinessGeneratedBody(firstNonEmpty(input.Why, input.Reason))
		input.Detail = cleanBusinessGeneratedBody(firstNonEmpty(input.Detail, input.Content, input.Body))
		if sameBusinessText(input.Why, input.Detail) {
			input.Why = ""
		}
		brief := cleanBusinessGeneratedBody(firstNonEmpty(input.Desc, input.Why))
		full := cleanBusinessGeneratedBody(firstNonEmpty(input.Detail, input.Content, input.Body, input.Why, input.Reason))
		input.Source = cleanBusinessGeneratedTitle(firstNonEmpty(input.Source, "OpenClaw Agent"))
		input.Type = cleanBusinessGeneratedTitle(firstNonEmpty(input.Type, "行业推荐"))
		input.Desc = firstNonEmpty(brief, newsBrief(firstNonEmpty(input.Why, full), 80))
		if !validRecommendationResultItem(input) {
			continue
		}
		input.Desc = firstNonEmpty(input.Desc, "该推荐由 OpenClaw Agent 生成，详情请结合来源复核。")
		full = firstNonEmpty(full, input.Detail, input.Why, input.Desc)
		input.Time = firstNonEmpty(input.Time, currentText())
		record := Recommendation{
			ID:    generatedBusinessID("REC", input.Title, input.Source),
			Title: input.Title, Desc: input.Desc, Why: input.Why, Detail: input.Detail, Content: full, Source: input.Source, Type: input.Type,
			SourceURL: input.SourceURL, Time: input.Time, UpdatedAt: currentText(), GeneratedBy: job.AgentID, DataStatus: "generated",
		}
		for _, old := range existing {
			if old.Title == record.Title && old.Source == record.Source {
				record.ID = old.ID
				break
			}
		}
		if err := a.putOrCreate(ctx, "recommendations", record.ID, record); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (a *businessAPI) displayableRecommendationCount(ctx context.Context) (int, error) {
	items, err := listRecords[Recommendation](ctx, a.store, "recommendations")
	if err != nil {
		return 0, err
	}
	return len(filterDisplayableRecommendations(items)), nil
}

func (a *businessAPI) persistRecommendationNarrative(ctx context.Context, job Job, summary string) (int, error) {
	narrative := cleanBusinessGeneratedBody(stripSTA100ResultBlocks(summary))
	if !usableRecommendationNarrative(narrative) {
		return 0, nil
	}
	brief := strings.Join(strings.Fields(narrative), " ")
	if len([]rune(brief)) > 100 {
		brief = string([]rune(brief)[:100]) + "…"
	}
	record := Recommendation{
		ID:          generatedBusinessID("REC", "推荐汇总", job.AgentID),
		Title:       "推荐汇总",
		Desc:        brief,
		Why:         "OpenClaw 推荐 Agent 汇总本轮关键变化。",
		Detail:      narrative,
		Content:     narrative,
		Source:      "OpenClaw 推荐 Agent",
		Type:        "智能推荐摘要",
		Time:        currentText(),
		UpdatedAt:   currentText(),
		GeneratedBy: job.AgentID,
		DataStatus:  "generated_narrative",
	}
	if err := a.putOrCreate(ctx, "recommendations", record.ID, record); err != nil {
		return 0, err
	}
	return 1, nil
}

func usableRecommendationNarrative(value string) bool {
	text := strings.TrimSpace(value)
	if len([]rune(text)) < 80 {
		return false
	}
	lower := strings.ToLower(text)
	for _, marker := range []string{
		"模型调用超时",
		"request timed out",
		"network connection error",
		"failovererror",
		"没有返回可自动",
		"未生成新闻",
		"空跑",
		"无输入数据",
	} {
		if strings.Contains(lower, strings.ToLower(marker)) {
			return false
		}
	}
	return true
}

func (a *businessAPI) runBuiltInBusinessJobNow(ctx context.Context, id string) (map[string]any, error) {
	persistCtx := context.Background()
	preferences := defaultPreferences()
	_ = a.store.getSetting(ctx, "preferences", &preferences)
	_, _ = a.syncRecommendationJobs(ctx, preferences)
	var item Job
	if err := a.store.get(ctx, "jobs", id, &item); err != nil {
		return nil, err
	}
	item = applyBuiltInJobDefaults(item)
	if item.OpenClawID == "" {
		updated, err := a.syncOneJob(ctx, item)
		if err != nil {
			return nil, err
		}
		item = updated
		if item.OpenClawID == "" {
			return nil, fmt.Errorf("任务尚未同步到 OpenClaw，暂时不能手动刷新")
		}
	}
	startedAtMs := time.Now().Add(-5 * time.Second).UnixMilli()
	submission, err := a.openClaw.service.CronRun(ctx, item.OpenClawID)
	if err != nil {
		return nil, err
	}
	item.Status = "running"
	item.SyncStatus = "synced"
	item.SyncMessage = fmt.Sprintf("已提交 OpenClaw，正在等待%s返回结果", businessJobAgentLabel(item.OutputTarget))
	item.UpdatedAt = currentText()
	_ = a.store.put(persistCtx, "jobs", item.ID, item)

	run, waitErr := a.waitForCronRun(ctx, item.OpenClawID, startedAtMs)
	if waitErr != nil {
		if item.OutputTarget == "recommendations" {
			if result, ok := loadLatestRecommendationWorkspaceResultSince(startedAtMs); ok {
				count, persistErr := a.persistRecommendations(persistCtx, item, result.Items)
				if persistErr != nil {
					item.BusinessStatus = "failed"
					item.BusinessMessage = "推荐 Agent 已写出结果文件，但写入本机推荐库失败：" + persistErr.Error()
					item.BusinessResultVersion = businessResultVersion
					item.UpdatedAt = currentText()
					_ = a.store.put(persistCtx, "jobs", item.ID, item)
					return nil, persistErr
				}
				visibleCount, _ := a.displayableRecommendationCount(persistCtx)
				if visibleCount > 0 {
					count = visibleCount
				}
				item.BusinessStatus = "updated"
				item.BusinessMessage = fmt.Sprintf("OpenClaw 状态仍在同步中，已从推荐 Agent 输出文件恢复 %d 条推荐。", count)
				item.BusinessUpdatedAt = currentText()
				item.BusinessResultVersion = businessResultVersion
				item.LastResult = "running"
				item.UpdatedAt = currentText()
				_ = a.store.put(persistCtx, "jobs", item.ID, item)
				return map[string]any{"submitted": true, "completed": false, "fallbackUsed": true, "job": item, "run": sanitizeCronRunSubmission(submission), "message": item.BusinessMessage, "updatedCount": count}, nil
			}
		}
		if item.OutputTarget == "news" {
			recovered, recoveredRun, recoverErr := a.recoverLatestSuccessfulNewsRun(persistCtx, item, "")
			if recoverErr != nil {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = waitErr.Error() + "；同时读取最近成功新闻记录失败：" + userJobSyncError(recoverErr)
				item.LastResult = "running"
				item.UpdatedAt = currentText()
				_ = a.store.put(persistCtx, "jobs", item.ID, item)
				return map[string]any{"submitted": true, "completed": false, "job": item, "run": sanitizeCronRunSubmission(submission), "message": item.BusinessMessage}, nil
			}
			if recovered > 0 {
				item.BusinessStatus = "updated"
				item.BusinessMessage = fmt.Sprintf("本次行业新闻刷新尚未返回有效结果，已恢复最近一次成功运行的 %d 条新闻。", recovered)
				item.BusinessUpdatedAt = currentText()
				item.BusinessResultVersion = businessResultVersion
				if recoveredRun != nil {
					item.LastRunID = recoveredRun.RunID
					item.LastProcessedRunAtMs = firstPositiveInt64(recoveredRun.RunAtMs, recoveredRun.EndedAtMs, recoveredRun.StartedAtMs)
				}
				item.UpdatedAt = currentText()
				_ = a.store.put(persistCtx, "jobs", item.ID, item)
				return map[string]any{"submitted": true, "completed": false, "fallbackUsed": true, "job": item, "run": sanitizeCronRunSubmission(submission), "message": item.BusinessMessage, "updatedCount": recovered}, nil
			}
		}
		item.BusinessStatus = "syncing"
		item.BusinessMessage = waitErr.Error()
		item.LastResult = "running"
		item.UpdatedAt = currentText()
		_ = a.store.put(persistCtx, "jobs", item.ID, item)
		return map[string]any{"submitted": true, "completed": false, "job": item, "run": sanitizeCronRunSubmission(submission), "message": waitErr.Error()}, nil
	}
	item.LastRunID = run.RunID
	item.LastRun = currentText()
	item.LastResult = firstNonEmpty(run.Status, "ok")
	item.Status = firstNonEmpty(run.Status, "success")
	item.LastProcessedRunAtMs = firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs)
	if isCronFailureStatus(run.Status) || strings.TrimSpace(run.Error) != "" {
		if item.OutputTarget == "news" {
			recovered, recoveredRun, recoverErr := a.recoverLatestSuccessfulNewsRun(persistCtx, item, run.RunID)
			if recoverErr == nil && recovered > 0 {
				item.BusinessStatus = "updated"
				item.BusinessMessage = fmt.Sprintf("本次 OpenClaw 行业新闻任务失败：%s；已恢复最近一次成功运行的 %d 条新闻。", cronFailureMessage(run), recovered)
				item.BusinessUpdatedAt = currentText()
				item.BusinessResultVersion = businessResultVersion
				if recoveredRun != nil {
					item.LastRunID = recoveredRun.RunID
					item.LastProcessedRunAtMs = firstPositiveInt64(recoveredRun.RunAtMs, recoveredRun.EndedAtMs, recoveredRun.StartedAtMs)
				}
				item.UpdatedAt = currentText()
				_ = a.store.put(persistCtx, "jobs", item.ID, item)
				return map[string]any{"submitted": true, "completed": true, "fallbackUsed": true, "job": item, "run": sanitizeCronRun(run), "message": item.BusinessMessage, "updatedCount": recovered}, nil
			}
		}
		item.BusinessStatus = "failed"
		item.BusinessMessage = cronFailureMessage(run)
		item.BusinessResultVersion = businessResultVersion
		item.UpdatedAt = currentText()
		_ = a.store.put(persistCtx, "jobs", item.ID, item)
		return map[string]any{"submitted": true, "completed": true, "job": item, "run": sanitizeCronRun(run), "message": item.BusinessMessage}, nil
	}
	item.Error = ""
	var result cronBusinessResult
	var ok bool
	if item.OutputTarget == "news" {
		result, ok = parseNewsBusinessResult(run.Summary, firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs))
	} else if item.OutputTarget == "recommendations" {
		result, ok = parseRecommendationBusinessResult(run.Summary, firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs))
	} else {
		result, ok = parseCronBusinessResult(run.Summary)
		if !ok {
			result, ok = parseCronMarkdownBusinessResult(item.OutputTarget, run.Summary)
		}
	}
	count := 0
	switch item.OutputTarget {
	case "recommendations":
		if !ok || result.Type != "recommendations" {
			item.BusinessStatus = "needs_review"
			item.BusinessMessage = "推荐 Agent 已执行，但没有返回可自动入库的推荐数据块；原始结果请在运行记录中复核。"
			item.BusinessUpdatedAt = currentText()
			item.BusinessResultVersion = businessResultVersion
			item.UpdatedAt = currentText()
			_ = a.store.put(persistCtx, "jobs", item.ID, item)
			return map[string]any{"submitted": true, "completed": true, "job": item, "run": sanitizeCronRun(run), "message": item.BusinessMessage}, nil
		}
		count, err = a.persistRecommendations(persistCtx, item, result.Items)
		if err != nil {
			item.BusinessStatus = "failed"
			item.BusinessMessage = "推荐 Agent 已返回，但写入本机推荐库失败：" + err.Error()
			item.BusinessResultVersion = businessResultVersion
			item.UpdatedAt = currentText()
			_ = a.store.put(persistCtx, "jobs", item.ID, item)
			return nil, err
		}
		item.BusinessUpdatedAt = currentText()
		item.BusinessResultVersion = businessResultVersion
		if count == 0 {
			item.BusinessStatus = "needs_review"
			item.BusinessMessage = "推荐 Agent 已执行，但没有返回可展示的有效推荐；页面继续展示本机缓存推荐。"
		} else {
			if visibleCount, err := a.displayableRecommendationCount(persistCtx); err == nil {
				count = visibleCount
			}
			if count == 0 {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = "推荐结果已写入，但没有符合页面展示规则的有效推荐；请复核 Agent 输出格式。"
				break
			}
			item.BusinessStatus = "updated"
			item.BusinessMessage = fmt.Sprintf("已从 OpenClaw 推荐 Agent 同步 %d 条推荐。", count)
		}
	case "news":
		if ok && result.Type == "news" {
			count, err = a.persistNews(persistCtx, item, result.Items)
			if err != nil {
				item.BusinessStatus = "failed"
				item.BusinessMessage = "新闻 Agent 已返回，但写入本机新闻库失败：" + err.Error()
				item.BusinessResultVersion = businessResultVersion
				item.UpdatedAt = currentText()
				_ = a.store.put(persistCtx, "jobs", item.ID, item)
				return nil, err
			}
		}
		if count > 0 {
			item.BusinessStatus = "updated"
			item.BusinessMessage = fmt.Sprintf("已从 OpenClaw 行业新闻 Agent 同步 %d 条新闻。", count)
		} else {
			recovered, recoveredRun, recoverErr := a.recoverLatestSuccessfulNewsRun(persistCtx, item, run.RunID)
			if recoverErr != nil {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = "行业新闻已执行，但没有返回可展示的有效新闻；读取最近成功新闻记录也失败：" + userJobSyncError(recoverErr)
			} else if recovered > 0 {
				count = recovered
				item.BusinessStatus = "updated"
				item.BusinessMessage = fmt.Sprintf("本次行业新闻未返回有效条目，已恢复最近一次成功运行的 %d 条新闻。", recovered)
				if recoveredRun != nil {
					item.LastRunID = recoveredRun.RunID
					item.LastProcessedRunAtMs = firstPositiveInt64(recoveredRun.RunAtMs, recoveredRun.EndedAtMs, recoveredRun.StartedAtMs)
				}
			} else {
				item.BusinessStatus = "needs_review"
				item.BusinessMessage = "行业新闻已执行，但没有返回可展示的有效新闻；原始结果保留在 OpenClaw 运行记录中。"
			}
		}
		item.BusinessUpdatedAt = currentText()
		item.BusinessResultVersion = businessResultVersion
	case "weekly_report":
		if !ok || result.Type != "weekly_report" || strings.TrimSpace(result.Summary) == "" {
			item.BusinessStatus = "needs_review"
			item.BusinessMessage = "周报已完成，但没有返回可展示的摘要。"
			item.BusinessResultVersion = businessResultVersion
			item.UpdatedAt = currentText()
			_ = a.store.put(persistCtx, "jobs", item.ID, item)
			return map[string]any{"submitted": true, "completed": true, "job": item, "run": sanitizeCronRun(run), "message": item.BusinessMessage}, nil
		}
		item.BusinessStatus = "updated"
		item.BusinessUpdatedAt = currentText()
		item.BusinessMessage = "周报已生成，完整内容请通过运行记录查看。"
		item.LastResult = result.Summary
		item.BusinessResultVersion = businessResultVersion
	case "private_files":
		count, err = a.refreshPrivateFileMetadata(persistCtx)
		if err != nil {
			item.BusinessStatus = "failed"
			item.BusinessMessage = "文件元数据扫描失败：" + err.Error()
			item.BusinessResultVersion = businessResultVersion
			item.UpdatedAt = currentText()
			_ = a.store.put(persistCtx, "jobs", item.ID, item)
			return nil, err
		}
		item.BusinessStatus = "updated"
		item.BusinessUpdatedAt = currentText()
		item.BusinessMessage = fmt.Sprintf("已完成 %d 个本地文件的元数据复核；正文解析和向量索引仍按数据格式待补充。", count)
		item.BusinessResultVersion = businessResultVersion
	default:
		item.BusinessStatus = "needs_review"
		item.BusinessMessage = "任务已完成，但未配置业务结果接收目标。"
		item.BusinessResultVersion = businessResultVersion
	}
	item.UpdatedAt = currentText()
	_ = a.store.put(persistCtx, "jobs", item.ID, item)
	return map[string]any{"submitted": true, "completed": true, "job": item, "run": sanitizeCronRun(run), "message": item.BusinessMessage, "updatedCount": count}, nil
}

func businessJobAgentLabel(target string) string {
	switch strings.TrimSpace(target) {
	case "news":
		return "行业新闻 Agent "
	case "recommendations":
		return "推荐 Agent "
	case "weekly_report":
		return "周报 Agent "
	default:
		return "业务 Agent "
	}
}

func (a *businessAPI) waitForCronRun(ctx context.Context, openClawID string, startedAtMs int64) (orchestrator.CronRun, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	deadline := time.NewTimer(150 * time.Second)
	defer deadline.Stop()
	for {
		result, err := a.openClaw.service.CronRuns(ctx, openClawID, 10)
		if err != nil {
			return orchestrator.CronRun{}, err
		}
		if run := latestCronRunSince(result.Entries, startedAtMs); run != nil {
			status := strings.ToLower(strings.TrimSpace(run.Status))
			if status != "running" && status != "active" && status != "queued" && status != "pending" && status != "waiting" {
				return *run, nil
			}
		}
		select {
		case <-ctx.Done():
			return orchestrator.CronRun{}, fmt.Errorf("OpenClaw Agent 仍在执行，稍后会由定时任务状态同步结果。")
		case <-deadline.C:
			return orchestrator.CronRun{}, fmt.Errorf("OpenClaw Agent 执行时间较长，已提交刷新；稍后请再次刷新页面查看结果。")
		case <-ticker.C:
		}
	}
}

func (a *businessAPI) recoverLatestSuccessfulNewsRun(ctx context.Context, item Job, skipRunID string) (int, *orchestrator.CronRun, error) {
	if strings.TrimSpace(item.OpenClawID) == "" {
		return 0, nil, nil
	}
	recoverCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	runs, err := a.openClaw.service.CronRuns(recoverCtx, item.OpenClawID, 50)
	if err != nil {
		return 0, nil, err
	}
	for index := range runs.Entries {
		run := runs.Entries[index]
		if skipRunID != "" && run.RunID == skipRunID {
			continue
		}
		if isCronFailureStatus(run.Status) || strings.TrimSpace(run.Error) != "" || strings.TrimSpace(run.Summary) == "" {
			continue
		}
		result, ok := parseNewsBusinessResult(run.Summary, firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs))
		if !ok || result.Type != "news" {
			continue
		}
		count, err := a.persistNews(ctx, item, result.Items)
		if err != nil {
			return 0, nil, err
		}
		if count > 0 {
			return count, &runs.Entries[index], nil
		}
	}
	return 0, nil, nil
}

func latestCronRunSince(entries []orchestrator.CronRun, startedAtMs int64) *orchestrator.CronRun {
	var latest *orchestrator.CronRun
	var latestScore int64
	for index := range entries {
		run := entries[index]
		score := firstPositiveInt64(run.RunAtMs, run.EndedAtMs, run.StartedAtMs)
		if score < startedAtMs {
			continue
		}
		if latest == nil || score > latestScore {
			latest = &entries[index]
			latestScore = score
		}
	}
	return latest
}

func (a *businessAPI) persistNews(ctx context.Context, job Job, rawItems []json.RawMessage) (int, error) {
	existing, err := listRecords[NewsItem](ctx, a.store, "news")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, raw := range rawItems {
		var input newsResultItem
		if err := json.Unmarshal(raw, &input); err != nil || strings.TrimSpace(input.Title) == "" {
			continue
		}
		input.Category, input.Title, input.Summary, input.Content = normalizeNewsRecordFields(input.Category, input.Title, input.Summary, input.Content)
		input.Source = cleanBusinessGeneratedTitle(firstNonEmpty(input.Source, "OpenClaw Agent"))
		if shouldSkipNewsResult(input) {
			continue
		}
		// A list item is only valid when the Agent supplied complete, separated
		// fields. Never infer sourceUrl from summary/content: doing so turns a
		// truncated URL-only response into a fake article.
		if !validNewsResultItem(input) {
			continue
		}
		input.Summary = firstNonEmpty(input.Summary, newsBrief(input.Content, 120))
		input.Time = firstNonEmpty(input.Time, currentText())
		input.Relevance = firstNonEmpty(input.Relevance, "待复核")
		record := NewsItem{
			ID:       generatedBusinessID("NEWS", input.Title, input.Source),
			Category: input.Category, Title: input.Title, Summary: input.Summary, Content: input.Content,
			Source: input.Source, SourceURL: input.SourceURL, Time: input.Time,
			Relevance: input.Relevance, UpdatedAt: currentText(), GeneratedBy: job.AgentID, DataStatus: "generated",
		}
		for _, old := range existing {
			if old.Title == record.Title && old.Source == record.Source {
				record.ID = old.ID
				break
			}
		}
		if err := a.putOrCreate(ctx, "news", record.ID, record); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func newsBrief(value string, limit int) string {
	text := strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if text == "" {
		return "该资讯由 OpenClaw Agent 整理，详情请结合来源复核。"
	}
	runes := []rune(text)
	if limit > 0 && len(runes) > limit {
		return string(runes[:limit]) + "…"
	}
	return text
}

func normalizeNewsText(value string) string {
	text := cleanBusinessGeneratedBody(stripSTA100ResultBlocks(value))
	if text == "" {
		return ""
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	// Agent responses sometimes repeat field labels inside the body. Remove
	// labels before deduplicating so "摘要:" and "完整内容:" never become
	// visible article text.
	text = regexp.MustCompile(`(?i)(摘要|完整内容|详细内容|正文|内容|详情)\s*[:：]\s*`).ReplaceAllString(text, "\n")
	seen := map[string]bool{}
	lines := make([]string, 0)
	for _, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			if len(lines) > 0 && lines[len(lines)-1] != "" {
				lines = append(lines, "")
			}
			continue
		}
		key := strings.Join(strings.Fields(line), " ")
		if seen[key] {
			continue
		}
		seen[key] = true
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func firstURL(value string) string {
	match := markdownBusinessURLPattern.FindString(value)
	return strings.TrimRight(match, "，。.,)")
}

func (a *businessAPI) putOrCreate(ctx context.Context, kind, id string, value any) error {
	if err := a.store.put(ctx, kind, id, value); err == nil {
		return nil
	} else if !errors.Is(err, errRecordNotFound) {
		return err
	}
	return a.store.create(ctx, kind, id, value)
}

func generatedBusinessID(prefix, title, source string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(title) + "\x00" + strings.TrimSpace(source)))
	return prefix + "-AI-" + hex.EncodeToString(sum[:])[:16]
}

func (a *businessAPI) refreshPrivateFileMetadata(ctx context.Context) (int, error) {
	files, err := listRecords[PrivateFile](ctx, a.store, "private_files")
	if err != nil {
		return 0, err
	}
	available := 0
	for _, file := range files {
		if _, err := privateFileStoragePath(file); err == nil {
			available++
		}
	}
	return available, nil
}

func (a *businessAPI) syncRecommendationJobs(ctx context.Context, preferences UserPreferences) (string, error) {
	items, err := listRecords[Job](ctx, a.store, "jobs")
	if err != nil {
		return "", err
	}
	frequency := strings.TrimSpace(preferences.NewsFrequency)
	hours := strings.TrimSuffix(frequency, "小时")
	scheduleValue := hours + "h"
	for index := range items {
		if !items[index].BuiltIn || (items[index].ID != "JOB-RECOMMEND" && items[index].ID != "JOB-NEWS") {
			continue
		}
		item := applyBuiltInJobDefaults(items[index])
		item.Enabled = preferences.RecommendationEnabled
		item.ScheduleKind = "every"
		item.ScheduleValue = scheduleValue
		item.Schedule = "每 " + frequency
		if item.ID == "JOB-RECOMMEND" {
			limit := preferences.RecommendationShowLimit
			if limit < 1 || limit > 20 {
				limit = 5
			}
			item.Prompt = fmt.Sprintf("根据关注国家：%s；关注主题：%s；读取本地业务数据和可用公开来源，最多生成 %d 条有效推荐。", preferences.NewsCountries, preferences.NewsTopics, limit)
		} else {
			fetchLimit := preferences.NewsFetchLimit
			if fetchLimit < 1 || fetchLimit > 100 {
				fetchLimit = 5
			}
			item.Prompt = fmt.Sprintf("按以下关注国家：%s；关注主题：%s；指定来源：%s，整理骑行行业新闻，每次最多输出 %d 条有效新闻。", preferences.NewsCountries, preferences.NewsTopics, preferences.NewsSources, fetchLimit)
		}
		updated, syncErr := a.syncOneJob(ctx, item)
		if syncErr != nil {
			return "", fmt.Errorf("%s：%w", item.Name, syncErr)
		}
		if err := a.store.put(ctx, "jobs", item.ID, updated); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("推荐和行业新闻任务已按每 %s 更新%s", frequency, map[bool]string{true: "", false: "，当前已停用"}[preferences.RecommendationEnabled]), nil
}

func cronInputFromJob(item Job) orchestrator.CronJobInput {
	kind, value, timezone := jobScheduleSpec(item)
	return orchestrator.CronJobInput{
		ID:             item.OpenClawID,
		Name:           item.Name,
		Description:    item.Description,
		DeclarationKey: "sta100:" + item.ID,
		Enabled:        item.Enabled,
		AgentID:        item.AgentID,
		SessionTarget:  "isolated",
		ScheduleKind:   kind,
		ScheduleValue:  value,
		Timezone:       firstNonEmpty(item.Timezone, timezone),
		Message:        cronPrompt(item),
	}
}

func cronPrompt(item Job) string {
	prompt := strings.TrimSpace(item.Prompt)
	switch item.OutputTarget {
	case "recommendations":
		return prompt + "\n\n" + structuredBusinessPrompt(structuredBusinessPromptSpec{
			ModuleName: "推荐内容",
			Scope:      "必须基于推荐设置中的关注国家、关注主题、本机业务数据和可用公开来源生成真实推荐内容，不要返回演示数据。",
			Fields:     []string{"标题", "为什么推荐", "详情", "来源", "来源链接", "类型", "时间", "摘要"},
			ResultType: "recommendations",
			CountHint:  "请严格遵守上方给出的数量上限，只输出有效推荐条目，不要输出总览壳子。",
			Notes: []string{
				"标题要短；为什么推荐写 1-2 句说明命中原因；详情写完整展开内容，至少 120 个中文字符，说明影响什么业务、建议下一步怎么做。",
				"来源写来源名称；来源链接写原文地址，没有则留空；类型写推荐类别；时间写更新时间；摘要只写 30-60 字概览。",
				"没有可靠结果时 items 输出空数组，不要编造事实。不要把证据缺失说明、提示语、工作区文件清单或本格式示例写入 items。",
			},
		})
	case "news":
		return prompt + "\n\n" + newsStructuredResultPrompt()
	case "weekly_report":
		return prompt + "\n\n请在最终结果末尾输出机器可读结果块，格式必须为：[STA100_RESULT]{\"type\":\"weekly_report\",\"summary\":\"周报摘要\"}[/STA100_RESULT]。"
	default:
		return prompt
	}
}

func newsStructuredResultPrompt() string {
	return structuredBusinessPrompt(structuredBusinessPromptSpec{
		ModuleName: "行业新闻",
		Scope:      "只基于新闻设置中的关注国家、关注主题、指定来源、OpenClaw 可用检索结果和 Agent 明确拿到的来源内容生成；不要使用演示数据，不要编造来源、发布时间、事实或链接。",
		Fields:     []string{"category", "title", "summary", "content", "source", "sourceUrl", "time", "relevance"},
		ResultType: "news",
		CountHint:  "你必须遵守上文任务 Prompt 中的“每次最多输出 N 条有效新闻”数量上限，只输出前 N 条完整新闻，不要补足演示数据。",
		Notes: []string{
			"每条 item 必须是一条独立新闻；summary 写 30-80 个中文字符；content 写完整展开内容，至少 120 个中文字符，说明事实背景、业务影响和建议复核点。",
			"sourceUrl 只能写原文地址，summary 和 content 里不要出现 URL；分类 category 只能优先使用：欧洲市场、法规、智能骑行、渠道、产品；无法归类时写行业资讯。",
			"relevance 使用高、中、低或待复核，并简短说明原因。",
			"如果你会写入工作区文件，请把同一份 JSON 写入 workspaces/sta100-news-curator/output/sta100-news-YYYY-MM-DD.json；文件 items 中也必须包含 content 字段。",
		},
		Workspace: newsWorkspaceOutputDir + "/sta100-news-YYYY-MM-DD.json",
		Schema:    structuredBusinessResultSchema,
	})
}

type structuredBusinessPromptSpec struct {
	ModuleName string
	Scope      string
	Fields     []string
	ResultType string
	CountHint  string
	Notes      []string
	Workspace  string
	Schema     string
}

func structuredBusinessPrompt(spec structuredBusinessPromptSpec) string {
	lines := []string{fmt.Sprintf("请按 STA-100 结构化结果协议输出%s。Go 侧只解析固定 JSON，不再从自由 Markdown 中猜业务字段。", spec.ModuleName)}
	if strings.TrimSpace(spec.Scope) != "" {
		lines = append(lines, "输入边界："+spec.Scope)
	}
	if len(spec.Fields) > 0 {
		lines = append(lines, "输出字段："+strings.Join(spec.Fields, "、"))
	}
	if strings.TrimSpace(spec.CountHint) != "" {
		lines = append(lines, spec.CountHint)
	}
	if len(spec.Notes) > 0 {
		lines = append(lines, spec.Notes...)
	}
	if strings.TrimSpace(spec.Workspace) != "" {
		lines = append(lines, "如果你会写入工作区文件，请把同一份 JSON 写入 "+spec.Workspace+"。")
	}
	if strings.TrimSpace(spec.ResultType) != "" {
		lines = append(lines, fmt.Sprintf("最终结果末尾必须只输出一个机器可读结果块，格式必须为：[STA100_RESULT]{\"schema\":\"%s\",\"type\":\"%s\",\"items\":[]}[/STA100_RESULT]。", firstNonEmpty(strings.TrimSpace(spec.Schema), structuredBusinessResultSchema), spec.ResultType))
	}
	lines = append(lines, "没有可靠结果时 items 输出空数组，不要编造事实。不要把证据缺失说明、提示语、工作区文件清单或本格式示例写入 items。")
	return strings.Join(lines, "\n\n")
}

func jobScheduleSpec(item Job) (string, string, string) {
	kind := strings.ToLower(strings.TrimSpace(item.ScheduleKind))
	value := strings.TrimSpace(firstNonEmpty(item.ScheduleValue, item.Schedule))
	if kind != "" && value != "" {
		return kind, value, strings.TrimSpace(item.Timezone)
	}
	if match := jobEveryMinutesPattern.FindStringSubmatch(value); len(match) == 2 {
		minutes, _ := strconv.Atoi(match[1])
		if minutes > 0 {
			return "every", strconv.Itoa(minutes) + "m", ""
		}
	}
	if match := jobEveryHoursPattern.FindStringSubmatch(value); len(match) == 2 {
		hours, _ := strconv.Atoi(match[1])
		if hours > 0 {
			return "every", strconv.Itoa(hours) + "h", ""
		}
	}
	switch value {
	case "每周":
		return "cron", "0 9 * * 1", ""
	case "每天", "每天 08:00":
		return "cron", "0 8 * * *", ""
	case "每周一 09:00":
		return "cron", "0 9 * * 1", ""
	}
	fields := strings.Fields(value)
	if len(fields) == 5 || len(fields) == 6 {
		return "cron", value, strings.TrimSpace(item.Timezone)
	}
	return "cron", "0 8 * * *", strings.TrimSpace(item.Timezone)
}

func cronJobMatches(item Job, runtime orchestrator.CronJob) bool {
	kind, value, timezone := jobScheduleSpec(item)
	scheduleMatches := runtime.Schedule.Kind == kind
	switch kind {
	case "every":
		scheduleMatches = scheduleMatches && runtime.Schedule.EveryMs == durationToMilliseconds(value)
	case "cron":
		scheduleMatches = scheduleMatches && runtime.Schedule.Expr == value && runtime.Schedule.TZ == timezone
	case "at":
		scheduleMatches = scheduleMatches && runtime.Schedule.AtMs > 0
	}
	return runtime.Name == item.Name &&
		runtime.Description == item.Description &&
		runtime.Enabled == item.Enabled &&
		runtime.AgentID == item.AgentID &&
		runtime.Payload.Message == cronPrompt(item) &&
		scheduleMatches
}

func durationToMilliseconds(value string) int64 {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return 0
	}
	var multiplier int64 = 1
	switch value[len(value)-1] {
	case 'm':
		multiplier = 60 * 1000
	case 'h':
		multiplier = 60 * 60 * 1000
	case 'd':
		multiplier = 24 * 60 * 60 * 1000
	case 's':
		multiplier = 1000
	default:
		return 0
	}
	number, err := strconv.ParseInt(strings.TrimSpace(value[:len(value)-1]), 10, 64)
	if err != nil || number <= 0 {
		return 0
	}
	return number * multiplier
}

func applyCronRuntime(item Job, runtime orchestrator.CronJob) Job {
	item.OpenClawID = runtime.ID
	item.Enabled = runtime.Enabled
	item.SyncStatus = "synced"
	item.SyncMessage = "已写入 OpenClaw Cron，并完成任务状态复核"
	item.Status = runtimeJobStatus(runtime)
	item.LastRun = formatJobTime(runtime.State.LastRunAtMs)
	item.NextRun = formatJobTime(runtime.State.NextRunAtMs)
	if isCronFailureStatus(cronRuntimeStatus(runtime)) {
		item.Error = userJobExecutionError(firstNonEmpty(runtime.State.LastError, runtime.State.LastRunError))
	} else {
		item.Error = ""
	}
	item.LastResult = firstNonEmpty(runtime.State.LastRunSummary, item.Error, runtime.Status)
	return item
}

func mergeCronRuntime(items []Job, runtime []orchestrator.CronJob) []Job {
	byID := make(map[string]orchestrator.CronJob, len(runtime))
	byDeclaration := make(map[string]orchestrator.CronJob, len(runtime))
	for _, job := range runtime {
		byID[job.ID] = job
		if job.DeclarationKey != "" {
			byDeclaration[job.DeclarationKey] = job
		}
	}
	for index := range items {
		job, ok := byID[items[index].OpenClawID]
		if !ok {
			job, ok = byDeclaration["sta100:"+items[index].ID]
		}
		if ok {
			items[index] = applyCronRuntime(items[index], job)
		} else if items[index].OpenClawID != "" {
			items[index].SyncStatus = "missing"
			items[index].SyncMessage = "本地任务记录存在，但 OpenClaw Cron 中未找到对应任务"
			items[index].Status = "unsynced"
		} else if items[index].Enabled {
			items[index].SyncStatus = firstNonEmpty(items[index].SyncStatus, "pending")
			items[index].Status = "unsynced"
		}
	}
	return items
}

func runtimeJobStatus(runtime orchestrator.CronJob) string {
	if !runtime.Enabled {
		return "disabled"
	}
	switch cronRuntimeStatus(runtime) {
	case "running":
		return "running"
	case "ok", "success":
		return "success"
	case "error", "failed":
		return "failed"
	case "skipped":
		return "skipped"
	default:
		return "waiting"
	}
}

func cronRuntimeStatus(runtime orchestrator.CronJob) string {
	if runtime.State.RunningAtMs > 0 {
		return "running"
	}
	return strings.ToLower(strings.TrimSpace(firstNonEmpty(
		runtime.State.LastRunStatus,
		runtime.State.LastStatus,
		runtime.Status,
	)))
}

func isCronFailureStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "error", "failed", "failure", "timeout", "aborted":
		return true
	default:
		return false
	}
}

func humanCronStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "error", "failed", "failure":
		return "执行失败"
	case "timeout":
		return "执行超时"
	case "aborted":
		return "执行被中止"
	case "success", "ok":
		return "执行成功"
	case "running":
		return "正在执行"
	default:
		return "状态异常"
	}
}

func cronFailureMessage(run orchestrator.CronRun) string {
	providerModel := ""
	if strings.TrimSpace(run.Provider) != "" || strings.TrimSpace(run.Model) != "" {
		providerModel = "（" + strings.TrimSpace(strings.TrimSpace(run.Provider)+" / "+strings.TrimSpace(run.Model)) + "）"
	}
	reason := firstNonEmpty(run.ErrorReason, run.Cause)
	raw := firstNonEmpty(run.Error, reason)
	return userJobExecutionErrorWithContext(raw, providerModel, humanCronStatus(firstNonEmpty(run.Status, "failed")))
}

func userJobExecutionError(raw string) string {
	return userJobExecutionErrorWithContext(raw, "", "执行失败")
}

func userJobExecutionErrorWithContext(raw, providerModel, status string) string {
	text := strings.TrimSpace(raw)
	lower := strings.ToLower(text)
	if strings.Contains(lower, "network connection error") || strings.Contains(lower, "econnreset") || strings.Contains(lower, "socket hang up") {
		return fmt.Sprintf("OpenClaw 模型调用%s失败：网络连接异常。请检查 Gateway 网络、模型 API 地址和 API Key；本次未生成业务结果。", providerModel)
	}
	if strings.Contains(lower, "timeout") || strings.Contains(lower, "timed out") {
		return fmt.Sprintf("OpenClaw 模型调用%s超时。请检查模型服务响应和网络连接；本次未生成业务结果。", providerModel)
	}
	if strings.Contains(lower, "unauthorized") || strings.Contains(lower, "authentication") || strings.Contains(lower, "api key") {
		return fmt.Sprintf("OpenClaw 模型鉴权%s失败。请在模型设置中重新测试 API Key；本次未生成业务结果。", providerModel)
	}
	if text == "" {
		return fmt.Sprintf("OpenClaw 最近一次%s，业务数据未更新。", status)
	}
	return fmt.Sprintf("OpenClaw 最近一次%s：%s", status, text)
}

func formatJobTime(milliseconds int64) string {
	if milliseconds <= 0 {
		return ""
	}
	return time.UnixMilli(milliseconds).In(time.Local).Format("2006-01-02 15:04:05")
}

func userJobSyncError(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "OpenClaw Cron 查询或同步超时，请检查 Gateway 是否正在运行"
	}
	if errors.Is(err, orchestrator.ErrUnavailable) {
		return "OpenClaw CLI 不可用，任务已保留为未同步状态"
	}
	return "OpenClaw Cron 同步失败：" + err.Error()
}
