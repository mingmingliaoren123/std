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
if [[ -f "$ROOT_DIR/openclaw/config/openclaw.json" ]]; then
  export OPENCLAW_CONFIG_PATH="${OPENCLAW_CONFIG_PATH:-$ROOT_DIR/openclaw/config/openclaw.json}"
fi
export STA100_AGENT_MANIFEST="$MANIFEST"

mkdir -p "$BUILD_DIR" "$RUN_DIR" "$LOG_DIR"

is_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null
}

start() {
  if is_running; then
    echo "STA-100 web already running: pid $(cat "$PID_FILE")"
    return 0
  fi
  echo "building STA-100 web..."
  (cd "$WEB_DIR" && go build -trimpath -o "$BIN" .)
  echo "starting STA-100 web on $ADDR..."
  local -a env_args
  env_args=(OPENCLAW_BIN="$OPENCLAW_BIN" STA100_AGENT_MANIFEST="$STA100_AGENT_MANIFEST")
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
