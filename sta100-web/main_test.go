package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMutationRequiresSTA100Header(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/openclaw/agents/sync", nil)
	recorder := httptest.NewRecorder()
	if allowMutation(recorder, request, http.MethodPost) {
		t.Fatal("mutation was allowed without X-STA100-Request")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestMutationAllowsSTA100Header(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/openclaw/agents/sync", nil)
	request.Header.Set("X-STA100-Request", "1")
	recorder := httptest.NewRecorder()
	if !allowMutation(recorder, request, http.MethodPost) {
		t.Fatalf("mutation rejected with required header: status %d", recorder.Code)
	}
}

func TestAuthManagerDefaultsAndPersistsHashedPassword(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth", "auth.json")
	t.Setenv("STA100_AUTH_FILE", path)
	manager, err := newAuthManager()
	if err != nil {
		t.Fatal(err)
	}
	if manager.username != "admin" || !verifyPassword("admin", manager.password) {
		t.Fatal("default admin/admin credentials were not initialized")
	}
	if manager.password == "admin" {
		t.Fatal("password was stored as plaintext")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("auth file permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestAuthManagerUpdatesCredentials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("STA100_AUTH_FILE", path)
	manager, err := newAuthManager()
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.updateCredentials("operator", "cycle2026"); err != nil {
		t.Fatal(err)
	}
	reloaded, err := newAuthManager()
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.username != "operator" || !verifyPassword("cycle2026", reloaded.password) {
		t.Fatal("updated credentials were not persisted")
	}
	if verifyPassword("admin", reloaded.password) {
		t.Fatal("old password remained valid")
	}
}
