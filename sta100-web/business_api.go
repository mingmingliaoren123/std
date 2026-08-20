package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type businessAPI struct {
	store    *businessStore
	openClaw *openClawService
}

func newBusinessAPI(store *businessStore, openClaw *openClawService) *businessAPI {
	return &businessAPI{store: store, openClaw: openClaw}
}

func (a *businessAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/"), "/")
	parts := strings.Split(path, "/")
	if path == "" || parts[0] == "" {
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "接口不存在")
		return
	}
	switch parts[0] {
	case "bootstrap":
		a.bootstrapHandler(w, r)
	case "search":
		a.searchHandler(w, r)
	case "assistant":
		a.assistantRouter(w, r, parts[1:])
	case "agent-token-usage":
		a.tokenUsageRouter(w, r, parts[1:])
	case "accounts":
		a.accountsRouter(w, r, parts[1:])
	case "products":
		a.productsRouter(w, r, parts[1:])
	case "suppliers":
		a.suppliersRouter(w, r, parts[1:])
	case "quotes":
		a.quotesRouter(w, r, parts[1:])
	case "orders":
		a.ordersRouter(w, r, parts[1:])
	case "documents":
		a.documentsRouter(w, r, parts[1:])
	case "private-files":
		a.privateFilesRouter(w, r, parts[1:])
	case "news":
		a.newsRouter(w, r, parts[1:])
	case "overview":
		a.overviewRouter(w, r, parts[1:])
	case "preferences":
		a.preferencesRouter(w, r, parts[1:])
	case "recommendations":
		a.overviewRecommendations(w, r)
	case "tasks":
		a.tasksRouter(w, r, parts[1:])
	case "agents":
		a.agentsRouter(w, r, parts[1:])
	case "plugins":
		a.pluginsRouter(w, r, parts[1:])
	case "jobs":
		a.jobsRouter(w, r, parts[1:])
	case "system":
		a.systemRouter(w, r, parts[1:])
	case "templates":
		a.templatesRouter(w, r, parts[1:])
	case "settings":
		a.settingsRouter(w, r, parts[1:])
	case "agent-backups":
		a.agentBackupsHandler(w, r)
	default:
		writeAPIError(w, http.StatusNotFound, "API_NOT_FOUND", "接口不存在")
	}
}

func (a *businessAPI) bootstrapHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx := r.Context()
	customers, err := a.listCustomers(ctx)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	quotes, err := listRecords[Quote](ctx, a.store, "quotes")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	orders, err := listRecords[Order](ctx, a.store, "orders")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	documents, err := listRecords[Document](ctx, a.store, "documents")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	products, err := listRecords[Product](ctx, a.store, "products")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	suppliers, err := listRecords[Supplier](ctx, a.store, "suppliers")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	files, err := listRecords[PrivateFile](ctx, a.store, "private_files")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	news, err := listRecords[NewsItem](ctx, a.store, "news")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	recommendations, err := listRecords[Recommendation](ctx, a.store, "recommendations")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	news = filterDisplayableNews(news)
	recommendations = filterDisplayableRecommendations(recommendations)
	jobs, err := listRecords[Job](ctx, a.store, "jobs")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	plugins, err := listRecords[Plugin](ctx, a.store, "plugins")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	preferences := defaultPreferences()
	if err := a.store.getSetting(ctx, "preferences", &preferences); err != nil && !errors.Is(err, errRecordNotFound) {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"schemaVersion":   businessSchemaVersion,
		"dataStatus":      "seeded_demo_data",
		"customers":       customers,
		"quotes":          quotes,
		"orders":          orders,
		"documents":       documents,
		"products":        products,
		"suppliers":       suppliers,
		"files":           files,
		"news":            news,
		"recommendations": recommendations,
		"jobs":            jobs,
		"plugins":         plugins,
		"overview":        a.overviewSummaryData(ctx),
		"automation":      a.overviewAutomationData(ctx),
		"preferences":     preferences,
	})
}

func (a *businessAPI) searchHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "total": 0})
		return
	}
	results := make([]map[string]string, 0)
	customers, _ := a.listCustomers(r.Context())
	for _, item := range customers {
		if containsFold(item.ID+" "+item.Name+" "+item.Contact+" "+item.Country, query) {
			results = append(results, map[string]string{"type": "客户", "id": item.ID, "label": item.Name, "sub": item.Country + " · " + item.Contact, "page": "customers"})
		}
	}
	orders, _ := listRecords[Order](r.Context(), a.store, "orders")
	for _, item := range orders {
		if containsFold(item.ID+" "+item.Customer+" "+item.Products, query) {
			results = append(results, map[string]string{"type": "订单", "id": item.ID, "label": item.ID, "sub": item.Customer + " · " + item.Value, "page": "orders"})
		}
	}
	products, _ := listRecords[Product](r.Context(), a.store, "products")
	for _, item := range products {
		if containsFold(item.ID+" "+item.Name+" "+item.Category, query) {
			results = append(results, map[string]string{"type": "产品", "id": item.ID, "label": item.Name, "sub": item.ID + " · " + item.Category, "page": "products"})
		}
	}
	files, _ := listRecords[PrivateFile](r.Context(), a.store, "private_files")
	for _, item := range files {
		if containsFold(item.ID+" "+item.Name+" "+item.Category+" "+strings.Join(item.Tags, " "), query) {
			results = append(results, map[string]string{"type": "文件", "id": item.ID, "label": item.Name, "sub": item.Category + " · " + item.Source, "page": "database"})
		}
	}
	if len(results) > 50 {
		results = results[:50]
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": results, "total": len(results)})
}

func writeBusinessError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errRecordNotFound):
		writeAPIError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "记录不存在或已归档")
	case errors.Is(err, errRecordConflict):
		writeAPIError(w, http.StatusConflict, "RECORD_CONFLICT", "记录编号或唯一字段已存在")
	default:
		writeAPIError(w, http.StatusInternalServerError, "BUSINESS_STORE_FAILED", "本地业务数据处理失败")
	}
}

func writeTODO(w http.ResponseWriter, code, message string, missing []string) {
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"error": map[string]any{"code": code, "message": message, "missing": missing},
	})
}

func requestOperator(r *http.Request) string {
	if username, ok := r.Context().Value(authUsernameContextKey{}).(string); ok && username != "" {
		return username
	}
	return "local-user"
}

func containsFold(value, lowerQuery string) bool {
	return strings.Contains(strings.ToLower(value), lowerQuery)
}

func currentText() string { return time.Now().Format("2006-01-02 15:04") }

func moneyNumber(value string) float64 {
	cleaned := strings.NewReplacer("EUR", "", "USD", "", "CNY", "", "GBP", "", ",", "", " ", "").Replace(value)
	number, _ := strconv.ParseFloat(cleaned, 64)
	return number
}

func formatMoney(value float64, currency string) string {
	if currency == "" {
		currency = "EUR"
	}
	return fmt.Sprintf("%s %.2f", currency, value)
}

func lineSummary(lines []BusinessLine) string {
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		parts = append(parts, fmt.Sprintf("%s x %s", line.ProductName, strconv.FormatFloat(line.Quantity, 'f', -1, 64)))
	}
	return strings.Join(parts, "；")
}

func listResponse[T any](w http.ResponseWriter, items []T) {
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func filterPage[T any](items []T, r *http.Request) []T {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page <= 0 || pageSize <= 0 {
		return items
	}
	if pageSize > 200 {
		pageSize = 200
	}
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []T{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func sortStrings(values []string) []string {
	result := append([]string{}, values...)
	sort.Strings(result)
	return result
}

func writeCSV(w http.ResponseWriter, filename string, rows [][]string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(w)
	_ = writer.WriteAll(rows)
}
