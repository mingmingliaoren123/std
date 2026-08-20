package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

func (s *Service) Plugins(ctx context.Context) ([]Plugin, error) {
	_ = ctx
	enabled := s.pluginEnabledState()
	items := fixedPlugins()
	for index := range items {
		if configured, ok := enabled[items[index].ID]; ok {
			items[index].Enabled = configured
		}
		items[index].Status = map[bool]string{true: "enabled", false: "disabled"}[items[index].Enabled]
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

func (s *Service) Channels(ctx context.Context) ([]Channel, error) {
	_ = ctx
	items := fixedChannels()
	enabled := s.pluginEnabledState()
	s.pluginMu.Lock()
	inventory := append([]pluginInventoryEntry(nil), s.pluginCache...)
	s.pluginMu.Unlock()
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	for index := range items {
		pluginID, installSpec, bindingMode := channelIntegrationMetadata(items[index].ID)
		items[index].PluginID = pluginID
		items[index].InstallSpec = installSpec
		items[index].BindingMode = bindingMode
		items[index].CanInstall = installSpec != "" || pluginID != ""
		items[index].CanUninstall = pluginID != ""
		items[index].Enabled = enabled[items[index].ID]
		if items[index].Enabled {
			items[index].Installed = true
		}
		if len(inventory) > 0 {
			if existing, ok := findChannelPlugin(inventory, items[index].ID, pluginID); ok {
				items[index].PluginID = existing.ID
				items[index].Installed = true
				items[index].Enabled = existing.Enabled || existing.Status == "loaded"
				items[index].Origin = existing.Origin
			} else {
				items[index].Installed = false
				items[index].Enabled = false
				items[index].Origin = "installable"
			}
		}
		if items[index].Installed {
			items[index].Configured, items[index].AccountCount = s.configuredChannelState(items[index].ID)
		}
		if items[index].Configured || items[index].AccountCount > 0 {
			items[index].Status = "configured"
		} else if items[index].Installed {
			items[index].Status = "installed"
		} else {
			items[index].Status = "installable"
		}
	}
	return items, nil
}

type pluginInventoryEntry struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Origin      string   `json:"origin"`
	Status      string   `json:"status"`
	Enabled     bool     `json:"enabled"`
	ChannelIDs  []string `json:"channelIds"`
}

func channelIntegrationMetadata(channel string) (pluginID, installSpec, bindingMode string) {
	metadata := map[string][3]string{
		"wecom":                {"wecom-openclaw-plugin", "@wecom/wecom-openclaw-plugin@2026.5.7", "config"},
		"yuanbao":              {"openclaw-plugin-yuanbao", "openclaw-plugin-yuanbao@2.15.0", "config"},
		"openclaw-weixin":      {"openclaw-weixin", "@tencent-weixin/openclaw-weixin@2.4.6", "login"},
		"openclaw-zaloclawbot": {"openclaw-zaloclawbot", "@zalo-platforms/openclaw-zaloclawbot@0.1.4", "login"},
		"clickclack":           {"clickclack", "@openclaw/clickclack", "config"},
		"discord":              {"discord", "@openclaw/discord", "config"},
		"feishu":               {"feishu", "@openclaw/feishu", "qr"},
		"googlechat":           {"googlechat", "@openclaw/googlechat", "config"},
		"irc":                  {"irc", "@openclaw/irc", "config"},
		"line":                 {"line", "@openclaw/line", "config"},
		"matrix":               {"matrix", "@openclaw/matrix", "config"},
		"mattermost":           {"mattermost", "@openclaw/mattermost", "config"},
		"msteams":              {"msteams", "@openclaw/msteams", "config"},
		"nextcloud-talk":       {"nextcloud-talk", "@openclaw/nextcloud-talk", "config"},
		"nostr":                {"nostr", "@openclaw/nostr", "config"},
		"qqbot":                {"qqbot", "@openclaw/qqbot", "config"},
		"raft":                 {"raft", "@openclaw/raft", "config"},
		"signal":               {"signal", "@openclaw/signal", "config"},
		"slack":                {"slack", "@openclaw/slack", "config"},
		"sms":                  {"sms", "@openclaw/sms", "config"},
		"synology-chat":        {"synology-chat", "@openclaw/synology-chat", "config"},
		"telegram":             {"telegram", "@openclaw/telegram", "token"},
		"tlon":                 {"tlon", "@openclaw/tlon", "config"},
		"twitch":               {"twitch", "@openclaw/twitch", "config"},
		"whatsapp":             {"whatsapp", "@openclaw/whatsapp", "login"},
		"zalo":                 {"zalo", "@openclaw/zalo", "config"},
		"zalouser":             {"zalouser", "@openclaw/zalouser", "config"},
	}
	channel = strings.TrimSpace(strings.ToLower(channel))
	if value, ok := metadata[channel]; ok {
		return value[0], value[1], value[2]
	}
	return channel, "", "config"
}

func (s *Service) pluginInventory(ctx context.Context) ([]pluginInventoryEntry, error) {
	s.pluginMu.Lock()
	if len(s.pluginCache) > 0 && time.Since(s.pluginAt) < 5*time.Second {
		items := append([]pluginInventoryEntry(nil), s.pluginCache...)
		s.pluginMu.Unlock()
		return items, nil
	}
	s.pluginMu.Unlock()

	out, err := s.run(ctx, nil, "plugins", "list", "--json")
	if err != nil {
		return nil, err
	}
	var response struct {
		Plugins []pluginInventoryEntry `json:"plugins"`
	}
	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("%w: plugins list: %v", ErrInvalidResponse, err)
	}
	s.pluginMu.Lock()
	s.pluginCache = append([]pluginInventoryEntry(nil), response.Plugins...)
	s.pluginAt = time.Now()
	s.pluginMu.Unlock()
	return response.Plugins, nil
}

func (s *Service) invalidatePluginInventory() {
	s.pluginMu.Lock()
	s.pluginCache = nil
	s.pluginAt = time.Time{}
	s.pluginMu.Unlock()
}

func findChannelPlugin(inventory []pluginInventoryEntry, channel string, pluginID string) (pluginInventoryEntry, bool) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	for _, item := range inventory {
		if item.ID == pluginID {
			return item, true
		}
		for _, channelID := range item.ChannelIDs {
			if strings.EqualFold(channelID, channel) {
				return item, true
			}
		}
	}
	return pluginInventoryEntry{}, false
}

func (s *Service) InstallChannel(ctx context.Context, channel, packageSpec string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	pluginID, defaultSpec, _ := channelIntegrationMetadata(channel)
	if pluginID == "" {
		return nil, errors.New("通道 ID 不能为空")
	}
	inventory, err := s.pluginInventory(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取 OpenClaw 插件目录失败: %w", err)
	}
	if existing, ok := findChannelPlugin(inventory, channel, pluginID); ok {
		if existing.Enabled || existing.Status == "loaded" {
			return channelPluginResult(channel, existing, "already_enabled", "通道插件已经安装并启用"), nil
		}
		if _, err := s.run(ctx, nil, "plugins", "enable", existing.ID); err != nil {
			return nil, fmt.Errorf("启用 OpenClaw 通道插件失败: %w", err)
		}
		verified, err := s.verifyChannelPlugin(ctx, channel, existing.ID, true)
		if err != nil {
			return nil, err
		}
		return channelPluginResult(channel, verified, "enabled", "通道插件已启用并通过 OpenClaw 目录复核"), nil
	}
	spec := strings.TrimSpace(packageSpec)
	if spec == "" {
		spec = defaultSpec
	}
	if spec == "" {
		return nil, fmt.Errorf("OpenClaw 当前未提供 %s 的默认安装包，请填写插件包名或 npm spec", channel)
	}
	if _, err := s.run(ctx, nil, "plugins", "install", spec); err != nil {
		return nil, fmt.Errorf("安装 OpenClaw 通道插件失败: %w", err)
	}
	s.invalidatePluginInventory()
	installed, err := s.verifyChannelPlugin(ctx, channel, pluginID, false)
	if err != nil {
		return nil, fmt.Errorf("OpenClaw 安装命令完成，但安装结果复核失败: %w", err)
	}
	if _, err := s.run(ctx, nil, "plugins", "enable", installed.ID); err != nil {
		return nil, fmt.Errorf("安装后启用 OpenClaw 通道插件失败: %w", err)
	}
	verified, err := s.verifyChannelPlugin(ctx, channel, installed.ID, true)
	if err != nil {
		return nil, err
	}
	result := channelPluginResult(channel, verified, "installed", "通道插件已安装、启用并通过 OpenClaw 目录复核")
	result["installSpec"] = spec
	return result, nil
}

func (s *Service) UninstallChannel(ctx context.Context, channel string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	pluginID, _, _ := channelIntegrationMetadata(channel)
	inventory, err := s.pluginInventory(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取 OpenClaw 插件目录失败: %w", err)
	}
	existing, ok := findChannelPlugin(inventory, channel, pluginID)
	if !ok {
		return map[string]any{"updated": true, "channel": channel, "action": "already_uninstalled", "message": "通道插件当前未安装"}, nil
	}
	if existing.Origin == "bundled" {
		if _, err := s.run(ctx, nil, "plugins", "disable", existing.ID); err != nil {
			return nil, fmt.Errorf("停用 OpenClaw 内置通道插件失败: %w", err)
		}
		verified, err := s.verifyChannelPlugin(ctx, channel, existing.ID, false)
		if err != nil {
			return nil, err
		}
		result := channelPluginResult(channel, verified, "disabled", "内置通道不能删除，已停用并通过 OpenClaw 目录复核")
		result["verifiedDisabled"] = !verified.Enabled && verified.Status != "loaded"
		return result, nil
	}
	if _, err := s.run(ctx, nil, "plugins", "uninstall", existing.ID, "--force"); err != nil {
		return nil, fmt.Errorf("卸载 OpenClaw 通道插件失败: %w", err)
	}
	s.invalidatePluginInventory()
	after, err := s.pluginInventory(ctx)
	if err != nil {
		return nil, fmt.Errorf("卸载命令完成，但无法复核 OpenClaw 插件目录: %w", err)
	}
	if _, ok := findChannelPlugin(after, channel, existing.ID); ok {
		return nil, fmt.Errorf("卸载命令完成，但 OpenClaw 插件目录仍显示该通道已安装")
	}
	result := channelPluginResult(channel, existing, "uninstalled", "通道插件已卸载并通过 OpenClaw 目录复核")
	result["verified"] = true
	return result, nil
}

func (s *Service) verifyChannelPlugin(ctx context.Context, channel, pluginID string, requireEnabled bool) (pluginInventoryEntry, error) {
	s.invalidatePluginInventory()
	inventory, err := s.pluginInventory(ctx)
	if err != nil {
		return pluginInventoryEntry{}, fmt.Errorf("复核 OpenClaw 插件目录失败: %w", err)
	}
	existing, ok := findChannelPlugin(inventory, channel, pluginID)
	if !ok {
		return pluginInventoryEntry{}, fmt.Errorf("未在 OpenClaw 插件目录中找到通道 %s", channel)
	}
	if requireEnabled && !existing.Enabled && existing.Status != "loaded" {
		return pluginInventoryEntry{}, fmt.Errorf("通道 %s 已安装但未启用，OpenClaw 返回状态为 %s", channel, existing.Status)
	}
	return existing, nil
}

func channelPluginResult(channel string, plugin pluginInventoryEntry, action, message string) map[string]any {
	return map[string]any{
		"updated": true, "verified": true, "channel": channel, "pluginId": plugin.ID,
		"pluginName": plugin.Name, "version": plugin.Version, "origin": plugin.Origin,
		"enabled": plugin.Enabled || plugin.Status == "loaded", "status": plugin.Status,
		"action": action, "message": message,
	}
}

func summarizeChannelStatus(status map[string]any, channel string) (configured, running bool, accountCount int, lastError string) {
	if channels, ok := status["channels"].(map[string]any); ok {
		if item, ok := channels[channel].(map[string]any); ok {
			configured, _ = item["configured"].(bool)
			running, _ = item["running"].(bool)
			lastError, _ = item["lastError"].(string)
		}
	}
	if accounts, ok := status["channelAccounts"].(map[string]any); ok {
		if list, ok := accounts[channel].([]any); ok {
			for _, raw := range list {
				account, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				accountConfigured, _ := account["configured"].(bool)
				if accountConfigured {
					accountCount++
				}
				if lastError == "" {
					lastError, _ = account["lastError"].(string)
				}
				if accountConfigured {
					configured = true
				}
				if accountRunning, _ := account["running"].(bool); accountRunning {
					running = true
				}
			}
		}
	}
	return configured, running, accountCount, strings.TrimSpace(lastError)
}

func (s *Service) SetPluginEnabled(ctx context.Context, pluginID string, enabled bool) error {
	pluginID = strings.TrimSpace(strings.ToLower(pluginID))
	if !providerPattern.MatchString(pluginID) {
		return errors.New("invalid plugin id")
	}
	command := "disable"
	if enabled {
		command = "enable"
	}
	if _, err := s.run(ctx, nil, "plugins", command, pluginID); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChannelStatus(ctx context.Context, channel string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	if !providerPattern.MatchString(channel) {
		return nil, errors.New("invalid channel id")
	}
	out, err := s.run(ctx, nil, "channels", "status", "--channel", channel, "--json")
	if err != nil {
		configured, accountCount := s.configuredChannelState(channel)
		queryState := "failed"
		queryMessage := "OpenClaw 通道状态查询失败，请检查 OpenClaw 服务后重试。"
		if errors.Is(err, context.DeadlineExceeded) {
			queryState = "timeout"
			queryMessage = "OpenClaw 通道状态查询超时，暂时无法确认通道真实状态，请稍后重试。"
		} else if errors.Is(err, context.Canceled) {
			queryState = "canceled"
			queryMessage = "通道状态查询已取消，请稍后重试。"
		}
		return map[string]any{
			"queryState":     queryState,
			"queryMessage":   queryMessage,
			"queryError":     err.Error(),
			"queryErrorText": tail(err.Error(), 300),
			"channels": map[string]any{
				channel: map[string]any{
					"configured": configured,
					"running":    false,
				},
			},
			"channelAccounts": map[string]any{
				channel: channelAccountPlaceholders(accountCount, configured),
			},
		}, nil
	}
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("%w: channels status: %v", ErrInvalidResponse, err)
	}
	return result, nil
}

func (s *Service) configuredChannelState(channel string) (bool, int) {
	if s.configPath == "" {
		return false, 0
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return false, 0
	}
	var config map[string]any
	if json.Unmarshal(raw, &config) != nil {
		return false, 0
	}
	channels, ok := config["channels"].(map[string]any)
	if !ok {
		return false, 0
	}
	channelConfig, ok := channels[channel].(map[string]any)
	if !ok {
		return false, 0
	}
	topConfigured := configuredChannelValue(channelConfig["appId"]) && configuredChannelValue(channelConfig["appSecret"])
	count := 0
	if topConfigured {
		count = 1
	}
	accounts, _ := channelConfig["accounts"].(map[string]any)
	for _, rawAccount := range accounts {
		account, ok := rawAccount.(map[string]any)
		if !ok {
			continue
		}
		appID := account["appId"]
		secret := account["appSecret"]
		if !configuredChannelValue(appID) {
			appID = channelConfig["appId"]
		}
		if !configuredChannelValue(secret) {
			secret = channelConfig["appSecret"]
		}
		if configuredChannelValue(appID) && configuredChannelValue(secret) {
			count++
		}
	}
	return count > 0, count
}

func configuredChannelValue(value any) bool {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed) != ""
	case map[string]any:
		if strings.EqualFold(strings.TrimSpace(fmt.Sprint(typed["source"])), "env") {
			name := strings.TrimSpace(fmt.Sprint(typed["id"]))
			return name != "" && os.Getenv(name) != ""
		}
		return len(typed) > 0
	default:
		return value != nil
	}
}

func channelAccountPlaceholders(count int, configured bool) []any {
	if count < 1 && !configured {
		return []any{map[string]any{"accountId": "default", "configured": false, "running": false, "lastError": "not configured"}}
	}
	accounts := make([]any, 0, count)
	for index := 0; index < count; index++ {
		accounts = append(accounts, map[string]any{"accountId": fmt.Sprintf("account-%d", index+1), "configured": true, "running": false})
	}
	return accounts
}

func (s *Service) AddChannelAccount(ctx context.Context, request ChannelAccountRequest) (map[string]any, error) {
	request.Channel = strings.TrimSpace(strings.ToLower(request.Channel))
	request.Account = strings.TrimSpace(request.Account)
	request.Name = strings.TrimSpace(request.Name)
	if !providerPattern.MatchString(request.Channel) {
		return nil, errors.New("invalid channel id")
	}
	args := []string{"channels", "add", "--channel", request.Channel}
	addPair := func(flag, value string) {
		value = strings.TrimSpace(value)
		if value != "" {
			args = append(args, flag, value)
		}
	}
	addPair("--account", request.Account)
	addPair("--name", request.Name)
	addPair("--url", request.URL)
	addPair("--base-url", request.BaseURL)
	addPair("--http-url", request.HTTPURL)
	addPair("--cli-path", request.CLIPath)
	addPair("--db-path", request.DBPath)
	addPair("--service", request.Service)
	addPair("--region", request.Region)
	addPair("--password", request.Password)
	if request.UseEnv {
		args = append(args, "--use-env")
	}

	cleanup := func() {}
	if strings.TrimSpace(request.Token) != "" {
		path, remove, err := writeSecretTempFile("openclaw-channel-token-", request.Token)
		if err != nil {
			return nil, err
		}
		cleanup = remove
		args = append(args, "--token-file", path)
	}
	if strings.TrimSpace(request.Secret) != "" {
		path, remove, err := writeSecretTempFile("openclaw-channel-secret-", request.Secret)
		if err != nil {
			cleanup()
			return nil, err
		}
		previous := cleanup
		cleanup = func() { previous(); remove() }
		args = append(args, "--secret-file", path)
	}
	defer cleanup()

	out, err := s.run(ctx, nil, args...)
	if err != nil {
		return nil, err
	}
	result := map[string]any{"updated": true, "channel": request.Channel}
	if trimmed := strings.TrimSpace(string(out)); trimmed != "" {
		var parsed any
		if json.Unmarshal(out, &parsed) == nil {
			result["openclaw"] = parsed
		} else {
			result["output"] = trimmed
		}
	}
	return result, nil
}

func (s *Service) LoginChannel(ctx context.Context, channel, account string) (map[string]any, error) {
	channel = strings.TrimSpace(strings.ToLower(channel))
	account = strings.TrimSpace(account)
	if !providerPattern.MatchString(channel) {
		return nil, errors.New("invalid channel id")
	}
	args := []string{"channels", "login", "--channel", channel}
	if account != "" {
		args = append(args, "--account", account)
	}
	out, err := s.run(ctx, nil, args...)
	if err != nil {
		return nil, err
	}
	result := map[string]any{"updated": true, "channel": channel}
	if account != "" {
		result["account"] = account
	}
	if trimmed := strings.TrimSpace(string(out)); trimmed != "" {
		result["output"] = trimmed
	}
	return result, nil
}

func writeSecretTempFile(pattern, value string) (string, func(), error) {
	file, err := os.CreateTemp("", pattern+"*.txt")
	if err != nil {
		return "", func() {}, fmt.Errorf("create secret temp file: %w", err)
	}
	path := file.Name()
	remove := func() { _ = os.Remove(path) }
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		remove()
		return "", func() {}, fmt.Errorf("secure secret temp file: %w", err)
	}
	if _, err := io.WriteString(file, strings.TrimSpace(value)); err != nil {
		file.Close()
		remove()
		return "", func() {}, fmt.Errorf("write secret temp file: %w", err)
	}
	if err := file.Close(); err != nil {
		remove()
		return "", func() {}, fmt.Errorf("close secret temp file: %w", err)
	}
	return path, remove, nil
}

func (s *Service) pluginEnabledState() map[string]bool {
	result := map[string]bool{}
	if s.configPath == "" {
		return result
	}
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return result
	}
	var config struct {
		Plugins struct {
			Entries map[string]struct {
				Enabled *bool `json:"enabled"`
			} `json:"entries"`
		} `json:"plugins"`
	}
	if json.Unmarshal(data, &config) != nil {
		return result
	}
	for id, entry := range config.Plugins.Entries {
		if entry.Enabled != nil {
			result[id] = *entry.Enabled
		}
	}
	return result
}
