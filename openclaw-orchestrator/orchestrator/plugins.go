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
	items := fixedChannels()
	enabled := s.pluginEnabledState()
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	for index := range items {
		items[index].Enabled = enabled[items[index].ID]
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
		return map[string]any{
			"queryError": err.Error(),
			"channels": map[string]any{
				channel: map[string]any{
					"configured": configured,
					"running":    false,
					"lastError":  tail(err.Error(), 300),
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
