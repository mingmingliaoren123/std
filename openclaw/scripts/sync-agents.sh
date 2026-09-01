#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_ROOT="$(cd "$ROOT/.." && pwd)"
ORCHESTRATOR_ROOT="${OPENCLAW_ORCHESTRATOR_ROOT:-$ROOT/../openclaw-orchestrator}"
ORCHESTRATOR_BIN="${OPENCLAW_ORCHESTRATOR_BIN:-$ORCHESTRATOR_ROOT/bin/openclaw-operator}"
MANIFEST="${OPENCLAW_AGENT_MANIFEST:-$ORCHESTRATOR_ROOT/../sta100-web/config/sta100-agents.json}"
OPENCLAW_HOME="${OPENCLAW_HOME:-$RELEASE_ROOT/data/openclaw-home}"
OPENCLAW_STATE_DIR="${OPENCLAW_STATE_DIR:-$RELEASE_ROOT/data/openclaw-state}"
OPENCLAW_PROFILE="${OPENCLAW_PROFILE:-sta100}"
OPENCLAW_GATEWAY_PORT="${OPENCLAW_GATEWAY_PORT:-18789}"

mkdir -p "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR"

if [[ -x "$ROOT/scripts/bootstrap-state.sh" ]]; then
  "$ROOT/scripts/bootstrap-state.sh"
fi

if [[ -x "$ORCHESTRATOR_BIN" ]]; then
exec env OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_STATE_DIR/openclaw.json" OPENCLAW_PROFILE="$OPENCLAW_PROFILE" OPENCLAW_GATEWAY_PORT="$OPENCLAW_GATEWAY_PORT" "$ORCHESTRATOR_BIN" --bin "$ROOT/bin/openclaw" --config "$OPENCLAW_STATE_DIR/openclaw.json" --manifest "$MANIFEST" agents sync "$@"
fi

command -v go >/dev/null 2>&1 || { echo "Go is required to run the OpenClaw orchestrator" >&2; exit 1; }
exec env OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_STATE_DIR/openclaw.json" OPENCLAW_PROFILE="$OPENCLAW_PROFILE" OPENCLAW_GATEWAY_PORT="$OPENCLAW_GATEWAY_PORT" go run "$ORCHESTRATOR_ROOT/cmd/openclaw-operator" --bin "$ROOT/bin/openclaw" --config "$OPENCLAW_STATE_DIR/openclaw.json" --manifest "$MANIFEST" agents sync "$@"
