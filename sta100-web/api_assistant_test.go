package main

import (
	"strings"
	"testing"
)

func TestParseAssistantStructuredCustomerDiscoveryItems(t *testing.T) {
	items := parseAssistantStructuredItems("customer_discovery", `[STA100_RESULT]{"schema":"sta100.business.v1","type":"customer_discovery","items":[{"name":"Berlin Bike Hub","country":"德国","city":"柏林","type":"Distributor","business":"E-bike 经销","contact":"info@example.test","source":"公开网站","sourceUrl":"https://example.test","updatedAt":"2026-08-19","score":92,"reason":"符合国家、城市和类型筛选"}]}[/STA100_RESULT]`)
	if len(items) != 1 {
		t.Fatalf("expected one customer discovery item, got %d", len(items))
	}
	if items[0]["name"] != "Berlin Bike Hub" || items[0]["score"] != 92 {
		t.Fatalf("unexpected normalized customer item: %+v", items[0])
	}
}

func TestParseAssistantStructuredCustomerDiscoveryFractionalScore(t *testing.T) {
	items := parseAssistantStructuredItems("customer_discovery", `[STA100_RESULT]{"schema":"sta100.business.v1","type":"customer_discovery","items":[{"name":"Fahrradstation GmbH","country":"Germany","city":"Berlin","type":"Distributor","score":0.82,"source":"公开网站"}]}[/STA100_RESULT]`)
	if len(items) != 1 {
		t.Fatalf("expected one customer discovery item, got %d", len(items))
	}
	if items[0]["score"] != 82 {
		t.Fatalf("expected fractional score to normalize to percentage, got %+v", items[0]["score"])
	}
}

func TestParseAssistantStructuredOEMItems(t *testing.T) {
	items := parseAssistantStructuredItems("oem-match", `[STA100_RESULT]{"schema":"sta100.business.v1","type":"oem_match","items":[{"title":"华东电池 OEM","category":"E-bike 电池","reason":"产品能力匹配","detail":"具备欧洲 E-bike 电池 Pack 经验","source":"供应商库","score":"88","capacity":"5000 组/月","moq":"100 组"}]}[/STA100_RESULT]`)
	if len(items) != 1 {
		t.Fatalf("expected one oem item, got %d", len(items))
	}
	if items[0]["title"] != "华东电池 OEM" || items[0]["score"] != 88 {
		t.Fatalf("unexpected normalized oem item: %+v", items[0])
	}
}

func TestParseAssistantStructuredCustomerSearchItems(t *testing.T) {
	items := parseAssistantStructuredItems("customer_search", `[STA100_RESULT]{"schema":"sta100.business.v1","type":"customer_search","items":[{"name":"Berlin Bike GmbH","country":"德国","city":"柏林","type":"Distributor","source":"公开网站","sourceUrl":"https://example.test/customer","updatedAt":"2026-08-31","score":86}]}[/STA100_RESULT]`)
	if len(items) != 1 || items[0]["source"] != "公开网站" || items[0]["sourceUrl"] != "https://example.test/customer" {
		t.Fatalf("customer search source fields were not retained: %+v", items)
	}
}

func TestParseAssistantStructuredResultAcceptsEmptyItems(t *testing.T) {
	result, ok := parseAssistantStructuredResult(`[STA100_RESULT]{"schema":"sta100.business.v1","type":"customer_discovery","items":[]}[/STA100_RESULT]`)
	if !ok {
		t.Fatal("empty structured result block should still be recognized")
	}
	if result.Type != "customer_discovery" || len(result.Items) != 0 {
		t.Fatalf("unexpected empty structured result: %+v", result)
	}
}

func TestCustomerDiscoveryAccessIssue(t *testing.T) {
	text := "来源访问情况：web_search：本轮再次禁用（disabled or no provider is available），web_fetch：请求被拦截（resolves to private/internal/special-use IP address）。本 Agent 无自有客户库。"
	issue, message := customerDiscoveryAccessIssue(text, 0)
	if !issue {
		t.Fatal("expected unavailable public-source capability to be treated as an access issue")
	}
	if !strings.Contains(message, "公开来源") {
		t.Fatalf("unexpected issue message: %q", message)
	}
}

func TestCustomerDiscoveryAccessIssueIgnoresToolWarningsWhenItemsExist(t *testing.T) {
	issue, _ := customerDiscoveryAccessIssue("web_search disabled but items already parsed", 1)
	if issue {
		t.Fatal("structured customer items should win over tool warnings")
	}
}
