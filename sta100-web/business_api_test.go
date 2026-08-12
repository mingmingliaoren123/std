package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"openclaw-orchestrator/orchestrator"
)

func newTestBusinessAPI(t *testing.T) (*businessAPI, *businessStore) {
	t.Helper()
	t.Setenv("STA100_DB_PATH", filepath.Join(t.TempDir(), "sta100.db"))
	store, err := newBusinessStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return newBusinessAPI(store, newOpenClawService()), store
}

func businessRequest(t *testing.T, api *businessAPI, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var input bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		input = *bytes.NewReader(data)
	} else {
		input = *bytes.NewReader(nil)
	}
	request := httptest.NewRequest(method, path, &input)
	request = request.WithContext(context.WithValue(request.Context(), authUsernameContextKey{}, "tester"))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if method != http.MethodGet {
		request.Header.Set("X-STA100-Request", "1")
	}
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	return recorder
}

func decodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(recorder.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode response (%d): %v\n%s", recorder.Code, err, recorder.Body.String())
	}
	return value
}

func TestBootstrapReturnsPersistentBusinessData(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	recorder := businessRequest(t, api, http.MethodGet, "/api/v1/bootstrap", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", recorder.Code, recorder.Body.String())
	}
	response := decodeResponse[struct {
		SchemaVersion int        `json:"schemaVersion"`
		Customers     []Customer `json:"customers"`
		Products      []Product  `json:"products"`
		Orders        []Order    `json:"orders"`
	}](t, recorder)
	if response.SchemaVersion != businessSchemaVersion || len(response.Customers) != 5 || len(response.Products) != 4 || len(response.Orders) != 4 {
		t.Fatalf("unexpected bootstrap response: %+v", response)
	}
	if response.Customers[0].Total == "" {
		t.Fatal("customer aggregate total was not calculated")
	}
}

func TestAccountCRUDPersistsAndArchives(t *testing.T) {
	api, store := newTestBusinessAPI(t)
	created := businessRequest(t, api, http.MethodPost, "/api/v1/accounts", Customer{Name: "API Test Cycling", Type: "Distributor", Country: "德国", Phone: "+49 1", Rating: "Prospect"})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d: %s", created.Code, created.Body.String())
	}
	account := decodeResponse[Customer](t, created)
	if account.ID == "" || account.Total != "EUR 0.00" {
		t.Fatalf("invalid created account: %+v", account)
	}
	account.City = "柏林"
	updated := businessRequest(t, api, http.MethodPatch, "/api/v1/accounts/"+account.ID, account)
	if updated.Code != http.StatusOK || decodeResponse[Customer](t, updated).City != "柏林" {
		t.Fatalf("update failed: %d %s", updated.Code, updated.Body.String())
	}
	var stored Customer
	if err := store.get(context.Background(), "accounts", account.ID, &stored); err != nil || stored.City != "柏林" {
		t.Fatalf("record not persisted: %+v %v", stored, err)
	}
	archived := businessRequest(t, api, http.MethodDelete, "/api/v1/accounts/"+account.ID, nil)
	if archived.Code != http.StatusOK {
		t.Fatalf("archive failed: %d %s", archived.Code, archived.Body.String())
	}
	if err := store.get(context.Background(), "accounts", account.ID, &stored); err != nil || !stored.Archived {
		t.Fatalf("archive flag not persisted: %+v %v", stored, err)
	}
}

func TestCustomerCommunicationsAreAppendOnlyAndRetained(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	const customerID = "ACC-0001"
	created := businessRequest(t, api, http.MethodPost, "/api/v1/accounts/"+customerID+"/communications", CustomerCommunication{
		Type: "电话", Subject: "德国渠道跟进", Content: "客户确认下周评估样机。", Contact: "Anna", OccurredAt: "2026-08-12T09:30",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("communication create = %d: %s", created.Code, created.Body.String())
	}
	record := decodeResponse[CustomerCommunication](t, created)
	if record.ID == "" || record.CustomerID != customerID || record.CreatedBy != "tester" || record.CreatedAt == "" {
		t.Fatalf("invalid communication: %+v", record)
	}
	listed := businessRequest(t, api, http.MethodGet, "/api/v1/accounts/"+customerID+"/communications", nil)
	items := decodeResponse[struct {
		Items []CustomerCommunication `json:"items"`
	}](t, listed).Items
	if listed.Code != http.StatusOK || len(items) != 1 || items[0].ID != record.ID {
		t.Fatalf("communication list = %d: %s", listed.Code, listed.Body.String())
	}
	for _, method := range []string{http.MethodPatch, http.MethodDelete} {
		response := businessRequest(t, api, method, "/api/v1/accounts/"+customerID+"/communications", map[string]any{})
		if response.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s communication status = %d, want 405", method, response.Code)
		}
	}
	archived := businessRequest(t, api, http.MethodDelete, "/api/v1/accounts/"+customerID, nil)
	if archived.Code != http.StatusOK {
		t.Fatalf("archive customer = %d: %s", archived.Code, archived.Body.String())
	}
	retained := businessRequest(t, api, http.MethodGet, "/api/v1/accounts/"+customerID+"/communications", nil)
	if retained.Code != http.StatusOK || len(decodeResponse[struct {
		Items []CustomerCommunication `json:"items"`
	}](t, retained).Items) != 1 {
		t.Fatalf("communication history not retained: %d %s", retained.Code, retained.Body.String())
	}
	rejected := businessRequest(t, api, http.MethodPost, "/api/v1/accounts/"+customerID+"/communications", CustomerCommunication{Type: "邮件", Content: "new", OccurredAt: "2026-08-12T10:00"})
	if rejected.Code != http.StatusConflict {
		t.Fatalf("archived customer append = %d, want 409: %s", rejected.Code, rejected.Body.String())
	}
}

func TestAssistantQueryRejectsSystemAgentTarget(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	response := businessRequest(t, api, http.MethodPost, "/api/v1/assistant/query", map[string]any{
		"page": "agents", "feature": "agent-chat", "message": "test", "sessionKey": "sta100-test",
		"context": map[string]any{"targetAgent": coordinatorAgentID},
	})
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "INVALID_TARGET_AGENT") {
		t.Fatalf("system agent target = %d: %s", response.Code, response.Body.String())
	}
}

func TestTokenUsageSummaryAndClear(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	api.recordTokenUsage(context.Background(), "AI-1", "agents", "agent-chat", "domain-agent", "export-agent", "ok", orchestrator.AgentMessageResult{
		RunID: "run-1", Usage: orchestrator.TokenUsage{Input: 100, Output: 20, CacheRead: 5, Total: 125},
	})
	api.recordTokenUsage(context.Background(), "AI-1", "agents", "agent-chat", "coordinator-agent", coordinatorAgentID, "failed", orchestrator.AgentMessageResult{})
	response := businessRequest(t, api, http.MethodGet, "/api/v1/agent-token-usage?request_id=AI-1", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("token summary = %d: %s", response.Code, response.Body.String())
	}
	summary := decodeResponse[tokenUsageSummary](t, response)
	if summary.Total != 125 || summary.CurrentRequestTotal != 125 || summary.Calls != 2 || summary.MeasuredCalls != 1 || summary.UnavailableCalls != 1 {
		t.Fatalf("unexpected token summary: %+v", summary)
	}
	latest := decodeResponse[tokenUsageSummary](t, businessRequest(t, api, http.MethodGet, "/api/v1/agent-token-usage", nil))
	if latest.CurrentRequestID != "AI-1" || latest.CurrentRequestTotal != 125 {
		t.Fatalf("latest request usage not retained: %+v", latest)
	}
	cleared := businessRequest(t, api, http.MethodDelete, "/api/v1/agent-token-usage", nil)
	if cleared.Code != http.StatusOK {
		t.Fatalf("clear token usage = %d: %s", cleared.Code, cleared.Body.String())
	}
	empty := decodeResponse[tokenUsageSummary](t, businessRequest(t, api, http.MethodGet, "/api/v1/agent-token-usage", nil))
	if empty.Calls != 0 || empty.Total != 0 {
		t.Fatalf("token usage was not cleared: %+v", empty)
	}
}

func TestCustomerDiscoveryDoesNotSubstituteUnmatchedLocalCustomer(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	request := assistantQueryRequest{Feature: "customer-discovery", Context: map[string]any{
		"country": "德国", "city": "柏林", "type": "Distributor",
	}}
	evidence := api.collectLocalEvidence(context.Background(), assistantQueryRequest{Message: "德国 客户"})
	if items := api.localAssistantItems(request, evidence); len(items) != 0 {
		t.Fatalf("unmatched discovery returned local customer: %+v", items)
	}
}

func TestQuoteCalculationConversionAndDocumentGeneration(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	quoteInput := Quote{
		Subject: "API 自动计算报价", Customer: "VeloTrade GmbH", Valid: "2026-12-31", Currency: "EUR",
		Freight: 20, Tax: 5,
		Lines: []BusinessLine{{ProductID: "STA-100-EU", Quantity: 2, UnitPrice: 100, Discount: 10}},
	}
	created := businessRequest(t, api, http.MethodPost, "/api/v1/quotes", quoteInput)
	if created.Code != http.StatusCreated {
		t.Fatalf("quote create = %d: %s", created.Code, created.Body.String())
	}
	quote := decodeResponse[Quote](t, created)
	if quote.Value != "EUR 205.00" || len(quote.Lines) != 1 || quote.Lines[0].Amount != 180 {
		t.Fatalf("server calculation incorrect: %+v", quote)
	}
	quote.Status = "Accepted"
	accepted := businessRequest(t, api, http.MethodPatch, "/api/v1/quotes/"+quote.ID, quote)
	if accepted.Code != http.StatusOK {
		t.Fatalf("quote accept = %d: %s", accepted.Code, accepted.Body.String())
	}
	converted := businessRequest(t, api, http.MethodPost, "/api/v1/quotes/"+quote.ID+"/convert-order", map[string]any{})
	if converted.Code != http.StatusCreated {
		t.Fatalf("convert = %d: %s", converted.Code, converted.Body.String())
	}
	order := decodeResponse[Order](t, converted)
	if order.Quote != quote.ID || order.Value != "EUR 180.00" || order.Lines[0].UnitPrice != 90 {
		t.Fatalf("quote snapshot was not converted correctly: %+v", order)
	}
	documents := businessRequest(t, api, http.MethodPost, "/api/v1/orders/"+order.ID+"/documents", map[string]any{"types": []string{"PI", "CI"}, "language": "英文"})
	if documents.Code != http.StatusCreated {
		t.Fatalf("generate documents = %d: %s", documents.Code, documents.Body.String())
	}
	result := decodeResponse[struct {
		Items []Document `json:"items"`
	}](t, documents)
	if len(result.Items) != 2 || result.Items[0].Order != order.ID || len(result.Items[0].Lines) != 1 {
		t.Fatalf("invalid documents: %+v", result.Items)
	}
}

func TestOrderRejectsInsufficientStock(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	order := Order{Customer: "VeloTrade GmbH", Delivery: "2026-12-31", Currency: "EUR", Lines: []BusinessLine{{ProductID: "EBK-CITY-03", Quantity: 20, UnitPrice: 1000}}}
	recorder := businessRequest(t, api, http.MethodPost, "/api/v1/orders", order)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400: %s", recorder.Code, recorder.Body.String())
	}
}

func TestReferencedProductCannotBeDeleted(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	recorder := businessRequest(t, api, http.MethodDelete, "/api/v1/products/STA-100-EU", nil)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRawDataDependentEndpointsAreExplicitTODO(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	for _, test := range []struct {
		path string
		body any
	}{
		{"/api/v1/overview/oem/match", map[string]any{"query": "battery"}},
		{"/api/v1/overview/customer-discovery", map[string]any{"country": "德国", "city": "柏林", "type": "Distributor"}},
		{"/api/v1/news/refresh", map[string]any{}},
		{"/api/v1/templates/upload", map[string]any{}},
		{"/api/v1/system/upgrade/import", map[string]any{}},
	} {
		recorder := businessRequest(t, api, http.MethodPost, test.path, test.body)
		if recorder.Code != http.StatusNotImplemented {
			t.Errorf("%s status = %d, want 501: %s", test.path, recorder.Code, recorder.Body.String())
		}
	}
}

func TestBusinessMutationRequiresRequestHeader(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	data, _ := json.Marshal(Customer{Name: "No Header", Country: "德国"})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/accounts", bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}

func TestCustomerRenameUpdatesBusinessSnapshots(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	created := businessRequest(t, api, http.MethodPost, "/api/v1/quotes", Quote{
		Subject: "改名关联测试", Customer: "VeloTrade GmbH", Valid: "2026-12-31", Currency: "EUR",
		Lines: []BusinessLine{{ProductID: "STA-100-EU", Quantity: 1, UnitPrice: 159}},
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("quote create = %d: %s", created.Code, created.Body.String())
	}
	quote := decodeResponse[Quote](t, created)
	accountResponse := businessRequest(t, api, http.MethodGet, "/api/v1/accounts", nil)
	accounts := decodeResponse[struct {
		Items []Customer `json:"items"`
	}](t, accountResponse).Items
	var account Customer
	for _, item := range accounts {
		if item.Name == "VeloTrade GmbH" {
			account = item
			break
		}
	}
	if account.ID == "" {
		t.Fatal("seed customer not found")
	}
	account.Name = "VeloTrade GmbH Renamed"
	updated := businessRequest(t, api, http.MethodPatch, "/api/v1/accounts/"+account.ID, account)
	if updated.Code != http.StatusOK {
		t.Fatalf("account rename = %d: %s", updated.Code, updated.Body.String())
	}
	quoteResponse := businessRequest(t, api, http.MethodGet, "/api/v1/quotes/"+quote.ID, nil)
	if got := decodeResponse[Quote](t, quoteResponse).Customer; got != account.Name {
		t.Fatalf("quote customer snapshot = %q, want %q", got, account.Name)
	}
}

func TestExcelExportsAreXLSX(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	for _, path := range []string{"/api/v1/accounts/export", "/api/v1/suppliers/export"} {
		recorder := businessRequest(t, api, http.MethodGet, path, nil)
		if recorder.Code != http.StatusOK || recorder.Header().Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			t.Fatalf("%s response = %d %q", path, recorder.Code, recorder.Header().Get("Content-Type"))
		}
		if !bytes.HasPrefix(recorder.Body.Bytes(), []byte("PK")) {
			t.Fatalf("%s is not a zip based xlsx", path)
		}
		archive, err := zip.NewReader(bytes.NewReader(recorder.Body.Bytes()), int64(recorder.Body.Len()))
		if err != nil {
			t.Fatalf("%s xlsx cannot be opened: %v", path, err)
		}
		foundWorksheet := false
		for _, file := range archive.File {
			if file.Name != "xl/worksheets/sheet1.xml" {
				continue
			}
			foundWorksheet = true
			reader, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(reader)
			_ = reader.Close()
			if err != nil {
				t.Fatal(err)
			}
			var worksheet struct{ XMLName xml.Name }
			if err := xml.Unmarshal(data, &worksheet); err != nil || worksheet.XMLName.Local != "worksheet" {
				t.Fatalf("%s worksheet XML is invalid: %v", path, err)
			}
		}
		if !foundWorksheet {
			t.Fatalf("%s does not contain sheet1.xml", path)
		}
	}
}

func TestCollectionPatchRoutes(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	job := Job{ID: "JOB-RECOMMEND", Enabled: false, Schedule: "每 2 小时"}
	updatedJob := businessRequest(t, api, http.MethodPatch, "/api/v1/jobs", job)
	if updatedJob.Code != http.StatusOK {
		t.Fatalf("jobs collection patch = %d: %s", updatedJob.Code, updatedJob.Body.String())
	}
	plugin := Plugin{ID: "wechat", Enabled: true}
	updatedPlugin := businessRequest(t, api, http.MethodPatch, "/api/v1/plugins", plugin)
	if updatedPlugin.Code != http.StatusOK {
		t.Fatalf("plugins collection patch = %d: %s", updatedPlugin.Code, updatedPlugin.Body.String())
	}
}

func TestPrivateFileUploadAndRawDataTODO(t *testing.T) {
	t.Setenv("STA100_PRIVATE_DATA_DIR", filepath.Join(t.TempDir(), "private-files"))
	api, _ := newTestBusinessAPI(t)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "customer.md")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, "customer data pending schema")
	_ = writer.WriteField("category", "客户资料")
	_ = writer.WriteField("tags", "客户,待确认")
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/private-files/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-STA100-Request", "1")
	request = request.WithContext(context.WithValue(request.Context(), authUsernameContextKey{}, "tester"))
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("upload = %d: %s", recorder.Code, recorder.Body.String())
	}
	response := decodeResponse[struct {
		Item PrivateFile `json:"item"`
	}](t, recorder)
	if response.Item.Status != "PendingParse" || response.Item.Path != "" {
		t.Fatalf("unexpected uploaded metadata: %+v", response.Item)
	}
	content := businessRequest(t, api, http.MethodGet, "/api/v1/private-files/"+response.Item.ID+"/content", nil)
	if content.Code != http.StatusOK || string(content.Body.Bytes()) != "customer data pending schema" {
		t.Fatalf("content = %d %q", content.Code, content.Body.String())
	}
	summary := businessRequest(t, api, http.MethodGet, "/api/v1/private-files/"+response.Item.ID+"/summary", nil)
	if summary.Code != http.StatusNotImplemented || !strings.Contains(summary.Body.String(), "TODO_PRIVATE_FILE_PARSER") {
		t.Fatalf("summary = %d: %s", summary.Code, summary.Body.String())
	}
}

func TestDesignTODOAliasesDoNotFallThroughTo404(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	for _, path := range []string{
		"/api/v1/overview/oem-matches",
		"/api/v1/overview/oem-matches/export",
		"/api/v1/overview/oem-matches/FACTORY-1",
		"/api/v1/news/refresh",
		"/api/v1/templates/upload",
		"/api/v1/system/upgrade/import",
		"/api/v1/products/import",
		"/api/v1/tasks",
	} {
		recorder := businessRequest(t, api, http.MethodPost, path, map[string]any{})
		if recorder.Code == http.StatusNotFound {
			t.Errorf("%s fell through to 404: %s", path, recorder.Body.String())
		}
	}
}
