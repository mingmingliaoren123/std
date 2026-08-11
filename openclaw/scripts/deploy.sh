#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUNTIME="$ROOT/runtime"
APP="$ROOT/app"
PACKAGE="$ROOT/downloads/openclaw-2026.7.1-2.tgz"
CLI="$ROOT/bin/openclaw"
UNIT_SOURCE="$ROOT/systemd/openclaw-gateway.service"
UNIT_TARGET="$HOME/.config/systemd/user/openclaw-gateway.service"
EXPECTED_NODE="v22.22.3"
EXPECTED_OPENCLAW="2026.7.1-2"

export PATH="$RUNTIME/bin:$PATH"

[[ -x "$RUNTIME/bin/node" ]] || { echo "Missing Node runtime" >&2; exit 1; }
[[ -f "$PACKAGE" ]] || { echo "Missing OpenClaw package" >&2; exit 1; }
[[ "$($RUNTIME/bin/node --version)" == "$EXPECTED_NODE" ]] || { echo "Unexpected Node version" >&2; exit 1; }

printf '%s  %s\n' \
  '1c4a9933a5e45bc88f54f70b5f91232c127ec49f1a5989d23fb85824c7adf9b7' \
  "$ROOT/downloads/node-v22.22.3-linux-arm64.tar.xz" | sha256sum -c -
printf '%s  %s\n' \
  '5bb525f36f471a41239615d321c441778c7e1c007018ed6d84b795be77803276' \
  "$PACKAGE" | sha256sum -c -

systemctl --user stop openclaw-gateway.service 2>/dev/null || true
"$ROOT/scripts/backup-state.sh"

INSTALLED_VERSION=""
if [[ -f "$APP/node_modules/openclaw/package.json" ]]; then
  INSTALLED_VERSION="$($RUNTIME/bin/node -p "require('$APP/node_modules/openclaw/package.json').version")"
fi
if [[ "$INSTALLED_VERSION" != "$EXPECTED_OPENCLAW" ]] || [[ "${FORCE_REINSTALL:-0}" == "1" ]]; then
  "$RUNTIME/bin/npm" install \
    --prefix "$APP" \
    "$PACKAGE" \
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
  -e "s|Environment=PATH=.*|Environment=PATH=$RUNTIME/bin:$HOME/.local/bin:/usr/local/bin:/usr/bin:/bin|" \
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
