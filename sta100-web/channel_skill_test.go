package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func jsonBody(t *testing.T, value any) *bytes.Reader {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(data)
}

func TestChannelSkillMatchAndMissing(t *testing.T) {
	candidates := matchChannelSkills("请生成报价单")
	if len(candidates) != 1 || candidates[0] != "quote.create" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	skill, ok := findChannelSkill("quote.create")
	if !ok {
		t.Fatal("quote skill not found")
	}
	missing := channelSkillMissing(skill, map[string]string{"customer": "Acme"})
	if len(missing) != 1 || missing[0] != "products" {
		t.Fatalf("unexpected missing parameters: %#v", missing)
	}
}

func TestChannelSkillRouteRequiresExplicitAllowlist(t *testing.T) {
	route := ChannelSkillRoute{Channel: "feishu", Account: "bot", AgentID: "export-agent", Enabled: true}
	if err := validateChannelSkillRoute(&route); err == nil {
		t.Fatal("route without sender/chat allowlist was accepted")
	}
	route.AllowedSenders = []string{"user-1"}
	if err := validateChannelSkillRoute(&route); err != nil {
		t.Fatal(err)
	}
	if !channelSkillRouteAllows(route, "user-1", "chat-1") || channelSkillRouteAllows(route, "user-2", "chat-1") {
		t.Fatal("route allowlist behavior is incorrect")
	}
}

func TestChannelSkillInboundRejectsWithoutRoute(t *testing.T) {
	t.Setenv("STA100_CHANNEL_SKILL_TOKEN", "test-channel-token")
	t.Setenv("STA100_DB_PATH", filepath.Join(t.TempDir(), "sta100.db"))
	store, err := newBusinessStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	a := newBusinessAPI(store, &openClawService{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/openclaw/inbound/feishu", jsonBody(t, map[string]any{
		"account": "bot", "sender": "user-1", "conversation": "chat-1", "message": "查询库存",
	}))
	request.RemoteAddr = "127.0.0.1:12345"
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-STA100-Channel-Token", "test-channel-token")
	recorder := httptest.NewRecorder()
	a.channelSkillInboundHandler(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", recorder.Code)
	}
}

func TestChannelSkillSessionPersists(t *testing.T) {
	t.Setenv("STA100_DB_PATH", filepath.Join(t.TempDir(), "sta100.db"))
	store, err := newBusinessStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	a := newBusinessAPI(store, &openClawService{})
	session := ChannelSkillSession{ID: "css-test", Channel: "feishu", Account: "bot", Sender: "user", Conversation: "chat", State: "collecting", Parameters: map[string]string{"query": "库存"}}
	if err := a.saveChannelSkillSession(context.Background(), session); err != nil {
		t.Fatal(err)
	}
	loaded, err := a.loadChannelSkillSession(context.Background(), "css-test")
	if err != nil || loaded.Parameters["query"] != "库存" {
		t.Fatalf("session was not persisted: %#v, %v", loaded, err)
	}
}
