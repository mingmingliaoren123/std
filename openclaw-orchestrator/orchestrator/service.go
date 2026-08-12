package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrUnavailable      = errors.New("openclaw unavailable")
	ErrInvalidResponse  = errors.New("invalid openclaw response")
	ErrModelUnavailable = errors.New("model unavailable")
	ErrInvalidModel     = errors.New("invalid model id")
	ErrInvalidProvider  = errors.New("invalid provider id")
	ErrInvalidAPIKey    = errors.New("invalid api key")
	ErrManifestMissing  = errors.New("agent manifest unavailable")
	ErrInvalidAgent     = errors.New("invalid agent id")
	ErrAgentUnavailable = errors.New("agent unavailable")
	ErrInvalidMessage   = errors.New("invalid agent message")
	ErrInvalidSession   = errors.New("invalid session key")
	ErrInvalidSource    = errors.New("invalid source policy")
)

var (
	modelIDPattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}/[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
	providerPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,39}$`)
	agentIDPattern  = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
	sessionPattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
	domainPattern   = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+)$`)
)

type Config struct {
	BinaryPath string
	ConfigPath string
	Manifest   string
}

type Service struct {
	bin          string
	configPath   string
	manifestPath string
}

type AuditIssue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

type GatewayStatus struct {
	Available     bool         `json:"available"`
	Version       string       `json:"version"`
	ServiceLoaded bool         `json:"serviceLoaded"`
	ServiceStatus string       `json:"serviceStatus"`
	ServiceState  string       `json:"serviceState"`
	ServiceDetail string       `json:"serviceDetail"`
	RPCOK         bool         `json:"rpcOK"`
	BindHost      string       `json:"bindHost"`
	Port          int          `json:"port"`
	AuditOK       bool         `json:"auditOK"`
	AuditIssues   []AuditIssue `json:"auditIssues"`
}

type CredentialTypes struct {
	APIKey int `json:"apiKey"`
	OAuth  int `json:"oauth"`
	Token  int `json:"token"`
}

type ProviderStatus struct {
	Provider        string          `json:"provider"`
	Configured      bool            `json:"configured"`
	ProfileCount    int             `json:"profileCount"`
	CredentialTypes CredentialTypes `json:"credentialTypes"`
}

type Model struct {
	Key           string   `json:"key"`
	Name          string   `json:"name"`
	Input         string   `json:"input"`
	ContextWindow int      `json:"contextWindow"`
	Local         bool     `json:"local"`
	Available     bool     `json:"available"`
	Tags          []string `json:"tags"`
	Missing       bool     `json:"missing"`
}

type Models struct {
	DefaultModel     string           `json:"defaultModel"`
	ResolvedDefault  string           `json:"resolvedDefault"`
	Fallbacks        []string         `json:"fallbacks"`
	Configured       bool             `json:"configured"`
	MissingProviders []string         `json:"missingProviders"`
	Providers        []ProviderStatus `json:"providers"`
	Models           []Model          `json:"models"`
}

type Agent struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IdentityName  string `json:"identityName"`
	IdentityEmoji string `json:"identityEmoji"`
	Model         string `json:"model"`
	Bindings      int    `json:"bindings"`
	IsDefault     bool   `json:"isDefault"`
	Role          string `json:"role,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
}

type AgentMessageInput struct {
	AgentID    string
	Message    string
	SessionKey string
	Sources    []string
	Allowlist  []string
}

type AgentMessageResult struct {
	AgentID    string     `json:"agentId"`
	SessionKey string     `json:"sessionKey"`
	Text       string     `json:"text"`
	RunID      string     `json:"runId,omitempty"`
	Status     string     `json:"status"`
	Usage      TokenUsage `json:"usage"`
}

type TokenUsage struct {
	Input      int64 `json:"input"`
	Output     int64 `json:"output"`
	CacheRead  int64 `json:"cacheRead"`
	CacheWrite int64 `json:"cacheWrite"`
	Total      int64 `json:"total"`
}

func (u TokenUsage) Available() bool {
	return u.Input > 0 || u.Output > 0 || u.CacheRead > 0 || u.CacheWrite > 0 || u.Total > 0
}

type commandError struct {
	args   []string
	stderr string
	err    error
}

func (e *commandError) Error() string {
	return fmt.Sprintf("openclaw %s failed: %v: %s", strings.Join(e.args, " "), e.err, e.stderr)
}

func (e *commandError) Unwrap() error { return e.err }

func New(config Config) *Service {
	config.BinaryPath = strings.TrimSpace(config.BinaryPath)
	config.ConfigPath = strings.TrimSpace(config.ConfigPath)
	config.Manifest = strings.TrimSpace(config.Manifest)
	return &Service{
		bin:          config.BinaryPath,
		configPath:   config.ConfigPath,
		manifestPath: config.Manifest,
	}
}

func Discover(manifestPath string) *Service {
	bin := strings.TrimSpace(os.Getenv("OPENCLAW_BIN"))
	if bin == "" {
		for _, root := range rootCandidates() {
			candidate := filepath.Clean(filepath.Join(root, "bin", "openclaw"))
			if isExecutable(candidate) {
				bin = candidate
				break
			}
		}
	}
	if bin == "" {
		if path, err := exec.LookPath("openclaw"); err == nil {
			bin = path
		}
	}
	configPath := strings.TrimSpace(os.Getenv("OPENCLAW_CONFIG_PATH"))
	if configPath == "" {
		if home, err := os.UserHomeDir(); err == nil {
			configPath = filepath.Join(home, ".openclaw", "openclaw.json")
		}
	}
	if envManifest := strings.TrimSpace(os.Getenv("OPENCLAW_AGENT_MANIFEST")); envManifest != "" {
		manifestPath = envManifest
	}
	return New(Config{BinaryPath: bin, ConfigPath: configPath, Manifest: manifestPath})
}

func (s *Service) BinaryPath() string { return s.bin }

func (s *Service) ManifestPath() string { return s.manifestPath }

func (s *Service) Status(ctx context.Context) (GatewayStatus, error) {
	var result GatewayStatus
	out, err := s.run(ctx, nil, "gateway", "status", "--json", "--require-rpc", "--timeout", "5000")
	if err != nil {
		return result, err
	}
	var raw struct {
		CLI struct {
			Version string `json:"version"`
		} `json:"cli"`
		Service struct {
			Loaded  bool `json:"loaded"`
			Runtime struct {
				Status   string `json:"status"`
				State    string `json:"state"`
				SubState string `json:"subState"`
			} `json:"runtime"`
			ConfigAudit struct {
				OK     bool         `json:"ok"`
				Issues []AuditIssue `json:"issues"`
			} `json:"configAudit"`
		} `json:"service"`
		Gateway struct {
			BindHost string `json:"bindHost"`
			Port     int    `json:"port"`
			Version  string `json:"version"`
		} `json:"gateway"`
		RPC struct {
			OK      bool   `json:"ok"`
			Version string `json:"version"`
		} `json:"rpc"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return result, fmt.Errorf("%w: gateway status: %v", ErrInvalidResponse, err)
	}
	result = GatewayStatus{
		Available:     true,
		Version:       firstNonEmpty(raw.Gateway.Version, raw.RPC.Version, raw.CLI.Version),
		ServiceLoaded: raw.Service.Loaded,
		ServiceStatus: raw.Service.Runtime.Status,
		ServiceState:  raw.Service.Runtime.State,
		ServiceDetail: raw.Service.Runtime.SubState,
		RPCOK:         raw.RPC.OK,
		BindHost:      raw.Gateway.BindHost,
		Port:          raw.Gateway.Port,
		AuditOK:       raw.Service.ConfigAudit.OK,
		AuditIssues:   raw.Service.ConfigAudit.Issues,
	}
	return result, nil
}

func (s *Service) Models(ctx context.Context) (Models, error) {
	type commandResult struct {
		stdout []byte
		err    error
	}
	statusChannel := make(chan commandResult, 1)
	listChannel := make(chan commandResult, 1)
	go func() {
		out, err := s.run(ctx, nil, "models", "status", "--json")
		statusChannel <- commandResult{stdout: out, err: err}
	}()
	go func() {
		out, err := s.run(ctx, nil, "models", "list", "--json")
		listChannel <- commandResult{stdout: out, err: err}
	}()
	statusResult, listResult := <-statusChannel, <-listChannel
	if statusResult.err != nil {
		return Models{}, statusResult.err
	}
	if listResult.err != nil {
		return Models{}, listResult.err
	}
	var status struct {
		DefaultModel    string   `json:"defaultModel"`
		ResolvedDefault string   `json:"resolvedDefault"`
		Fallbacks       []string `json:"fallbacks"`
		Auth            struct {
			MissingProvidersInUse []string `json:"missingProvidersInUse"`
			Providers             []struct {
				Provider string `json:"provider"`
				Profiles struct {
					Count  int `json:"count"`
					APIKey int `json:"apiKey"`
					OAuth  int `json:"oauth"`
					Token  int `json:"token"`
				} `json:"profiles"`
			} `json:"providers"`
		} `json:"auth"`
	}
	var list struct {
		Models []Model `json:"models"`
	}
	if json.Unmarshal(statusResult.stdout, &status) != nil || json.Unmarshal(listResult.stdout, &list) != nil {
		return Models{}, ErrInvalidResponse
	}
	providers := make([]ProviderStatus, 0, len(status.Auth.Providers))
	for _, provider := range status.Auth.Providers {
		providers = append(providers, ProviderStatus{
			Provider:     provider.Provider,
			Configured:   provider.Profiles.Count > 0,
			ProfileCount: provider.Profiles.Count,
			CredentialTypes: CredentialTypes{
				APIKey: provider.Profiles.APIKey,
				OAuth:  provider.Profiles.OAuth,
				Token:  provider.Profiles.Token,
			},
		})
	}
	return Models{
		DefaultModel:     status.DefaultModel,
		ResolvedDefault:  status.ResolvedDefault,
		Fallbacks:        status.Fallbacks,
		Configured:       status.ResolvedDefault != "" && len(status.Auth.MissingProvidersInUse) == 0,
		MissingProviders: status.Auth.MissingProvidersInUse,
		Providers:        providers,
		Models:           list.Models,
	}, nil
}

func (s *Service) SetDefaultModel(ctx context.Context, modelID string) error {
	modelID = strings.TrimSpace(modelID)
	if !modelIDPattern.MatchString(modelID) {
		return ErrInvalidModel
	}
	models, err := s.Models(ctx)
	if err != nil {
		return err
	}
	for _, model := range models.Models {
		if model.Key == modelID && model.Available && !model.Missing {
			_, err = s.run(ctx, nil, "models", "set", modelID)
			return err
		}
	}
	return ErrModelUnavailable
}

func (s *Service) SaveAPIKey(ctx context.Context, provider, apiKey string) (string, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	apiKey = strings.TrimSpace(apiKey)
	if !providerPattern.MatchString(provider) {
		return "", ErrInvalidProvider
	}
	if len(apiKey) < 8 || len(apiKey) > 8192 || strings.ContainsAny(apiKey, "\r\n\x00") {
		return "", ErrInvalidAPIKey
	}
	profileID := provider + ":default"
	_, err := s.run(ctx, strings.NewReader(apiKey+"\n"), "models", "auth", "paste-api-key", "--provider", provider, "--profile-id", profileID)
	if err != nil {
		return "", err
	}
	return profileID, nil
}

func (s *Service) ListAgents(ctx context.Context) ([]Agent, error) {
	if agents, err := s.listConfiguredAgents(); err == nil {
		return s.decorateAgents(agents), nil
	}
	out, err := s.run(ctx, nil, "agents", "list", "--json")
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := json.Unmarshal(out, &agents); err != nil {
		return nil, fmt.Errorf("%w: agents list: %v", ErrInvalidResponse, err)
	}
	return s.decorateAgents(agents), nil
}

func (s *Service) SendAgentMessage(ctx context.Context, input AgentMessageInput) (AgentMessageResult, error) {
	input.AgentID = strings.TrimSpace(strings.ToLower(input.AgentID))
	input.Message = strings.TrimSpace(input.Message)
	input.SessionKey = strings.TrimSpace(input.SessionKey)
	if !agentIDPattern.MatchString(input.AgentID) {
		return AgentMessageResult{}, ErrInvalidAgent
	}
	if input.Message == "" || len(input.Message) > 32<<10 || strings.ContainsRune(input.Message, '\x00') {
		return AgentMessageResult{}, ErrInvalidMessage
	}
	if !sessionPattern.MatchString(input.SessionKey) {
		return AgentMessageResult{}, ErrInvalidSession
	}
	agents, err := s.ListAgents(ctx)
	if err != nil {
		return AgentMessageResult{}, err
	}
	available := false
	for _, agent := range agents {
		if agent.ID == input.AgentID {
			available = true
			break
		}
	}
	if !available {
		return AgentMessageResult{}, ErrAgentUnavailable
	}
	if len(input.Sources) == 0 || len(input.Sources) > 8 || len(input.Allowlist) > 50 {
		return AgentMessageResult{}, ErrInvalidSource
	}
	for _, source := range input.Sources {
		if source != "本地业务数据库" && source != "客户私有知识库" && source != "联网检索" {
			return AgentMessageResult{}, ErrInvalidSource
		}
	}
	for _, domain := range input.Allowlist {
		if !domainPattern.MatchString(strings.TrimSpace(domain)) {
			return AgentMessageResult{}, ErrInvalidSource
		}
	}

	contents := buildAgentMessage(input)
	messageFile, err := os.CreateTemp("", "sta100-agent-message-*.txt")
	if err != nil {
		return AgentMessageResult{}, fmt.Errorf("create agent message file: %w", err)
	}
	messagePath := messageFile.Name()
	defer os.Remove(messagePath)
	if err := messageFile.Chmod(0o600); err != nil {
		messageFile.Close()
		return AgentMessageResult{}, fmt.Errorf("secure agent message file: %w", err)
	}
	if _, err := messageFile.WriteString(contents); err != nil {
		messageFile.Close()
		return AgentMessageResult{}, fmt.Errorf("write agent message file: %w", err)
	}
	if err := messageFile.Close(); err != nil {
		return AgentMessageResult{}, fmt.Errorf("close agent message file: %w", err)
	}

	out, err := s.run(ctx, nil, "agent", "--agent", input.AgentID, "--message-file", messagePath, "--session-key", input.SessionKey, "--json", "--timeout", "90")
	if err != nil {
		return AgentMessageResult{}, err
	}
	result, err := parseAgentMessageResult(out)
	if err != nil {
		return AgentMessageResult{}, err
	}
	result.AgentID = input.AgentID
	result.SessionKey = input.SessionKey
	return result, nil
}

func buildAgentMessage(input AgentMessageInput) string {
	var builder strings.Builder
	builder.WriteString("[STA-100 本次来源约束]\n允许来源类别：")
	builder.WriteString(strings.Join(input.Sources, "、"))
	if len(input.Allowlist) > 0 {
		builder.WriteString("\n联网检索仅允许以下域名：")
		builder.WriteString(strings.Join(input.Allowlist, "、"))
	}
	builder.WriteString("\n回答时请区分事实来源；无法访问所选来源时必须明确说明，不得假设已检索。\n\n[用户消息]\n")
	builder.WriteString(input.Message)
	return builder.String()
}

func parseAgentMessageResult(out []byte) (AgentMessageResult, error) {
	type rawUsage struct {
		Input       int64 `json:"input"`
		Output      int64 `json:"output"`
		CacheRead   int64 `json:"cacheRead"`
		CacheWrite  int64 `json:"cacheWrite"`
		Total       int64 `json:"total"`
		TotalTokens int64 `json:"totalTokens"`
	}
	type rawMeta struct {
		AgentMeta struct {
			Usage rawUsage `json:"usage"`
		} `json:"agentMeta"`
	}
	var raw struct {
		RunID   string `json:"runId"`
		Status  string `json:"status"`
		Summary string `json:"summary"`
		Result  struct {
			Payloads []struct {
				Text string `json:"text"`
			} `json:"payloads"`
			Meta rawMeta `json:"meta"`
		} `json:"result"`
		Payloads []struct {
			Text string `json:"text"`
		} `json:"payloads"`
		Meta rawMeta `json:"meta"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return AgentMessageResult{}, fmt.Errorf("%w: agent result: %v", ErrInvalidResponse, err)
	}
	texts := make([]string, 0, len(raw.Result.Payloads)+len(raw.Payloads))
	for _, payload := range append(raw.Result.Payloads, raw.Payloads...) {
		if text := strings.TrimSpace(payload.Text); text != "" {
			texts = append(texts, text)
		}
	}
	text := strings.Join(texts, "\n\n")
	if text == "" {
		text = strings.TrimSpace(raw.Summary)
	}
	if text == "" {
		return AgentMessageResult{}, fmt.Errorf("%w: agent result has no text", ErrInvalidResponse)
	}
	status := strings.TrimSpace(raw.Status)
	if status == "" {
		status = "ok"
	}
	usage := raw.Result.Meta.AgentMeta.Usage
	if usage.Input == 0 && usage.Output == 0 && usage.CacheRead == 0 && usage.CacheWrite == 0 && usage.Total == 0 && usage.TotalTokens == 0 {
		usage = raw.Meta.AgentMeta.Usage
	}
	total := usage.Total
	if total == 0 {
		total = usage.TotalTokens
	}
	if total == 0 {
		total = usage.Input + usage.Output + usage.CacheRead + usage.CacheWrite
	}
	return AgentMessageResult{
		Text: text, RunID: raw.RunID, Status: status,
		Usage: TokenUsage{Input: usage.Input, Output: usage.Output, CacheRead: usage.CacheRead, CacheWrite: usage.CacheWrite, Total: total},
	}, nil
}

func (s *Service) listConfiguredAgents() ([]Agent, error) {
	if s.configPath == "" {
		return nil, errors.New("openclaw config path unavailable")
	}
	contents, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, fmt.Errorf("read openclaw config: %w", err)
	}
	var config struct {
		Agents struct {
			Defaults struct {
				Model json.RawMessage `json:"model"`
			} `json:"defaults"`
			List []struct {
				ID       string          `json:"id"`
				Name     string          `json:"name"`
				Model    json.RawMessage `json:"model"`
				Identity struct {
					Name  string `json:"name"`
					Emoji string `json:"emoji"`
				} `json:"identity"`
			} `json:"list"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("parse openclaw config: %w", err)
	}
	if len(config.Agents.List) == 0 {
		return nil, errors.New("openclaw agent list is empty")
	}
	defaultModel := parseModelReference(config.Agents.Defaults.Model)
	agents := make([]Agent, 0, len(config.Agents.List))
	for index, configured := range config.Agents.List {
		if configured.ID == "" {
			continue
		}
		model := parseModelReference(configured.Model)
		if model == "" {
			model = defaultModel
		}
		agents = append(agents, Agent{
			ID:            configured.ID,
			Name:          configured.Name,
			IdentityName:  configured.Identity.Name,
			IdentityEmoji: configured.Identity.Emoji,
			Model:         model,
			IsDefault:     configured.ID == "main" || index == 0,
		})
	}
	return agents, nil
}

func parseModelReference(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var value string
	if json.Unmarshal(raw, &value) == nil {
		return value
	}
	var reference struct {
		Primary string `json:"primary"`
	}
	if json.Unmarshal(raw, &reference) == nil {
		return reference.Primary
	}
	return ""
}

func (s *Service) run(ctx context.Context, stdin io.Reader, args ...string) ([]byte, error) {
	if s.bin == "" || !isExecutable(s.bin) {
		return nil, ErrUnavailable
	}
	command := exec.CommandContext(ctx, s.bin, args...)
	command.Stdin = stdin
	command.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, &commandError{args: append([]string(nil), args...), stderr: tail(stderr.String(), 800), err: err}
	}
	return stdout.Bytes(), nil
}

func rootCandidates() []string {
	var candidates []string
	if executable, err := os.Executable(); err == nil {
		dir := filepath.Dir(executable)
		candidates = append(candidates, filepath.Join(dir, "..", "..", "openclaw"), filepath.Join(dir, "..", "openclaw"))
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "..", "openclaw"), filepath.Join(cwd, "openclaw"))
	}
	return candidates
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func tail(value string, limit int) string {
	value = strings.TrimSpace(value)
	if len(value) <= limit {
		return value
	}
	return value[len(value)-limit:]
}
