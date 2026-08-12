#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP="$ROOT/app"
CLI="$ROOT/bin/openclaw"
UNIT_SOURCE="$ROOT/systemd/openclaw-gateway.service"
UNIT_TARGET="$HOME/.config/systemd/user/openclaw-gateway.service"
EXPECTED_OPENCLAW="2026.7.1-2"

if [[ -n "${OPENCLAW_NODE:-}" ]]; then
  NODE="$OPENCLAW_NODE"
elif [[ -x /usr/bin/node ]]; then
  NODE=/usr/bin/node
else
  NODE="$(command -v node || true)"
fi

if [[ -n "${OPENCLAW_NPM:-}" ]]; then
  NPM="$OPENCLAW_NPM"
elif [[ -x /usr/bin/npm ]]; then
  NPM=/usr/bin/npm
else
  NPM="$(command -v npm || true)"
fi

[[ -x "$NODE" ]] || { echo "Node.js is not available in the current box environment" >&2; exit 1; }
[[ -x "$NPM" ]] || { echo "npm is not available in the current box environment" >&2; exit 1; }
"$NODE" -e '
  const [major, minor, patch] = process.versions.node.split(".").map(Number);
  const compatible =
    (major === 22 && (minor > 22 || (minor === 22 && patch >= 3))) ||
    (major === 24 && minor >= 15) ||
    (major >= 25 && (major > 25 || minor >= 9));
  process.exit(compatible ? 0 : 1);
' || { echo "OpenClaw requires Node.js >=22.22.3 <23, >=24.15.0 <25, or >=25.9.0; found $($NODE --version)" >&2; exit 1; }

systemctl --user stop openclaw-gateway.service 2>/dev/null || true
"$ROOT/scripts/backup-state.sh"

INSTALLED_VERSION=""
if [[ -f "$APP/node_modules/openclaw/package.json" ]]; then
  INSTALLED_VERSION="$($NODE -p "require('$APP/node_modules/openclaw/package.json').version")"
fi
if [[ "$INSTALLED_VERSION" != "$EXPECTED_OPENCLAW" ]] || [[ "${FORCE_REINSTALL:-0}" == "1" ]]; then
  "$NPM" install \
    --prefix "$APP" \
    "openclaw@$EXPECTED_OPENCLAW" \
    --omit=dev \
    --registry="${NPM_CONFIG_REGISTRY:-https://registry.npmmirror.com}"
fi

ACTUAL_VERSION="$($CLI --version | sed -n 's/^OpenClaw \([^ ]*\).*/\1/p')"
[[ "$ACTUAL_VERSION" == "$EXPECTED_OPENCLAW" ]] || {
  echo "OpenClaw version mismatch: $ACTUAL_VERSION" >&2
  exit 1
}

mkdir -p "$(dirname "$UNIT_TARGET")"
sed \
  -e "s|Environment=HOME=.*|Environment=HOME=$HOME|" \
  -e "s|Environment=PATH=.*|Environment=PATH=$(dirname "$NODE"):$HOME/.local/bin:/usr/local/bin:/usr/bin:/bin|" \
  "$UNIT_SOURCE" > "$UNIT_TARGET"
chmod 0644 "$UNIT_TARGET"
systemctl --user daemon-reload
systemctl --user enable --now openclaw-gateway.service

for _ in $(seq 1 30); do
  if "$CLI" gateway status --json --timeout 2000 >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

"$CLI" gateway status --json --timeout 5000 >/dev/null

if [[ "${REMOVE_LEGACY:-1}" == "1" ]] && [[ -e /usr/local/lib/node_modules/openclaw ]]; then
  if sudo -n true 2>/dev/null; then
    sudo -n /usr/local/bin/npm uninstall -g openclaw
  else
    echo "Legacy global OpenClaw remains because passwordless sudo is unavailable." >&2
  fi
fi

echo "OpenClaw $EXPECTED_OPENCLAW deployed from $ROOT"
