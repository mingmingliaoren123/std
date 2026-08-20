package orchestrator

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
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
	PluginDir  string
}

const agentMessageTimeoutSeconds = 300

type Service struct {
	bin          string
	configPath   string
	manifestPath string
	pluginDir    string
	qrMu         sync.Mutex
	qrSessions   map[string]*feishuQRSession
	pluginMu     sync.Mutex
	pluginCache  []pluginInventoryEntry
	pluginAt     time.Time
	messageGate  chan struct{}
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
	Provider         string          `json:"provider"`
	Configured       bool            `json:"configured"`
	APIKeyConfigured bool            `json:"apiKeyConfigured"`
	APIKeySupported  bool            `json:"apiKeySupported"`
	ProfileCount     int             `json:"profileCount"`
	CredentialTypes  CredentialTypes `json:"credentialTypes"`
}

type Model struct {
	Key             string   `json:"key"`
	Name            string   `json:"name"`
	Input           string   `json:"input"`
	ContextWindow   int      `json:"contextWindow"`
	Local           bool     `json:"local"`
	Available       bool     `json:"available"`
	Source          string   `json:"source,omitempty"`
	Tags            []string `json:"tags"`
	Missing         bool     `json:"missing"`
	LastTestStatus  string   `json:"lastTestStatus,omitempty"`
	LastTestMessage string   `json:"lastTestMessage,omitempty"`
	LastTestAt      string   `json:"lastTestAt,omitempty"`
}

type Models struct {
	CatalogVersion   string            `json:"catalogVersion"`
	DefaultModel     string            `json:"defaultModel"`
	ResolvedDefault  string            `json:"resolvedDefault"`
	Fallbacks        []string          `json:"fallbacks"`
	Configured       bool              `json:"configured"`
	MissingProviders []string          `json:"missingProviders"`
	Providers        []ProviderStatus  `json:"providers"`
	Models           []Model           `json:"models"`
	ConfiguredModels []Model           `json:"configuredModels,omitempty"`
	CatalogModels    []Model           `json:"catalogModels,omitempty"`
	ProviderBaseURLs map[string]string `json:"providerBaseUrls,omitempty"`
}

func (m *Models) LimitToConfiguredSelection() {
	if len(m.CatalogModels) == 0 {
		m.CatalogModels = append([]Model(nil), m.Models...)
	}
	selected := map[string]bool{}
	if strings.TrimSpace(m.ResolvedDefault) != "" {
		selected[strings.TrimSpace(m.ResolvedDefault)] = true
	}
	if strings.TrimSpace(m.DefaultModel) != "" {
		selected[strings.TrimSpace(m.DefaultModel)] = true
	}
	if len(selected) == 0 && len(m.ConfiguredModels) > 0 {
		for _, model := range m.ConfiguredModels {
			selected[model.Key] = true
		}
	}
	if len(selected) == 0 {
		m.ConfiguredModels = nil
		return
	}
	filtered := make([]Model, 0, len(selected))
	seen := map[string]bool{}
	for _, model := range m.Models {
		if selected[model.Key] && !seen[model.Key] {
			filtered = append(filtered, model)
			seen[model.Key] = true
		}
	}
	m.ConfiguredModels = filtered
}

type Plugin struct {
	ID           string   `json:"id"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Version      string   `json:"version,omitempty"`
	Enabled      bool     `json:"enabled"`
	Status       string   `json:"status"`
	Source       string   `json:"source"`
	Providers    []string `json:"providers,omitempty"`
	Channels     []string `json:"channels,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
	Configurable bool     `json:"configurable"`
}

type Channel struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	PluginID     string `json:"pluginId,omitempty"`
	InstallSpec  string `json:"installSpec,omitempty"`
	BindingMode  string `json:"bindingMode,omitempty"`
	CanInstall   bool   `json:"canInstall"`
	CanUninstall bool   `json:"canUninstall"`
	Installed    bool   `json:"installed"`
	Enabled      bool   `json:"enabled"`
	Configured   bool   `json:"configured"`
	Running      bool   `json:"running"`
	LastError    string `json:"lastError,omitempty"`
	Origin       string `json:"origin"`
	AccountCount int    `json:"accountCount"`
	Status       string `json:"status"`
}

type ChannelAccountRequest struct {
	Channel  string `json:"channel"`
	Account  string `json:"account,omitempty"`
	Name     string `json:"name,omitempty"`
	Token    string `json:"token,omitempty"`
	Secret   string `json:"secret,omitempty"`
	URL      string `json:"url,omitempty"`
	BaseURL  string `json:"baseUrl,omitempty"`
	HTTPURL  string `json:"httpUrl,omitempty"`
	CLIPath  string `json:"cliPath,omitempty"`
	DBPath   string `json:"dbPath,omitempty"`
	Service  string `json:"service,omitempty"`
	Region   string `json:"region,omitempty"`
	Password string `json:"password,omitempty"`
	UseEnv   bool   `json:"useEnv,omitempty"`
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
	AgentID     string
	Model       string
	Message     string
	SessionKey  string
	Sources     []string
	Allowlist   []string
	Attachments []MessageAttachment
}

type MessageAttachment struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Mime string `json:"mime,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type AgentMessageResult struct {
	AgentID    string     `json:"agentId"`
	SessionKey string     `json:"sessionKey"`
	Text       string     `json:"text"`
	RunID      string     `json:"runId,omitempty"`
	Status     string     `json:"status"`
	Usage      TokenUsage `json:"usage"`
}

// SessionHistoryMessage is the user-facing subset of an OpenClaw session.
// Internal prompts, tool calls and thinking events are intentionally omitted.
type SessionHistoryMessage struct {
	Role      string `json:"role"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt,omitempty"`
	Model     string `json:"model,omitempty"`
	Provider  string `json:"provider,omitempty"`
	RunID     string `json:"runId,omitempty"`
	Error     string `json:"error,omitempty"`
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
	config.PluginDir = strings.TrimSpace(config.PluginDir)
	if config.PluginDir == "" && config.BinaryPath != "" {
		config.PluginDir = filepath.Clean(filepath.Join(filepath.Dir(config.BinaryPath), "..", "..", "openclaw", "app", "node_modules", "openclaw", "dist", "extensions"))
	}
	return &Service{
		bin:          config.BinaryPath,
		configPath:   config.ConfigPath,
		manifestPath: config.Manifest,
		pluginDir:    config.PluginDir,
		qrSessions:   make(map[string]*feishuQRSession),
		messageGate:  make(chan struct{}, 2),
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
	pluginDir := ""
	if bin != "" {
		pluginDir = filepath.Clean(filepath.Join(filepath.Dir(bin), "..", "..", "openclaw", "app", "node_modules", "openclaw", "dist", "extensions"))
	}
	return New(Config{BinaryPath: bin, ConfigPath: configPath, Manifest: manifestPath, PluginDir: pluginDir})
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
	local, localOK := s.modelsFromConfig()
	statusContext, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancel()
	statusOutput, statusErr := s.run(statusContext, nil, "models", "status", "--json")
	if statusErr != nil {
		if localOK {
			return local, nil
		}
		return Models{}, statusErr
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
	if json.Unmarshal(statusOutput, &status) != nil {
		return Models{}, ErrInvalidResponse
	}
	list := struct {
		Models []Model `json:"models"`
	}{Models: fixedModels()}
	if localOK {
		known := make(map[string]bool, len(list.Models))
		for _, model := range list.Models {
			known[model.Key] = true
		}
		for _, model := range local.Models {
			if !known[model.Key] {
				list.Models = append(list.Models, model)
				known[model.Key] = true
			}
		}
	}
	providerStatuses := make(map[string]ProviderStatus, len(status.Auth.Providers))
	for _, provider := range status.Auth.Providers {
		providerStatuses[provider.Provider] = ProviderStatus{
			Provider:         provider.Provider,
			Configured:       provider.Profiles.Count > 0,
			APIKeyConfigured: provider.Profiles.APIKey > 0,
			APIKeySupported:  fixedAPIKeyProvider(provider.Provider),
			ProfileCount:     provider.Profiles.Count,
			CredentialTypes: CredentialTypes{
				APIKey: provider.Profiles.APIKey,
				OAuth:  provider.Profiles.OAuth,
				Token:  provider.Profiles.Token,
			},
		}
	}
	for _, model := range list.Models {
		provider, _, ok := strings.Cut(model.Key, "/")
		if ok && providerPattern.MatchString(provider) {
			if _, exists := providerStatuses[provider]; !exists {
				providerStatuses[provider] = ProviderStatus{Provider: provider, APIKeySupported: supportsAPIKeyProvider(provider)}
			}
		}
	}
	for index := range list.Models {
		provider, _, _ := strings.Cut(list.Models[index].Key, "/")
		configured := providerStatuses[provider].Configured
		if !configured && localOK {
			for _, localModel := range local.Models {
				if localModel.Key == list.Models[index].Key {
					configured = localModel.Available
					break
				}
			}
		}
		list.Models[index].Available = configured
		if list.Models[index].Source == "" {
			list.Models[index].Source = "fixed_catalog"
		}
		list.Models[index].Missing = !configured
	}
	providers := make([]ProviderStatus, 0, len(providerStatuses))
	for _, provider := range providerStatuses {
		providers = append(providers, provider)
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].Provider < providers[j].Provider })
	result := Models{
		CatalogVersion:   ModelCatalogVersion,
		DefaultModel:     status.DefaultModel,
		ResolvedDefault:  status.ResolvedDefault,
		Fallbacks:        status.Fallbacks,
		Configured:       status.ResolvedDefault != "" && len(status.Auth.MissingProvidersInUse) == 0,
		MissingProviders: status.Auth.MissingProvidersInUse,
		Providers:        providers,
		Models:           list.Models,
		CatalogModels:    append([]Model(nil), list.Models...),
	}
	if localOK {
		result.ProviderBaseURLs = local.ProviderBaseURLs
		if result.DefaultModel == "" {
			result.DefaultModel = local.DefaultModel
		}
		if result.ResolvedDefault == "" {
			result.ResolvedDefault = local.ResolvedDefault
		}
		if local.Configured {
			result.Configured = true
		}
		selected := map[string]bool{}
		if result.ResolvedDefault != "" {
			selected[result.ResolvedDefault] = true
		}
		if result.DefaultModel != "" {
			selected[result.DefaultModel] = true
		}
		for _, model := range local.Models {
			selected[model.Key] = true
		}
		filtered := make([]Model, 0, len(selected))
		for _, model := range result.CatalogModels {
			if selected[model.Key] {
				filtered = append(filtered, model)
			}
		}
		result.Models = filtered
	} else {
		result.LimitToConfiguredSelection()
	}
	return result, nil
}

func (s *Service) ModelSnapshot() Models {
	if local, ok := s.modelsFromConfig(); ok {
		return local
	}
	models := fixedModels()
	providers := make([]ProviderStatus, 0)
	seen := map[string]bool{}
	for _, model := range models {
		provider, _, _ := strings.Cut(model.Key, "/")
		if seen[provider] {
			continue
		}
		seen[provider] = true
		providers = append(providers, ProviderStatus{Provider: provider, APIKeySupported: fixedAPIKeyProvider(provider)})
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].Provider < providers[j].Provider })
	return Models{
		CatalogVersion: ModelCatalogVersion,
		Providers:      providers,
		Models:         models,
	}
}

func (s *Service) modelsFromConfig() (Models, bool) {
	if s.configPath == "" {
		return Models{}, false
	}
	contents, err := os.ReadFile(s.configPath)
	if err != nil {
		return Models{}, false
	}
	var config struct {
		Agents struct {
			Defaults struct {
				Model  json.RawMessage        `json:"model"`
				Models map[string]interface{} `json:"models"`
			} `json:"defaults"`
		} `json:"agents"`
		Models struct {
			Providers map[string]struct {
				BaseURL string `json:"baseUrl"`
				Models  []struct {
					ID            string   `json:"id"`
					Name          string   `json:"name"`
					Input         []string `json:"input"`
					ContextWindow int      `json:"contextWindow"`
				} `json:"models"`
			} `json:"providers"`
		} `json:"models"`
		Auth struct {
			Profiles map[string]struct {
				Provider string `json:"provider"`
				Mode     string `json:"mode"`
			} `json:"profiles"`
		} `json:"auth"`
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		return Models{}, false
	}
	models := fixedModels()
	known := make(map[string]bool, len(models))
	for _, model := range models {
		known[model.Key] = true
	}
	providerConfigured := map[string]ProviderStatus{}
	providerBaseURLs := map[string]string{}
	for _, model := range models {
		provider, _, _ := strings.Cut(model.Key, "/")
		if _, exists := providerConfigured[provider]; !exists {
			providerConfigured[provider] = ProviderStatus{Provider: provider, APIKeySupported: supportsAPIKeyProvider(provider)}
		}
	}
	for provider, definition := range config.Models.Providers {
		if baseURL := strings.TrimSpace(definition.BaseURL); baseURL != "" {
			providerBaseURLs[provider] = baseURL
		}
		for _, item := range definition.Models {
			key := provider + "/" + strings.TrimSpace(item.ID)
			if item.ID == "" || known[key] {
				continue
			}
			input := strings.Join(item.Input, "+")
			if input == "" {
				input = "text"
			}
			models = append(models, Model{Key: key, Name: firstNonEmpty(item.Name, item.ID), Input: input, ContextWindow: item.ContextWindow, Source: "configured", Available: true})
			known[key] = true
		}
	}
	for profileID, profile := range config.Auth.Profiles {
		provider := strings.TrimSpace(profile.Provider)
		if provider == "" {
			provider = strings.SplitN(profileID, ":", 2)[0]
		}
		if provider == "" {
			continue
		}
		status := providerConfigured[provider]
		status.Provider = provider
		status.Configured = true
		if strings.EqualFold(profile.Mode, "api_key") || profile.Mode == "" {
			status.APIKeyConfigured = true
			status.CredentialTypes.APIKey++
		}
		status.ProfileCount++
		status.APIKeySupported = supportsAPIKeyProvider(provider)
		providerConfigured[provider] = status
	}
	for index := range models {
		provider, _, _ := strings.Cut(models[index].Key, "/")
		status := providerConfigured[provider]
		models[index].Available = status.Configured
		models[index].Missing = !models[index].Available
		if models[index].Source == "" {
			models[index].Source = "fixed_catalog"
		}
	}
	providers := make([]ProviderStatus, 0, len(providerConfigured))
	for _, status := range providerConfigured {
		providers = append(providers, status)
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].Provider < providers[j].Provider })
	defaultModel := parseModelReference(config.Agents.Defaults.Model)
	if defaultModel == "" {
		for model := range config.Agents.Defaults.Models {
			defaultModel = model
			break
		}
	}
	selectedModels := map[string]bool{}
	if defaultModel != "" {
		selectedModels[defaultModel] = true
	}
	for model := range config.Agents.Defaults.Models {
		if strings.TrimSpace(model) != "" {
			selectedModels[strings.TrimSpace(model)] = true
		}
	}
	configuredModels := make([]Model, 0, len(selectedModels))
	for _, model := range models {
		if selectedModels[model.Key] {
			configuredModels = append(configuredModels, model)
		}
	}
	configured := false
	for _, provider := range providerConfigured {
		if provider.Configured {
			configured = true
			break
		}
	}
	return Models{
		CatalogVersion:   ModelCatalogVersion,
		DefaultModel:     defaultModel,
		ResolvedDefault:  defaultModel,
		Configured:       configured,
		Providers:        providers,
		Models:           append([]Model(nil), models...),
		ConfiguredModels: configuredModels,
		CatalogModels:    append([]Model(nil), models...),
		ProviderBaseURLs: providerBaseURLs,
	}, true
}

func supportsAPIKeyProvider(provider string) bool {
	provider = strings.ToLower(strings.TrimSpace(provider))
	return fixedAPIKeyProvider(provider)
}

func (s *Service) SetDefaultModel(ctx context.Context, modelID string) error {
	modelID = strings.TrimSpace(modelID)
	if !modelIDPattern.MatchString(modelID) {
		return ErrInvalidModel
	}
	if !s.modelAvailable(ctx, modelID) {
		return ErrModelUnavailable
	}
	_, err := s.run(ctx, nil, "models", "set", modelID)
	return err
}

func (s *Service) SetConfiguredDefaultModel(modelID string) error {
	modelID = strings.TrimSpace(modelID)
	if !modelIDPattern.MatchString(modelID) {
		return ErrInvalidModel
	}
	if !s.configuredModelAvailable(modelID) {
		return ErrModelUnavailable
	}
	return s.applyDefaultModelConfig(modelID)
}

func (s *Service) SetConfiguredModelSelection(defaultModel string, selectedModels []string) error {
	defaultModel = strings.TrimSpace(defaultModel)
	if !modelIDPattern.MatchString(defaultModel) {
		return ErrInvalidModel
	}
	seen := map[string]bool{}
	models := make([]string, 0, len(selectedModels)+1)
	for _, modelID := range append([]string{defaultModel}, selectedModels...) {
		modelID = strings.TrimSpace(modelID)
		if modelID == "" || seen[modelID] {
			continue
		}
		if !modelIDPattern.MatchString(modelID) {
			return ErrInvalidModel
		}
		if !s.configuredModelAvailable(modelID) {
			return ErrModelUnavailable
		}
		seen[modelID] = true
		models = append(models, modelID)
	}
	if len(models) == 0 {
		return ErrInvalidModel
	}
	return s.applyModelSelectionConfig(defaultModel, models)
}

func (s *Service) ClearConfiguredModelSelection() error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	agents := ensureObject(config, "agents")
	defaults := ensureObject(agents, "defaults")
	delete(defaults, "model")
	defaults["models"] = map[string]any{}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func (s *Service) SetProviderBaseURL(provider, baseURL string) error {
	provider = strings.TrimSpace(strings.ToLower(provider))
	baseURL = strings.TrimSpace(baseURL)
	if !providerPattern.MatchString(provider) {
		return ErrInvalidProvider
	}
	if baseURL != "" {
		parsed, err := url.Parse(baseURL)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return ErrInvalidResponse
		}
	}
	return s.applyProviderBaseURLConfig(provider, baseURL)
}

func (s *Service) modelAvailable(ctx context.Context, modelID string) bool {
	if provider, _, ok := strings.Cut(modelID, "/"); ok && s.providerConfigured(provider) {
		for _, model := range fixedModelCatalog {
			if model.Key == modelID {
				return true
			}
		}
	}
	models, err := s.Models(ctx)
	if err == nil {
		for _, model := range append(append([]Model(nil), models.Models...), models.CatalogModels...) {
			if model.Key == modelID {
				return model.Available
			}
		}
	}
	return s.configuredModel(modelID)
}

func (s *Service) configuredModelAvailable(modelID string) bool {
	_, _, ok := strings.Cut(modelID, "/")
	if !ok {
		return false
	}
	for _, model := range fixedModelCatalog {
		if model.Key == modelID {
			return true
		}
	}
	return s.configuredModel(modelID)
}

func (s *Service) providerConfigured(provider string) bool {
	if s.configPath == "" {
		return false
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return false
	}
	var config struct {
		Auth struct {
			Profiles map[string]struct {
				Provider string `json:"provider"`
			} `json:"profiles"`
		} `json:"auth"`
	}
	if json.Unmarshal(raw, &config) != nil {
		return false
	}
	for profileID, profile := range config.Auth.Profiles {
		value := strings.TrimSpace(profile.Provider)
		if value == "" {
			value = strings.SplitN(profileID, ":", 2)[0]
		}
		if value == provider {
			return true
		}
	}
	return false
}

func (s *Service) configuredModel(modelID string) bool {
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return false
	}
	var config struct {
		Models struct {
			Providers map[string]struct {
				Models []struct {
					ID string `json:"id"`
				} `json:"models"`
			} `json:"providers"`
		} `json:"models"`
	}
	if json.Unmarshal(raw, &config) != nil {
		return false
	}
	provider, model, ok := strings.Cut(modelID, "/")
	if !ok {
		return false
	}
	for _, entry := range config.Models.Providers[provider].Models {
		if entry.ID == model {
			return true
		}
	}
	return false
}

func (s *Service) SaveAPIKey(ctx context.Context, provider, apiKey string) (string, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	apiKey = strings.TrimSpace(apiKey)
	if !providerPattern.MatchString(provider) {
		return "", ErrInvalidProvider
	}
	if !supportsAPIKeyProvider(provider) {
		return "", ErrInvalidProvider
	}
	if len(apiKey) < 8 || len(apiKey) > 8192 || strings.ContainsAny(apiKey, "\r\n\x00") {
		return "", ErrInvalidAPIKey
	}
	profileID := provider + ":default"
	if err := s.saveAPIKeyProfile(ctx, provider, profileID, apiKey); err != nil {
		return "", err
	}
	return profileID, nil
}

func (s *Service) ConfigureWebSearchForModelProvider(ctx context.Context, provider, apiKey string) error {
	provider = strings.TrimSpace(strings.ToLower(provider))
	webSearchProvider, ok := webSearchProviderForModelProvider(provider)
	if !ok {
		return nil
	}
	if !providerPattern.MatchString(provider) {
		return ErrInvalidProvider
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		savedAPIKey, err := s.LoadAPIKey(ctx, provider)
		if err != nil {
			return err
		}
		apiKey = savedAPIKey
	}
	if len(apiKey) < 8 || len(apiKey) > 8192 || strings.ContainsAny(apiKey, "\r\n\x00") {
		return ErrInvalidAPIKey
	}
	return s.applyWebSearchProviderConfig(provider, webSearchProvider, apiKey)
}

func (s *Service) LoadAPIKey(ctx context.Context, provider string) (string, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	if !providerPattern.MatchString(provider) {
		return "", ErrInvalidProvider
	}
	agentDir, err := s.defaultAgentDir()
	if err != nil {
		return "", err
	}
	return readAuthProfileAPIKey(ctx, filepath.Join(agentDir, "openclaw-agent.sqlite"), provider, provider+":default")
}

func (s *Service) saveAPIKeyProfile(ctx context.Context, provider, profileID, apiKey string) error {
	agentDir, err := s.defaultAgentDir()
	if err != nil {
		return err
	}
	if err := writeAuthProfileStore(ctx, filepath.Join(agentDir, "openclaw-agent.sqlite"), provider, profileID, apiKey); err != nil {
		return err
	}
	if err := s.applyAuthProfileConfig(provider, profileID, "api_key"); err != nil {
		return err
	}
	go s.refreshGatewayAuthState(context.Background())
	return nil
}

func writeAuthProfileStore(ctx context.Context, databasePath, provider, profileID, apiKey string) error {
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o700); err != nil {
		return err
	}
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.ExecContext(ctx, "PRAGMA busy_timeout = 5000"); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS auth_profile_store (
  store_key TEXT NOT NULL PRIMARY KEY,
  store_json TEXT NOT NULL,
  updated_at INTEGER NOT NULL
)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS auth_profile_state (
  state_key TEXT NOT NULL PRIMARY KEY,
  state_json TEXT NOT NULL,
  updated_at INTEGER NOT NULL
)`); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var store authProfileStore
	var raw string
	err = tx.QueryRowContext(ctx, `SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`).Scan(&raw)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		store.Version = 1
		store.Profiles = map[string]map[string]any{}
	case err != nil:
		return err
	default:
		if err := json.Unmarshal([]byte(raw), &store); err != nil {
			return fmt.Errorf("%w: auth profile store: %v", ErrInvalidResponse, err)
		}
		if store.Version == 0 {
			store.Version = 1
		}
		if store.Profiles == nil {
			store.Profiles = map[string]map[string]any{}
		}
	}
	store.Profiles[profileID] = map[string]any{"type": "api_key", "provider": provider, "key": apiKey}
	payload, err := json.Marshal(store)
	if err != nil {
		return err
	}
	updatedAt := time.Now().UnixMilli()
	if _, err := tx.ExecContext(ctx, `INSERT INTO auth_profile_store (store_key, store_json, updated_at)
VALUES ('primary', ?, ?)
ON CONFLICT(store_key) DO UPDATE SET store_json = excluded.store_json, updated_at = excluded.updated_at`, string(payload), updatedAt); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	_ = os.Chmod(databasePath, 0o600)
	return nil
}

type authProfileStore struct {
	Version    int                       `json:"version"`
	Profiles   map[string]map[string]any `json:"profiles"`
	Order      map[string][]string       `json:"order,omitempty"`
	LastGood   map[string]string         `json:"lastGood,omitempty"`
	UsageStats map[string]any            `json:"usageStats,omitempty"`
}

func (s *Service) applyAuthProfileConfig(provider, profileID, mode string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	auth := ensureObject(config, "auth")
	profiles := ensureObject(auth, "profiles")
	profiles[profileID] = map[string]any{"provider": provider, "mode": mode}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func webSearchProviderForModelProvider(provider string) (string, bool) {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "minimax":
		return "minimax", true
	case "xai":
		return "grok", true
	case "google":
		return "gemini", true
	default:
		return "", false
	}
}

func (s *Service) applyWebSearchProviderConfig(provider, webSearchProvider, apiKey string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	tools := ensureObject(config, "tools")
	web := ensureObject(tools, "web")
	search := ensureObject(web, "search")
	search["enabled"] = true
	search["provider"] = webSearchProvider
	if _, ok := search["maxResults"]; !ok {
		search["maxResults"] = 10
	}
	if _, ok := search["timeoutSeconds"]; !ok {
		search["timeoutSeconds"] = 30
	}

	plugins := ensureObject(config, "plugins")
	entries := ensureObject(plugins, "entries")
	plugin := ensureObject(entries, provider)
	plugin["enabled"] = true
	pluginConfig := ensureObject(plugin, "config")
	webSearch := ensureObject(pluginConfig, "webSearch")
	webSearch["apiKey"] = apiKey
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err := writeFileAtomic(s.configPath, append(data, '\n'), 0o600); err != nil {
		return err
	}
	go s.refreshGatewayAuthState(context.Background())
	return nil
}

func (s *Service) applyDefaultModelConfig(modelID string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	agents := ensureObject(config, "agents")
	defaults := ensureObject(agents, "defaults")
	defaults["model"] = map[string]any{"primary": modelID}
	defaultModels := ensureObject(defaults, "models")
	if _, ok := defaultModels[modelID]; !ok {
		defaultModels[modelID] = map[string]any{}
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func (s *Service) applyModelSelectionConfig(defaultModel string, selectedModels []string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	agents := ensureObject(config, "agents")
	defaults := ensureObject(agents, "defaults")
	defaults["model"] = map[string]any{"primary": defaultModel}
	defaultModels := make(map[string]any, len(selectedModels))
	for _, modelID := range selectedModels {
		modelID = strings.TrimSpace(modelID)
		if modelID != "" {
			defaultModels[modelID] = map[string]any{}
		}
	}
	defaultModels[defaultModel] = map[string]any{}
	defaults["models"] = defaultModels
	ensureConfiguredProviderModels(config, selectedModels)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func (s *Service) ensureAgentDefaultsTimeout(seconds int) error {
	if s.configPath == "" || seconds <= 0 {
		return nil
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	agents := ensureObject(config, "agents")
	defaults := ensureObject(agents, "defaults")
	if configuredAgentTimeoutSeconds(defaults["timeoutSeconds"]) >= seconds {
		return nil
	}
	defaults["timeoutSeconds"] = seconds
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func configuredAgentTimeoutSeconds(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		parsed, _ := typed.Int64()
		return int(parsed)
	case string:
		parsed, _ := strconv.Atoi(strings.TrimSpace(typed))
		return parsed
	default:
		return 0
	}
}

// ensureConfiguredProviderModels bridges the application catalog and OpenClaw's
// provider configuration. Some OpenClaw providers expose their model catalog
// only after a model row is present in models.providers. MiniMax is one such
// provider in the pinned version, so saving a catalog model must also register
// its provider-side definition.
func ensureConfiguredProviderModels(config map[string]any, selectedModels []string) {
	modelsObject := ensureObject(config, "models")
	providers := ensureObject(modelsObject, "providers")
	for _, modelID := range selectedModels {
		provider, modelName, ok := strings.Cut(strings.TrimSpace(modelID), "/")
		if !ok || provider != "minimax" {
			continue
		}
		catalogModel, found := fixedCatalogModel(modelID)
		if !found {
			continue
		}
		providerObject := ensureObject(providers, provider)
		existing, _ := providerObject["models"].([]any)
		alreadyPresent := false
		for _, item := range existing {
			entry, _ := item.(map[string]any)
			if strings.TrimSpace(stringValue(entry["id"])) == modelName {
				alreadyPresent = true
				break
			}
		}
		if alreadyPresent {
			continue
		}
		input := []any{"text"}
		if strings.Contains(catalogModel.Input, "image") {
			input = append(input, "image")
		}
		existing = append(existing, map[string]any{
			"api":           "anthropic-messages",
			"id":            modelName,
			"name":          catalogModel.Name,
			"input":         input,
			"contextWindow": catalogModel.ContextWindow,
			"reasoning":     true,
		})
		providerObject["models"] = existing
	}
}

func fixedCatalogModel(modelID string) (Model, bool) {
	for _, model := range fixedModelCatalog {
		if model.Key == modelID {
			return model, true
		}
	}
	return Model{}, false
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}

func (s *Service) applyProviderBaseURLConfig(provider, baseURL string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if len(strings.TrimSpace(string(raw))) == 0 {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	models := ensureObject(config, "models")
	providers := ensureObject(models, "providers")
	entry := ensureObject(providers, provider)
	if baseURL == "" {
		delete(entry, "baseUrl")
	} else {
		entry["baseUrl"] = baseURL
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func readAuthProfileAPIKey(ctx context.Context, databasePath, provider, profileID string) (string, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return "", err
	}
	defer db.Close()
	var raw string
	if err := db.QueryRowContext(ctx, `SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`).Scan(&raw); err != nil {
		return "", err
	}
	var store authProfileStore
	if err := json.Unmarshal([]byte(raw), &store); err != nil {
		return "", fmt.Errorf("%w: auth profile store: %v", ErrInvalidResponse, err)
	}
	profile, ok := store.Profiles[profileID]
	if !ok {
		return "", ErrInvalidAPIKey
	}
	if storedProvider, _ := profile["provider"].(string); strings.TrimSpace(storedProvider) != "" && strings.TrimSpace(storedProvider) != provider {
		return "", ErrInvalidProvider
	}
	apiKey, _ := profile["key"].(string)
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", ErrInvalidAPIKey
	}
	return apiKey, nil
}

func ensureObject(parent map[string]any, key string) map[string]any {
	if existing, ok := parent[key].(map[string]any); ok {
		return existing
	}
	next := map[string]any{}
	parent[key] = next
	return next
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func (s *Service) refreshGatewayAuthState(ctx context.Context) {
	if s.bin == "" || !isExecutable(s.bin) {
		return
	}
	refreshCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, _ = s.run(refreshCtx, nil, "gateway", "call", "models.authStatus", "--params", `{"refresh":true}`)
}

func (s *Service) defaultAgentDir() (string, error) {
	stateDir, err := s.stateDir()
	if err != nil {
		return "", err
	}
	if s.configPath == "" {
		return filepath.Join(stateDir, "agents", "main", "agent"), nil
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return "", err
	}
	var config struct {
		Agents struct {
			List []struct {
				ID       string `json:"id"`
				Default  bool   `json:"default"`
				AgentDir string `json:"agentDir"`
			} `json:"list"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		return "", fmt.Errorf("%w: openclaw config: %v", ErrInvalidResponse, err)
	}
	if len(config.Agents.List) == 0 {
		return filepath.Join(stateDir, "agents", "main", "agent"), nil
	}
	chosen := config.Agents.List[0]
	for _, agent := range config.Agents.List {
		if agent.Default {
			chosen = agent
			break
		}
	}
	if strings.TrimSpace(chosen.AgentDir) != "" {
		return expandUserPath(strings.TrimSpace(chosen.AgentDir)), nil
	}
	id := strings.TrimSpace(strings.ToLower(chosen.ID))
	if id == "" {
		id = "main"
	}
	return filepath.Join(stateDir, "agents", id, "agent"), nil
}

func (s *Service) stateDir() (string, error) {
	if value := strings.TrimSpace(os.Getenv("OPENCLAW_STATE_DIR")); value != "" {
		return expandUserPath(value), nil
	}
	if s.configPath != "" {
		return filepath.Dir(s.configPath), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".openclaw"), nil
}

func expandUserPath(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
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
	if err := s.ensureAgentDefaultsTimeout(agentMessageTimeoutSeconds); err != nil {
		return AgentMessageResult{}, err
	}
	release, err := s.acquireMessageSlot(ctx)
	if err != nil {
		return AgentMessageResult{}, err
	}
	defer release()

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

	args := []string{"agent", "--agent", input.AgentID}
	model := s.agentMessageModelOverride(input.Model)
	if model != "" {
		if !modelIDPattern.MatchString(model) {
			return AgentMessageResult{}, ErrInvalidModel
		}
		args = append(args, "--model", model)
	}
	args = append(args, "--message-file", messagePath, "--session-key", input.SessionKey, "--json", "--timeout", strconv.Itoa(agentMessageTimeoutSeconds))
	out, err := s.run(ctx, nil, args...)
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

func (s *Service) acquireMessageSlot(ctx context.Context) (func(), error) {
	if s == nil || s.messageGate == nil {
		return func() {}, nil
	}
	select {
	case s.messageGate <- struct{}{}:
		return func() { <-s.messageGate }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// SessionHistory reads persisted messages from OpenClaw's JSONL session store.
// sessionKey is the application session key, for example
// "sta100-export-agent". The caller may pass the stage-specific key as well.
func (s *Service) SessionHistory(ctx context.Context, agentID, sessionKey string, limit int) ([]SessionHistoryMessage, error) {
	agentID = strings.TrimSpace(strings.ToLower(agentID))
	sessionKey = strings.TrimSpace(sessionKey)
	if !agentIDPattern.MatchString(agentID) {
		return nil, ErrInvalidAgent
	}
	if !sessionPattern.MatchString(sessionKey) {
		return nil, ErrInvalidSession
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	records, err := s.listSessionRecords(ctx, agentID)
	if err != nil {
		return nil, err
	}
	targets := map[string]bool{
		sessionKey: true,
		strings.TrimSuffix(sessionKey, "-coordinator") + "-coordinator": true,
	}
	var selected sessionRecord
	found := false
	for _, record := range records {
		if sessionRecordMatches(record.Key, agentID, targets) {
			if !found || record.UpdatedAt > selected.UpdatedAt {
				selected = record
				found = true
			}
		}
	}
	if !found {
		return []SessionHistoryMessage{}, nil
	}
	path, err := s.safeSessionFilePath(agentID, selected.SessionFile)
	if err != nil {
		return nil, err
	}
	messages, err := readSessionHistoryFile(path, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

type sessionRecord struct {
	Key         string `json:"key"`
	UpdatedAt   int64  `json:"updatedAt"`
	SessionFile string `json:"sessionFile"`
}

func (s *Service) listSessionRecords(ctx context.Context, agentID string) ([]sessionRecord, error) {
	if stateDir, stateErr := s.stateDir(); stateErr == nil {
		indexPath := filepath.Join(stateDir, "agents", agentID, "sessions", "sessions.json")
		if data, readErr := os.ReadFile(indexPath); readErr == nil {
			if records, parseErr := parseSessionIndex(data); parseErr == nil {
				return records, nil
			}
		}
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	out, err := s.run(ctx, nil, "sessions", "--agent", agentID, "--json", "--limit", "500")
	if err != nil {
		return nil, err
	}
	var response struct {
		Sessions []sessionRecord `json:"sessions"`
	}
	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("%w: sessions list: %v", ErrInvalidResponse, err)
	}
	return response.Sessions, nil
}

func parseSessionIndex(data []byte) ([]sessionRecord, error) {
	var response struct {
		Sessions []sessionRecord `json:"sessions"`
	}
	if err := json.Unmarshal(data, &response); err == nil && response.Sessions != nil {
		return response.Sessions, nil
	}
	var entries map[string]sessionRecord
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	records := make([]sessionRecord, 0, len(entries))
	for key, record := range entries {
		record.Key = key
		records = append(records, record)
	}
	return records, nil
}

func sessionRecordMatches(key, agentID string, targets map[string]bool) bool {
	key = strings.TrimSpace(key)
	prefix := "agent:" + agentID + ":"
	if strings.HasPrefix(key, prefix) {
		key = strings.TrimPrefix(key, prefix)
	}
	return targets[key]
}

func (s *Service) safeSessionFilePath(agentID, sessionFile string) (string, error) {
	sessionFile = expandUserPath(strings.TrimSpace(sessionFile))
	if sessionFile == "" {
		return "", fmt.Errorf("%w: session file is empty", ErrInvalidResponse)
	}
	root, err := s.stateDir()
	if err != nil {
		return "", err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(sessionFile)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: session file is outside OpenClaw state directory", ErrInvalidResponse)
	}
	expectedPrefix := filepath.Join("agents", agentID, "sessions") + string(filepath.Separator)
	if !strings.HasPrefix(relative, expectedPrefix) {
		return "", fmt.Errorf("%w: session file is outside the agent session directory", ErrInvalidResponse)
	}
	return path, nil
}

func readSessionHistoryFile(path string, limit int) ([]SessionHistoryMessage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read OpenClaw session: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	messages := make([]SessionHistoryMessage, 0, limit)
	for {
		var event struct {
			Type      string `json:"type"`
			Timestamp string `json:"timestamp"`
			Message   struct {
				Role         string          `json:"role"`
				Content      json.RawMessage `json:"content"`
				Provider     string          `json:"provider"`
				Model        string          `json:"model"`
				RunID        string          `json:"runId"`
				StopReason   string          `json:"stopReason"`
				ErrorMessage string          `json:"errorMessage"`
			} `json:"message"`
		}
		if err := decoder.Decode(&event); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("%w: session JSONL: %v", ErrInvalidResponse, err)
		}
		if event.Type != "message" {
			continue
		}
		role := strings.TrimSpace(event.Message.Role)
		if role != "user" && role != "assistant" {
			continue
		}
		text := extractSessionText(event.Message.Content)
		if role == "user" {
			text = extractOriginalUserMessage(text)
		}
		if text == "" && event.Message.ErrorMessage != "" {
			text = "OpenClaw 调用失败：" + strings.TrimSpace(event.Message.ErrorMessage)
		}
		if text == "" {
			continue
		}
		item := SessionHistoryMessage{
			Role: role, Text: text, CreatedAt: event.Timestamp,
			Model: event.Message.Model, Provider: event.Message.Provider,
			RunID: event.Message.RunID, Error: event.Message.ErrorMessage,
		}
		if role == "assistant" && item.Error != "" {
			item.Text = "OpenClaw 调用失败：" + strings.TrimSpace(item.Error)
		}
		messages = append(messages, item)
	}
	if len(messages) > limit {
		messages = append([]SessionHistoryMessage(nil), messages[len(messages)-limit:]...)
	}
	return messages, nil
}

func extractSessionText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return strings.TrimSpace(text)
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) == nil {
		parts := make([]string, 0, len(blocks))
		for _, block := range blocks {
			if block.Type == "" || block.Type == "text" {
				if value := strings.TrimSpace(block.Text); value != "" {
					parts = append(parts, value)
				}
			}
		}
		return strings.TrimSpace(strings.Join(parts, "\n"))
	}
	var object struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &object) == nil {
		return strings.TrimSpace(object.Text)
	}
	return ""
}

func extractOriginalUserMessage(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	const marker = "[用户消息]"
	if index := strings.Index(text, marker); index >= 0 {
		candidate := strings.TrimSpace(text[index+len(marker):])
		if jsonStart := strings.Index(candidate, "{"); jsonStart >= 0 {
			var envelope struct {
				UserMessage string `json:"userMessage"`
			}
			if err := json.Unmarshal([]byte(candidate[jsonStart:]), &envelope); err == nil && strings.TrimSpace(envelope.UserMessage) != "" {
				return strings.TrimSpace(envelope.UserMessage)
			}
		}
		return candidate
	}
	return text
}

func (s *Service) agentMessageModelOverride(model string) string {
	model = strings.TrimSpace(model)
	if model == "" {
		return ""
	}
	snapshot := s.ModelSnapshot()
	defaultModel := strings.TrimSpace(firstNonEmpty(snapshot.ResolvedDefault, snapshot.DefaultModel))
	if defaultModel != "" && model == defaultModel {
		return ""
	}
	return model
}

func buildAgentMessage(input AgentMessageInput) string {
	var builder strings.Builder
	builder.WriteString("[STA-100 本次来源约束]\n允许来源类别：")
	builder.WriteString(strings.Join(input.Sources, "、"))
	if len(input.Allowlist) > 0 {
		builder.WriteString("\n联网检索仅允许以下域名：")
		builder.WriteString(strings.Join(input.Allowlist, "、"))
	}
	if len(input.Attachments) > 0 {
		builder.WriteString("\n本次附件（由 Go 应用保存于本机临时目录，Agent 需要使用可用文件读取能力读取；如无法读取必须明确说明，不得假设已查看）：")
		for _, attachment := range input.Attachments {
			builder.WriteString("\n- ")
			builder.WriteString(attachment.Name)
			builder.WriteString(" | ")
			builder.WriteString(attachment.Path)
			if attachment.Mime != "" {
				builder.WriteString(" | ")
				builder.WriteString(attachment.Mime)
			}
		}
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
	command.Env = s.commandEnv()
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

func (s *Service) commandEnv() []string {
	env := os.Environ()
	if s.configPath != "" && os.Getenv("OPENCLAW_CONFIG_PATH") == "" {
		env = append(env, "OPENCLAW_CONFIG_PATH="+s.configPath)
	}
	if os.Getenv("OPENCLAW_GATEWAY_TOKEN") == "" {
		if token := s.gatewayTokenFromConfig(); token != "" {
			env = append(env, "OPENCLAW_GATEWAY_TOKEN="+token)
		}
	}
	return env
}

func (s *Service) gatewayTokenFromConfig() string {
	if s.configPath == "" {
		return ""
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return ""
	}
	var config struct {
		Gateway struct {
			Auth struct {
				Mode  string `json:"mode"`
				Token string `json:"token"`
			} `json:"auth"`
		} `json:"gateway"`
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		return ""
	}
	if strings.EqualFold(strings.TrimSpace(config.Gateway.Auth.Mode), "token") {
		return strings.TrimSpace(config.Gateway.Auth.Token)
	}
	return ""
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
