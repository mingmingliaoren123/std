#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ORCHESTRATOR_ROOT="${OPENCLAW_ORCHESTRATOR_ROOT:-$ROOT/../openclaw-orchestrator}"
ORCHESTRATOR_BIN="${OPENCLAW_ORCHESTRATOR_BIN:-$ORCHESTRATOR_ROOT/bin/openclaw-operator}"
MANIFEST="${OPENCLAW_AGENT_MANIFEST:-$ORCHESTRATOR_ROOT/../sta100-web/config/sta100-agents.json}"

if [[ -x "$ORCHESTRATOR_BIN" ]]; then
  exec "$ORCHESTRATOR_BIN" --bin "$ROOT/bin/openclaw" --manifest "$MANIFEST" agents sync "$@"
fi

command -v go >/dev/null 2>&1 || { echo "Go is required to run the OpenClaw orchestrator" >&2; exit 1; }
exec go run "$ORCHESTRATOR_ROOT/cmd/openclaw-operator" --bin "$ROOT/bin/openclaw" --manifest "$MANIFEST" agents sync "$@"
