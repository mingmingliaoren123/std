package orchestrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// TestPortableAuthProfileRoundTrip covers the read/write helpers used by
// propagateMainAuthToSubagents: a profile written through writeAuthProfileStore
// must be readable by readPortableAuthProfiles and re-writable to a fresh
// database, preserving the provider / api_key fields. Non-portable OAuth
// entries are intentionally dropped.
func TestPortableAuthProfileRoundTrip(t *testing.T) {
	mainDir := t.TempDir()
	subDir := t.TempDir()
	mainDB := filepath.Join(mainDir, "openclaw-agent.sqlite")
	subDB := filepath.Join(subDir, "openclaw-agent.sqlite")
	if err := writeAuthProfileStore(context.Background(), mainDB, "minimax", "minimax:default", "test-key-12345678"); err != nil {
		t.Fatalf("seed main store: %v", err)
	}
	// Inject a non-portable OAuth entry directly to ensure the reader skips it.
	if err := appendOAuthProfile(t, mainDB, "minimax:oauth", "minimax"); err != nil {
		t.Fatalf("inject oauth profile: %v", err)
	}
	profiles, err := readPortableAuthProfiles(mainDB)
	if err != nil {
		t.Fatalf("read main store: %v", err)
	}
	if got := profiles["minimax:default"]; got == nil || got["key"] != "test-key-12345678" || got["provider"] != "minimax" {
		t.Fatalf("missing or wrong portable profile: %#v", profiles)
	}
	if _, ok := profiles["minimax:oauth"]; ok {
		t.Fatalf("non-portable oauth profile leaked into portable set: %#v", profiles)
	}
	if err := writePortableAuthProfiles(subDB, profiles); err != nil {
		t.Fatalf("write sub store: %v", err)
	}
	if info, err := os.Stat(subDB); err != nil {
		t.Fatalf("sub db missing: %v", err)
	} else if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("sub db permission = %o, want 0o600", perm)
	}
	again, err := readPortableAuthProfiles(subDB)
	if err != nil {
		t.Fatalf("read sub store: %v", err)
	}
	if again["minimax:default"]["key"] != "test-key-12345678" {
		t.Fatalf("round-tripped key mismatch: %#v", again)
	}
	if _, ok := again["minimax:oauth"]; ok {
		t.Fatalf("oauth profile survived propagation: %#v", again)
	}
}

// TestPropagateMainAuthToSubagentsMergesExisting verifies that propagation
// seeds missing stores and updates stale static API key profiles in existing
// stores. Non-portable OAuth/token state is still preserved.
func TestPropagateMainAuthToSubagentsMergesExisting(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state")
	mainAgentDir := filepath.Join(stateDir, "agents", "main", "agent")
	if err := os.MkdirAll(mainAgentDir, 0o700); err != nil {
		t.Fatal(err)
	}
	mainDB := filepath.Join(mainAgentDir, "openclaw-agent.sqlite")
	if err := writeAuthProfileStore(context.Background(), mainDB, "minimax", "minimax:default", "main-key-12345678"); err != nil {
		t.Fatalf("seed main store: %v", err)
	}

	alphaDir := filepath.Join(stateDir, "agents", "alpha", "agent")
	bravoDir := filepath.Join(stateDir, "agents", "bravo", "agent")
	if err := os.MkdirAll(alphaDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(bravoDir, 0o700); err != nil {
		t.Fatal(err)
	}
	bravoDB := filepath.Join(bravoDir, "openclaw-agent.sqlite")
	if err := writeAuthProfileStore(context.Background(), bravoDB, "minimax", "minimax:default", "bravo-key-87654321"); err != nil {
		t.Fatalf("seed bravo store: %v", err)
	}
	if err := appendOAuthProfile(t, bravoDB, "feishu:oauth", "feishu"); err != nil {
		t.Fatalf("seed bravo oauth profile: %v", err)
	}

	configPath := filepath.Join(stateDir, "openclaw.json")
	if err := os.WriteFile(configPath, []byte(`{"agents":{"list":[{"id":"main"}]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	service := New(Config{ConfigPath: configPath, StateDir: stateDir})
	t.Setenv("OPENCLAW_STATE_DIR", stateDir)
	if err := service.propagateMainAuthToSubagents(); err != nil {
		t.Fatalf("propagate: %v", err)
	}

	alphaDB := filepath.Join(alphaDir, "openclaw-agent.sqlite")
	if _, err := os.Stat(alphaDB); err != nil {
		t.Fatalf("expected alpha db to be created: %v", err)
	}
	got, err := readPortableAuthProfiles(alphaDB)
	if err != nil {
		t.Fatalf("read alpha: %v", err)
	}
	if got["minimax:default"]["key"] != "main-key-12345678" {
		t.Fatalf("alpha copy missing api key: %#v", got)
	}
	updated, err := readPortableAuthProfiles(bravoDB)
	if err != nil {
		t.Fatalf("read bravo: %v", err)
	}
	if updated["minimax:default"]["key"] != "main-key-12345678" {
		t.Fatalf("bravo static api key was not refreshed: %#v", updated)
	}
	if !authProfileExists(t, bravoDB, "feishu:oauth") {
		t.Fatalf("bravo oauth profile was not preserved")
	}
}

// TestPropagateMainAuthSeedsSubagentsFromOpenclawConfig exercises the real
// post-sync state: openclaw.json lists subagents whose on-disk dirs do not
// exist yet. Propagation must mkdir each one and drop a sqlite store there.
func TestPropagateMainAuthSeedsSubagentsFromOpenclawConfig(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state")
	mainAgentDir := filepath.Join(stateDir, "agents", "main", "agent")
	if err := os.MkdirAll(mainAgentDir, 0o700); err != nil {
		t.Fatal(err)
	}
	mainDB := filepath.Join(mainAgentDir, "openclaw-agent.sqlite")
	if err := writeAuthProfileStore(context.Background(), mainDB, "minimax", "minimax:default", "main-key-12345678"); err != nil {
		t.Fatalf("seed main store: %v", err)
	}

	alphaDir := filepath.Join(stateDir, "agents", "alpha", "agent")
	bravoDir := filepath.Join(stateDir, "agents", "bravo", "agent")
	configPath := filepath.Join(stateDir, "openclaw.json")
	config := `{
  "agents": {
    "list": [
      {"id":"main","agentDir":"` + mainAgentDir + `"},
      {"id":"alpha","agentDir":"` + alphaDir + `"},
      {"id":"bravo","agentDir":"` + bravoDir + `"}
    ]
  }
}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	service := New(Config{ConfigPath: configPath, StateDir: stateDir})
	t.Setenv("OPENCLAW_STATE_DIR", stateDir)
	if err := service.propagateMainAuthToSubagents(); err != nil {
		t.Fatalf("propagate: %v", err)
	}
	for _, dir := range []string{alphaDir, bravoDir} {
		db := filepath.Join(dir, "openclaw-agent.sqlite")
		info, err := os.Stat(db)
		if err != nil {
			t.Fatalf("expected seeded db at %s: %v", db, err)
		}
		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Fatalf("seeded db permission = %o, want 0o600", perm)
		}
		got, err := readPortableAuthProfiles(db)
		if err != nil {
			t.Fatalf("read seeded db %s: %v", db, err)
		}
		if got["minimax:default"]["key"] != "main-key-12345678" {
			t.Fatalf("seeded db missing api key: %#v", got)
		}
	}
	if _, err := os.Stat(mainDB); err != nil {
		t.Fatalf("main db should remain: %v", err)
	}
}

func appendOAuthProfile(t *testing.T, databasePath, profileID, provider string) error {
	t.Helper()
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	var raw string
	if err := db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key='primary'`).Scan(&raw); err != nil {
		return err
	}
	var envelope struct {
		Version  int                       `json:"version"`
		Profiles map[string]map[string]any `json:"profiles"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return err
	}
	if envelope.Profiles == nil {
		envelope.Profiles = map[string]map[string]any{}
	}
	envelope.Profiles[profileID] = map[string]any{
		"type":     "oauth",
		"provider": provider,
		"token":    "secret-token",
	}
	updated, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE auth_profile_store SET store_json=?, updated_at=? WHERE store_key='primary'`,
		string(updated), time.Now().UnixMilli())
	return err
}

func authProfileExists(t *testing.T, databasePath, profileID string) bool {
	t.Helper()
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var raw string
	if err := db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key='primary'`).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Profiles map[string]map[string]any `json:"profiles"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		t.Fatal(err)
	}
	_, ok := envelope.Profiles[profileID]
	return ok
}
