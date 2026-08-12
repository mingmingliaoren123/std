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
  "models status") printf '%s' '{"defaultModel":"demo/old","resolvedDefault":"demo/old","auth":{"providers":[]}}' ;;
  "models list") printf '%s' '{"models":[{"key":"demo/new","name":"New","available":true,"missing":false}]}' ;;
  "models set") printf '%s\n' "$*" >> "`+logPath+`" ;;
  *) exit 2 ;;
esac
`)
	service := New(Config{BinaryPath: script})
	if err := service.SetDefaultModel(context.Background(), "demo/new"); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(contents)) != "models set demo/new" {
		t.Fatalf("unexpected command: %q", contents)
	}
	if err := service.SetDefaultModel(context.Background(), "demo/missing"); err != ErrModelUnavailable {
		t.Fatalf("error = %v, want ErrModelUnavailable", err)
	}
}

func TestSaveAPIKeyUsesStdinAndDoesNotPutSecretInArguments(t *testing.T) {
	dir := t.TempDir()
	argsPath := filepath.Join(dir, "args.log")
	stdinPath := filepath.Join(dir, "stdin.log")
	script := writeExecutable(t, dir, `#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" > "`+argsPath+`"
IFS= read -r secret
printf '%s' "$secret" > "`+stdinPath+`"
`)
	service := New(Config{BinaryPath: script})
	secret := "secret-key-value"
	profile, err := service.SaveAPIKey(context.Background(), "deepseek", secret)
	if err != nil {
		t.Fatal(err)
	}
	if profile != "deepseek:default" {
		t.Fatalf("profile = %q", profile)
	}
	args, _ := os.ReadFile(argsPath)
	stdin, _ := os.ReadFile(stdinPath)
	if strings.Contains(string(args), secret) {
		t.Fatal("secret appeared in command arguments")
	}
	if string(stdin) != secret {
		t.Fatalf("stdin = %q, want secret", stdin)
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
