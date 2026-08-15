package orchestrator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseModelReference(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "string", raw: `"provider/model"`, want: "provider/model"},
		{name: "object", raw: `{"primary":"provider/default"}`, want: "provider/default"},
		{name: "null", raw: `null`, want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseModelReference(json.RawMessage(test.raw)); got != test.want {
				t.Fatalf("parseModelReference() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestListConfiguredAgentsDoesNotExposeSecrets(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "openclaw.json")
	config := `{
		"gateway":{"auth":{"token":"must-not-leak"}},
		"agents":{"defaults":{"model":{"primary":"deepseek/default"}},"list":[
			{"id":"main"},
			{"id":"export-agent","name":"export-agent","identity":{"name":"出口业务助手","emoji":"🛒"}}
		]}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	service := New(Config{ConfigPath: configPath})
	agents, err := service.ListAgents(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 2 || agents[1].IdentityEmoji != "🛒" || agents[1].Model != "deepseek/default" {
		t.Fatalf("unexpected agents: %#v", agents)
	}
	encoded, err := json.Marshal(agents)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "must-not-leak") {
		t.Fatal("sensitive configuration leaked into agent response")
	}
}

func TestListAgentsDecoratesSystemRoleFromManifest(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	manifestPath := filepath.Join(dir, "agents.json")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"main"},{"id":"sta100-coordinator"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, []byte(`{"schema_version":1,"agents":[{"id":"sta100-coordinator","name":"Coordinator","workspace":"workspace","role":"coordinator","visibility":"system"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	agents, err := New(Config{ConfigPath: configPath, Manifest: manifestPath}).ListAgents(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 2 || agents[1].Role != "coordinator" || agents[1].Visibility != "system" {
		t.Fatalf("system metadata not applied: %+v", agents)
	}
}

func TestSetDefaultModelChecksAvailability(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "commands.log")
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
	"models status") printf '%s' '{"defaultModel":"deepseek/deepseek-chat","resolvedDefault":"deepseek/deepseek-chat","auth":{"providers":[{"provider":"deepseek","profiles":{"count":1,"apiKey":1,"oauth":0,"token":0}},{"provider":"custom-old","profiles":{"count":1,"apiKey":1,"oauth":0,"token":0}},{"provider":"github-copilot","profiles":{"count":1,"apiKey":0,"oauth":1,"token":0}}]}}' ;;
	"models set") printf '%s\n' "$*" >> "`+logPath+`" ;;
  *) exit 2 ;;
esac
`)
	service := New(Config{BinaryPath: script})
	if err := service.SetDefaultModel(context.Background(), "deepseek/deepseek-reasoner"); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(contents)) != "models set deepseek/deepseek-reasoner" {
		t.Fatalf("unexpected command: %q", contents)
	}
	if err := service.SetDefaultModel(context.Background(), "minimax/MiniMax-M3"); err != ErrModelUnavailable {
		t.Fatalf("error = %v, want ErrModelUnavailable", err)
	}
}

func TestModelsReturnsFullCatalogProviders(t *testing.T) {
	dir := t.TempDir()
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
	"models status") printf '%s' '{"defaultModel":"deepseek/deepseek-chat","resolvedDefault":"deepseek/deepseek-chat","auth":{"providers":[{"provider":"deepseek","profiles":{"count":1,"apiKey":1,"oauth":0,"token":0}}]}}' ;;
	*) exit 2 ;;
esac
`)
	models, err := New(Config{BinaryPath: script}).Models(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(fixedModelCatalog) != 91 || len(models.Models) != len(fixedModelCatalog) || len(models.Providers) != 18 {
		t.Fatalf("catalog = %d models, %d providers: %+v", len(models.Models), len(models.Providers), models)
	}
	if models.CatalogVersion != ModelCatalogVersion {
		t.Fatalf("catalog version = %q", models.CatalogVersion)
	}
	var minimaxModels int
	for _, model := range models.Models {
		if strings.HasPrefix(model.Key, "minimax/") {
			minimaxModels++
			if model.Available || !model.Missing {
				t.Fatalf("unconfigured MiniMax model has wrong state: %+v", model)
			}
		}
	}
	if minimaxModels != 3 {
		t.Fatalf("MiniMax model count = %d", minimaxModels)
	}
}

func TestConfiguredModelSelectionDoesNotRequireCredentialAndRegistersMiniMaxModel(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"models":{"providers":{}},"agents":{"defaults":{}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	service := New(Config{ConfigPath: configPath})
	if err := service.SetConfiguredModelSelection("minimax/MiniMax-M3", []string{"minimax/MiniMax-M3"}); err != nil {
		t.Fatalf("selection should be saved before credential testing: %v", err)
	}
	var config struct {
		Models struct {
			Providers map[string]struct {
				Models []struct {
					ID  string `json:"id"`
					API string `json:"api"`
				} `json:"models"`
			} `json:"providers"`
		} `json:"models"`
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		t.Fatal(err)
	}
	entries := config.Models.Providers["minimax"].Models
	if len(entries) != 1 || entries[0].ID != "MiniMax-M3" || entries[0].API != "anthropic-messages" {
		t.Fatalf("MiniMax model was not registered in OpenClaw config: %+v", entries)
	}
}

func TestChannelsReturnPinnedCatalogWithoutCLI(t *testing.T) {
	channels, err := New(Config{}).Channels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var feishu Channel
	for _, channel := range channels {
		if channel.ID == "feishu" {
			feishu = channel
			break
		}
	}
	if !feishu.Installed || feishu.Origin != "available" || feishu.AccountCount != 0 || feishu.Status != "installed" {
		t.Fatalf("pinned Feishu channel state is wrong: %+v", feishu)
	}
}

func TestApplyFeishuRegistrationWritesDefaultAccountAndRestartsGateway(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":{"entries":{"feishu":{"enabled":true}}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
if [[ "$1 $2" != "gateway restart" ]]; then
  exit 2
fi
`)
	service := New(Config{BinaryPath: script, ConfigPath: configPath})
	if err := service.applyFeishuRegistration(context.Background(), "default", "feishu", "cli_xxx", "secret_xxx", "ou_owner"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]any
	if err := json.Unmarshal(raw, &config); err != nil {
		t.Fatal(err)
	}
	feishu := config["channels"].(map[string]any)["feishu"].(map[string]any)
	if feishu["appId"] != "cli_xxx" || feishu["appSecret"] != "secret_xxx" || feishu["domain"] != "feishu" {
		t.Fatalf("Feishu credentials were not persisted: %#v", feishu)
	}
	if feishu["dmPolicy"] != "allowlist" {
		t.Fatalf("Feishu DM policy was not restricted after QR registration: %#v", feishu)
	}
}

func TestSupportsAPIKeyProviderExcludesOAuthAndCLIOnlyRoutes(t *testing.T) {
	if supportsAPIKeyProvider("github-copilot") || supportsAPIKeyProvider("claude-cli") || supportsAPIKeyProvider("minimax-portal") || supportsAPIKeyProvider("lmstudio") {
		t.Fatal("OAuth/CLI/local-only providers must not appear in API Key configuration")
	}
	if !supportsAPIKeyProvider("deepseek") || !supportsAPIKeyProvider("anthropic") || !supportsAPIKeyProvider("minimax") {
		t.Fatal("API Key providers must remain available")
	}
}

func TestSaveAPIKeyWritesOpenClawStoreWithoutInteractiveCLI(t *testing.T) {
	dir := t.TempDir()
	argsPath := filepath.Join(dir, "args.log")
	configPath := filepath.Join(dir, "openclaw.json")
	agentDir := filepath.Join(dir, "agents", "main", "agent")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"main","agentDir":"`+filepath.ToSlash(agentDir)+`"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" > "`+argsPath+`"
`)
	service := New(Config{BinaryPath: script, ConfigPath: configPath})
	secret := "secret-key-value"
	profile, err := service.SaveAPIKey(context.Background(), "minimax", secret)
	if err != nil {
		t.Fatal(err)
	}
	if profile != "minimax:default" {
		t.Fatalf("profile = %q", profile)
	}
	args, _ := os.ReadFile(argsPath)
	if strings.Contains(string(args), secret) {
		t.Fatal("secret appeared in command arguments")
	}
	if strings.Contains(string(args), "paste-api-key") {
		t.Fatalf("SaveAPIKey used interactive OpenClaw auth command: %s", args)
	}
	if info, err := os.Stat(filepath.Join(agentDir, "openclaw-agent.sqlite")); err != nil || info.Size() == 0 {
		t.Fatalf("auth database was not created: info=%v err=%v", info, err)
	}
	models := service.ModelSnapshot()
	var minimaxConfigured bool
	for _, provider := range models.Providers {
		if provider.Provider == "minimax" {
			minimaxConfigured = provider.APIKeyConfigured
			break
		}
	}
	if !minimaxConfigured {
		t.Fatalf("minimax provider was not marked configured: %+v", models.Providers)
	}
}

func TestSendAgentMessageUsesPrivateFileAndParsesReply(t *testing.T) {
	dir := t.TempDir()
	argsPath := filepath.Join(dir, "args.log")
	messagePath := filepath.Join(dir, "message.log")
	modePath := filepath.Join(dir, "mode.log")
	configPath := filepath.Join(dir, "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"export-agent"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" > "`+argsPath+`"
while [[ $# -gt 0 ]]; do
  if [[ "$1" == "--message-file" ]]; then
    shift
    cat "$1" > "`+messagePath+`"
    stat -c '%a' "$1" > "`+modePath+`"
    break
  fi
  shift
done
printf '%s' '{"runId":"run-1","status":"ok","result":{"payloads":[{"text":"真实回复"}],"meta":{"agentMeta":{"usage":{"input":120,"output":30,"cacheRead":10,"cacheWrite":5,"total":165}}}}}'
`)
	service := New(Config{BinaryPath: script, ConfigPath: configPath})
	secretMessage := "检查欧洲报价"
	result, err := service.SendAgentMessage(context.Background(), AgentMessageInput{
		AgentID: "export-agent", Message: secretMessage, SessionKey: "sta100-export-agent",
		Sources: []string{"本地业务数据库", "联网检索"}, Allowlist: []string{"eur-lex.europa.eu"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Text != "真实回复" || result.RunID != "run-1" || result.Usage.Total != 165 || result.Usage.Input != 120 || result.Usage.Output != 30 {
		t.Fatalf("unexpected result: %#v", result)
	}
	args, _ := os.ReadFile(argsPath)
	if strings.Contains(string(args), secretMessage) || strings.Contains(string(args), "[用户消息]") {
		t.Fatal("message content appeared in command arguments")
	}
	contents, _ := os.ReadFile(messagePath)
	if !strings.Contains(string(contents), secretMessage) || !strings.Contains(string(contents), "eur-lex.europa.eu") {
		t.Fatalf("message file missing content or source policy: %s", contents)
	}
	mode, _ := os.ReadFile(modePath)
	if strings.TrimSpace(string(mode)) != "600" {
		t.Fatalf("message file mode = %q, want 600", mode)
	}
}

func TestSendAgentMessageRejectsUnregisteredAgentAndInvalidPolicy(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"export-agent"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	service := New(Config{ConfigPath: configPath})
	base := AgentMessageInput{AgentID: "missing-agent", Message: "hello", SessionKey: "session-1", Sources: []string{"本地业务数据库"}}
	if _, err := service.SendAgentMessage(context.Background(), base); err != ErrAgentUnavailable {
		t.Fatalf("error = %v, want ErrAgentUnavailable", err)
	}
	base.AgentID = "export-agent"
	base.Allowlist = []string{"https://example.com/path"}
	if _, err := service.SendAgentMessage(context.Background(), base); err != ErrInvalidSource {
		t.Fatalf("error = %v, want ErrInvalidSource", err)
	}
}

func TestSendAgentMessageOnlyPassesNonDefaultModelOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	argsPath := filepath.Join(dir, "args.txt")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"defaults":{"model":{"primary":"deepseek/deepseek-v4-flash"}},"list":[{"id":"export-agent"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "`+argsPath+`"
printf '%s' '{"runId":"run-1","status":"ok","result":{"payloads":[{"text":"ok"}]}}'
`)
	service := New(Config{BinaryPath: script, ConfigPath: configPath})
	base := AgentMessageInput{AgentID: "export-agent", Message: "hello", SessionKey: "session-1", Sources: []string{"本地业务数据库"}}
	base.Model = "deepseek/deepseek-v4-flash"
	if _, err := service.SendAgentMessage(context.Background(), base); err != nil {
		t.Fatal(err)
	}
	base.SessionKey = "session-2"
	base.Model = "minimax/MiniMax-M3"
	if _, err := service.SendAgentMessage(context.Background(), base); err != nil {
		t.Fatal(err)
	}
	args, _ := os.ReadFile(argsPath)
	lines := strings.Split(strings.TrimSpace(string(args)), "\n")
	if len(lines) != 2 {
		t.Fatalf("captured args = %q", args)
	}
	if strings.Contains(lines[0], "--model") {
		t.Fatalf("default model should not be passed as override: %s", lines[0])
	}
	if !strings.Contains(lines[1], "--model minimax/MiniMax-M3") {
		t.Fatalf("non-default model should be passed as override: %s", lines[1])
	}
}

func TestSessionHistoryReadsPersistedCoordinatorMessages(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "openclaw.json")
	sessionDir := filepath.Join(dir, "agents", "sta100-coordinator", "sessions")
	sessionPath := filepath.Join(sessionDir, "history.jsonl")
	if err := os.MkdirAll(sessionDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"sta100-coordinator"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	session := strings.Join([]string{
		`{"type":"session","version":3,"id":"history-1","timestamp":"2026-08-14T08:00:00Z"}`,
		`{"type":"thinking_level_change","thinkingLevel":"high"}`,
		`{"type":"message","timestamp":"2026-08-14T08:01:00Z","message":{"role":"user","content":"[STA-100 本次来源约束]\n\n[用户消息]\n请查找德国客户\n\n{\"userMessage\":\"德国客户\"}"}}`,
		`{"type":"message","timestamp":"2026-08-14T08:01:02Z","message":{"role":"assistant","content":[{"type":"text","text":"已找到 2 条客户记录。"}],"provider":"minimax","model":"MiniMax-M3","runId":"run-1"}}`,
		`{"type":"message","timestamp":"2026-08-14T08:02:00Z","message":{"role":"tool","content":"internal"}}`,
	}, "\n") + "\n"
	if err := os.WriteFile(sessionPath, []byte(session), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
if [[ "$1" == "sessions" ]]; then
  printf '%s' '{"sessions":[{"key":"agent:sta100-coordinator:sta100-export-agent-coordinator","updatedAt":1786693496064,"sessionFile":"`+filepath.ToSlash(sessionPath)+`"}]}'
  exit 0
fi
exit 2
`)
	service := New(Config{BinaryPath: script, ConfigPath: configPath})
	messages, err := service.SessionHistory(context.Background(), "sta100-coordinator", "sta100-export-agent-coordinator", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 {
		t.Fatalf("messages = %#v, want user and assistant messages", messages)
	}
	if messages[0].Role != "user" || messages[0].Text != "德国客户" {
		t.Fatalf("user message was not normalized: %#v", messages[0])
	}
	if messages[1].Role != "assistant" || messages[1].Text != "已找到 2 条客户记录。" || messages[1].Model != "MiniMax-M3" {
		t.Fatalf("assistant message was not parsed: %#v", messages[1])
	}
}

func TestParseSessionIndexSupportsOpenClawMapFormat(t *testing.T) {
	records, err := parseSessionIndex([]byte(`{
		"agent:sta100-coordinator:sta100-export-agent-coordinator": {
			"updatedAt": 1786713496064,
			"sessionFile": "/home/User/.openclaw/agents/sta100-coordinator/sessions/example.jsonl"
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Key != "agent:sta100-coordinator:sta100-export-agent-coordinator" {
		t.Fatalf("unexpected records: %#v", records)
	}
}

func TestSyncAgentsUsesApplicationManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "config", "agents.json")
	workspaceRoot := filepath.Join(dir, "workspaces")
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schema_version":1,"workspace_root":"../workspaces","agents":[{"id":"export-agent","name":"出口业务助手","technical_name":"ExportAgent","emoji":"🛒","workspace":"export-agent","role":"domain","visibility":"business","instructions":"# Operating rules\\n\\nUse supplied evidence only."}]}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
if [[ "$1 $2" == "agents list" ]]; then
  printf '%s' '[{"id":"export-agent","identityName":"出口业务助手","identityEmoji":"🛒","model":"demo/model"}]'
  exit 0
fi
exit 9
`)
	service := New(Config{BinaryPath: script, Manifest: manifestPath})
	agents, err := service.SyncAgents(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 1 || agents[0].ID != "export-agent" {
		t.Fatalf("unexpected agents: %#v", agents)
	}
	identity, err := os.ReadFile(filepath.Join(workspaceRoot, "export-agent", "IDENTITY.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(identity), "ExportAgent") {
		t.Fatalf("identity file missing technical name: %s", identity)
	}
	instructions, err := os.ReadFile(filepath.Join(workspaceRoot, "export-agent", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(instructions), "Use supplied evidence only") {
		t.Fatalf("instructions were not synchronized: %s", instructions)
	}
}

func TestHTTPMutationRequiresOperatorHeader(t *testing.T) {
	handler := NewHTTPHandler(New(Config{}), HTTPOptions{})
	request := httptest.NewRequest(http.MethodPost, "/v1/agents/sync", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func writeExecutable(t *testing.T, dir, contents string) string {
	t.Helper()
	path := filepath.Join(dir, "openclaw")
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}
