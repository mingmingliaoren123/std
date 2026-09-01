package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AgentManifest struct {
	SchemaVersion int             `json:"schema_version"`
	WorkspaceRoot string          `json:"workspace_root,omitempty"`
	Agents        []ManifestAgent `json:"agents"`
}

type ManifestAgent struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TechnicalName string `json:"technical_name"`
	Emoji         string `json:"emoji"`
	Workspace     string `json:"workspace"`
	Model         string `json:"model,omitempty"`
	Role          string `json:"role,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
	Instructions  string `json:"instructions,omitempty"`
}

func LoadManifest(path string) (AgentManifest, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return AgentManifest{}, fmt.Errorf("%w: %v", ErrManifestMissing, err)
	}
	var manifest AgentManifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return AgentManifest{}, fmt.Errorf("invalid agent manifest: %w", err)
	}
	if len(manifest.Agents) == 0 {
		return AgentManifest{}, errors.New("agent manifest contains no agents")
	}
	for _, agent := range manifest.Agents {
		if strings.TrimSpace(agent.ID) == "" || strings.TrimSpace(agent.Name) == "" || strings.TrimSpace(agent.Workspace) == "" {
			return AgentManifest{}, fmt.Errorf("agent manifest has incomplete agent %q", agent.ID)
		}
	}
	return manifest, nil
}

func (s *Service) decorateAgents(agents []Agent) []Agent {
	manifest, err := LoadManifest(s.manifestPath)
	if err != nil {
		return agents
	}
	metadata := make(map[string]ManifestAgent, len(manifest.Agents))
	for _, item := range manifest.Agents {
		metadata[item.ID] = item
	}
	for index := range agents {
		desired, ok := metadata[agents[index].ID]
		if !ok {
			continue
		}
		agents[index].Role = firstManifestValue(desired.Role, "domain")
		agents[index].Visibility = firstManifestValue(desired.Visibility, "business")
	}
	return agents
}

func firstManifestValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func (s *Service) SyncAgents(ctx context.Context) ([]Agent, error) {
	if s.manifestPath == "" {
		return nil, ErrManifestMissing
	}
	manifest, err := LoadManifest(s.manifestPath)
	if err != nil {
		return nil, err
	}
	base := filepath.Dir(s.manifestPath)
	workspaceRoot := base
	if manifest.WorkspaceRoot != "" {
		workspaceRoot = manifest.WorkspaceRoot
		if !filepath.IsAbs(workspaceRoot) {
			workspaceRoot = filepath.Join(base, workspaceRoot)
		}
		workspaceRoot = filepath.Clean(workspaceRoot)
	}
	for _, desired := range manifest.Agents {
		workspace, err := resolveWorkspace(workspaceRoot, desired.Workspace)
		if err != nil {
			return nil, fmt.Errorf("agent %s workspace: %w", desired.ID, err)
		}
		if err := os.MkdirAll(workspace, 0o755); err != nil {
			return nil, fmt.Errorf("create agent %s workspace: %w", desired.ID, err)
		}
		identityPath := filepath.Join(workspace, "IDENTITY.md")
		if _, err := os.Stat(identityPath); errors.Is(err, os.ErrNotExist) {
			content := fmt.Sprintf("# %s\n\n- Technical name: %s\n- Emoji: %s\n- Role: cycling industry business agent\n", desired.Name, desired.TechnicalName, desired.Emoji)
			if err := os.WriteFile(identityPath, []byte(content), 0o644); err != nil {
				return nil, fmt.Errorf("write agent %s identity: %w", desired.ID, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("inspect agent %s identity: %w", desired.ID, err)
		}
		if strings.TrimSpace(desired.Instructions) != "" {
			instructions := strings.TrimSpace(desired.Instructions) + "\n"
			if err := os.WriteFile(filepath.Join(workspace, "AGENTS.md"), []byte(instructions), 0o644); err != nil {
				return nil, fmt.Errorf("write agent %s instructions: %w", desired.ID, err)
			}
		}
	}
	if err := s.applyManifestAgentsConfig(manifest, workspaceRoot); err != nil {
		return nil, err
	}
	// 把 main agent 的便携静态 auth profile（api_key 等）复制到本次新增的 subagent 目录里，
	// 否则子 Agent 找不到 minimax 等 provider 的 API Key 会报 ProviderAuthError。
	if err := s.propagateMainAuthToSubagents(); err != nil {
		// 这里不返回错误，auth 复制失败不应阻塞 Agent 同步本身。
		// 真实场景下用户也可以在页面上重新保存一次 API Key 来主动重发。
		_ = err
	}
	return s.ListAgents(ctx)
}

// propagateMainAuthToSubagents reads the main agent's openclaw-agent.sqlite
// and merges its portable static api_key profiles into each subagent dir.
// Non-portable profiles (OAuth tokens etc.) are skipped per OpenClaw guidance:
// "copy only portable static auth profiles from the main agentDir".
//
// Subagent targets are derived from the configured openclaw.json agents.list
// (falls back to the on-disk state dir if the config is unreadable). The
// function also creates the <state>/agents/<id>/agent directory if missing,
// so a fresh sync without prior agent invocations still seeds the stores.
func (s *Service) propagateMainAuthToSubagents() error {
	mainDir, err := s.defaultAgentDir()
	if err != nil {
		return err
	}
	mainDB := filepath.Join(mainDir, "openclaw-agent.sqlite")
	if _, statErr := os.Stat(mainDB); statErr != nil {
		return nil
	}
	profiles, err := readPortableAuthProfiles(mainDB)
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		return nil
	}
	targets, err := s.subAgentDirs()
	if err != nil {
		return err
	}
	for _, targetDir := range targets {
		if targetDir == "" {
			continue
		}
		targetDB := filepath.Join(targetDir, "openclaw-agent.sqlite")
		if err := os.MkdirAll(targetDir, 0o700); err != nil {
			continue
		}
		if err := writePortableAuthProfiles(targetDB, profiles); err != nil {
			continue
		}
		_ = os.Chmod(targetDB, 0o600)
	}
	return nil
}

// subAgentDirs returns the agentDir for every agent in openclaw.json except
// the default agent. Falls back to walking <stateDir>/agents when the config
// is missing or unreadable, so a stale state can still be reseeded.
func (s *Service) subAgentDirs() ([]string, error) {
	if s.configPath != "" {
		if dirs, err := subAgentDirsFromConfig(s.configPath); err == nil && len(dirs) > 0 {
			return dirs, nil
		}
	}
	stateDir, err := s.stateDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(stateDir, "agents"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "main" {
			continue
		}
		out = append(out, filepath.Join(stateDir, "agents", entry.Name(), "agent"))
	}
	return out, nil
}

func subAgentDirsFromConfig(configPath string) ([]string, error) {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config struct {
		Agents struct {
			Defaults struct {
				AgentDir string `json:"agentDir"`
			} `json:"defaults"`
			List []struct {
				ID       string `json:"id"`
				AgentDir string `json:"agentDir"`
				Default  bool   `json:"default"`
			} `json:"list"`
		} `json:"agents"`
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, err
	}
	defaultAgentDir := strings.TrimSpace(config.Agents.Defaults.AgentDir)
	out := make([]string, 0, len(config.Agents.List))
	for _, item := range config.Agents.List {
		dir := strings.TrimSpace(item.AgentDir)
		if dir == "" {
			continue
		}
		// Skip the default agent — that's where the source profile lives.
		if defaultAgentDir != "" && expandUserPath(dir) == expandUserPath(defaultAgentDir) {
			continue
		}
		if item.Default && defaultAgentDir == "" {
			continue
		}
		out = append(out, expandUserPath(dir))
	}
	return out, nil
}

func (s *Service) applyManifestAgentsConfig(manifest AgentManifest, workspaceRoot string) error {
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
	defaultModel := parseModelReference(rawMessageFromValue(defaults["model"]))
	if defaultModel == "" {
		if models, ok := defaults["models"].(map[string]any); ok {
			for model := range models {
				defaultModel = model
				break
			}
		}
	}
	stateDir, err := s.stateDir()
	if err != nil {
		return err
	}
	existingList, _ := agents["list"].([]any)
	merged := make([]map[string]any, 0, len(existingList)+len(manifest.Agents))
	seen := map[string]bool{}
	for _, item := range existingList {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id := strings.TrimSpace(stringValue(entry["id"]))
		if id == "" || seen[id] {
			continue
		}
		if id != "main" {
			continue
		}
		merged = append(merged, entry)
		seen[id] = true
	}
	for _, desired := range manifest.Agents {
		workspace, err := resolveWorkspace(workspaceRoot, desired.Workspace)
		if err != nil {
			return fmt.Errorf("agent %s workspace: %w", desired.ID, err)
		}
		id := strings.TrimSpace(desired.ID)
		if id == "" || seen[id] {
			continue
		}
		agentDir := filepath.Join(stateDir, "agents", id, "agent")
		entry := map[string]any{
			"id":        id,
			"name":      desired.TechnicalName,
			"workspace": workspace,
			"agentDir":  agentDir,
			"identity": map[string]any{
				"name":  desired.Name,
				"emoji": desired.Emoji,
			},
		}
		if strings.TrimSpace(desired.Model) != "" && strings.TrimSpace(desired.Model) != defaultModel {
			entry["model"] = map[string]any{"primary": strings.TrimSpace(desired.Model)}
		}
		merged = append(merged, entry)
		seen[id] = true
	}
	agents["list"] = mapSliceToAny(merged)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.configPath, append(data, '\n'), 0o600)
}

func rawMessageFromValue(value any) json.RawMessage {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return data
}

func mapSliceToAny(values []map[string]any) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

func (s *Service) runConfigMutationWithRetry(ctx context.Context, args ...string) error {
	if _, err := s.run(ctx, nil, args...); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "config changed since last load") {
			time.Sleep(250 * time.Millisecond)
			if _, retryErr := s.run(ctx, nil, args...); retryErr == nil {
				return nil
			} else {
				return retryErr
			}
		}
		return err
	}
	return nil
}

func (s *Service) listAgentsFromCLI(ctx context.Context) ([]Agent, error) {
	out, err := s.run(ctx, nil, "agents", "list", "--json")
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := json.Unmarshal(out, &agents); err != nil {
		return nil, fmt.Errorf("%w: agents list: %v", ErrInvalidResponse, err)
	}
	return agents, nil
}

func resolveWorkspace(base, workspace string) (string, error) {
	if filepath.IsAbs(workspace) {
		return filepath.Clean(workspace), nil
	}
	base = filepath.Clean(base)
	resolved := filepath.Clean(filepath.Join(base, workspace))
	relative, err := filepath.Rel(base, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("workspace must stay inside manifest directory")
	}
	return resolved, nil
}
