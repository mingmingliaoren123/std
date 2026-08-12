package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"openclaw-orchestrator/orchestrator"
)

func defaultPreferences() UserPreferences {
	return UserPreferences{
		RecommendationEnabled: true,
		NewsShowLimit:         20,
		NewsFrequency:         "1小时",
		NewsCountries:         "德国、法国、波兰、瑞典",
		NewsTopics:            "E-bike、智能骑行、经销商、欧盟法规",
		NewsSources:           "EUR-Lex\nBike Europe\nCycling Industry News\nEurobike",
		AgentAllowlists:       map[string][]string{},
	}
}

func (a *businessAPI) privateFilesRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.privateFilesCollection(w, r)
		return
	}
	if parts[0] == "upload" {
		a.privateFileUpload(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "reindex" {
		a.privateFileReindex(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "download" {
		a.privateFileDownload(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "content" {
		a.privateFileContent(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "summary" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		var item PrivateFile
		if err := a.store.get(r.Context(), "private_files", parts[0], &item); err != nil {
			writeBusinessError(w, err)
			return
		}
		writeTODO(w, "TODO_PRIVATE_FILE_PARSER", "文件摘要需要在客户原始数据格式和解析规则确认后生成", []string{"原始文件格式", "解析器", "摘要字段", "索引策略"})
		return
	}
	a.privateFileItem(w, r, parts[0])
}

func (a *businessAPI) privateFilesCollection(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := listRecords[PrivateFile](r.Context(), a.store, "private_files")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	category := r.URL.Query().Get("category")
	filtered := make([]PrivateFile, 0, len(items))
	for _, item := range items {
		if query != "" && !containsFold(item.ID+" "+item.Name+" "+item.Category+" "+strings.Join(item.Tags, " "), query) {
			continue
		}
		if category != "" && category != "all" && item.Category != category {
			continue
		}
		filtered = append(filtered, item)
	}
	listResponse(w, filtered)
}

func (a *businessAPI) privateFileUpload(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_UPLOAD", "文件上传格式无效或超过 50 MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "FILE_REQUIRED", "必须上传 file 字段")
		return
	}
	defer file.Close()
	extension := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".pdf": true, ".docx": true, ".xlsx": true, ".csv": true, ".txt": true, ".md": true, ".jpg": true, ".jpeg": true, ".png": true}
	if !allowed[extension] {
		writeAPIError(w, http.StatusUnsupportedMediaType, "FILE_TYPE_UNSUPPORTED", "文件格式暂不支持")
		return
	}
	root, err := privateDataRoot()
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	temporary, err := os.CreateTemp(root, ".upload-*")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(temporary, hash), file)
	closeErr := temporary.Close()
	if err != nil || closeErr != nil {
		writeAPIError(w, http.StatusInternalServerError, "UPLOAD_WRITE_FAILED", "文件写入本地私有区失败")
		return
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	existing, _ := listRecords[PrivateFile](r.Context(), a.store, "private_files")
	for _, item := range existing {
		if item.SHA256 != "" && item.SHA256 == digest {
			writeAPIError(w, http.StatusConflict, "FILE_DUPLICATE", "相同内容文件已经存在")
			return
		}
	}
	id, err := a.store.nextSequence(r.Context(), "private_files", "FILE", 4)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	target := filepath.Join(root, id+extension)
	if err := os.Rename(temporaryName, target); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "UPLOAD_COMMIT_FAILED", "文件保存失败")
		return
	}
	_ = os.Chmod(target, 0600)
	category := strings.TrimSpace(r.FormValue("category"))
	if category == "" || category == "自动识别" {
		category = "待分类"
	}
	tags := splitTags(r.FormValue("tags"))
	item := PrivateFile{ID: id, Name: filepath.Base(header.Filename), Category: category, Tags: tags, Size: humanBytes(written), Bytes: written, SHA256: digest, Mime: header.Header.Get("Content-Type"), Source: "客户上传", Status: "PendingParse", Path: target, Updated: currentText()}
	if err := a.store.create(r.Context(), "private_files", id, item); err != nil {
		_ = os.Remove(target)
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "upload", "private_file", id, requestOperator(r), map[string]any{"sha256": digest, "bytes": written})
	writeJSON(w, http.StatusCreated, map[string]any{"item": item, "todo": "等待客户提供原始数据格式后接入解析、分类和索引"})
}

func (a *businessAPI) privateFileItem(w http.ResponseWriter, r *http.Request, id string) {
	var item PrivateFile
	if err := a.store.get(r.Context(), "private_files", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, item)
	case http.MethodPatch:
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
		var request PrivateFile
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID, request.Path, request.SHA256, request.Bytes, request.Mime, request.Size = id, item.Path, item.SHA256, item.Bytes, item.Mime, item.Size
		if request.Tags == nil {
			request.Tags = item.Tags
		}
		if strings.TrimSpace(request.Source) == "" {
			request.Source = item.Source
		}
		if strings.TrimSpace(request.Status) == "" {
			request.Status = item.Status
		}
		request.Updated = currentText()
		if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Category) == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_FILE_METADATA", "文件名和主分类不能为空")
			return
		}
		if err := a.store.put(r.Context(), "private_files", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "private_file", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if path, err := privateFileStoragePath(item); err == nil {
			_ = os.Remove(path)
		}
		if err := a.store.softDelete(r.Context(), "private_files", id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "delete", "private_file", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) privateFileReindex(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var item PrivateFile
	if err := a.store.get(r.Context(), "private_files", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	item.Status, item.Updated = "WaitingDataFormat", currentText()
	if err := a.store.put(r.Context(), "private_files", id, item); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "reindex_requested", "private_file", id, requestOperator(r), nil)
	writeJSON(w, http.StatusAccepted, map[string]any{"item": item, "todo": "原始数据格式和索引策略待提供"})
}

func (a *businessAPI) privateFileDownload(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item PrivateFile
	if err := a.store.get(r.Context(), "private_files", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	path, err := privateFileStoragePath(item)
	if err != nil {
		writeAPIError(w, http.StatusNotFound, "FILE_CONTENT_UNAVAILABLE", "原始文件不存在或已移动")
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(item.Name, `"`, "")+`"`)
	http.ServeFile(w, r, path)
}

func (a *businessAPI) privateFileContent(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item PrivateFile
	if err := a.store.get(r.Context(), "private_files", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	path, err := privateFileStoragePath(item)
	if err != nil {
		writeAPIError(w, http.StatusNotFound, "FILE_CONTENT_UNAVAILABLE", "原始文件不存在或已移动")
		return
	}
	file, err := os.Open(path)
	if err != nil {
		writeAPIError(w, http.StatusNotFound, "FILE_CONTENT_UNAVAILABLE", "原始文件不存在或已移动")
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	w.Header().Set("Content-Disposition", `inline; filename="`+strings.ReplaceAll(item.Name, `"`, "")+`"`)
	if item.Mime != "" {
		w.Header().Set("Content-Type", item.Mime)
	}
	http.ServeContent(w, r, item.Name, info.ModTime(), file)
}

func (a *businessAPI) tasksRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	writeTODO(w, "TODO_TASK_SOURCE", "待办写入接口已保留，需确认待办来源、字段、提醒和完成规则", []string{"待办字段", "提醒规则", "完成状态", "Agent 同步范围"})
}

func privateDataRoot() (string, error) {
	root := strings.TrimSpace(os.Getenv("STA100_PRIVATE_DATA_DIR"))
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(home, ".local", "share", "sta100", "private-files")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	return root, nil
}

func privateFileStoragePath(item PrivateFile) (string, error) {
	if item.Path != "" {
		if _, err := os.Stat(item.Path); err == nil {
			return item.Path, nil
		}
	}
	root, err := privateDataRoot()
	if err != nil {
		return "", err
	}
	matches, err := filepath.Glob(filepath.Join(root, item.ID+".*"))
	if err != nil || len(matches) != 1 {
		return "", os.ErrNotExist
	}
	info, err := os.Stat(matches[0])
	if err != nil || !info.Mode().IsRegular() {
		return "", os.ErrNotExist
	}
	return matches[0], nil
}

func splitTags(value string) []string {
	fields := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == '，' || r == '/' || r == '、' || r == '\n' })
	result, seen := make([]string, 0, len(fields)), map[string]bool{}
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" && !seen[field] {
			seen[field] = true
			result = append(result, field)
		}
	}
	return result
}

func humanBytes(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d B", value)
	}
	if value < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(value)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(value)/(1024*1024))
}

func (a *businessAPI) newsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) > 0 && parts[0] == "refresh" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_NEWS_SOURCES", "行业新闻刷新接口已保留，需客户确认来源白名单、授权和抓取规则", []string{"来源白名单", "授权方式", "抓取频率", "失败重试", "留存周期"})
		return
	}
	if len(parts) > 0 && parts[0] != "" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		var item NewsItem
		if err := a.store.get(r.Context(), "news", parts[0], &item); err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := listRecords[NewsItem](r.Context(), a.store, "news")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	category := r.URL.Query().Get("category")
	filtered := make([]NewsItem, 0, len(items))
	for _, item := range items {
		if category == "" || category == "全部" || item.Category == category {
			filtered = append(filtered, item)
		}
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit > 0 && limit < len(filtered) {
		filtered = filtered[:limit]
	}
	listResponse(w, filtered)
}

func (a *businessAPI) overviewRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	path := strings.Join(parts, "/")
	switch path {
	case "":
		a.overviewFull(w, r)
	case "summary":
		a.overviewSummary(w, r)
	case "recommendations":
		a.overviewRecommendations(w, r)
	case "subscription":
		a.overviewSubscription(w, r)
	case "oem/match", "oem-matches", "oem-matches/export":
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_OEM_DATA", "OEM 匹配接口和查询参数已保留，需客户提供正式工厂数据、分类和评分规则", []string{"工厂原始数据", "骑行类目", "评分权重", "TOP 排序口径", "报告格式"})
	case "customers/search":
		a.overviewCustomerSearch(w, r)
	case "customer-discovery":
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_DISCOVERY_DATA", "本地客户发现链路已保留，需客户提供本地数据格式和允许的公开来源", []string{"本地客户数据", "公开来源白名单", "字段映射", "结果保存规则"})
	default:
		if len(parts) == 2 && parts[0] == "oem-matches" {
			if !allowMethod(w, r, http.MethodGet) {
				return
			}
			writeTODO(w, "TODO_OEM_DATA", "OEM 工厂详情接口已保留，正式工厂数据待提供", []string{"工厂原始数据", "来源链接", "认证字段", "抓取时间"})
			return
		}
		if len(parts) == 2 && parts[0] == "customers" {
			if !allowMethod(w, r, http.MethodGet) {
				return
			}
			var item Customer
			if err := a.store.get(r.Context(), "accounts", parts[1], &item); err != nil {
				writeBusinessError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "概览接口不存在")
	}
}

func (a *businessAPI) overviewFull(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	news, _ := listRecords[NewsItem](r.Context(), a.store, "news")
	recommendations, _ := listRecords[Recommendation](r.Context(), a.store, "recommendations")
	preferences := defaultPreferences()
	_ = a.store.getSetting(r.Context(), "preferences", &preferences)
	visibleNews := news
	if len(visibleNews) > 3 {
		visibleNews = visibleNews[:3]
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"summary":     a.overviewSummaryData(r.Context()),
		"preferences": preferences, "recommendations": recommendations, "news": visibleNews,
		"dataStatus": "partial_real_data", "todo": []string{"待办来源", "日程来源", "Agent 会话统计", "OEM 原始数据"},
	})
}

func (a *businessAPI) preferencesRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) != 1 || parts[0] != "recommendation" {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "偏好接口不存在")
		return
	}
	if r.Method == http.MethodGet {
		a.overviewSubscription(w, r)
		return
	}
	if !allowMutation(w, r, http.MethodPatch) {
		return
	}
	var request struct {
		Enabled   *bool  `json:"enabled"`
		Frequency string `json:"frequency"`
		Countries string `json:"countries"`
		Topics    string `json:"topics"`
		Sources   string `json:"sources"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	preferences := defaultPreferences()
	_ = a.store.getSetting(r.Context(), "preferences", &preferences)
	if request.Enabled != nil {
		preferences.RecommendationEnabled = *request.Enabled
	}
	if request.Frequency != "" {
		if !validNewsFrequency(request.Frequency) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_NEWS_FREQUENCY", "获取频率不受支持")
			return
		}
		preferences.NewsFrequency = request.Frequency
	}
	if request.Countries != "" {
		preferences.NewsCountries = request.Countries
	}
	if request.Topics != "" {
		preferences.NewsTopics = request.Topics
	}
	if request.Sources != "" {
		preferences.NewsSources = request.Sources
	}
	if err := a.store.putSetting(r.Context(), "preferences", preferences); err != nil {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, preferences)
}

func (a *businessAPI) overviewSummary(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, a.overviewSummaryData(r.Context()))
}

func (a *businessAPI) overviewSummaryData(ctx context.Context) map[string]any {
	orders, _ := listRecords[Order](ctx, a.store, "orders")
	documents, _ := listRecords[Document](ctx, a.store, "documents")
	news, _ := listRecords[NewsItem](ctx, a.store, "news")
	messages, _ := listRecords[AgentMessageRecord](ctx, a.store, "agent_messages")
	activeOrders, reviewDocuments := 0, 0
	for _, item := range orders {
		if item.Status != "Completed" && item.Status != "Cancelled" {
			activeOrders++
		}
	}
	for _, item := range documents {
		if item.Status == "Draft" || item.Status == "Review" {
			reviewDocuments++
		}
	}
	return map[string]any{"tasks": 0, "meetings": 0, "documents": reviewDocuments, "orders": activeOrders, "chats": len(messages) / 2, "news": len(news), "dataStatus": "business_and_agent_aggregates", "todo": []string{"待办来源", "日程来源"}}
}

func (a *businessAPI) overviewRecommendations(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := listRecords[Recommendation](r.Context(), a.store, "recommendations")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	listResponse(w, items)
}

func (a *businessAPI) overviewSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		preferences := defaultPreferences()
		_ = a.store.getSetting(r.Context(), "preferences", &preferences)
		writeJSON(w, http.StatusOK, preferences)
		return
	}
	if !allowMutation(w, r, http.MethodPatch) {
		return
	}
	var request struct {
		Enabled bool `json:"enabled"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	preferences := defaultPreferences()
	_ = a.store.getSetting(r.Context(), "preferences", &preferences)
	preferences.RecommendationEnabled = request.Enabled
	if err := a.store.putSetting(r.Context(), "preferences", preferences); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "update", "setting", "recommendation_enabled", requestOperator(r), request)
	writeJSON(w, http.StatusOK, preferences)
}

func (a *businessAPI) overviewCustomerSearch(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Query      string `json:"query"`
		Mode       string `json:"mode"`
		HasContact bool   `json:"hasContact"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.Mode == "" {
		request.Mode = "local"
	}
	if request.Mode != "local" && request.Mode != "rag" && request.Mode != "mixed" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_SEARCH_MODE", "搜索模式必须是 local、rag 或 mixed")
		return
	}
	if request.Mode != "local" {
		writeTODO(w, "TODO_RAG_DATA", "联网/RAG 客户搜索需要客户提供数据和来源规则；本地搜索已可用", []string{"RAG 数据格式", "联网来源", "去重规则", "字段来源展示"})
		return
	}
	customers, err := a.listCustomers(r.Context())
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	query := strings.ToLower(strings.TrimSpace(request.Query))
	results := make([]map[string]any, 0)
	for _, item := range customers {
		if item.Archived {
			continue
		}
		if query != "" && !containsFold(item.ID+" "+item.Name+" "+item.Country+" "+item.Type+" "+item.Contact+" "+item.Phone+" "+item.Email, query) {
			continue
		}
		contact := firstNonEmpty(item.Email, item.Phone, item.Website)
		if request.HasContact && contact == "" {
			continue
		}
		results = append(results, map[string]any{"id": item.ID, "name": item.Name, "country": item.Country, "type": item.Type, "contact": contact, "business": item.Description, "source": "本地知识库", "score": 100})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": results, "total": len(results), "mode": request.Mode})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (a *businessAPI) agentsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		a.openClaw.agentsHandler(w, r)
		return
	}
	if parts[0] == "weekly-report" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		a.agentWeeklyReport(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "chat" {
		a.agentChatByID(w, r, parts[0])
		return
	}
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	agents, err := a.openClaw.service.ListAgents(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	for _, agent := range businessVisibleAgents(agents) {
		if agent.ID == parts[0] {
			writeJSON(w, http.StatusOK, agent)
			return
		}
	}
	writeAPIError(w, http.StatusNotFound, "AGENT_UNAVAILABLE", "所选 Agent 未在 OpenClaw 中注册")
}

func (a *businessAPI) agentChatByID(w http.ResponseWriter, r *http.Request, agentID string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Message    string   `json:"message"`
		SessionKey string   `json:"sessionKey"`
		Sources    []string `json:"sources"`
		Allowlist  []string `json:"allowlist"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 105*time.Second)
	defer cancel()
	result, err := a.openClaw.service.SendAgentMessage(ctx, orchestrator.AgentMessageInput{AgentID: agentID, Message: request.Message, SessionKey: request.SessionKey, Sources: request.Sources, Allowlist: request.Allowlist})
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	createdAt := time.Now().Format(time.RFC3339)
	userID, sequenceErr := a.store.nextSequence(r.Context(), "agent_messages", "MSG", 8)
	if sequenceErr != nil {
		writeBusinessError(w, sequenceErr)
		return
	}
	userMessage := AgentMessageRecord{ID: userID, AgentID: agentID, SessionKey: request.SessionKey, Role: "user", Text: request.Message, Sources: request.Sources, CreatedAt: createdAt}
	if err := a.store.create(r.Context(), "agent_messages", userID, userMessage); err != nil {
		writeBusinessError(w, err)
		return
	}
	assistantID, sequenceErr := a.store.nextSequence(r.Context(), "agent_messages", "MSG", 8)
	if sequenceErr != nil {
		writeBusinessError(w, sequenceErr)
		return
	}
	assistantMessage := AgentMessageRecord{ID: assistantID, AgentID: agentID, SessionKey: request.SessionKey, Role: "assistant", Text: result.Text, Sources: request.Sources, RunID: result.RunID, CreatedAt: createdAt}
	if err := a.store.create(r.Context(), "agent_messages", assistantID, assistantMessage); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "chat", "agent", agentID, requestOperator(r), map[string]any{"sessionKey": request.SessionKey, "sources": request.Sources})
	writeJSON(w, http.StatusOK, result)
}

func (a *businessAPI) agentWeeklyReport(w http.ResponseWriter, r *http.Request) {
	var auditCount int
	_ = a.store.db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM audit_logs WHERE created_at >= datetime('now','-7 day')`).Scan(&auditCount)
	messages, _ := listRecords[AgentMessageRecord](r.Context(), a.store, "agent_messages")
	usage := map[string]int{}
	for _, message := range messages {
		if message.Role == "user" {
			usage[message.AgentID]++
		}
	}
	agentIDs := make([]string, 0, len(usage))
	for agentID := range usage {
		agentIDs = append(agentIDs, agentID)
	}
	sort.Strings(agentIDs)
	var usageLines strings.Builder
	if len(agentIDs) == 0 {
		usageLines.WriteString("- 本周尚无已保存的 Agent 对话。\n")
	}
	for _, agentID := range agentIDs {
		fmt.Fprintf(&usageLines, "- %s：%d 次对话\n", agentID, usage[agentID])
	}
	report := fmt.Sprintf("# STA-100 智能体周报\n\n生成时间：%s\n\n## 使用概览\n\n%s\n## 业务操作\n\n- 最近 7 天已记录操作：%d 次。\n\n## 完成事项\n\n- 客户、报价、订单、单据、产品、供应商和设置变更已写入本机审计日志。\n\n## 待跟进\n\n- 待办和日程来源尚待客户确认。\n- 客户原始知识数据格式尚待提供，私有知识库摘要暂不纳入。\n\n## 来源摘要\n\n- 本地业务数据库\n- 本机 Agent 会话记录\n", time.Now().Format(time.RFC3339), usageLines.String(), auditCount)
	a.store.audit(r.Context(), "weekly_report", "agent", "all", requestOperator(r), map[string]any{"auditCount": auditCount, "messageCount": len(messages)})
	writeJSON(w, http.StatusOK, map[string]any{"markdown": report, "generatedAt": time.Now().Format(time.RFC3339), "partial": true, "generationMode": "local_aggregate", "todo": "客户确认周报模板后可增加指定 OpenClaw Agent 的内容摘要"})
}

func (a *businessAPI) pluginsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		if r.Method == http.MethodGet {
			items, err := listRecords[Plugin](r.Context(), a.store, "plugins")
			if err != nil {
				writeBusinessError(w, err)
				return
			}
			listResponse(w, items)
			return
		}
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
		var request Plugin
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		if request.ID == "" {
			writeAPIError(w, http.StatusBadRequest, "PLUGIN_ID_REQUIRED", "插件 ID 不能为空")
			return
		}
		a.updatePlugin(w, r, request.ID, request)
		return
	}
	if !allowMutation(w, r, http.MethodPatch) {
		return
	}
	var request Plugin
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	a.updatePlugin(w, r, parts[0], request)
}

func (a *businessAPI) updatePlugin(w http.ResponseWriter, r *http.Request, id string, request Plugin) {
	var item Plugin
	if err := a.store.get(r.Context(), "plugins", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	request.ID, request.UpdatedAt = item.ID, currentText()
	if request.Name == "" {
		request.Name = item.Name
	}
	if request.Capabilities == nil {
		request.Capabilities = item.Capabilities
	}
	if request.Enabled {
		request.Status = "PendingBinding"
	} else {
		request.Status = "Unbound"
	}
	if err := a.store.put(r.Context(), "plugins", item.ID, request); err != nil {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"item": request, "todo": "插件绑定凭据、推送内容和同步范围待客户确认"})
}

func (a *businessAPI) jobsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		if r.Method == http.MethodGet {
			items, err := listRecords[Job](r.Context(), a.store, "jobs")
			if err != nil {
				writeBusinessError(w, err)
				return
			}
			listResponse(w, items)
			return
		}
		if r.Method == http.MethodPatch {
			if !allowMutation(w, r, http.MethodPatch) {
				return
			}
			var request Job
			if err := decodeJSONBody(w, r, &request); err != nil {
				return
			}
			if strings.TrimSpace(request.ID) == "" {
				writeAPIError(w, http.StatusBadRequest, "JOB_ID_REQUIRED", "任务 ID 不能为空")
				return
			}
			a.updateJob(w, r, request.ID, request)
			return
		}
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var item Job
		if err := decodeJSONBody(w, r, &item); err != nil {
			return
		}
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Kind) == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_JOB", "任务名称和类型不能为空")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "jobs", "JOB", 4)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		item.ID, item.Status, item.UpdatedAt = id, "Ready", currentText()
		if err := a.store.create(r.Context(), "jobs", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	var item Job
	if err := a.store.get(r.Context(), "jobs", parts[0], &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	if r.Method == http.MethodPatch {
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
		var request Job
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		a.updateJob(w, r, item.ID, request)
		return
	}
	if r.Method == http.MethodDelete {
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if item.BuiltIn {
			writeAPIError(w, http.StatusConflict, "BUILTIN_JOB", "内置任务不能删除，只能停用")
			return
		}
		if err := a.store.softDelete(r.Context(), "jobs", item.ID); err != nil {
			writeBusinessError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
		return
	}
	w.Header().Set("Allow", "PATCH, DELETE")
	writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
}

func (a *businessAPI) updateJob(w http.ResponseWriter, r *http.Request, id string, request Job) {
	var item Job
	if err := a.store.get(r.Context(), "jobs", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	request.ID, request.BuiltIn, request.UpdatedAt = item.ID, item.BuiltIn, currentText()
	if strings.TrimSpace(request.Name) == "" {
		request.Name = item.Name
	}
	if strings.TrimSpace(request.Kind) == "" {
		request.Kind = item.Kind
	}
	if strings.TrimSpace(request.Status) == "" {
		request.Status = item.Status
	}
	if strings.TrimSpace(request.Schedule) == "" {
		request.Schedule = item.Schedule
	}
	request.LastRun = firstNonEmpty(request.LastRun, item.LastRun)
	request.NextRun = firstNonEmpty(request.NextRun, item.NextRun)
	if err := a.store.put(r.Context(), "jobs", item.ID, request); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "update", "job", item.ID, requestOperator(r), nil)
	writeJSON(w, http.StatusOK, request)
}

func (a *businessAPI) settingsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) > 0 && parts[0] == "model" {
		a.settingsModelRouter(w, r, parts[1:])
		return
	}
	if len(parts) == 0 || parts[0] != "preferences" {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "设置接口不存在")
		return
	}
	if r.Method == http.MethodGet {
		preferences := defaultPreferences()
		_ = a.store.getSetting(r.Context(), "preferences", &preferences)
		writeJSON(w, http.StatusOK, preferences)
		return
	}
	if !allowMutation(w, r, http.MethodPatch) {
		return
	}
	var request UserPreferences
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.NewsShowLimit < 1 || request.NewsShowLimit > 100 {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NEWS_LIMIT", "每次展示数量必须是 1-100")
		return
	}
	if !validNewsFrequency(request.NewsFrequency) {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NEWS_FREQUENCY", "获取频率不受支持")
		return
	}
	if request.AgentAllowlists == nil {
		request.AgentAllowlists = map[string][]string{}
	}
	if err := a.store.putSetting(r.Context(), "preferences", request); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "update", "setting", "preferences", requestOperator(r), nil)
	writeJSON(w, http.StatusOK, request)
}

func validNewsFrequency(value string) bool {
	frequencies := map[string]bool{"1小时": true, "2小时": true, "3小时": true, "6小时": true, "8小时": true, "12小时": true, "24小时": true}
	return frequencies[value]
}

func (a *businessAPI) settingsModelRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 1 && parts[0] == "test" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		started := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		status, err := a.openClaw.service.Status(ctx)
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": status.RPCOK, "durationMs": time.Since(started).Milliseconds(), "version": status.Version})
		return
	}
	if r.Method == http.MethodGet {
		a.openClaw.modelsHandler(w, r)
		return
	}
	if !allowMutation(w, r, http.MethodPatch) {
		return
	}
	var request struct {
		Model    string `json:"model"`
		Provider string `json:"provider"`
		APIKey   string `json:"apiKey"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	if request.APIKey != "" {
		if _, err := a.openClaw.service.SaveAPIKey(ctx, request.Provider, request.APIKey); err != nil {
			writeOpenClawError(w, err)
			return
		}
	}
	if request.Model != "" {
		if err := a.openClaw.service.SetDefaultModel(ctx, request.Model); err != nil {
			writeOpenClawError(w, err)
			return
		}
	}
	models, err := a.openClaw.service.Models(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	a.store.audit(r.Context(), "update", "setting", "model", requestOperator(r), map[string]any{"model": request.Model, "provider": request.Provider, "apiKeyUpdated": request.APIKey != ""})
	writeJSON(w, http.StatusOK, models)
}

func (a *businessAPI) systemRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 1 && parts[0] == "health" {
		a.systemHealth(w, r)
		return
	}
	if len(parts) >= 2 && parts[0] == "upgrade" {
		if parts[1] == "history" {
			if !allowMethod(w, r, http.MethodGet) {
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "total": 0})
			return
		}
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_UPGRADE_PACKAGE_SPEC", "离线升级接口已保留，需提供升级包清单、签名、兼容和回滚规范", []string{"manifest 格式", "签名算法", "兼容版本", "迁移脚本", "回滚验收"})
		return
	}
	writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "系统接口不存在")
}

func (a *businessAPI) systemHealth(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	databaseOK := a.store.db.PingContext(r.Context()) == nil
	var disk syscall.Statfs_t
	diskOK := syscall.Statfs(filepath.Dir(a.store.path), &disk) == nil
	available := uint64(0)
	if diskOK {
		available = disk.Bavail * uint64(disk.Bsize)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	openClawStatus, openClawErr := a.openClaw.service.Status(ctx)
	databaseBytes := int64(0)
	if info, err := os.Stat(a.store.path); err == nil {
		databaseBytes = info.Size()
	}
	privateFiles, _ := listRecords[PrivateFile](r.Context(), a.store, "private_files")
	privateBytes, indexedFiles := int64(0), 0
	for _, item := range privateFiles {
		privateBytes += item.Bytes
		if item.Status == "Indexed" {
			indexedFiles++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": map[bool]string{true: "ok", false: "degraded"}[databaseOK && openClawErr == nil], "database": map[string]any{"ok": databaseOK, "schemaVersion": businessSchemaVersion, "bytes": databaseBytes}, "openclaw": map[string]any{"ok": openClawErr == nil, "status": openClawStatus}, "runtime": map[string]any{"go": runtime.Version(), "os": runtime.GOOS, "arch": runtime.GOARCH}, "storage": map[string]any{"ok": diskOK, "availableBytes": available, "privateFileBytes": privateBytes}, "index": map[string]any{"files": len(privateFiles), "indexed": indexedFiles, "status": "waiting_raw_data_format"}, "rawData": map[string]any{"status": "todo", "reason": "客户原始数据格式待提供"}})
}

func (a *businessAPI) templatesRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "total": 0, "status": "waiting_spec"})
		return
	}
	if r.Method == http.MethodPatch {
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
	} else if !allowMutation(w, r, http.MethodPost) {
		return
	}
	writeTODO(w, "TODO_TEMPLATE_SPEC", "模板接口已保留，需客户提供正式模板、字段映射和 OCR 规则", []string{"模板样例", "占位符规范", "OCR 引擎", "纸张语言签章", "发布规则"})
}

func (a *businessAPI) agentBackupsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	source := filepath.Join(home, ".openclaw", "agents")
	if _, err := os.Stat(source); err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, "AGENT_DATA_UNAVAILABLE", "未找到 OpenClaw Agent 数据目录")
		return
	}
	backupRoot := filepath.Join(home, ".local", "share", "sta100", "agent-backups")
	if err := os.MkdirAll(backupRoot, 0700); err != nil {
		writeBusinessError(w, err)
		return
	}
	target := filepath.Join(backupRoot, "agents-"+time.Now().Format("20060102-150405")+".tar.gz")
	if err := createTarGZ(source, target); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "AGENT_BACKUP_FAILED", "Agent 备份创建失败")
		return
	}
	info, _ := os.Stat(target)
	a.store.audit(r.Context(), "backup", "agent", "all", requestOperator(r), map[string]any{"path": target})
	writeJSON(w, http.StatusCreated, map[string]any{"path": target, "bytes": info.Size(), "createdAt": currentText(), "restoreSupported": false})
}

func createTarGZ(source, target string) error {
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer output.Close()
	gzipWriter := gzip.NewWriter(output)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	base := filepath.Dir(source)
	return filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(base, path)
		if err != nil {
			return err
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tarWriter, file)
		return err
	})
}

func (a *businessAPI) openClawErrorStatus(err error) int {
	if errors.Is(err, orchestrator.ErrUnavailable) {
		return http.StatusServiceUnavailable
	}
	return http.StatusBadGateway
}

func sortedKeys(values map[string][]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
