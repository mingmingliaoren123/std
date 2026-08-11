package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func (s *Service) SyncAgents(ctx context.Context) ([]Agent, error) {
	if s.manifestPath == "" {
		return nil, ErrManifestMissing
	}
	manifest, err := LoadManifest(s.manifestPath)
	if err != nil {
		return nil, err
	}
	existing, err := s.listAgentsFromCLI(ctx)
	if err != nil {
		return nil, err
	}
	existingByID := make(map[string]Agent, len(existing))
	for _, agent := range existing {
		existingByID[agent.ID] = agent
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

		current, exists := existingByID[desired.ID]
		if !exists {
			args := []string{"agents", "add", desired.ID, "--workspace", workspace, "--non-interactive", "--json"}
			if desired.Model != "" {
				args = append(args, "--model", desired.Model)
			}
			if _, err := s.run(ctx, nil, args...); err != nil {
				return nil, err
			}
		}
		if !exists || current.IdentityName != desired.Name || current.IdentityEmoji != desired.Emoji {
			if _, err := s.run(ctx, nil, "agents", "set-identity", "--agent", desired.ID, "--name", desired.Name, "--emoji", desired.Emoji, "--identity-file", identityPath, "--json"); err != nil {
				return nil, err
			}
		}
	}
	return s.ListAgents(ctx)
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
