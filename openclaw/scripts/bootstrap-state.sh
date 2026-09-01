#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_ROOT="$(cd "$ROOT/.." && pwd)"
OPENCLAW_HOME="${OPENCLAW_HOME:-$RELEASE_ROOT/data/openclaw-home}"
OPENCLAW_STATE_DIR="${OPENCLAW_STATE_DIR:-$RELEASE_ROOT/data/openclaw-state}"
OPENCLAW_PROFILE="${OPENCLAW_PROFILE:-sta100}"
OPENCLAW_GATEWAY_PORT="${OPENCLAW_GATEWAY_PORT:-18789}"
CONFIG_PATH="${OPENCLAW_CONFIG_PATH:-$OPENCLAW_STATE_DIR/openclaw.json}"
RELEASE_ROOT="$(cd "$ROOT/.." && pwd)"
CHANNEL_SKILL_TOKEN="$RELEASE_ROOT/data/channel-skill.token"

if [[ ! "$OPENCLAW_GATEWAY_PORT" =~ ^[0-9]+$ ]] || (( OPENCLAW_GATEWAY_PORT < 1 || OPENCLAW_GATEWAY_PORT > 65535 )); then
  echo "invalid OPENCLAW_GATEWAY_PORT: $OPENCLAW_GATEWAY_PORT" >&2
  exit 2
fi

mkdir -p \
  "$OPENCLAW_HOME" \
  "$OPENCLAW_HOME/workspace" \
  "$OPENCLAW_STATE_DIR" \
  "$OPENCLAW_STATE_DIR/state" \
  "$OPENCLAW_STATE_DIR/agents/main/agent" \
  "$OPENCLAW_STATE_DIR/agents/main/sessions" \
  "$OPENCLAW_STATE_DIR/logs"

if [[ ! -s "$CHANNEL_SKILL_TOKEN" ]]; then
  umask 077
  if command -v node >/dev/null 2>&1; then
    node -e 'console.log(require("crypto").randomBytes(32).toString("hex"))' >"$CHANNEL_SKILL_TOKEN"
  else
    head -c 32 /dev/urandom | od -An -tx1 | tr -d " \n" >"$CHANNEL_SKILL_TOKEN"
    printf '\n' >>"$CHANNEL_SKILL_TOKEN"
  fi
  chmod 0600 "$CHANNEL_SKILL_TOKEN"
fi

ensure_channel_skill_plugin() {
  [[ -s "$CONFIG_PATH" ]] || return 0
  command -v node >/dev/null 2>&1 || return 0
  node - "$CONFIG_PATH" "$ROOT/extensions/sta100-channel-skill" <<'NODE'
const fs = require("fs");
const [configPath, pluginPath] = process.argv.slice(2);
try {
  const config = JSON.parse(fs.readFileSync(configPath, "utf8"));
  config.plugins ??= {};
  config.plugins.entries ??= {};
  config.plugins.entries["sta100-channel-skill"] ??= { enabled: true };
  config.plugins.load ??= {};
  const paths = Array.isArray(config.plugins.load.paths) ? config.plugins.load.paths : [];
  const pluginEntry = `${pluginPath}/index.js`;
  const normalizedPaths = paths.filter((entry) => entry !== pluginPath);
  if (!normalizedPaths.includes(pluginEntry)) normalizedPaths.push(pluginEntry);
  if (JSON.stringify(paths) !== JSON.stringify(normalizedPaths)) {
    config.plugins.load.paths = normalizedPaths;
    fs.writeFileSync(configPath, JSON.stringify(config, null, 2) + "\n", { mode: 0o600 });
  }
} catch (error) {
  console.error(`channel skill plugin config update failed: ${error.message}`);
  process.exitCode = 1;
}
NODE
}

if [[ -s "$CONFIG_PATH" ]]; then
  ensure_channel_skill_plugin
  exit 0
fi

if command -v node >/dev/null 2>&1; then
  OPENCLAW_GATEWAY_TOKEN="${OPENCLAW_GATEWAY_TOKEN:-$(node -e 'console.log(require("crypto").randomBytes(24).toString("hex"))')}"
else
  OPENCLAW_GATEWAY_TOKEN="${OPENCLAW_GATEWAY_TOKEN:-$(head -c 24 /dev/urandom | od -An -tx1 | tr -d " \n")}"
fi

cat > "$CONFIG_PATH" <<EOF
{
  "agents": {
    "defaults": {
      "workspace": "$OPENCLAW_HOME/workspace",
      "models": {}
    },
    "list": [
      {
        "id": "main",
        "workspace": "$OPENCLAW_HOME/workspace",
        "agentDir": "$OPENCLAW_STATE_DIR/agents/main/agent",
        "identity": {
          "name": "STA-100 主智能体",
          "emoji": "🧭"
        }
      }
    ]
  },
  "gateway": {
    "mode": "local",
    "auth": {
      "mode": "token",
      "token": "$OPENCLAW_GATEWAY_TOKEN"
    },
    "port": $OPENCLAW_GATEWAY_PORT,
    "bind": "loopback",
    "tailscale": {
      "mode": "off",
      "resetOnExit": false
    },
    "controlUi": {
      "allowInsecureAuth": true
    },
    "nodes": {
      "denyCommands": [
        "camera.snap",
        "camera.clip",
        "screen.record",
        "contacts.add",
        "calendar.add",
        "reminders.add",
        "sms.send",
        "sms.search"
      ]
    }
  },
  "session": {
    "dmScope": "per-channel-peer"
  },
  "tools": {
    "profile": "coding"
  },
  "plugins": {
    "entries": {
      "sta100-channel-skill": { "enabled": true }
    },
    "load": {
      "paths": ["$ROOT/extensions/sta100-channel-skill"]
    }
  },
  "models": {
    "mode": "merge",
    "providers": {}
  },
  "auth": {
    "profiles": {}
  },
  "wizard": {
    "lastRunAt": "$(date -Iseconds)",
    "lastRunVersion": "2026.7.1-2",
    "lastRunCommand": "bootstrap",
    "lastRunMode": "local"
  },
  "meta": {
    "lastTouchedVersion": "2026.7.1-2",
    "lastTouchedAt": "$(date -Iseconds)"
  }
}
EOF

cp "$CONFIG_PATH" "$CONFIG_PATH.last-good"
cp "$CONFIG_PATH" "$CONFIG_PATH.bak"
