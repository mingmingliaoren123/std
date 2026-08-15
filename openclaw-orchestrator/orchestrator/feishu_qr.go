package orchestrator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	feishuQRDefaultInterval = 5
	feishuQRDefaultExpiry   = 600
)

type feishuQRSession struct {
	ID         string
	Channel    string
	Account    string
	Domain     string
	DeviceCode string
	QRURL      string
	QRDataURL  string
	Interval   int
	ExpiresAt  time.Time
	NextPoll   time.Time
	Status     string
	Message    string
	Polling    bool
}

// StartFeishuQR begins the same device-code registration flow used by the
// pinned OpenClaw Feishu plugin. It intentionally does not run the interactive
// `channels login` wizard because that wizard requires a TTY.
func (s *Service) StartFeishuQR(ctx context.Context, channel, account, domain string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	account = strings.TrimSpace(account)
	domain = strings.TrimSpace(strings.ToLower(domain))
	if channel != "feishu" {
		return nil, errors.New("当前仅支持飞书扫码绑定")
	}
	if account == "" {
		account = "default"
	}
	if domain != "lark" {
		domain = "feishu"
	}

	if _, err := s.feishuRegistrationRequest(ctx, domain, url.Values{"action": {"init"}}); err != nil {
		return nil, fmt.Errorf("飞书扫码环境检查失败: %w", err)
	}
	begin, err := s.feishuRegistrationRequest(ctx, domain, url.Values{
		"action":            {"begin"},
		"archetype":         {"PersonalAgent"},
		"auth_method":       {"client_secret"},
		"request_user_info": {"open_id"},
	})
	if err != nil {
		return nil, fmt.Errorf("飞书二维码生成失败: %w", err)
	}

	deviceCode := jsonString(begin["device_code"])
	verificationURI := jsonString(begin["verification_uri_complete"])
	if deviceCode == "" || verificationURI == "" {
		return nil, fmt.Errorf("%w: 飞书未返回 device_code 或二维码地址", ErrInvalidResponse)
	}
	qrURL, err := url.Parse(verificationURI)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Feishu QR URL", ErrInvalidResponse)
	}
	query := qrURL.Query()
	query.Set("from", "oc_onboard")
	query.Set("tp", "ob_cli_app")
	qrURL.RawQuery = query.Encode()

	interval := jsonInt(begin["interval"])
	if interval <= 0 {
		interval = feishuQRDefaultInterval
	}
	expireIn := jsonInt(begin["expires_in"])
	if expireIn <= 0 {
		expireIn = feishuQRDefaultExpiry
	}
	sessionID, err := randomSessionID()
	if err != nil {
		return nil, err
	}
	session := &feishuQRSession{
		ID:         sessionID,
		Channel:    channel,
		Account:    account,
		Domain:     domain,
		DeviceCode: deviceCode,
		QRURL:      qrURL.String(),
		QRDataURL:  renderQRDataURL(qrURL.String()),
		Interval:   interval,
		ExpiresAt:  time.Now().Add(time.Duration(expireIn) * time.Second),
		Status:     "pending",
		Message:    "请使用飞书移动端扫描二维码并确认授权",
	}
	s.qrMu.Lock()
	s.qrSessions[sessionID] = session
	s.qrMu.Unlock()

	return s.feishuQRResponse(session), nil
}

// PollFeishuQR performs one non-blocking poll. The browser owns the polling
// interval so a page request never blocks for the whole 10-minute QR lifetime.
func (s *Service) PollFeishuQR(ctx context.Context, channel, sessionID string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	sessionID = strings.TrimSpace(sessionID)
	if channel != "feishu" || sessionID == "" {
		return nil, errors.New("invalid Feishu QR session")
	}
	s.qrMu.Lock()
	session, ok := s.qrSessions[sessionID]
	if !ok || session.Channel != channel {
		s.qrMu.Unlock()
		return nil, errors.New("飞书二维码会话不存在或已失效")
	}
	if session.Status != "pending" {
		result := s.feishuQRResponse(session)
		s.qrMu.Unlock()
		return result, nil
	}
	if time.Now().After(session.ExpiresAt) {
		session.Status = "expired"
		session.Message = "二维码已过期，请重新生成"
		result := s.feishuQRResponse(session)
		s.qrMu.Unlock()
		return result, nil
	}
	if session.Polling {
		result := s.feishuQRResponse(session)
		s.qrMu.Unlock()
		return result, nil
	}
	if !session.NextPoll.IsZero() && time.Now().Before(session.NextPoll) {
		result := s.feishuQRResponse(session)
		s.qrMu.Unlock()
		return result, nil
	}
	session.Polling = true
	deviceCode := session.DeviceCode
	domain := session.Domain
	interval := session.Interval
	account := session.Account
	s.qrMu.Unlock()

	result, err := s.feishuRegistrationRequest(ctx, domain, url.Values{
		"action":      {"poll"},
		"device_code": {deviceCode},
		"tp":          {"ob_cli_app"},
	})

	s.qrMu.Lock()
	defer s.qrMu.Unlock()
	session.Polling = false
	session.NextPoll = time.Now().Add(time.Duration(interval) * time.Second)
	if err != nil {
		session.Message = "暂时无法连接飞书授权服务，将自动重试"
		return s.feishuQRResponse(session), nil
	}

	if tenant := nestedString(result, "user_info", "tenant_brand"); tenant == "lark" && session.Domain != "lark" {
		session.Domain = "lark"
		session.Message = "已识别为 Lark，正在切换授权区域"
		session.NextPoll = time.Now()
		return s.feishuQRResponse(session), nil
	}
	appID := jsonString(result["client_id"])
	appSecret := jsonString(result["client_secret"])
	if appID != "" && appSecret != "" {
		openID := nestedString(result, "user_info", "open_id")
		if err := s.applyFeishuRegistration(ctx, account, session.Domain, appID, appSecret, openID); err != nil {
			session.Status = "error"
			session.Message = "授权成功，但写入 OpenClaw 配置或重启网关失败: " + err.Error()
			return s.feishuQRResponse(session), nil
		}
		session.Status = "success"
		session.Message = "飞书已绑定，OpenClaw 配置和网关已更新"
		return s.feishuQRResponse(session), nil
	}

	switch jsonString(result["error"]) {
	case "access_denied":
		session.Status = "denied"
		session.Message = "用户拒绝了飞书授权"
	case "expired_token":
		session.Status = "expired"
		session.Message = "二维码已过期，请重新生成"
	case "slow_down":
		session.Interval += 5
		session.Message = "授权服务要求降低轮询频率，正在继续等待"
	case "", "authorization_pending":
		session.Message = "等待飞书扫码确认"
	default:
		session.Status = "error"
		session.Message = jsonString(result["error"]) + ": " + firstNonEmpty(jsonString(result["error_description"]), "飞书授权失败")
	}
	return s.feishuQRResponse(session), nil
}

func (s *Service) feishuQRResponse(session *feishuQRSession) map[string]any {
	return map[string]any{
		"status":    session.Status,
		"message":   session.Message,
		"sessionId": session.ID,
		"channel":   session.Channel,
		"account":   session.Account,
		"domain":    session.Domain,
		"qrUrl":     session.QRURL,
		"qrDataUrl": session.QRDataURL,
		"expiresAt": session.ExpiresAt.Format(time.RFC3339),
		"interval":  session.Interval,
	}
}

func (s *Service) feishuRegistrationRequest(ctx context.Context, domain string, values url.Values) (map[string]any, error) {
	base := "https://accounts.feishu.cn"
	if domain == "lark" {
		base = "https://accounts.larksuite.com"
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/oauth/v1/app/registration", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 12 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("%w: 飞书返回无效 JSON", ErrInvalidResponse)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if values.Get("action") == "poll" {
			return result, nil
		}
		return nil, fmt.Errorf("HTTP %d: %s", response.StatusCode, firstNonEmpty(jsonString(result["error_description"]), jsonString(result["message"]), "飞书请求失败"))
	}
	if supported, ok := result["supported_auth_methods"].([]any); ok && len(supported) > 0 {
		hasClientSecret := false
		for _, method := range supported {
			if jsonString(method) == "client_secret" {
				hasClientSecret = true
				break
			}
		}
		if !hasClientSecret {
			return nil, errors.New("当前飞书环境不支持 client_secret 扫码注册")
		}
	}
	return result, nil
}

func (s *Service) applyFeishuRegistration(ctx context.Context, account, domain, appID, appSecret, openID string) error {
	if s.configPath == "" {
		return ErrUnavailable
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}
	var config map[string]any
	if strings.TrimSpace(string(raw)) == "" {
		config = map[string]any{}
	} else if err := json.Unmarshal(raw, &config); err != nil {
		return fmt.Errorf("%w: OpenClaw 配置格式无效", ErrInvalidResponse)
	}
	channels := ensureObject(config, "channels")
	feishu := ensureObject(channels, "feishu")
	feishu["enabled"] = true
	if domain != "" {
		feishu["domain"] = domain
	}
	target := feishu
	if account != "" && account != "default" {
		accounts := ensureObject(feishu, "accounts")
		target = ensureObject(accounts, account)
		target["enabled"] = true
		target["appId"] = appID
		target["appSecret"] = appSecret
		target["domain"] = domain
	} else {
		target["appId"] = appID
		target["appSecret"] = appSecret
	}
	if openID != "" {
		target["dmPolicy"] = "allowlist"
		target["allowFrom"] = []any{openID}
	}
	if _, exists := target["groupPolicy"]; !exists {
		target["groupPolicy"] = "allowlist"
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err := writeFileAtomic(s.configPath, append(data, '\n'), 0o600); err != nil {
		return err
	}
	if s.bin != "" && isExecutable(s.bin) {
		restartCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		if _, err := s.run(restartCtx, nil, "gateway", "restart"); err != nil {
			return fmt.Errorf("OpenClaw 网关重启失败: %v", err)
		}
	}
	return nil
}

func randomSessionID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func renderQRDataURL(value string) string {
	qrencode, err := exec.LookPath("qrencode")
	if err != nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, qrencode, "-t", "PNG", "-o", "-", "-s", "6", "-m", "2", value).Output()
	if err != nil || len(output) == 0 {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(output)
}

func jsonString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return fmt.Sprintf("%v", typed)
	default:
		return ""
	}
}

func jsonInt(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case json.Number:
		value, _ := typed.Int64()
		return int(value)
	default:
		return 0
	}
}

func nestedString(value map[string]any, keys ...string) string {
	var current any = value
	for _, key := range keys {
		object, ok := current.(map[string]any)
		if !ok {
			return ""
		}
		current = object[key]
	}
	return jsonString(current)
}
