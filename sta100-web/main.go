package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"openclaw-orchestrator/orchestrator"
)

//go:embed index.html assets/*
var ui embed.FS

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Prototype bool   `json:"prototype"`
	Time      string `json:"time"`
}

// openClawService is the STA-100 HTTP adapter. OpenClaw discovery, CLI calls,
// model validation and manifest synchronization live in the reusable package.
type openClawService struct {
	service *orchestrator.Service
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	openClaw := newOpenClawService()
	auth, err := newAuthManager()
	if err != nil {
		log.Fatalf("initialize STA-100 authentication: %v", err)
	}
	store, err := newBusinessStore()
	if err != nil {
		log.Fatalf("initialize STA-100 business database: %v", err)
	}
	defer store.Close()
	business := newBusinessAPI(store, openClaw)
	files := http.FileServer(http.FS(ui))
	staticFiles := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The UI is embedded and changes are deployed by restarting this local
		// service, so stale browser assets must not survive a refresh.
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		files.ServeHTTP(w, r)
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/v1/auth/status", auth.statusHandler)
	mux.HandleFunc("/api/v1/auth/login", auth.loginHandler)
	mux.HandleFunc("/api/v1/auth/logout", auth.logoutHandler)
	mux.HandleFunc("/api/v1/auth/verify-master", auth.verifyMasterHandler)
	mux.HandleFunc("/api/v1/auth/reset", auth.resetHandler)
	protectedAPI := http.NewServeMux()
	protectedAPI.HandleFunc("/api/v1/openclaw/status", openClaw.statusHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/models", openClaw.modelsHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/models/default", openClaw.defaultModelHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/models/auth", openClaw.modelAuthHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/plugins", openClaw.pluginsHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/plugins/", openClaw.pluginHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/channels", openClaw.channelsHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/channels/", openClaw.channelHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/agents", openClaw.agentsHandler)
	protectedAPI.HandleFunc("/api/v1/openclaw/agents/sync", openClaw.syncAgentsHandler)
	protectedAPI.HandleFunc("/api/v1/agents/chat", openClaw.chatHandler)
	protectedAPI.HandleFunc("/api/v1/auth/account", auth.accountHandler)
	protectedAPI.Handle("/api/v1/", business)
	mux.Handle("/api/v1/", auth.requireSession(protectedAPI))
	mux.Handle("/", staticFiles)

	log.Printf("STA-100 web service listening on %s", *addr)
	log.Printf("OpenClaw CLI: %s", valueOrUnavailable(openClaw.service.BinaryPath()))
	log.Printf("STA-100 database: %s", store.path)
	log.Fatal(http.ListenAndServe(*addr, requestLog(securityHeaders(mux))))
}

func newOpenClawService() *openClawService {
	manifest := strings.TrimSpace(os.Getenv("STA100_AGENT_MANIFEST"))
	if manifest == "" {
		manifest = findAgentManifest()
	}
	return &openClawService{service: orchestrator.Discover(manifest)}
}

func findAgentManifest() string {
	candidates := []string{}
	if executable, err := os.Executable(); err == nil {
		dir := filepath.Dir(executable)
		candidates = append(candidates, filepath.Join(dir, "..", "config", "sta100-agents.json"))
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "config", "sta100-agents.json"))
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return filepath.Clean(candidate)
		}
	}
	return ""
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{
		Status:    "ok",
		Service:   "sta100-web",
		Prototype: false,
		Time:      time.Now().Format(time.RFC3339),
	})
}

func (s *openClawService) statusHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	status, err := s.service.Status(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *openClawService) modelsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	models, err := s.service.Models(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, models)
}

func (s *openClawService) defaultModelHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPut) {
		return
	}
	var request struct {
		Model string `json:"model"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if err := s.service.SetConfiguredDefaultModel(request.Model); err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": true, "defaultModel": strings.TrimSpace(request.Model)})
}

func (s *openClawService) modelAuthHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Provider string `json:"provider"`
		APIKey   string `json:"apiKey"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	profile, err := s.service.SaveAPIKey(ctx, request.Provider, request.APIKey)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"updated":  true,
		"provider": strings.ToLower(strings.TrimSpace(request.Provider)),
		"profile":  profile,
	})
}

func (s *openClawService) pluginsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	plugins, err := s.service.Plugins(r.Context())
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"catalogVersion": orchestrator.PluginCatalogVersion, "count": len(plugins), "plugins": plugins})
}

func (s *openClawService) pluginHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	pluginID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/openclaw/plugins/"), "/")
	var request struct {
		Enabled bool `json:"enabled"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := s.service.SetPluginEnabled(ctx, pluginID, request.Enabled); err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": true, "id": pluginID, "enabled": request.Enabled})
}

func (s *openClawService) channelsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 75*time.Second)
	defer cancel()
	channels, err := s.service.Channels(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"catalogVersion": orchestrator.ChannelCatalogVersion, "count": len(channels), "channels": channels})
}

func (s *openClawService) channelHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/openclaw/channels/"), "/"), "/")
	if len(parts) == 2 && parts[1] == "install" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request struct {
			PackageSpec string `json:"packageSpec"`
		}
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()
		result, err := s.service.InstallChannel(ctx, parts[0], request.PackageSpec)
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if len(parts) == 2 && parts[1] == "uninstall" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()
		result, err := s.service.UninstallChannel(ctx, parts[0])
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if len(parts) == 3 && parts[1] == "qr" && parts[2] == "start" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request struct {
			Account string `json:"account"`
			Domain  string `json:"domain"`
		}
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()
		result, err := s.service.StartFeishuQR(ctx, parts[0], request.Account, request.Domain)
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if len(parts) == 4 && parts[1] == "qr" && parts[3] == "status" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		result, err := s.service.PollFeishuQR(ctx, parts[0], parts[2])
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if len(parts) == 2 && parts[1] == "status" {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		// `openclaw channels status` performs a live probe and can take more
		// than 30 seconds for a healthy Feishu connection. Keep this timeout
		// separate from quick gateway/status requests so a slow probe is not
		// misreported as a channel failure.
		ctx, cancel := context.WithTimeout(r.Context(), 55*time.Second)
		defer cancel()
		status, err := s.service.ChannelStatus(ctx, parts[0])
		if err != nil {
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, normalizeChannelStatusResponse(parts[0], status))
		return
	}
	if len(parts) == 2 && parts[1] == "login" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request struct {
			Account string `json:"account"`
		}
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		result, err := s.service.LoginChannel(ctx, parts[0], request.Account)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				writeAPIError(w, http.StatusGatewayTimeout, "CHANNEL_LOGIN_TIMEOUT", "OpenClaw 通道登录等待超时；该通道可能需要扫码或网页登录确认，但当前 OpenClaw 未在限定时间内返回可展示的授权信息。")
				return
			}
			writeOpenClawError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}
	if len(parts) != 1 || parts[0] == "" {
		writeAPIError(w, http.StatusNotFound, "CHANNEL_ENDPOINT_NOT_FOUND", "通道接口不存在")
		return
	}
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request orchestrator.ChannelAccountRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	request.Channel = parts[0]
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	result, err := s.service.AddChannelAccount(ctx, request)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func normalizeChannelStatusResponse(channel string, raw map[string]any) map[string]any {
	channel = strings.TrimSpace(strings.ToLower(channel))
	item := map[string]any{}
	if channels, ok := raw["channels"].(map[string]any); ok {
		if value, ok := channels[channel].(map[string]any); ok {
			item = value
		}
	}
	configured, _ := item["configured"].(bool)
	running, _ := item["running"].(bool)
	connected, _ := item["connected"].(bool)
	lastError, _ := item["lastError"].(string)
	queryState, _ := raw["queryState"].(string)
	queryMessage, _ := raw["queryMessage"].(string)
	accounts := make([]map[string]any, 0)
	accountConnected := false
	accountError := ""
	if channelAccounts, ok := raw["channelAccounts"].(map[string]any); ok {
		if values, ok := channelAccounts[channel].([]any); ok {
			for _, value := range values {
				account, ok := value.(map[string]any)
				if !ok {
					continue
				}
				accountLabel := "未配置"
				accountIsConnected, _ := account["connected"].(bool)
				accountConfigured, _ := account["configured"].(bool)
				accountLastError, _ := account["lastError"].(string)
				switch {
				case strings.TrimSpace(accountLastError) != "":
					accountLabel = "异常"
				case accountIsConnected:
					accountLabel = "已连接"
					accountConnected = true
				case accountConfigured:
					accountLabel = "已配置"
				}
				if strings.TrimSpace(accountLastError) != "" && accountLastError != "not configured" && accountLastError != "not connected" {
					if accountError == "" {
						accountError = accountLastError
					}
				}
				accounts = append(accounts, map[string]any{
					"accountId": account["accountId"], "label": accountLabel,
					"configured": accountConfigured, "connected": accountIsConnected,
					"lastError": strings.TrimSpace(accountLastError),
				})
			}
		}
	}
	if strings.TrimSpace(lastError) == "" {
		lastError = accountError
	}
	label := "未配置"
	message := "OpenClaw 中尚未配置该通道账号。"
	switch {
	case queryState == "timeout":
		label = "状态查询失败"
		message = firstNonEmptyString(queryMessage, "OpenClaw 状态查询超时，暂时无法确认通道真实状态。")
	case queryState != "" && queryState != "ok":
		label = "状态查询失败"
		message = firstNonEmptyString(queryMessage, "OpenClaw 状态查询失败，暂时无法确认通道真实状态。")
	case strings.TrimSpace(lastError) != "":
		label = "通道异常"
		message = "通道运行异常，请检查账号配置和 OpenClaw 网关日志。"
	case connected || accountConnected:
		label = "已连接"
		message = "通道账号已配置，并已确认与服务建立连接。"
	case running:
		label = "运行中"
		message = "通道正在运行，但 OpenClaw 尚未返回已连接状态。"
	case configured:
		label = "已配置"
		message = "通道账号已写入 OpenClaw，但当前尚未确认连接。"
	}
	return map[string]any{
		"channel":    channel,
		"label":      label,
		"message":    message,
		"queryState": queryStateOrOK(queryState),
		"configured": configured,
		"running":    running,
		"connected":  connected || accountConnected,
		"lastError":  strings.TrimSpace(lastError),
		"accounts":   accounts,
		"checkedAt":  time.Now().Format(time.RFC3339),
	}
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func queryStateOrOK(value string) string {
	if strings.TrimSpace(value) == "" {
		return "ok"
	}
	return strings.TrimSpace(value)
}

func (s *openClawService) agentsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	agents, err := s.service.ListAgents(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	agents = businessVisibleAgents(agents)
	writeJSON(w, http.StatusOK, map[string]any{"count": len(agents), "agents": agents})
}

func (s *openClawService) syncAgentsHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()
	agents, err := s.service.SyncAgents(ctx)
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	visible := businessVisibleAgents(agents)
	writeJSON(w, http.StatusOK, map[string]any{"synced": true, "count": len(visible), "agents": visible, "systemAgents": len(agents) - len(visible)})
}

func businessVisibleAgents(agents []orchestrator.Agent) []orchestrator.Agent {
	visible := make([]orchestrator.Agent, 0, len(agents))
	for _, agent := range agents {
		if agent.IsDefault || agent.Visibility == "system" {
			continue
		}
		visible = append(visible, agent)
	}
	return visible
}

func (s *openClawService) chatHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request struct {
		AgentID     string                           `json:"agentId"`
		Model       string                           `json:"model"`
		Message     string                           `json:"message"`
		SessionKey  string                           `json:"sessionKey"`
		Sources     []string                         `json:"sources"`
		Allowlist   []string                         `json:"allowlist"`
		Attachments []orchestrator.MessageAttachment `json:"attachments"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 105*time.Second)
	defer cancel()
	result, err := s.service.SendAgentMessage(ctx, orchestrator.AgentMessageInput{
		AgentID: request.AgentID, Model: request.Model, Message: request.Message, SessionKey: request.SessionKey,
		Sources: request.Sources, Allowlist: request.Allowlist, Attachments: request.Attachments,
	})
	if err != nil {
		writeOpenClawError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, destination any) error {
	if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		writeAPIError(w, http.StatusUnsupportedMediaType, "JSON_REQUIRED", "请求必须使用 application/json")
		return errors.New("json content type required")
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_JSON", "请求数据格式无效："+err.Error())
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeAPIError(w, http.StatusBadRequest, "INVALID_JSON", "请求只能包含一个 JSON 对象")
		return errors.New("multiple json objects")
	}
	return nil
}

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	return false
}

func allowMutation(w http.ResponseWriter, r *http.Request, method string) bool {
	if !allowMethod(w, r, method) {
		return false
	}
	if r.Header.Get("X-STA100-Request") != "1" {
		writeAPIError(w, http.StatusForbidden, "REQUEST_HEADER_REQUIRED", "缺少 STA-100 写操作请求头")
		return false
	}
	return true
}

func writeOpenClawError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		writeAPIError(w, http.StatusGatewayTimeout, "OPENCLAW_TIMEOUT", "OpenClaw 操作超时")
	case errors.Is(err, orchestrator.ErrUnavailable):
		writeAPIError(w, http.StatusServiceUnavailable, "OPENCLAW_UNAVAILABLE", "未找到 OpenClaw CLI")
	case errors.Is(err, orchestrator.ErrInvalidModel):
		writeAPIError(w, http.StatusBadRequest, "INVALID_MODEL", "模型 ID 格式无效")
	case errors.Is(err, orchestrator.ErrModelUnavailable):
		writeAPIError(w, http.StatusBadRequest, "MODEL_UNAVAILABLE", "所选模型不在 OpenClaw 可用模型列表中")
	case errors.Is(err, orchestrator.ErrInvalidProvider):
		writeAPIError(w, http.StatusBadRequest, "INVALID_PROVIDER", "模型提供商 ID 格式无效")
	case errors.Is(err, orchestrator.ErrInvalidAPIKey):
		writeAPIError(w, http.StatusBadRequest, "INVALID_API_KEY", "API Key 格式无效")
	case errors.Is(err, orchestrator.ErrManifestMissing):
		writeAPIError(w, http.StatusServiceUnavailable, "SYNC_MANIFEST_UNAVAILABLE", "未找到 Agent 编排清单")
	case errors.Is(err, orchestrator.ErrInvalidAgent):
		writeAPIError(w, http.StatusBadRequest, "INVALID_AGENT", "Agent ID 格式无效")
	case errors.Is(err, orchestrator.ErrAgentUnavailable):
		writeAPIError(w, http.StatusNotFound, "AGENT_UNAVAILABLE", "所选 Agent 未在 OpenClaw 中注册")
	case errors.Is(err, orchestrator.ErrInvalidMessage):
		writeAPIError(w, http.StatusBadRequest, "INVALID_MESSAGE", "消息不能为空且不能超过 32 KiB")
	case errors.Is(err, orchestrator.ErrInvalidSession):
		writeAPIError(w, http.StatusBadRequest, "INVALID_SESSION", "会话标识格式无效")
	case errors.Is(err, orchestrator.ErrInvalidSource):
		writeAPIError(w, http.StatusBadRequest, "INVALID_SOURCE_POLICY", "本次来源或联网白名单格式无效")
	case errors.Is(err, orchestrator.ErrInvalidResponse):
		writeAPIError(w, http.StatusBadGateway, "OPENCLAW_INVALID_RESPONSE", "OpenClaw 返回了无法解析的信息")
	default:
		writeAPIError(w, http.StatusBadGateway, "OPENCLAW_COMMAND_FAILED", "OpenClaw 调用失败，请检查 Gateway、模型配置或账户额度")
	}
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started).Round(time.Millisecond))
	})
}

func valueOrUnavailable(value string) string {
	if value == "" {
		return "unavailable"
	}
	return value
}
