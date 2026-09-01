#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEB_DIR="$ROOT_DIR/sta100-web"
BUILD_DIR="$ROOT_DIR/.build"
RUN_DIR="$ROOT_DIR/.run"
LOG_DIR="$ROOT_DIR/logs"
BIN="$BUILD_DIR/sta100-web"
PID_FILE="$RUN_DIR/sta100-web.pid"
LOG_FILE="$LOG_DIR/sta100-web.log"
ADDR="${STA100_ADDR:-127.0.0.1:18080}"
MANIFEST="${STA100_AGENT_MANIFEST:-$WEB_DIR/config/sta100-agents.json}"

export OPENCLAW_BIN="${OPENCLAW_BIN:-$ROOT_DIR/openclaw/bin/openclaw}"
export OPENCLAW_HOME="${OPENCLAW_HOME:-$ROOT_DIR/data/openclaw-home}"
export OPENCLAW_STATE_DIR="${OPENCLAW_STATE_DIR:-$ROOT_DIR/data/openclaw-state}"
export OPENCLAW_CONFIG_PATH="${OPENCLAW_CONFIG_PATH:-$OPENCLAW_STATE_DIR/openclaw.json}"
export OPENCLAW_GATEWAY_PORT="${OPENCLAW_GATEWAY_PORT:-18789}"
export OPENCLAW_GATEWAY_PORT
export STA100_AGENT_MANIFEST="$MANIFEST"

mkdir -p "$BUILD_DIR" "$RUN_DIR" "$LOG_DIR" "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR"

bootstrap_openclaw_state() {
  if [[ -x "$ROOT_DIR/openclaw/scripts/bootstrap-state.sh" ]]; then
    "$ROOT_DIR/openclaw/scripts/bootstrap-state.sh"
  fi
}

start_openclaw_gateway() {
  local gateway_pid_file="$RUN_DIR/openclaw-gateway.pid"
  local gateway_log="$LOG_DIR/openclaw-gateway.log"
  if [[ -f "$gateway_pid_file" ]]; then
    local pid
    pid="$(cat "$gateway_pid_file" 2>/dev/null || true)"
    if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
      return 0
    fi
    rm -f "$gateway_pid_file" 2>/dev/null || true
  fi
  mkdir -p "$LOG_DIR"
  printf '%s start request at %s\n' "openclaw-gateway" "$(date -Iseconds)" >>"$gateway_log"
  nohup setsid "$OPENCLAW_BIN" gateway --port "${OPENCLAW_GATEWAY_PORT:-18789}" >>"$gateway_log" 2>&1 </dev/null &
  echo $! > "$gateway_pid_file"
  sleep 1
  if ! kill -0 "$(cat "$gateway_pid_file")" 2>/dev/null; then
    echo "warning: OpenClaw gateway failed to start. Recent log:" >&2
    tail -40 "$gateway_log" >&2 || true
    rm -f "$gateway_pid_file"
    return 1
  fi
}

wait_for_openclaw_gateway() {
  local attempts="${1:-30}"
  local gateway_url="ws://127.0.0.1:${OPENCLAW_GATEWAY_PORT:-18789}"
  for _ in $(seq 1 "$attempts"); do
    if env OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_CONFIG_PATH" OPENCLAW_GATEWAY_PORT="$OPENCLAW_GATEWAY_PORT" "$OPENCLAW_BIN" gateway status --url "$gateway_url" --json --timeout 2000 2>/dev/null | tr -d '\n' | grep -Eq '"rpc"[[:space:]]*:[[:space:]]*\{[^}]*"ok"[[:space:]]*:[[:space:]]*true'; then
      return 0
    fi
    sleep 1
  done
  return 1
}

sync_openclaw_agents() {
  if [[ "${STA100_SYNC_AGENTS_ON_START:-1}" != "1" ]] || [[ ! -x "$ROOT_DIR/openclaw/scripts/sync-agents.sh" ]]; then
    return 0
  fi
  local sync_log="$LOG_DIR/openclaw-sync.log"
  mkdir -p "$LOG_DIR"
  if command -v timeout >/dev/null 2>&1; then
    if timeout "${STA100_SYNC_TIMEOUT:-120}"s "$ROOT_DIR/openclaw/scripts/sync-agents.sh" >>"$sync_log" 2>&1; then
      return 0
    fi
  else
    if "$ROOT_DIR/openclaw/scripts/sync-agents.sh" >>"$sync_log" 2>&1; then
      return 0
    fi
  fi
  echo "warning: OpenClaw Agent sync failed or timed out; STA-100 will still start, but agent calls may be unavailable." >&2
  tail -40 "$sync_log" >&2 || true
  return 1
}

is_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null
}

openclaw_gateway_running() {
  local gateway_url="ws://127.0.0.1:${OPENCLAW_GATEWAY_PORT:-18789}"
  env OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_CONFIG_PATH" OPENCLAW_GATEWAY_PORT="$OPENCLAW_GATEWAY_PORT" "$OPENCLAW_BIN" gateway status --url "$gateway_url" --json --timeout 2000 2>/dev/null | tr -d '\n' | grep -Eq '"rpc"[[:space:]]*:[[:space:]]*\{[^}]*"ok"[[:space:]]*:[[:space:]]*true' && return 0
  local gateway_pid_file="$RUN_DIR/openclaw-gateway.pid"
  [[ -f "$gateway_pid_file" ]] || return 1
  local pid
  pid="$(cat "$gateway_pid_file" 2>/dev/null || true)"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null
}

start() {
  if is_running; then
    echo "STA-100 web already running: pid $(cat "$PID_FILE")"
    return 0
  fi
  bootstrap_openclaw_state
  sync_openclaw_agents || true
  start_openclaw_gateway || echo "warning: OpenClaw gateway is unavailable; STA-100 will still start, but OpenClaw features may be unavailable." >&2
  if ! wait_for_openclaw_gateway "${STA100_GATEWAY_READY_TIMEOUT:-120}"; then
    echo "warning: OpenClaw gateway did not become ready in time; STA-100 will still start, but agent calls may be unavailable." >&2
  fi
  echo "building STA-100 web..."
  (cd "$WEB_DIR" && go build -trimpath -o "$BIN" .)
  echo "starting STA-100 web on $ADDR..."
  local -a env_args
  env_args=(OPENCLAW_BIN="$OPENCLAW_BIN" OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_CONFIG_PATH" STA100_AGENT_MANIFEST="$STA100_AGENT_MANIFEST")
  if [[ -n "${OPENCLAW_CONFIG_PATH:-}" ]]; then
    env_args+=(OPENCLAW_CONFIG_PATH="$OPENCLAW_CONFIG_PATH")
  fi
  setsid -f env "${env_args[@]}" "$BIN" -addr "$ADDR" >>"$LOG_FILE" 2>&1 </dev/null
  for _ in {1..20}; do
    if ss -ltn 2>/dev/null | awk '{print $4}' | grep -Fq "$ADDR"; then
      local pid
      pid="$(ps -eo pid,cmd | grep -F "$BIN -addr $ADDR" | grep -v grep | awk 'NR==1{print $1}')"
      if [[ -n "${pid:-}" ]]; then
        echo "$pid" >"$PID_FILE"
      fi
      echo "STA-100 web started: pid ${pid:-unknown}"
      echo "log: $LOG_FILE"
      return 0
    fi
    sleep 0.5
  done
  echo "STA-100 web failed to start"
  tail -n 40 "$LOG_FILE" || true
  return 1
}

stop() {
  if ! [[ -f "$PID_FILE" ]]; then
    echo "STA-100 web is not running"
    return 0
  fi
  local pid
  pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  if [[ -z "$pid" ]]; then
    rm -f "$PID_FILE"
    echo "STA-100 web is not running"
    return 0
  fi
  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid"
    for _ in {1..20}; do
      if ! kill -0 "$pid" 2>/dev/null; then
        break
      fi
      sleep 0.2
    done
    if kill -0 "$pid" 2>/dev/null; then
      kill -9 "$pid" 2>/dev/null || true
    fi
  fi
  rm -f "$PID_FILE"
  echo "STA-100 web stopped"
}

status() {
  if is_running; then
    echo "STA-100 web running: pid $(cat "$PID_FILE"), addr $ADDR"
    return 0
  fi
  echo "STA-100 web not running"
  return 1
}

logs() {
  tail -n "${1:-80}" "$LOG_FILE"
}

case "${1:-start}" in
  start) start ;;
  stop) stop ;;
  restart) stop || true; start ;;
  status) status ;;
  logs) logs "${2:-80}" ;;
  *)
    echo "usage: $0 {start|stop|restart|status|logs [lines]}"
    exit 2
    ;;
esac
