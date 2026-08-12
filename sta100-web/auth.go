package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type authUsernameContextKey struct{}

const (
	authCookieName = "sta100_session"
	masterPassword = "sta100"
	passwordRounds = 150000
)

type authManager struct {
	mu       sync.Mutex
	path     string
	username string
	password string
	sessions map[string]time.Time
}

type authFile struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func newAuthManager() (*authManager, error) {
	path := strings.TrimSpace(os.Getenv("STA100_AUTH_FILE"))
	if path == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(configDir, "sta100", "auth.json")
	}
	manager := &authManager{path: path, sessions: make(map[string]time.Time)}
	if err := manager.load(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (a *authManager) load() error {
	data, err := os.ReadFile(a.path)
	if errors.Is(err, os.ErrNotExist) {
		hash, hashErr := hashPassword("admin")
		if hashErr != nil {
			return hashErr
		}
		a.username, a.password = "admin", hash
		return a.persist()
	}
	if err != nil {
		return err
	}
	var saved authFile
	if err := json.Unmarshal(data, &saved); err != nil {
		return err
	}
	if strings.TrimSpace(saved.Username) == "" || saved.PasswordHash == "" {
		return errors.New("STA-100 auth file is incomplete")
	}
	a.username, a.password = strings.TrimSpace(saved.Username), saved.PasswordHash
	return nil
}

func (a *authManager) persist() error {
	if err := os.MkdirAll(filepath.Dir(a.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(authFile{Username: a.username, PasswordHash: a.password}, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(a.path), ".auth-*")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryName, a.path)
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	digest := derivePassword(password, salt, passwordRounds)
	return "$sta100$" + strconv.Itoa(passwordRounds) + "$" + hex.EncodeToString(salt) + "$" + hex.EncodeToString(digest), nil
}

func verifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 || parts[1] != "sta100" {
		return false
	}
	rounds, err := strconv.Atoi(parts[2])
	salt, saltErr := hex.DecodeString(parts[3])
	expected, expectedErr := hex.DecodeString(parts[4])
	if err != nil || rounds < 100000 || rounds > 1000000 || saltErr != nil || expectedErr != nil || len(expected) != sha256.Size {
		return false
	}
	actual := derivePassword(password, salt, rounds)
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func derivePassword(password string, salt []byte, rounds int) []byte {
	input := append(append([]byte{}, salt...), []byte(password)...)
	digest := sha256.Sum256(input)
	for index := 1; index < rounds; index++ {
		block := make([]byte, 0, len(digest)+len(password))
		block = append(block, digest[:]...)
		block = append(block, []byte(password)...)
		digest = sha256.Sum256(block)
	}
	return digest[:]
}

func (a *authManager) statusHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	a.mu.Lock()
	username := a.username
	a.mu.Unlock()
	_, authenticated := a.sessionUsername(r)
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": authenticated, "username": username})
}

func (a *authManager) loginHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	a.mu.Lock()
	valid := strings.TrimSpace(request.Username) == a.username && verifyPassword(request.Password, a.password)
	username := a.username
	a.mu.Unlock()
	if !valid {
		writeAPIError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码不正确")
		return
	}
	if err := a.issueSession(w); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "SESSION_FAILED", "登录会话创建失败")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "username": username})
}

func (a *authManager) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	a.revokeSession(r)
	http.SetCookie(w, &http.Cookie{Name: authCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": false})
}

func (a *authManager) verifyMasterHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		MasterPassword string `json:"masterPassword"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.MasterPassword != masterPassword {
		writeAPIError(w, http.StatusUnauthorized, "INVALID_MASTER_PASSWORD", "万能密码不正确")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"verified": true})
}

func (a *authManager) resetHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var request struct {
		MasterPassword string `json:"masterPassword"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.MasterPassword != masterPassword {
		writeAPIError(w, http.StatusUnauthorized, "INVALID_MASTER_PASSWORD", "万能密码不正确")
		return
	}
	if err := a.updateCredentials(strings.TrimSpace(request.Username), request.Password); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ACCOUNT", err.Error())
		return
	}
	if err := a.issueSession(w); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "SESSION_FAILED", "账户重置后会话创建失败")
		return
	}
	a.mu.Lock()
	username := a.username
	a.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "username": username, "reset": true})
}

func (a *authManager) accountHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPatch) {
		return
	}
	var request struct {
		CurrentPassword string `json:"currentPassword"`
		Username        string `json:"username"`
		Password        string `json:"password"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	a.mu.Lock()
	valid := verifyPassword(request.CurrentPassword, a.password)
	a.mu.Unlock()
	if !valid {
		writeAPIError(w, http.StatusUnauthorized, "INVALID_CURRENT_PASSWORD", "当前密码不正确")
		return
	}
	if err := a.updateCredentials(strings.TrimSpace(request.Username), request.Password); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ACCOUNT", err.Error())
		return
	}
	a.mu.Lock()
	username := a.username
	a.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "username": username})
}

func (a *authManager) updateCredentials(username, password string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}
	if len([]rune(password)) < 4 {
		return errors.New("密码至少需要 4 个字符")
	}
	hash, err := hashPassword(password)
	if err != nil {
		return errors.New("密码保存失败")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	previousUsername, previousPassword := a.username, a.password
	a.username, a.password = username, hash
	if err := a.persist(); err != nil {
		a.username, a.password = previousUsername, previousPassword
		return errors.New("账户保存失败")
	}
	return nil
}

func (a *authManager) issueSession(w http.ResponseWriter) error {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)
	a.mu.Lock()
	a.sessions[token] = time.Now().Add(12 * time.Hour)
	a.mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: authCookieName, Value: token, Path: "/", MaxAge: 12 * 60 * 60, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	return nil
}

func (a *authManager) revokeSession(r *http.Request) {
	if cookie, err := r.Cookie(authCookieName); err == nil {
		a.mu.Lock()
		delete(a.sessions, cookie.Value)
		a.mu.Unlock()
	}
}

func (a *authManager) sessionUsername(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}
	a.mu.Lock()
	expires, ok := a.sessions[cookie.Value]
	if ok && time.Now().After(expires) {
		delete(a.sessions, cookie.Value)
		ok = false
	}
	username := a.username
	a.mu.Unlock()
	return username, ok
}

func (a *authManager) requireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := a.sessionUsername(r)
		if !ok {
			writeAPIError(w, http.StatusUnauthorized, "AUTH_REQUIRED", "请先登录 STA-100")
			return
		}
		ctx := context.WithValue(r.Context(), authUsernameContextKey{}, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
