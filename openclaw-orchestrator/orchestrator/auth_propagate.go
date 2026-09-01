package orchestrator

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// readPortableAuthProfiles opens an openclaw-agent.sqlite and returns the
// map of portable static api_key profiles that can be safely copied into
// another agent's auth store. Non-portable profile types (OAuth/token etc.)
// are skipped because their validity is bound to the issuing agent's state.
func readPortableAuthProfiles(databasePath string) (map[string]map[string]any, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var raw string
	row := db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`)
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("read auth profile store: %w", err)
	}
	var envelope struct {
		Version  int                       `json:"version"`
		Profiles map[string]map[string]any `json:"profiles"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil, fmt.Errorf("parse auth profile store: %w", err)
	}
	result := map[string]map[string]any{}
	for profileID, profile := range envelope.Profiles {
		if profile == nil {
			continue
		}
		profileType := strings.ToLower(strings.TrimSpace(toString(profile["type"])))
		provider := strings.ToLower(strings.TrimSpace(toString(profile["provider"])))
		if provider == "" {
			if idx := strings.Index(profileID, ":"); idx > 0 {
				provider = strings.ToLower(strings.TrimSpace(profileID[:idx]))
			}
		}
		if profileType != "" && profileType != "api_key" {
			continue
		}
		key := strings.TrimSpace(toString(profile["key"]))
		if key == "" {
			continue
		}
		result[profileID] = map[string]any{
			"type":     "api_key",
			"provider": provider,
			"key":      key,
		}
	}
	return result, nil
}

// writePortableAuthProfiles creates or updates an openclaw-agent.sqlite at
// databasePath and merges the supplied portable profiles into the existing
// store. Existing non-portable profiles are preserved; matching static API key
// profiles are overwritten so saving a new key updates every subagent.
func writePortableAuthProfiles(databasePath string, profiles map[string]map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o700); err != nil {
		return err
	}
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS auth_profile_store (
  store_key TEXT NOT NULL PRIMARY KEY,
  store_json TEXT NOT NULL,
  updated_at INTEGER NOT NULL
)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS auth_profile_state (
  state_key TEXT NOT NULL PRIMARY KEY,
  state_json TEXT NOT NULL,
  updated_at INTEGER NOT NULL
)`); err != nil {
		return err
	}
	store := authProfileStore{Version: 1, Profiles: map[string]map[string]any{}}
	var raw string
	err = db.QueryRow(`SELECT store_json FROM auth_profile_store WHERE store_key = 'primary'`).Scan(&raw)
	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return err
	default:
		if err := json.Unmarshal([]byte(raw), &store); err != nil {
			return fmt.Errorf("parse existing auth profile store: %w", err)
		}
		if store.Version == 0 {
			store.Version = 1
		}
		if store.Profiles == nil {
			store.Profiles = map[string]map[string]any{}
		}
	}
	for profileID, profile := range profiles {
		if strings.TrimSpace(profileID) == "" || profile == nil {
			continue
		}
		store.Profiles[profileID] = profile
	}
	payload, err := json.Marshal(store)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO auth_profile_store (store_key, store_json, updated_at)
VALUES ('primary', ?, ?)
ON CONFLICT(store_key) DO UPDATE SET store_json = excluded.store_json, updated_at = excluded.updated_at`,
		string(payload), time.Now().UnixMilli())
	if err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	_ = os.Chmod(databasePath, 0o600)
	return nil
}

// toString converts an arbitrary map value into a trimmed string.
// Used to defensively read provider/type/key fields that may be encoded as
// either strings or numbers in some upstream code paths.
func toString(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}
