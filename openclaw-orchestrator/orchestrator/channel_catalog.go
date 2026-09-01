package orchestrator

const ChannelCatalogVersion = "openclaw-2026.7.1-2"

// fixedChannelCatalog is generated from `openclaw channels list --all --json` for OpenClaw 2026.7.1-2.
var fixedChannelCatalog = []Channel{
	{ID: "feishu", Name: "飞书", Installed: false, Origin: "installable", AccountCount: 0, Description: "飞书机器人/应用通道，需安装通道包后按 OpenClaw 引导绑定。"},
	{ID: "wecom", Name: "企业微信", Installed: false, Origin: "installable", AccountCount: 0, Description: "企业微信通道，需安装通道包并配置企业微信应用凭据。"},
	{ID: "openclaw-weixin", Name: "微信", Installed: false, Origin: "installable", AccountCount: 0, Description: "微信通道，需安装 openclaw-weixin 通道包后按 OpenClaw 登录/绑定流程处理。"},
	{ID: "telegram", Name: "Telegram", Installed: false, Origin: "installable", AccountCount: 0, Description: "Telegram Bot 通道，可通过 Bot Token 绑定。"},
	{ID: "whatsapp", Name: "WhatsApp", Installed: false, Origin: "installable", AccountCount: 0, Description: "WhatsApp Web 通道，通常需要网页登录。"},
	{ID: "slack", Name: "Slack", Installed: false, Origin: "installable", AccountCount: 0, Description: "Slack 通道，通常需要 OAuth 或应用 Token。"},
	{ID: "imessage", Name: "iMessage", Installed: false, Origin: "installable", AccountCount: 0, Description: "iMessage 本地通道，通常需要本机消息数据库权限。"},
	{ID: "clickclack", Name: "ClickClack", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "discord", Name: "Discord", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "googlechat", Name: "Google Chat", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "irc", Name: "IRC", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "line", Name: "LINE", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "matrix", Name: "Matrix", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "mattermost", Name: "Mattermost", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "msteams", Name: "Microsoft Teams", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "nextcloud-talk", Name: "Nextcloud Talk", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "nostr", Name: "Nostr", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "openclaw-zaloclawbot", Name: "Zalo Clawbot", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "qqbot", Name: "QQ Bot", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "raft", Name: "Raft", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "signal", Name: "Signal", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "sms", Name: "SMS", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "synology-chat", Name: "Synology Chat", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "tlon", Name: "Tlon", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "twitch", Name: "Twitch", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "yuanbao", Name: "元宝", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "zalo", Name: "Zalo", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
	{ID: "zalouser", Name: "Zalo User", Installed: false, Origin: "installable", AccountCount: 0, Description: "OpenClaw 聊天通道，具体绑定字段以 OpenClaw 当前通道实现为准。"},
}

func fixedChannels() []Channel {
	return append([]Channel(nil), fixedChannelCatalog...)
}
