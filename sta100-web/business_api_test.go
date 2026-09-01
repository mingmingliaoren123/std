package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
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
	seedTestBusinessFixtures(t, store)
	return newBusinessAPI(store, newOpenClawService()), store
}

func seedTestBusinessFixtures(t *testing.T, store *businessStore) {
	t.Helper()
	ctx := context.Background()
	customers := []Customer{
		{ID: "TACC-0001", Name: "Fixture Customer Alpha", Type: "Distributor", Country: "德国", City: "柏林", Contact: "Anna", Phone: "+49 30 1000 0001", Email: "alpha@example.test", Owner: "tester", Rating: "Active", Source: "展会", Updated: "2026-08-10 09:42"},
		{ID: "TACC-0002", Name: "Fixture Customer Beta", Type: "Importer", Country: "瑞典", City: "斯德哥尔摩", Contact: "Erik", Phone: "+46 8 1000 0002", Email: "beta@example.test", Owner: "tester", Rating: "Prospect", Source: "朋友介绍", Updated: "2026-08-09 16:20"},
		{ID: "TACC-0003", Name: "Fixture Customer Gamma", Type: "Customer", Country: "波兰", City: "华沙", Contact: "Marek", Phone: "+48 22 1000 0003", Email: "gamma@example.test", Owner: "tester", Rating: "Active", Source: "电话", Updated: "2026-08-08 11:15"},
		{ID: "TACC-0004", Name: "Fixture Customer Delta", Type: "Reseller", Country: "西班牙", City: "马德里", Contact: "Lucia", Phone: "+34 91 1000 0004", Email: "delta@example.test", Owner: "tester", Rating: "Prospect", Source: "拜访", Updated: "2026-08-07 14:02"},
		{ID: "TACC-0005", Name: "Fixture Customer Epsilon", Type: "Integrator", Country: "法国", City: "巴黎", Contact: "Louis", Phone: "+33 1 1000 0005", Email: "epsilon@example.test", Owner: "tester", Rating: "Active", Source: "互联网线索", Updated: "2026-08-06 10:30"},
	}
	products := []Product{
		{ID: "STA-100-EU", Name: "Fixture Edge Device", Category: "智能设备", Manufacturer: "Fixture Maker", Price: "EUR 159.00", Stock: 426, HS: "8471504090", Status: "Active", Description: "测试用边缘设备。", Updated: "2026-08-10 08:20"},
		{ID: "PM-DUAL-01", Name: "Fixture Power Meter", Category: "智能骑行", Manufacturer: "Fixture Maker", Price: "EUR 438.00", Stock: 68, HS: "9029209000", Status: "Active", Description: "测试用功率计。", Updated: "2026-08-09 08:20"},
		{ID: "GPS-PRO-02", Name: "Fixture GPS Computer", Category: "智能骑行", Manufacturer: "Fixture Maker", Price: "EUR 219.00", Stock: 112, HS: "8526919090", Status: "Active", Description: "测试用码表。", Updated: "2026-08-08 08:20"},
		{ID: "EBK-CITY-03", Name: "Fixture E-bike Kit", Category: "整车方案", Manufacturer: "Fixture Maker", Price: "EUR 1,280.00", Stock: 19, HS: "8711601000", Status: "Review", Description: "测试用整车方案。", Updated: "2026-08-07 08:20"},
	}
	quotes := []Quote{
		{ID: "QUO-2026-0188", Subject: "Fixture Quote Alpha", Customer: customers[0].Name, Value: "EUR 34,800.00", Valid: "2026-09-08", Status: "Delivered", Products: "Fixture Edge Device x 200", Owner: "tester", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "Fixture Edge Device", Quantity: 200, UnitPrice: 174, Amount: 34800}}, Updated: "2026-08-10 09:40"},
		{ID: "QUO-2026-0187", Subject: "Fixture Quote Beta", Customer: customers[1].Name, Value: "EUR 18,450.00", Valid: "2026-09-02", Status: "Draft", Products: "Fixture Power Meter x 50", Owner: "tester", Currency: "EUR", Lines: []BusinessLine{{ProductID: "PM-DUAL-01", ProductName: "Fixture Power Meter", Quantity: 50, UnitPrice: 369, Amount: 18450}}, Updated: "2026-08-09 09:40"},
		{ID: "QUO-2026-0185", Subject: "Fixture Quote Gamma", Customer: customers[2].Name, Value: "EUR 22,100.00", Valid: "2026-08-26", Status: "Accepted", Products: "Fixture Edge Device x 120", Owner: "tester", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "Fixture Edge Device", Quantity: 120, UnitPrice: 184.1666667, Amount: 22100}}, Updated: "2026-08-08 09:40"},
		{ID: "QUO-2026-0179", Subject: "Fixture Quote Delta", Customer: customers[3].Name, Value: "EUR 5,760.00", Valid: "2026-08-15", Status: "Rejected", Products: "Fixture Edge Device x 32", Owner: "tester", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "Fixture Edge Device", Quantity: 32, UnitPrice: 180, Amount: 5760}}, Updated: "2026-08-07 09:40"},
	}
	orders := []Order{
		{ID: "SO-2026-0106", Customer: customers[2].Name, Quote: "QUO-2026-0185", Products: "Fixture Edge Device x 120", Value: "EUR 22,100.00", Currency: "EUR", Status: "Production", Delivery: "2026-09-10", Progress: 55, Lines: quotes[2].Lines, Updated: "2026-08-10 10:20"},
		{ID: "SO-2026-0105", Customer: customers[0].Name, Quote: "QUO-2026-0188", Products: "Fixture Edge Device x 80", Value: "EUR 14,900.00", Currency: "EUR", Status: "Confirmed", Delivery: "2026-08-30", Progress: 30, Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "Fixture Edge Device", Quantity: 80, UnitPrice: 186.25, Amount: 14900}}, Updated: "2026-08-09 10:20"},
		{ID: "SO-2026-0102", Customer: customers[4].Name, Quote: "QUO-2026-0187", Products: "Fixture Power Meter x 40", Value: "EUR 31,600.00", Currency: "EUR", Status: "Shipped", Delivery: "2026-08-18", Progress: 82, Lines: []BusinessLine{{ProductID: "PM-DUAL-01", ProductName: "Fixture Power Meter", Quantity: 40, UnitPrice: 790, Amount: 31600}}, Updated: "2026-08-08 10:20"},
		{ID: "SO-2026-0098", Customer: customers[1].Name, Quote: "QUO-2026-0187", Products: "Fixture Edge Device x 25", Value: "EUR 4,650.00", Currency: "EUR", Status: "Completed", Delivery: "2026-08-03", Progress: 100, Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "Fixture Edge Device", Quantity: 25, UnitPrice: 186, Amount: 4650}}, Updated: "2026-08-07 10:20"},
	}
	documents := []Document{
		{ID: "PI-20260810-003", Type: "PI", Customer: customers[2].Name, Order: "SO-2026-0106", Template: "Fixture PI", Language: "英文", Status: "Draft", Value: orders[0].Value, Lines: orders[0].Lines, Updated: "2026-08-10 10:22"},
		{ID: "CI-20260809-012", Type: "CI", Customer: customers[4].Name, Order: "SO-2026-0102", Template: "Fixture CI", Language: "英文", Status: "Confirmed", Value: orders[2].Value, Lines: orders[2].Lines, Updated: "2026-08-09 17:45"},
		{ID: "PL-20260809-009", Type: "PL", Customer: customers[4].Name, Order: "SO-2026-0102", Template: "Fixture PL", Language: "英文", Status: "Confirmed", Value: orders[2].Value, Lines: orders[2].Lines, Updated: "2026-08-09 17:41"},
		{ID: "CD-20260808-006", Type: "报关单", Customer: customers[0].Name, Order: "SO-2026-0105", Template: "Fixture CD", Language: "中文 / 英文双语", Status: "Review", Value: orders[1].Value, Lines: orders[1].Lines, Updated: "2026-08-08 15:12"},
	}
	suppliers := []Supplier{
		{ID: "SUP-0001", Company: "Fixture Supplier Alpha", Phone: "0512-1000 0001", Contact: "Zhou", Email: "supplier-alpha@example.test", Product: "E-bike 电池", Specification: "36V / 48V", Quote: "EUR 82 / 组", Notes: "测试供应商", Source: "展会", Updated: "2026-08-10 08:50"},
		{ID: "SUP-0002", Company: "Fixture Supplier Beta", Phone: "0769-1000 0002", Contact: "Chen", Email: "supplier-beta@example.test", Product: "码表、功率计", Specification: "ANT+ / BLE", Quote: "EUR 38 / 件起", Notes: "测试供应商", Source: "朋友介绍", Updated: "2026-08-09 15:20"},
		{ID: "SUP-0003", Company: "Fixture Supplier Gamma", Phone: "0592-1000 0003", Contact: "Lin", Email: "supplier-gamma@example.test", Product: "头盔、骑行装备", Specification: "EN1078", Quote: "EUR 19 / 件起", Notes: "测试供应商", Source: "电话", Updated: "2026-08-08 10:05"},
	}
	privateFiles := []PrivateFile{
		{ID: "FILE-0001", Name: "Fixture_Product_Manual.pdf", Category: "产品手册", Tags: []string{"Fixture", "英文"}, Size: "6.8 MB", Source: "测试资料", Status: "Indexed", Updated: "2026-08-10 08:20"},
		{ID: "FILE-0002", Name: "Fixture_Agreement_Template.docx", Category: "合同", Tags: []string{"测试", "合同"}, Size: "1.2 MB", Source: "测试资料", Status: "Indexed", Updated: "2026-08-09 16:10"},
		{ID: "FILE-0003", Name: "Fixture_Regulation_Reference.pdf", Category: "法规", Tags: []string{"测试", "E-bike"}, Size: "3.1 MB", Source: "测试资料", Status: "Indexed", Updated: "2026-08-09 12:45"},
		{ID: "FILE-0004", Name: "Fixture_Compatibility_Matrix.xlsx", Category: "产品资料", Tags: []string{"兼容", "组件"}, Size: "846 KB", Source: "测试资料", Status: "Review", Updated: "2026-08-08 18:33"},
	}
	for _, entry := range []struct {
		kind  string
		items any
	}{
		{"accounts", customers}, {"products", products}, {"quotes", quotes}, {"orders", orders},
		{"documents", documents}, {"suppliers", suppliers}, {"private_files", privateFiles},
	} {
		if err := store.seedSlice(ctx, entry.kind, entry.items); err != nil {
			t.Fatalf("seed fixture %s: %v", entry.kind, err)
		}
	}
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
	const customerID = "TACC-0001"
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

func TestLeadCommunicationsConvertToCustomerCommunications(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	created := businessRequest(t, api, http.MethodPost, "/api/v1/leads", Lead{
		Name: "Berlin Lead", Type: "Distributor", Country: "德国", City: "柏林", Contact: "Mia",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("lead create = %d: %s", created.Code, created.Body.String())
	}
	lead := decodeResponse[Lead](t, created)
	appended := businessRequest(t, api, http.MethodPost, "/api/v1/leads/"+lead.ID+"/communications", LeadCommunication{
		Type: "邮件", Subject: "首次跟进", Content: "客户需要 STA-100 资料。", Contact: "Mia", OccurredAt: "2026-08-21T10:00",
	})
	if appended.Code != http.StatusCreated {
		t.Fatalf("lead communication append = %d: %s", appended.Code, appended.Body.String())
	}
	record := decodeResponse[LeadCommunication](t, appended)
	if record.ID == "" || record.LeadID != lead.ID || record.CreatedBy != "tester" {
		t.Fatalf("invalid lead communication: %+v", record)
	}
	converted := businessRequest(t, api, http.MethodPost, "/api/v1/leads/"+lead.ID+"/convert", map[string]any{"keepSource": true})
	if converted.Code != http.StatusCreated {
		t.Fatalf("lead convert = %d: %s", converted.Code, converted.Body.String())
	}
	result := decodeResponse[struct {
		Customer Customer `json:"customer"`
	}](t, converted)
	listed := businessRequest(t, api, http.MethodGet, "/api/v1/accounts/"+result.Customer.ID+"/communications", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("customer communications = %d: %s", listed.Code, listed.Body.String())
	}
	items := decodeResponse[struct {
		Items []CustomerCommunication `json:"items"`
	}](t, listed).Items
	if len(items) != 1 || items[0].Content != record.Content || items[0].CustomerID != result.Customer.ID {
		t.Fatalf("converted communication not retained: %+v", items)
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

func TestAssistantQueryRequiresUsableDefaultModel(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	api.openClaw = &openClawService{service: orchestrator.New(orchestrator.Config{})}
	started := time.Now()
	response := businessRequest(t, api, http.MethodPost, "/api/v1/assistant/query", map[string]any{
		"page": "overview", "feature": "customer-discovery", "message": "国家=德国，城市=柏林，客户类型=Distributor",
		"sessionKey": "sta100-test-model-required",
		"context":    map[string]any{"targetAgent": "customer-measurement-agent", "country": "德国", "cities": []string{"柏林"}, "types": []string{"Distributor"}},
	})
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "MODEL_NOT_READY") {
		t.Fatalf("missing default model = %d: %s", response.Code, response.Body.String())
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("missing model should fail before OpenClaw call, took %s", elapsed)
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
	evidence := api.collectLocalEvidence(context.Background(), request)
	if len(evidence) != 0 {
		t.Fatalf("discovery must not read local evidence: %+v", evidence)
	}
	if items := api.localAssistantItems(request, evidence); len(items) != 0 {
		t.Fatalf("unmatched discovery returned local customer: %+v", items)
	}
}

func TestLocalCustomerSearchItemsKeepLocalSource(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	items := api.localCustomerSearchItems(context.Background(), assistantQueryRequest{
		Feature: "customer-search", Message: "德国 柏林 Distributor", Context: map[string]any{"hasContact": true},
	})
	if len(items) != 1 {
		t.Fatalf("expected one local customer result, got %+v", items)
	}
	if items[0]["name"] != "Fixture Customer Alpha" || items[0]["sourceType"] != "local_business" || items[0]["sourceLabel"] != "本地业务数据库" {
		t.Fatalf("local customer source was not retained: %+v", items[0])
	}
}

func TestLocalCustomerSearchMatchesChineseChannelTermsAcrossCustomersAndLeads(t *testing.T) {
	api, store := newTestBusinessAPI(t)
	lead := Lead{ID: "TLEAD-0001", Name: "Fixture Bike Dealer", Type: "Dealer", Country: "Germany", City: "Berlin", Phone: "+49 30 1000 0009", Business: "自行车与 E-bike 经销商", Source: "本地客户发现", Updated: "2026-08-31T10:00:00Z"}
	if err := store.create(context.Background(), "leads", lead.ID, lead); err != nil {
		t.Fatal(err)
	}
	items := api.localCustomerSearchItems(context.Background(), assistantQueryRequest{Feature: "customer-search", Message: "骑行相关的经销商", Context: map[string]any{"hasContact": true}})
	foundCustomer, foundLead := false, false
	for _, item := range items {
		switch item["sourceType"] {
		case "local_business":
			if item["name"] == "Fixture Customer Alpha" {
				foundCustomer = true
			}
		case "local_lead":
			if item["name"] == "Fixture Bike Dealer" {
				foundLead = true
			}
		}
	}
	if !foundCustomer || !foundLead {
		t.Fatalf("Chinese channel query should match local customer and lead, got %+v", items)
	}
}

func TestAssistantAgentSourcesMatchOpenClawPolicy(t *testing.T) {
	allowed := map[string]bool{"本地业务数据库": true, "客户私有知识库": true, "Agent 知识库": true, "联网检索": true}
	for _, agentID := range []string{"export-agent", "payment-advisor", "sta100-news-curator", "unknown-agent"} {
		for _, source := range assistantAgentSources(agentID) {
			if !allowed[source] {
				t.Fatalf("agent %s returned unsupported OpenClaw source %q", agentID, source)
			}
		}
	}
	if sources := strings.Join(assistantAgentSources("customer-discovery"), ","); sources != "联网检索" {
		t.Fatalf("customer discovery should use web retrieval source, got %q", sources)
	}
	if sources := strings.Join(assistantAgentKnowledgeSources(), ","); sources != "Agent 知识库" {
		t.Fatalf("agent knowledge source changed unexpectedly: %q", sources)
	}
}

func TestAgentKnowledgeSyncFeedsAssistantEvidence(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	count, err := api.syncAgentKnowledgeFromLocalSources(context.Background())
	if err != nil {
		t.Fatalf("sync agent knowledge: %v", err)
	}
	if count == 0 {
		t.Fatal("expected synced agent knowledge")
	}
	evidence := api.collectLocalEvidence(context.Background(), assistantQueryRequest{
		Feature: "product",
		Message: "Shimano 兼容",
		Context: map[string]any{},
	})
	if len(evidence) == 0 {
		t.Fatal("expected product query to use synced agent knowledge")
	}
	for _, item := range evidence {
		if item.Entity != "agent-knowledge" {
			t.Fatalf("assistant evidence must come from agent knowledge, got %+v", item)
		}
	}
}

func TestAssistantMemoryExtractsReusableSessionContext(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	request := assistantQueryRequest{
		Feature:    "agent-chat",
		Message:    "请记住默认使用德国市场口径",
		Model:      "deepseek/deepseek-chat",
		SessionKey: "sta100-export-agent",
		Context: map[string]any{
			"targetAgent": "export-agent",
			"lang":        "zh",
			"country":     "德国",
		},
	}
	api.extractAndPersistAssistantMemory(context.Background(), request, "MSG-TEST")
	items, err := listRecords[AgentMemory](context.Background(), api.store, agentMemoryKind)
	if err != nil {
		t.Fatalf("list memories: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected extracted memories")
	}
	prompt := api.assistantMemoryPrompt(context.Background(), "export-agent", "sta100-export-agent")
	if !strings.Contains(prompt, "德国") || !strings.Contains(prompt, "language") {
		t.Fatalf("unexpected memory prompt: %s", prompt)
	}
}

func TestAssistantQueueLimitsConcurrentAgentTasks(t *testing.T) {
	queue := newAssistantQueue(2)
	ctx := context.Background()
	first, releaseFirst, err := queue.acquire(ctx, "export-agent", "agent-chat")
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer releaseFirst()
	if first.Queued {
		t.Fatalf("first task should not be queued: %+v", first)
	}
	second, releaseSecond, err := queue.acquire(ctx, "payment-advisor", "agent-chat")
	if err != nil {
		t.Fatalf("second acquire: %v", err)
	}
	defer releaseSecond()
	if second.Queued {
		t.Fatalf("second task should not be queued: %+v", second)
	}
	acquired := make(chan assistantQueueInfo, 1)
	go func() {
		info, releaseThird, err := queue.acquire(ctx, "shipping-eta", "agent-chat")
		if err != nil {
			t.Errorf("third acquire: %v", err)
			return
		}
		defer releaseThird()
		acquired <- info
	}()
	var status assistantQueueInfo
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case info := <-acquired:
			t.Fatalf("third task acquired before capacity was released: %+v", info)
		default:
		}
		status = queue.status()
		if len(status.Running) == 2 && status.QueueLength == 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(status.Running) != 2 || status.QueueLength != 1 {
		t.Fatalf("unexpected queue status: %+v", status)
	}
	releaseFirst()
	select {
	case info := <-acquired:
		if !info.Queued || info.WaitedMs <= 0 {
			t.Fatalf("third task should report queued wait: %+v", info)
		}
		if !strings.Contains(strings.Join(info.RunningNames, ","), "出口业务助手") || !strings.Contains(strings.Join(info.RunningNames, ","), "支付条款助手") {
			t.Fatalf("third task should report the two running blockers: %+v", info.RunningNames)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("third task did not acquire after capacity was released")
	}
	releaseFirst = func() {}
}

func TestQuoteCalculationConversionAndDocumentGeneration(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	quoteInput := Quote{
		Subject: "API 自动计算报价", Customer: "Fixture Customer Alpha", Valid: "2026-12-31", Currency: "EUR",
		Freight: 20, Tax: 5,
		Lines: []BusinessLine{{ProductID: "STA-100-EU", Quantity: 2, UnitPrice: 100, Discount: 10}},
	}
	created := businessRequest(t, api, http.MethodPost, "/api/v1/quotes", quoteInput)
	if created.Code != http.StatusCreated {
		t.Fatalf("quote create = %d: %s", created.Code, created.Body.String())
	}
	quote := decodeResponse[Quote](t, created)
	if quote.Value != "EUR 200.00" || len(quote.Lines) != 1 || quote.Lines[0].Amount != 180 {
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
	order := Order{Customer: "Fixture Customer Alpha", Delivery: "2026-12-31", Currency: "EUR", Lines: []BusinessLine{{ProductID: "EBK-CITY-03", Quantity: 20, UnitPrice: 1000}}}
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

func TestQuoteTemplateUploadIsLockedAndBuiltInDownloadWorks(t *testing.T) {
	t.Setenv("STA100_TEMPLATE_DIR", filepath.Join(t.TempDir(), "templates"))
	t.Setenv("STA100_PDFCPU_CONFIG_DIR", filepath.Join(t.TempDir(), "pdfcpu"))
	api, _ := newTestBusinessAPI(t)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("kind", "quote")
	part, err := writer.CreateFormFile("file", "quote-template.html")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, `<h1>{{quote.subject}}</h1><div>{{customer.name}}</div><section>{{record.lines}}</section><strong>{{record.total}}</strong>`)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/templates/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-STA100-Request", "1")
	request = request.WithContext(context.WithValue(request.Context(), authUsernameContextKey{}, "tester"))
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusConflict || !strings.Contains(recorder.Body.String(), "STANDARD_QUOTE_TEMPLATE_LOCKED") {
		t.Fatalf("quote template upload = %d: %s", recorder.Code, recorder.Body.String())
	}
	download := businessRequest(t, api, http.MethodGet, "/api/v1/quotes/QUO-2026-0188/download", nil)
	if download.Code != http.StatusOK {
		t.Fatalf("quote download = %d: %s", download.Code, download.Body.String())
	}
	if !bytes.HasPrefix(download.Body.Bytes(), []byte("%PDF")) {
		t.Fatalf("download is not a PDF: %q", download.Body.Bytes()[:min(8, download.Body.Len())])
	}
}

func TestOrderTemplatesAreDisabledBecauseDocumentsOwnTemplates(t *testing.T) {
	t.Setenv("STA100_TEMPLATE_DIR", filepath.Join(t.TempDir(), "templates"))
	t.Setenv("STA100_PDFCPU_CONFIG_DIR", filepath.Join(t.TempDir(), "pdfcpu"))
	api, _ := newTestBusinessAPI(t)
	templatePath := filepath.Join("..", "需要的文件", "PI-20260625-DW01.pdf")
	file, err := os.Open(templatePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("kind", "order")
	part, err := writer.CreateFormFile("file", filepath.Base(templatePath))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/templates/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-STA100-Request", "1")
	request = request.WithContext(context.WithValue(request.Context(), authUsernameContextKey{}, "tester"))
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusConflict || !strings.Contains(recorder.Body.String(), "STANDARD_ORDER_TEMPLATE_LOCKED") {
		t.Fatalf("order template upload should be disabled = %d: %s", recorder.Code, recorder.Body.String())
	}
	download := businessRequest(t, api, http.MethodGet, "/api/v1/orders/SO-2026-0105/download", nil)
	if download.Code != http.StatusOK {
		t.Fatalf("order download = %d: %s", download.Code, download.Body.String())
	}
	if contentType := download.Header().Get("Content-Type"); !strings.Contains(contentType, "text/html") {
		t.Fatalf("unexpected content type: %s", contentType)
	}
	if bytes.HasPrefix(download.Body.Bytes(), []byte("%PDF")) {
		t.Fatalf("order download should not use PI PDF template")
	}
}

func TestBuiltInQuotePDFTemplateIsAvailableByDefault(t *testing.T) {
	t.Setenv("STA100_TEMPLATE_DIR", filepath.Join(t.TempDir(), "templates"))
	t.Setenv("STA100_PDFCPU_CONFIG_DIR", filepath.Join(t.TempDir(), "pdfcpu"))
	api, _ := newTestBusinessAPI(t)
	if err := api.ensureBuiltInBusinessTemplates(context.Background()); err != nil {
		t.Fatalf("ensure built-in templates: %v", err)
	}
	listed := businessRequest(t, api, http.MethodGet, "/api/v1/templates?kind=quote", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("list templates = %d: %s", listed.Code, listed.Body.String())
	}
	templates := decodeResponse[struct {
		Items []BusinessTemplate `json:"items"`
	}](t, listed)
	found := false
	for _, item := range templates.Items {
		if item.ID == builtInQuotePDFTemplateID && item.Default && item.OutputFormat == "pdf" && len(item.PDFFields) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("built-in quote PDF template not registered: %+v", templates.Items)
	}
	var quote Quote
	if err := api.store.get(context.Background(), "quotes", "QUO-2026-0188", &quote); err != nil {
		t.Fatal(err)
	}
	tpl, err := api.templateForQuote(context.Background(), quote)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := api.renderBusinessTemplate(context.Background(), tpl, api.quoteTemplateValues(context.Background(), quote, tpl), "报价单-"+quote.ID); err != nil {
		t.Fatalf("render built-in quote PDF: %v", err)
	}
	download := businessRequest(t, api, http.MethodGet, "/api/v1/quotes/QUO-2026-0188/download", nil)
	if download.Code != http.StatusOK {
		t.Fatalf("download built-in pdf = %d: %s", download.Code, download.Body.String())
	}
	if !bytes.HasPrefix(download.Body.Bytes(), []byte("%PDF")) {
		t.Fatalf("download is not a PDF: %q", download.Body.Bytes()[:min(8, download.Body.Len())])
	}
}

func TestQuoteTemplatePreviewLeavesRealtimeFieldsEmpty(t *testing.T) {
	values := quoteTemplatePreviewValues(map[string]string{
		"customer.name":      "Should Not Render",
		"pdf.record.id":      "QUO-SAMPLE",
		"line.1.model":       "STA-100",
		"pdf.record.total":   "999.00",
		"quote.priceTerms":   "FOB Shenzhen",
		"supplier.company":   "STRATRONIX",
		"quote.validityDays": "30",
	})
	for _, key := range quoteTemplateRealtimeFields() {
		if got := strings.TrimSpace(values[key]); got != "" {
			t.Fatalf("preview realtime field %s = %q, want empty", key, got)
		}
	}
	if strings.TrimSpace(values["supplier.company"]) == "" || strings.TrimSpace(values["quote.priceTerms"]) == "" {
		t.Fatalf("template defaults should remain visible in preview: %+v", values)
	}
}

func TestQuoteTemplateValuesKeepLineItemsAndTotalsAligned(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	quote := Quote{
		ID:       "QUO-TEST-0001",
		Customer: "Fixture Customer Alpha",
		Value:    "EUR 34800.00",
		Currency: "USD",
		Freight:  25,
		Tax:      5,
		Lines: []BusinessLine{
			{ProductID: "STA-100-EU", ProductName: "Edge Device", Quantity: 2, UnitPrice: 195, Amount: 390},
			{ProductID: "PM-DUAL-01", ProductName: "Power Meter", Quantity: 3, UnitPrice: 80, Amount: 240},
			{ProductID: "EBK-CITY-03", ProductName: "City Kit", Quantity: 1, UnitPrice: 120, Amount: 120},
		},
	}
	values := api.quoteTemplateValues(context.Background(), quote, BusinessTemplate{DefaultValues: defaultQuoteTemplateValues()})
	checks := map[string]string{
		"line.1.no":            "1",
		"line.1.model":         "STA-100-EU",
		"line.1.description":   "Edge Device",
		"pdf.line.1.unitPrice": "195.00",
		"pdf.line.1.amount":    "390.00",
		"line.2.no":            "2",
		"line.2.model":         "PM-DUAL-01",
		"line.2.description":   "Power Meter",
		"pdf.line.2.unitPrice": "80.00",
		"pdf.line.2.amount":    "240.00",
		"line.3.no":            "3",
		"line.3.model":         "EBK-CITY-03",
		"line.3.description":   "City Kit",
		"pdf.line.3.unitPrice": "120.00",
		"pdf.line.3.amount":    "120.00",
		"pdf.record.subtotal":  "750.00",
		"pdf.record.freight":   "25.00",
		"pdf.record.total":     "775.00",
		"record.total":         "USD 775.00",
	}
	for key, want := range checks {
		if got := values[key]; got != want {
			t.Fatalf("quote template value %s = %q, want %q", key, got, want)
		}
	}
	if strings.Contains(strings.Join([]string{values["line.1.model"], values["line.2.model"], values["line.3.model"]}, " "), "MIXED") {
		t.Fatalf("quote line values should keep individual products, got %+v", values)
	}
}

func TestDynamicQuotePDFPaginatesLineItems(t *testing.T) {
	t.Setenv("STA100_PDFCPU_CONFIG_DIR", filepath.Join(t.TempDir(), "pdfcpu"))
	api, _ := newTestBusinessAPI(t)
	var quote Quote
	if err := api.store.get(context.Background(), "quotes", "QUO-2026-0188", &quote); err != nil {
		t.Fatal(err)
	}
	quote.Currency = "USD"
	quote.Freight = 0
	quote.Tax = 0
	quote.Lines = make([]BusinessLine, 0, 40)
	for i := 1; i <= 40; i++ {
		quote.Lines = append(quote.Lines, BusinessLine{
			ProductID:   fmt.Sprintf("P-%03d", i),
			ProductName: fmt.Sprintf("Product %03d", i),
			Quantity:    float64(i),
			UnitPrice:   float64(i) * 10,
			Amount:      float64(i) * float64(i) * 10,
		})
	}
	values := api.quoteTemplateValues(context.Background(), quote, BusinessTemplate{DefaultValues: defaultQuoteTemplateValues()})
	if got := values["line.40.model"]; got != "P-040" {
		t.Fatalf("line.40.model = %q, want dynamic placeholder value", got)
	}
	content, _, err := renderDynamicQuotePDF(quote, values, "报价单-"+quote.ID)
	if err != nil {
		t.Fatalf("render dynamic quote pdf: %v", err)
	}
	if !bytes.HasPrefix(content, []byte("%PDF")) {
		t.Fatalf("dynamic quote output is not PDF: %q", content[:min(8, len(content))])
	}
	path := filepath.Join(t.TempDir(), "dynamic-quote.pdf")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, err := pdfapi.ReadContextFile(path)
	if err != nil {
		t.Fatalf("read dynamic quote pdf: %v", err)
	}
	if ctx.PageCount < 2 {
		t.Fatalf("dynamic quote PDF should paginate 40 line items, got %d page", ctx.PageCount)
	}
}

func TestBuiltInDocumentTemplatesGenerateStandardFiles(t *testing.T) {
	t.Setenv("STA100_TEMPLATE_DIR", filepath.Join(t.TempDir(), "templates"))
	t.Setenv("STA100_PDFCPU_CONFIG_DIR", filepath.Join(t.TempDir(), "pdfcpu"))
	api, _ := newTestBusinessAPI(t)
	if err := api.ensureBuiltInBusinessTemplates(context.Background()); err != nil {
		t.Fatalf("ensure built-in templates: %v", err)
	}
	listed := businessRequest(t, api, http.MethodGet, "/api/v1/templates?kind=document", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("list document templates = %d: %s", listed.Code, listed.Body.String())
	}
	templates := decodeResponse[struct {
		Items []BusinessTemplate `json:"items"`
	}](t, listed)
	if len(templates.Items) != 5 {
		t.Fatalf("expected 5 document templates, got %d: %+v", len(templates.Items), templates.Items)
	}
	seen := map[string]bool{}
	for _, item := range templates.Items {
		seen[item.ID] = true
	}
	if !seen[builtInPIPDFTemplateID] || !seen[builtInDocumentCustomsTemplateID] || !seen[builtInDocumentContractTemplateID] || !seen[builtInDocumentInvoiceTemplateID] || !seen[builtInDocumentPackingTemplateID] {
		t.Fatalf("built-in document templates missing: %+v", templates.Items)
	}
	if seen["DOCUMENT-TPL-XLSM-20260826"] {
		t.Fatalf("old grouped document template should not be listed: %+v", templates.Items)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("kind", "document")
	part, err := writer.CreateFormFile("file", "document-template.html")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, `<h1>{{document.type}}</h1>`)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/templates/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-STA100-Request", "1")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusConflict || !strings.Contains(recorder.Body.String(), "STANDARD_DOCUMENT_TEMPLATE_LOCKED") {
		t.Fatalf("document template upload = %d: %s", recorder.Code, recorder.Body.String())
	}

	pi := businessRequest(t, api, http.MethodGet, "/api/v1/documents/PI-20260810-003/download", nil)
	if pi.Code != http.StatusOK {
		t.Fatalf("PI download = %d: %s", pi.Code, pi.Body.String())
	}
	if !bytes.HasPrefix(pi.Body.Bytes(), []byte("%PDF")) {
		t.Fatalf("PI download is not a PDF: %q", pi.Body.Bytes()[:min(8, pi.Body.Len())])
	}

	for _, id := range []string{"CI-20260809-012", "PL-20260809-009"} {
		download := businessRequest(t, api, http.MethodGet, "/api/v1/documents/"+id+"/download?format=excel", nil)
		if download.Code != http.StatusOK {
			t.Fatalf("%s download = %d: %s", id, download.Code, download.Body.String())
		}
		if contentType := download.Header().Get("Content-Type"); !strings.Contains(contentType, "spreadsheetml.sheet") {
			t.Fatalf("%s unexpected content type: %s", id, contentType)
		}
		if !strings.Contains(download.Header().Get("Content-Disposition"), ".xlsx") {
			t.Fatalf("%s should download as .xlsx: %s", id, download.Header().Get("Content-Disposition"))
		}
		contentTypes := readZipEntryForTest(t, download.Body.Bytes(), "[Content_Types].xml")
		if bytes.Contains(bytes.ToLower(contentTypes), []byte("macroenabled")) || bytes.Contains(bytes.ToLower(contentTypes), []byte("vbaproject")) {
			t.Fatalf("%s XLSX still contains macro metadata", id)
		}
		sheet := readZipEntryForTest(t, download.Body.Bytes(), "xl/worksheets/sheet1.xml")
		var worksheet struct{ XMLName xml.Name }
		if err := xml.Unmarshal(sheet, &worksheet); err != nil || worksheet.XMLName.Local != "worksheet" {
			t.Fatalf("%s worksheet XML is invalid: %v", id, err)
		}
		hasKnownProduct := bytes.Contains(sheet, []byte("8471504090")) || bytes.Contains(sheet, []byte("9029209000")) || bytes.Contains(sheet, []byte("Fixture Power Meter"))
		if !bytes.Contains(sheet, []byte("Fixture Customer")) || !hasKnownProduct {
			e11 := bytes.Index(sheet, []byte(`r="E11"`))
			from := max(0, e11-120)
			to := min(len(sheet), e11+320)
			t.Fatalf("%s workbook did not contain expected business values near E11=%d: %s", id, e11, sheet[from:to])
		}
	}

	for _, format := range []string{"pdf", "word", "excel"} {
		download := businessRequest(t, api, http.MethodGet, "/api/v1/documents/CI-20260809-012/download?format="+format, nil)
		if download.Code != http.StatusOK {
			t.Fatalf("CI %s download = %d: %s", format, download.Code, download.Body.String())
		}
		switch format {
		case "pdf":
			if !bytes.HasPrefix(download.Body.Bytes(), []byte("%PDF")) || !strings.Contains(download.Header().Get("Content-Type"), "application/pdf") {
				t.Fatalf("CI PDF output is invalid: content-type=%s prefix=%q", download.Header().Get("Content-Type"), download.Body.Bytes()[:min(8, download.Body.Len())])
			}
		case "word":
			if !strings.Contains(download.Header().Get("Content-Type"), "wordprocessingml.document") || !bytes.HasPrefix(download.Body.Bytes(), []byte("PK")) {
				t.Fatalf("CI Word output is invalid: content-type=%s", download.Header().Get("Content-Type"))
			}
			documentXML := readZipEntryForTest(t, download.Body.Bytes(), "word/document.xml")
			if !bytes.Contains(documentXML, []byte("Commercial Invoice")) || !bytes.Contains(documentXML, []byte("Fixture Customer")) {
				t.Fatalf("CI Word document is missing business values: %s", documentXML)
			}
		case "excel":
			if !strings.Contains(download.Header().Get("Content-Type"), "spreadsheetml.sheet") || !bytes.HasPrefix(download.Body.Bytes(), []byte("PK")) {
				t.Fatalf("CI Excel output is invalid: content-type=%s", download.Header().Get("Content-Type"))
			}
		}
	}

	customsPDF := businessRequest(t, api, http.MethodGet, "/api/v1/documents/CD-20260808-006/download?format=pdf", nil)
	if customsPDF.Code != http.StatusOK || !bytes.HasPrefix(customsPDF.Body.Bytes(), []byte("%PDF")) {
		t.Fatalf("customs PDF output = %d, prefix=%q: %s", customsPDF.Code, customsPDF.Body.Bytes()[:min(8, customsPDF.Body.Len())], customsPDF.Body.String())
	}
	customsXML := businessRequest(t, api, http.MethodGet, "/api/v1/documents/CD-20260808-006/download?format=xml", nil)
	if customsXML.Code != http.StatusOK || !strings.Contains(customsXML.Header().Get("Content-Type"), "application/xml") || !bytes.Contains(customsXML.Body.Bytes(), []byte("<CustomsDeclaration>")) || !bytes.Contains(customsXML.Body.Bytes(), []byte("Fixture Customer")) {
		t.Fatalf("customs XML output = %d, content-type=%s body=%q", customsXML.Code, customsXML.Header().Get("Content-Type"), customsXML.Body.String())
	}
	unsupported := businessRequest(t, api, http.MethodGet, "/api/v1/documents/CD-20260808-006/download?format=word", nil)
	if unsupported.Code != http.StatusBadRequest || !strings.Contains(unsupported.Body.String(), "仅支持下载格式") {
		t.Fatalf("customs unsupported format = %d: %s", unsupported.Code, unsupported.Body.String())
	}

	created := businessRequest(t, api, http.MethodPost, "/api/v1/orders/SO-2026-0105/documents", map[string]any{"types": []string{"合同"}, "template": "系统内置标准模板", "language": "中文 / 英文双语"})
	if created.Code != http.StatusCreated {
		t.Fatalf("contract document create = %d: %s", created.Code, created.Body.String())
	}
	response := decodeResponse[struct {
		Items []Document `json:"items"`
	}](t, created)
	if len(response.Items) != 1 || response.Items[0].Type != "合同" || !strings.HasPrefix(response.Items[0].ID, "CT-") {
		t.Fatalf("unexpected contract document: %+v", response.Items)
	}
	contract := businessRequest(t, api, http.MethodGet, "/api/v1/documents/"+response.Items[0].ID+"/download?format=excel", nil)
	if contract.Code != http.StatusOK || !bytes.HasPrefix(contract.Body.Bytes(), []byte("PK")) {
		t.Fatalf("contract download = %d, zip=%v: %s", contract.Code, bytes.HasPrefix(contract.Body.Bytes(), []byte("PK")), contract.Body.String())
	}
}

func TestUpdateWorkbookActiveSheetShowsOnlyRequestedDocument(t *testing.T) {
	source := []byte(`<workbook><bookViews><workbookView activeTab="4"/></bookViews><sheets><sheet name="报关资料录入" sheetId="1" state="hidden" r:id="rId1"/><sheet name="报关单" sheetId="2" r:id="rId2"/><sheet name="合同" sheetId="7" r:id="rId3"/><sheet name="发票" sheetId="6" r:id="rId4"/><sheet name="装箱单" sheetId="5" r:id="rId5"/></sheets><calcPr calcId="1"/></workbook>`)
	updated := string(updateWorkbookActiveSheet(source, "CI"))
	if !strings.Contains(updated, `activeTab="3"`) {
		t.Fatalf("invoice active tab was not selected: %s", updated)
	}
	if !strings.Contains(updated, `<sheet name="发票" sheetId="6" r:id="rId4"/>`) {
		t.Fatalf("invoice sheet should stay visible: %s", updated)
	}
	for _, name := range []string{"报关资料录入", "报关单", "合同", "装箱单"} {
		if !strings.Contains(updated, `name="`+name+`"`) || !strings.Contains(updated[strings.Index(updated, `name="`+name+`"`):], `state="hidden"`) {
			t.Fatalf("sheet %s should be hidden: %s", name, updated)
		}
	}
}

func TestKeepWorkbookDocumentSheetRetainsOnlyRequestedSheet(t *testing.T) {
	workbook := []byte(`<?xml version="1.0"?><workbook><bookViews><workbookView activeTab="4"/></bookViews><sheets><sheet name="报关资料录入" sheetId="1" state="hidden" r:id="rId1"/><sheet name="报关单" sheetId="2" r:id="rId2"/><sheet name="合同" sheetId="7" r:id="rId3"/><sheet name="发票" sheetId="6" r:id="rId4"/><sheet name="装箱单" sheetId="5" r:id="rId5"/></sheets><definedNames><definedName name="x">报关单!A1</definedName></definedNames></workbook>`)
	rels := []byte(`<?xml version="1.0"?><Relationships><Relationship Id="rId1" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Target="worksheets/sheet2.xml"/><Relationship Id="rId3" Target="worksheets/sheet3.xml"/><Relationship Id="rId4" Target="worksheets/sheet4.xml"/><Relationship Id="rId5" Target="worksheets/sheet5.xml"/></Relationships>`)
	contentTypes := []byte(`<?xml version="1.0"?><Types><Override PartName="/xl/worksheets/sheet1.xml"/><Override PartName="/xl/worksheets/sheet2.xml"/><Override PartName="/xl/worksheets/sheet3.xml"/><Override PartName="/xl/worksheets/sheet4.xml"/><Override PartName="/xl/worksheets/sheet5.xml"/></Types>`)
	worksheet := []byte(`<worksheet><sheetData><row><c><f>SUM(A1:A2)</f><v>3</v></c></row></sheetData></worksheet>`)
	var source bytes.Buffer
	writer := zip.NewWriter(&source)
	for entry, content := range map[string][]byte{
		"xl/workbook.xml": workbook, "xl/_rels/workbook.xml.rels": rels, "[Content_Types].xml": contentTypes,
		"xl/worksheets/sheet1.xml": worksheet, "xl/worksheets/sheet2.xml": worksheet, "xl/worksheets/sheet3.xml": worksheet, "xl/worksheets/sheet4.xml": worksheet, "xl/worksheets/sheet5.xml": worksheet,
		"xl/calcChain.xml": []byte("<calcChain/>"), "xl/externalLinks/externalLink1.xml": []byte("<externalLink/>"),
	} {
		file, err := writer.Create(entry)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := file.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := keepWorkbookDocumentSheet(source.Bytes(), "CI")
	if err != nil {
		t.Fatal(err)
	}
	if got := string(readZipEntryForTest(t, output, "xl/workbook.xml")); !strings.Contains(got, `name="发票"`) || strings.Contains(got, `name="合同"`) || !strings.Contains(got, `发票!$A$1:$G$27`) {
		t.Fatalf("unexpected workbook manifest: %s", got)
	}
	if got := string(readZipEntryForTest(t, output, "xl/worksheets/sheet4.xml")); strings.Contains(got, "<f>") || !strings.Contains(got, "<v>3</v>") {
		t.Fatalf("target worksheet formulas were not flattened: %s", got)
	}
	archive, err := zip.NewReader(bytes.NewReader(output), int64(len(output)))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range archive.File {
		if strings.Contains(file.Name, "sheet1.xml") || strings.Contains(file.Name, "sheet2.xml") || strings.Contains(file.Name, "sheet3.xml") || strings.Contains(file.Name, "sheet5.xml") || strings.Contains(file.Name, "externalLinks") || file.Name == "xl/calcChain.xml" {
			t.Fatalf("unexpected workbook part retained: %s", file.Name)
		}
	}
}

func readZipEntryForTest(t *testing.T, content []byte, name string) []byte {
	t.Helper()
	archive, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		t.Fatalf("cannot open zip: %v", err)
	}
	for _, file := range archive.File {
		if file.Name != name {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	t.Fatalf("zip entry %s not found", name)
	return nil
}

func TestRawDataDependentEndpointsAreExplicitTODO(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	for _, test := range []struct {
		path string
		body any
	}{
		{"/api/v1/overview/oem/match", map[string]any{"query": "battery"}},
		{"/api/v1/system/upgrade/import", map[string]any{}},
	} {
		recorder := businessRequest(t, api, http.MethodPost, test.path, test.body)
		if recorder.Code != http.StatusNotImplemented {
			t.Errorf("%s status = %d, want 501: %s", test.path, recorder.Code, recorder.Body.String())
		}
	}
}

func TestOverviewCustomerDiscoveryRequiresUsableDefaultModel(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	api.openClaw = &openClawService{service: orchestrator.New(orchestrator.Config{})}
	response := businessRequest(t, api, http.MethodPost, "/api/v1/overview/customer-discovery", map[string]any{
		"country": "德国", "city": "柏林", "type": "Distributor", "limit": 3,
	})
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "MODEL_NOT_READY") {
		t.Fatalf("overview customer discovery without model = %d: %s", response.Code, response.Body.String())
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
		Subject: "改名关联测试", Customer: "Fixture Customer Alpha", Valid: "2026-12-31", Currency: "EUR",
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
		if item.Name == "Fixture Customer Alpha" {
			account = item
			break
		}
	}
	if account.ID == "" {
		t.Fatal("seed customer not found")
	}
	account.Name = "Fixture Customer Alpha Renamed"
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
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "openclaw")
	script := `#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
  "cron add") printf '%s' '{"job":{"id":"oc-test","name":"每日推荐更新","description":"test","enabled":false,"agentId":"sta100-recommend-curator","sessionTarget":"isolated","schedule":{"kind":"every","everyMs":7200000},"payload":{"kind":"agentTurn","message":"test"},"state":{},"status":"disabled"}}' ;;
  *) exit 2 ;;
esac
`
	if err := os.WriteFile(scriptPath, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	api.openClaw = &openClawService{service: orchestrator.New(orchestrator.Config{BinaryPath: scriptPath})}
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

func TestOverviewAutomationMessageSeparatesFailedAndReviewCounts(t *testing.T) {
	message := overviewAutomationMessage("failed", 2, 1, "")
	if !strings.Contains(message, "1 个自动任务执行失败") || !strings.Contains(message, "1 个任务需要人工复核") {
		t.Fatalf("message did not separate failed and review counts: %q", message)
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
	if response.Item.Status != "Indexed" || response.Item.Path != "" {
		t.Fatalf("unexpected uploaded metadata: %+v", response.Item)
	}
	content := businessRequest(t, api, http.MethodGet, "/api/v1/private-files/"+response.Item.ID+"/content", nil)
	if content.Code != http.StatusOK || string(content.Body.Bytes()) != "customer data pending schema" {
		t.Fatalf("content = %d %q", content.Code, content.Body.String())
	}
	summary := businessRequest(t, api, http.MethodGet, "/api/v1/private-files/"+response.Item.ID+"/summary", nil)
	if summary.Code != http.StatusOK || !strings.Contains(summary.Body.String(), "已解析") {
		t.Fatalf("summary = %d: %s", summary.Code, summary.Body.String())
	}
}

func TestKnowledgeSyncContinuesAfterUnparseableSource(t *testing.T) {
	shared := t.TempDir()
	private := t.TempDir()
	t.Setenv("STA100_SHARED_KNOWLEDGE_DIR", shared)
	t.Setenv("STA100_KNOWLEDGE_DATA_DIR", filepath.Join(t.TempDir(), "knowledge"))
	t.Setenv("STA100_PRIVATE_DATA_DIR", private)
	if err := os.WriteFile(filepath.Join(shared, "one.md"), []byte("德国经销商和电动自行车渠道资料"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shared, "broken.docx"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shared, "two.md"), []byte("法国法规与产品合规资料"), 0600); err != nil {
		t.Fatal(err)
	}
	api, _ := newTestBusinessAPI(t)
	count, syncErr := api.syncKnowledgeSources(context.Background())
	if count != 3 || syncErr != nil {
		t.Fatalf("expected all files scanned, count=%d err=%v", count, syncErr)
	}
	status := api.knowledgeIndexStatus(context.Background())
	if status["documents"] != 2 || status["pending"] != 1 {
		t.Fatalf("unexpected knowledge status: %+v", status)
	}
	hits, err := api.searchKnowledge(context.Background(), []string{"market-analyzer"}, "德国经销商", 8)
	if err != nil || len(hits) == 0 {
		t.Fatalf("expected searchable indexed source, hits=%d err=%v", len(hits), err)
	}
}

func TestPrivateKnowledgeDeleteRemovesIndex(t *testing.T) {
	t.Setenv("STA100_PRIVATE_DATA_DIR", filepath.Join(t.TempDir(), "private-files"))
	api, _ := newTestBusinessAPI(t)
	path := filepath.Join(os.Getenv("STA100_PRIVATE_DATA_DIR"), "FILE-0009.md")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("专用产品资料"), 0600); err != nil {
		t.Fatal(err)
	}
	item := PrivateFile{ID: "FILE-0009", Name: "source.md", Category: "产品资料", Path: path, Bytes: 18, Updated: currentText()}
	if err := api.indexPrivateKnowledgeFile(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if err := api.removePrivateKnowledgeDocument(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	var documents int
	if err := api.store.db.QueryRow(`SELECT COUNT(*) FROM knowledge_documents`).Scan(&documents); err != nil {
		t.Fatal(err)
	}
	if documents != 0 {
		t.Fatalf("expected index cleanup, documents=%d", documents)
	}
}

func TestDesignTODOAliasesDoNotFallThroughTo404(t *testing.T) {
	api, _ := newTestBusinessAPI(t)
	for _, path := range []string{
		"/api/v1/overview/oem-matches",
		"/api/v1/overview/oem-matches/export",
		"/api/v1/overview/oem-matches/FACTORY-1",
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

func TestModelConnectionTestUsesRealProviderProbe(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("OPENCLAW_STATE_DIR", filepath.Join(dir, "state"))
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			http.Error(w, "unexpected path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer valid-key" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"message":"invalid api key","type":"authentication_error"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"pong"}}]}`))
	}))
	defer provider.Close()
	configPath := filepath.Join(dir, "openclaw.json")
	config := `{"agents":{"defaults":{"model":{"primary":"deepseek/deepseek-chat"},"models":{"deepseek/deepseek-chat":{}}},"list":[{"id":"sta100-coordinator"}]},"models":{"providers":{"deepseek":{"baseUrl":%q,"models":[{"id":"deepseek-chat","name":"DeepSeek Chat"}]}}},"auth":{"profiles":{}}}`
	if err := os.WriteFile(configPath, []byte(fmt.Sprintf(config, provider.URL)), 0o600); err != nil {
		t.Fatal(err)
	}
	scriptPath := filepath.Join(dir, "openclaw")
	script := `#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
  "gateway status") printf '%s' '{"cli":{"version":"test"},"service":{"loaded":true,"runtime":{"status":"running"},"configAudit":{"ok":true}},"gateway":{"version":"test"},"rpc":{"ok":true,"version":"test"}}' ;;
  "models status") printf '%s' '{"defaultModel":"demo/working","resolvedDefault":"demo/working","auth":{"missingProvidersInUse":[],"providers":[{"provider":"demo","profiles":{"count":1,"apiKey":1,"oauth":0,"token":0}}]}}' ;;
  "models list") printf '%s' '{"models":[{"key":"demo/working","name":"Working","available":true,"missing":false}]}' ;;
  *) exit 2 ;;
esac
`
	if err := os.WriteFile(scriptPath, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	api, _ := newTestBusinessAPI(t)
	api.openClaw = &openClawService{service: orchestrator.New(orchestrator.Config{BinaryPath: scriptPath, ConfigPath: configPath})}
	recorder := businessRequest(t, api, http.MethodPost, "/api/v1/settings/model/test", map[string]any{"model": "deepseek/deepseek-chat", "apiKey": "valid-key"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", recorder.Code, recorder.Body.String())
	}
	result := decodeResponse[struct {
		OK              bool   `json:"ok"`
		GatewayOK       bool   `json:"gatewayOK"`
		ConfigurationOK bool   `json:"configurationOK"`
		GenerationOK    bool   `json:"generationOK"`
		Model           string `json:"model"`
	}](t, recorder)
	if !result.OK || !result.ConfigurationOK || !result.GenerationOK || result.Model != "deepseek/deepseek-chat" {
		t.Fatalf("unexpected model test response: %+v\n%s", result, recorder.Body.String())
	}
	recorder = businessRequest(t, api, http.MethodPost, "/api/v1/settings/model/test", map[string]any{"model": "deepseek/deepseek-chat", "apiKey": "invalid-key"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", recorder.Code, recorder.Body.String())
	}
	failed := decodeResponse[struct {
		OK           bool   `json:"ok"`
		GenerationOK bool   `json:"generationOK"`
		Message      string `json:"message"`
	}](t, recorder)
	if failed.OK || failed.GenerationOK || !strings.Contains(failed.Message, "API Key") {
		t.Fatalf("invalid key should fail provider probe: %+v\n%s", failed, recorder.Body.String())
	}
}

func TestModelConfigurationSupportsSelectedVersionsAndSingleDefault(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"defaults":{"model":{"primary":"demo/working"},"models":{"demo/working":{},"demo/fast":{}}}},"models":{"providers":{"demo":{"models":[{"id":"working"},{"id":"fast"}]}}},"auth":{"profiles":{"demo:default":{"provider":"demo","mode":"api_key"}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	service := orchestrator.New(orchestrator.Config{ConfigPath: configPath})
	if err := service.SetConfiguredModelSelection("demo/fast", []string{"demo/working", "demo/fast"}); err != nil {
		t.Fatalf("set selected models: %v", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Agents struct {
			Defaults struct {
				Model  map[string]string `json:"model"`
				Models map[string]any    `json:"models"`
			} `json:"defaults"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if config.Agents.Defaults.Model["primary"] != "demo/fast" {
		t.Fatalf("default model = %q", config.Agents.Defaults.Model["primary"])
	}
	if len(config.Agents.Defaults.Models) != 2 {
		t.Fatalf("selected model count = %d", len(config.Agents.Defaults.Models))
	}
}

func TestSavingFailedModelKeepsPreviouslyTestedDefault(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	config := `{"agents":{"defaults":{"model":{"primary":"demo/working"},"models":{"demo/working":{}},"list":[{"id":"main","agentDir":"%s"}]}},"models":{"providers":{"demo":{"models":[{"id":"working"},{"id":"fast"}]}}},"auth":{"profiles":{}}}`
	if err := os.WriteFile(configPath, []byte(fmt.Sprintf(config, filepath.ToSlash(filepath.Join(dir, "main")))), 0o600); err != nil {
		t.Fatal(err)
	}
	api, _ := newTestBusinessAPI(t)
	api.openClaw = &openClawService{service: orchestrator.New(orchestrator.Config{ConfigPath: configPath})}
	api.saveModelTestState(context.Background(), "demo/working", true, modelRealProbeSuccessMessage)
	api.saveModelTestState(context.Background(), "demo/fast", false, "模型提供商鉴权失败")

	recorder := businessRequest(t, api, http.MethodPatch, "/api/v1/settings/model", map[string]any{
		"model": "demo/fast", "selectedModels": []string{"demo/working", "demo/fast"},
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("saving a failed non-default model returned %d: %s", recorder.Code, recorder.Body.String())
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var saved struct {
		Agents struct {
			Defaults struct {
				Model  map[string]string `json:"model"`
				Models map[string]any    `json:"models"`
			} `json:"defaults"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatal(err)
	}
	if got := saved.Agents.Defaults.Model["primary"]; got != "demo/working" {
		t.Fatalf("default changed to an untested model: %q", got)
	}
	if _, ok := saved.Agents.Defaults.Models["demo/fast"]; !ok {
		t.Fatalf("failed model was not saved: %#v", saved.Agents.Defaults.Models)
	}
}

func TestModelTestFailureMessageClassifiesProviderErrors(t *testing.T) {
	if got := modelTestFailureMessage(errors.New("provider: insufficient balance"), "fallback"); !strings.Contains(got, "余额") {
		t.Fatalf("balance error = %q", got)
	}
	if got := modelTestFailureMessage(errors.New("401 invalid api key"), "fallback"); !strings.Contains(got, "API Key") {
		t.Fatalf("authentication error = %q", got)
	}
	if got := modelTestFailureMessage(errors.New("The selected model was not found by the provider"), "fallback"); !strings.Contains(got, "无权使用") {
		t.Fatalf("model missing error = %q", got)
	}
}

func TestModelTestStateRequiresRealGeneration(t *testing.T) {
	legacy := modelTestState{OK: true, Message: "API Key 写入和模型配置检查通过"}
	if modelTestStatePassed(legacy) {
		t.Fatal("legacy configuration-only model test must not enable agent chat")
	}
	if got := modelTestStateMessage(legacy); !strings.Contains(got, "真实模型调用验证") {
		t.Fatalf("legacy message = %q", got)
	}
	current := modelTestState{OK: true, Message: modelRealProbeSuccessMessage}
	if !modelTestStatePassed(current) {
		t.Fatal("real provider probe success should enable agent chat")
	}
}
