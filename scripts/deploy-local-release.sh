#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_DIR="${RELEASE_DIR:-$ROOT/release/sta100-release}"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-arm64}"

STA100_WEB_DIR="$ROOT/sta100-web"
OPENCLAW_DIR="$ROOT/openclaw"
ORCH_DIR="$ROOT/openclaw-orchestrator"
OFFICE_RUNTIME_DIR="$RELEASE_DIR/office"
SHARED_KNOWLEDGE_SOURCE="${STA100_SHARED_KNOWLEDGE_SOURCE:-/home/User/tmp/0831}"

require_file() {
  [[ -e "$1" ]] || { echo "missing required path: $1" >&2; exit 1; }
}

write_release_launcher() {
  mkdir -p "$RELEASE_DIR/scripts"
  cat > "$RELEASE_DIR/scripts/sta100.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CMD="${1:-start}"
LOG_DIR="$ROOT/logs"
DATA_DIR="$ROOT/data"
RUN_DIR="$ROOT/.run"
ASSET_DIR="$ROOT/assets"
EMOJI_FONT="$ASSET_DIR/NotoColorEmoji.ttf"
OPENCLAW_HOME="${OPENCLAW_HOME:-$DATA_DIR/openclaw-home}"
OPENCLAW_STATE_DIR="${OPENCLAW_STATE_DIR:-$DATA_DIR/openclaw-state}"
OPENCLAW_CONFIG_PATH="${OPENCLAW_CONFIG_PATH:-$OPENCLAW_STATE_DIR/openclaw.json}"
OPENCLAW_PROFILE="${OPENCLAW_PROFILE:-sta100}"
OPENCLAW_GATEWAY_PORT="${OPENCLAW_GATEWAY_PORT:-18789}"
OPENCLAW_BIN="${OPENCLAW_BIN:-$ROOT/openclaw/bin/openclaw}"
STA100_ADDR="${STA100_ADDR:-0.0.0.0:18080}"
STA100_AGENT_MANIFEST="${STA100_AGENT_MANIFEST:-$ROOT/config/sta100-agents.json}"
STA100_DB_PATH="${STA100_DB_PATH:-$DATA_DIR/sta100.db}"
STA100_AUTH_FILE="${STA100_AUTH_FILE:-$DATA_DIR/auth.json}"
STA100_TEMPLATE_DIR="${STA100_TEMPLATE_DIR:-$DATA_DIR/templates}"
STA100_PRIVATE_DATA_DIR="${STA100_PRIVATE_DATA_DIR:-$DATA_DIR/private-files}"
STA100_KNOWLEDGE_DATA_DIR="${STA100_KNOWLEDGE_DATA_DIR:-$DATA_DIR/knowledge}"
STA100_SHARED_KNOWLEDGE_DIR="${STA100_SHARED_KNOWLEDGE_DIR:-$ROOT/knowledge/shared}"
STA100_OFFICE_BIN="${STA100_OFFICE_BIN:-$ROOT/office/bin/soffice}"
STA100_SYNC_AGENTS_ON_START="${STA100_SYNC_AGENTS_ON_START:-1}"
PID_FILE="$RUN_DIR/sta100-web.pid"
GATEWAY_PID_FILE="$RUN_DIR/openclaw-gateway.pid"
BIN="$ROOT/bin/sta100-web"
SERVICE_NAME="sta100-web.service"
GATEWAY_SERVICE_NAME="sta100-openclaw-gateway.service"
SYSTEMD_DIR="/etc/systemd/system"
CRON_MARKER="# STA100_AUTOSTART"

mkdir -p "$RUN_DIR" "$LOG_DIR" "$DATA_DIR" "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR" "$STA100_KNOWLEDGE_DATA_DIR" "$STA100_SHARED_KNOWLEDGE_DIR"

export OPENCLAW_BIN OPENCLAW_HOME OPENCLAW_STATE_DIR OPENCLAW_CONFIG_PATH OPENCLAW_PROFILE OPENCLAW_GATEWAY_PORT
export STA100_AGENT_MANIFEST STA100_DB_PATH STA100_AUTH_FILE STA100_TEMPLATE_DIR STA100_PRIVATE_DATA_DIR STA100_KNOWLEDGE_DATA_DIR STA100_SHARED_KNOWLEDGE_DIR STA100_OFFICE_BIN

gateway_url() {
  printf 'ws://127.0.0.1:%s' "$OPENCLAW_GATEWAY_PORT"
}

sta100_web_port() {
  local port="${STA100_ADDR##*:}"
  if [[ "$port" =~ ^[0-9]+$ ]]; then
    printf '%s' "$port"
  else
    printf '18080'
  fi
}

systemd_active() {
  command -v systemctl >/dev/null 2>&1 || return 1
  systemctl is-active --quiet "$1" 2>/dev/null
}

systemd_main_pid() {
  command -v systemctl >/dev/null 2>&1 || return 0
  systemctl show -p MainPID --value "$1" 2>/dev/null | awk '$1 != "0" {print $1; exit}'
}

openclaw_gateway_ready() {
  "$OPENCLAW_BIN" gateway status --url "$(gateway_url)" --json --timeout 2000 2>/dev/null | tr -d '\n' | grep -Eq '"rpc"[[:space:]]*:[[:space:]]*\{[^}]*"ok"[[:space:]]*:[[:space:]]*true'
}

openclaw_gateway_service_running() {
  systemd_active "$GATEWAY_SERVICE_NAME"
}

openclaw_gateway_running() {
  openclaw_gateway_ready && return 0
  openclaw_gateway_service_running && return 0
  [[ -f "$GATEWAY_PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$GATEWAY_PID_FILE" 2>/dev/null || true)"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  [[ -r "/proc/$pid/cmdline" ]] && tr '\0' ' ' < "/proc/$pid/cmdline" | grep -Eq "openclaw-gateway|openclaw.mjs gateway|$ROOT/openclaw/bin/openclaw"
}

openclaw_gateway_status_text() {
  if openclaw_gateway_ready; then
    local pid
    pid="$(systemd_main_pid "$GATEWAY_SERVICE_NAME")"
    if [[ -n "$pid" ]]; then
      echo "running: systemd pid $pid, ready"
    elif [[ -f "$GATEWAY_PID_FILE" ]]; then
      echo "running: pid $(cat "$GATEWAY_PID_FILE"), ready"
    else
      echo "running: ready"
    fi
    return 0
  fi
  if openclaw_gateway_service_running; then
    local pid
    pid="$(systemd_main_pid "$GATEWAY_SERVICE_NAME")"
    echo "starting: systemd${pid:+ pid $pid}, gateway not ready"
    return 0
  fi
  if [[ -f "$GATEWAY_PID_FILE" ]]; then
    local pid
    pid="$(cat "$GATEWAY_PID_FILE" 2>/dev/null || true)"
    if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
      echo "starting: pid $pid, gateway not ready"
      return 0
    fi
  fi
  echo "not running"
  return 1
}

bootstrap_openclaw_state() {
  if [[ -x "$ROOT/openclaw/scripts/bootstrap-state.sh" ]]; then
    "$ROOT/openclaw/scripts/bootstrap-state.sh"
  fi
}

sync_openclaw_agents() {
  if [[ "$STA100_SYNC_AGENTS_ON_START" != "1" ]] || [[ ! -x "$ROOT/openclaw/scripts/sync-agents.sh" ]]; then
    return 0
  fi
  local sync_log="$LOG_DIR/openclaw-sync.log"
  if command -v timeout >/dev/null 2>&1; then
    timeout "${STA100_SYNC_TIMEOUT:-120}"s "$ROOT/openclaw/scripts/sync-agents.sh" >>"$sync_log" 2>&1
  else
    "$ROOT/openclaw/scripts/sync-agents.sh" >>"$sync_log" 2>&1
  fi
}

start_openclaw_gateway() {
  if openclaw_gateway_running; then
    return 0
  fi
  bootstrap_openclaw_state
  printf '%s start request at %s\n' "openclaw-gateway" "$(date -Iseconds)" >>"$LOG_DIR/openclaw-gateway.log"
  nohup setsid "$OPENCLAW_BIN" gateway --port "$OPENCLAW_GATEWAY_PORT" >>"$LOG_DIR/openclaw-gateway.log" 2>&1 </dev/null &
  echo $! > "$GATEWAY_PID_FILE"
  sleep 0.8
  if ! kill -0 "$(cat "$GATEWAY_PID_FILE")" 2>/dev/null; then
    echo "warning: OpenClaw gateway failed to start. Recent log:" >&2
    tail -40 "$LOG_DIR/openclaw-gateway.log" >&2 || true
    rm -f "$GATEWAY_PID_FILE"
    return 1
  fi
}

wait_for_openclaw_gateway() {
  local attempts="${1:-30}"
  for _ in $(seq 1 "$attempts"); do
    openclaw_gateway_ready && return 0
    sleep 1
  done
  return 1
}

web_running() {
  web_pid_file_running && return 0
  systemd_active "$SERVICE_NAME" && return 0
  web_http_ready && return 0
  return 1
}

web_pid_file_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  [[ -r "/proc/$pid/cmdline" ]] && tr '\0' ' ' < "/proc/$pid/cmdline" | grep -q "$BIN"
}

web_http_ready() {
  command -v curl >/dev/null 2>&1 || return 1
  curl -fsS --max-time 2 "http://127.0.0.1:$(sta100_web_port)/api/health" >/dev/null 2>&1
}

web_status_text() {
  if web_pid_file_running; then
    echo "running: pid $(cat "$PID_FILE")"
    return 0
  fi
  if systemd_active "$SERVICE_NAME"; then
    local pid
    pid="$(systemd_main_pid "$SERVICE_NAME")"
    if web_http_ready; then
      echo "running: systemd${pid:+ pid $pid}, healthy"
    else
      echo "starting: systemd${pid:+ pid $pid}, health not ready"
    fi
    return 0
  fi
  if web_http_ready; then
    echo "running: http healthy"
    return 0
  fi
  echo "not running"
  return 1
}

start() {
  if web_running; then
    echo "sta100-web already running: $(cat "$PID_FILE")"
    return 0
  fi
  if [[ -f "$PID_FILE" ]]; then
    rm -f "$PID_FILE"
  fi
  bootstrap_openclaw_state
  sync_openclaw_agents || {
    echo "warning: OpenClaw Agent sync failed or timed out; STA-100 will still start, but agent calls may be unavailable." >&2
    tail -40 "$LOG_DIR/openclaw-sync.log" >&2 || true
  }
  start_openclaw_gateway || echo "warning: OpenClaw gateway is unavailable; STA-100 will still start, but OpenClaw features may be unavailable." >&2
  if ! wait_for_openclaw_gateway "${STA100_GATEWAY_READY_TIMEOUT:-120}"; then
    echo "warning: OpenClaw gateway did not become ready in time; STA-100 will still start, but agent calls may be unavailable." >&2
  fi
  printf '%s start request at %s\n' "sta100-web" "$(date -Iseconds)" >>"$LOG_DIR/sta100-web.log"
  nohup setsid "$BIN" -addr "$STA100_ADDR" >>"$LOG_DIR/sta100-web.log" 2>&1 </dev/null &
  echo $! > "$PID_FILE"
  sleep 0.5
  if ! kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
    echo "sta100-web failed to start. Recent log:" >&2
    tail -40 "$LOG_DIR/sta100-web.log" >&2 || true
    rm -f "$PID_FILE"
    exit 1
  fi
  echo "started: $(cat "$PID_FILE")"
}

stop_openclaw_gateway() {
  [[ -f "$GATEWAY_PID_FILE" ]] || return 0
  local pid
  pid="$(cat "$GATEWAY_PID_FILE" 2>/dev/null || true)"
  if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" || true
  fi
  rm -f "$GATEWAY_PID_FILE"
}

stop() {
  if [[ -f "$PID_FILE" ]]; then
    local pid
    pid="$(cat "$PID_FILE" 2>/dev/null || true)"
    if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
      kill "$pid" || true
    fi
    rm -f "$PID_FILE"
  fi
  stop_openclaw_gateway
  echo "stopped"
}

status() {
  local web_status gateway_status
  web_status="$(web_status_text)"
  gateway_status="$(openclaw_gateway_status_text)"
  echo "sta100-web: $web_status"
  echo "openclaw-gateway: $gateway_status"
  [[ "$web_status" == running:* || "$web_status" == starting:* ]] && [[ "$gateway_status" == running:* || "$gateway_status" == starting:* ]]
}

logs() {
  tail -n "${2:-80}" "$LOG_DIR/${1:-sta100-web}.log"
}

install_autostart() {
  command -v systemctl >/dev/null 2>&1 || { echo "systemctl not found; Ubuntu systemd is required for autostart" >&2; return 1; }
  local service_user service_group tmp_gateway tmp_web
  service_user="${STA100_SERVICE_USER:-${SUDO_USER:-$(id -un)}}"
  service_group="${STA100_SERVICE_GROUP:-$(id -gn "$service_user" 2>/dev/null || id -gn)}"
  mkdir -p "$RUN_DIR" "$LOG_DIR" "$DATA_DIR" "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR" "$STA100_PRIVATE_DATA_DIR" "$STA100_KNOWLEDGE_DATA_DIR" "$STA100_SHARED_KNOWLEDGE_DIR"
  chmod +x "$BIN" "$OPENCLAW_BIN" "$ROOT/openclaw/scripts/sync-agents.sh" "$ROOT/openclaw/scripts/bootstrap-state.sh" 2>/dev/null || true

  run_root() {
    if [[ "$(id -u)" -eq 0 ]]; then
      "$@"
    else
      sudo "$@"
    fi
  }

  run_root chown -R "$service_user:$service_group" "$RUN_DIR" "$LOG_DIR" "$DATA_DIR" "$STA100_PRIVATE_DATA_DIR" "$STA100_KNOWLEDGE_DATA_DIR" 2>/dev/null || true
  run_root chmod 0755 "$ROOT" "$RUN_DIR" "$LOG_DIR" "$DATA_DIR" "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR" "$STA100_PRIVATE_DATA_DIR" "$STA100_KNOWLEDGE_DATA_DIR" 2>/dev/null || true
  run_root chmod 0700 "$STA100_PRIVATE_DATA_DIR" "$STA100_KNOWLEDGE_DATA_DIR" 2>/dev/null || true
  if [[ -f "$EMOJI_FONT" ]]; then
    run_root install -d -m 0755 /usr/local/share/fonts/sta100 2>/dev/null || true
    run_root install -m 0644 "$EMOJI_FONT" /usr/local/share/fonts/sta100/NotoColorEmoji.ttf 2>/dev/null || true
    if command -v fc-cache >/dev/null 2>&1; then
      run_root fc-cache -f /usr/local/share/fonts/sta100 2>/dev/null || true
    fi
  fi

  tmp_gateway="$(mktemp)"
  tmp_web="$(mktemp)"
  trap 'rm -f "$tmp_gateway" "$tmp_web"' RETURN

  cat > "$tmp_gateway" <<SERVICE
[Unit]
Description=STA-100 bundled OpenClaw Gateway
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$service_user
Group=$service_group
WorkingDirectory=$ROOT
Environment=OPENCLAW_HOME=$ROOT/data/openclaw-home
Environment=OPENCLAW_STATE_DIR=$ROOT/data/openclaw-state
Environment=OPENCLAW_CONFIG_PATH=$ROOT/data/openclaw-state/openclaw.json
Environment=OPENCLAW_PROFILE=sta100
Environment=OPENCLAW_GATEWAY_PORT=${OPENCLAW_GATEWAY_PORT}
ExecStartPre=-$ROOT/openclaw/scripts/bootstrap-state.sh
ExecStartPre=-/usr/bin/env bash -lc 'if command -v timeout >/dev/null 2>&1; then timeout "${STA100_SYNC_TIMEOUT:-120}"s "$ROOT/openclaw/scripts/sync-agents.sh" >>"$ROOT/logs/openclaw-sync.log" 2>&1; else "$ROOT/openclaw/scripts/sync-agents.sh" >>"$ROOT/logs/openclaw-sync.log" 2>&1; fi'
ExecStart=$ROOT/openclaw/bin/openclaw gateway --port ${OPENCLAW_GATEWAY_PORT}
Restart=always
RestartSec=5
StandardOutput=append:$ROOT/logs/openclaw-gateway.log
StandardError=append:$ROOT/logs/openclaw-gateway.log

[Install]
WantedBy=multi-user.target
SERVICE
  cat > "$tmp_web" <<SERVICE
[Unit]
Description=STA-100 Web
After=network-online.target
Requires=$GATEWAY_SERVICE_NAME
After=$GATEWAY_SERVICE_NAME

[Service]
Type=simple
User=$service_user
Group=$service_group
WorkingDirectory=$ROOT
Environment=STA100_ADDR=${STA100_ADDR}
Environment=OPENCLAW_BIN=$ROOT/openclaw/bin/openclaw
Environment=OPENCLAW_HOME=$ROOT/data/openclaw-home
Environment=OPENCLAW_STATE_DIR=$ROOT/data/openclaw-state
Environment=OPENCLAW_CONFIG_PATH=$ROOT/data/openclaw-state/openclaw.json
Environment=OPENCLAW_PROFILE=sta100
Environment=OPENCLAW_GATEWAY_PORT=${OPENCLAW_GATEWAY_PORT}
Environment=STA100_AGENT_MANIFEST=$ROOT/config/sta100-agents.json
Environment=STA100_DB_PATH=$ROOT/data/sta100.db
Environment=STA100_AUTH_FILE=$ROOT/data/auth.json
Environment=STA100_TEMPLATE_DIR=$ROOT/data/templates
Environment=STA100_PRIVATE_DATA_DIR=${STA100_PRIVATE_DATA_DIR}
Environment=STA100_KNOWLEDGE_DATA_DIR=${STA100_KNOWLEDGE_DATA_DIR}
Environment=STA100_SHARED_KNOWLEDGE_DIR=${STA100_SHARED_KNOWLEDGE_DIR}
Environment=STA100_OFFICE_BIN=$ROOT/office/bin/soffice
ExecStart=$ROOT/bin/sta100-web -addr ${STA100_ADDR}
Restart=always
RestartSec=5
StandardOutput=append:$ROOT/logs/sta100-web.log
StandardError=append:$ROOT/logs/sta100-web.log

[Install]
WantedBy=multi-user.target
SERVICE

  stop || true
  run_root install -m 0644 "$tmp_gateway" "$SYSTEMD_DIR/$GATEWAY_SERVICE_NAME"
  run_root install -m 0644 "$tmp_web" "$SYSTEMD_DIR/$SERVICE_NAME"
  run_root systemctl daemon-reload
  run_root systemctl enable "$GATEWAY_SERVICE_NAME" "$SERVICE_NAME"
  run_root systemctl restart "$GATEWAY_SERVICE_NAME"
  run_root systemctl restart "$SERVICE_NAME"
  echo "autostart installed with systemd system services"
  echo "services: $GATEWAY_SERVICE_NAME, $SERVICE_NAME"
  echo "user: $service_user"
}

uninstall_autostart() {
  if command -v systemctl >/dev/null 2>&1; then
    if [[ "$(id -u)" -eq 0 ]]; then
      systemctl disable --now "$SERVICE_NAME" 2>/dev/null || true
      systemctl disable --now "$GATEWAY_SERVICE_NAME" 2>/dev/null || true
      rm -f "$SYSTEMD_DIR/$SERVICE_NAME" "$SYSTEMD_DIR/$GATEWAY_SERVICE_NAME"
      systemctl daemon-reload 2>/dev/null || true
    else
      sudo systemctl disable --now "$SERVICE_NAME" 2>/dev/null || true
      sudo systemctl disable --now "$GATEWAY_SERVICE_NAME" 2>/dev/null || true
      sudo rm -f "$SYSTEMD_DIR/$SERVICE_NAME" "$SYSTEMD_DIR/$GATEWAY_SERVICE_NAME"
      sudo systemctl daemon-reload 2>/dev/null || true
    fi
  fi
  if command -v crontab >/dev/null 2>&1; then
    crontab -l 2>/dev/null | grep -v "$CRON_MARKER" | crontab - 2>/dev/null || true
  fi
  echo "autostart uninstalled"
}

case "$CMD" in
  start) start ;;
  stop) stop ;;
  restart) stop || true; start ;;
  status) status ;;
  logs) logs "${2:-sta100-web}" "${3:-80}" ;;
  install-autostart) install_autostart ;;
  uninstall-autostart) uninstall_autostart ;;
  *) echo "usage: $0 {start|stop|restart|status|logs [name] [lines]|install-autostart|uninstall-autostart}" ; exit 2 ;;
esac
EOF
  chmod +x "$RELEASE_DIR/scripts/sta100.sh"
}

write_release_sync_script() {
  mkdir -p "$RELEASE_DIR/openclaw/scripts"
  cat > "$RELEASE_DIR/openclaw/scripts/sync-agents.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_ROOT="$(cd "$ROOT/.." && pwd)"
OPENCLAW_HOME="${OPENCLAW_HOME:-$RELEASE_ROOT/data/openclaw-home}"
OPENCLAW_STATE_DIR="${OPENCLAW_STATE_DIR:-$RELEASE_ROOT/data/openclaw-state}"
OPENCLAW_PROFILE="${OPENCLAW_PROFILE:-sta100}"
OPENCLAW_GATEWAY_PORT="${OPENCLAW_GATEWAY_PORT:-18789}"
OPERATOR="${OPENCLAW_ORCHESTRATOR_BIN:-$RELEASE_ROOT/bin/openclaw-operator}"
OPENCLAW="${OPENCLAW_BIN:-$ROOT/bin/openclaw}"
MANIFEST="${OPENCLAW_AGENT_MANIFEST:-$RELEASE_ROOT/config/sta100-agents.json}"

[[ -x "$OPERATOR" ]] || { echo "openclaw-operator missing: $OPERATOR" >&2; exit 1; }
[[ -x "$OPENCLAW" ]] || { echo "openclaw CLI missing: $OPENCLAW" >&2; exit 1; }
[[ -f "$MANIFEST" ]] || { echo "agent manifest missing: $MANIFEST" >&2; exit 1; }
mkdir -p "$OPENCLAW_HOME" "$OPENCLAW_STATE_DIR"

if [[ -x "$ROOT/scripts/bootstrap-state.sh" ]]; then
  "$ROOT/scripts/bootstrap-state.sh"
fi

exec env OPENCLAW_HOME="$OPENCLAW_HOME" OPENCLAW_STATE_DIR="$OPENCLAW_STATE_DIR" OPENCLAW_CONFIG_PATH="$OPENCLAW_STATE_DIR/openclaw.json" OPENCLAW_PROFILE="$OPENCLAW_PROFILE" OPENCLAW_GATEWAY_PORT="$OPENCLAW_GATEWAY_PORT" "$OPERATOR" --bin "$OPENCLAW" --config "$OPENCLAW_STATE_DIR/openclaw.json" --manifest "$MANIFEST" agents sync "$@"
EOF
  chmod +x "$RELEASE_DIR/openclaw/scripts/sync-agents.sh"
}

write_release_readme() {
  cat > "$RELEASE_DIR/README-DEPLOY.md" <<'EOF'
# STA-100 Release Deploy

启动：

```bash
./scripts/sta100.sh start
```

安装为 Ubuntu 开机自启动服务：

```bash
sudo ./scripts/sta100.sh install-autostart
```

常用命令：

```bash
./scripts/sta100.sh start
./scripts/sta100.sh stop
./scripts/sta100.sh restart
./scripts/sta100.sh status
./scripts/sta100.sh logs sta100-web 120
./scripts/sta100.sh logs openclaw-gateway 120
sudo ./scripts/sta100.sh install-autostart
sudo ./scripts/sta100.sh uninstall-autostart
sudo systemctl status sta100-web.service
sudo systemctl status sta100-openclaw-gateway.service
```

默认配置：

- Web 端口：`0.0.0.0:18080`
- OpenClaw Gateway：`127.0.0.1:18789`
- OpenClaw CLI：`./openclaw/bin/openclaw`
- OpenClaw 状态：`./data/openclaw-state`
- STA-100 SQLite：`./data/sta100.db`
- 共有知识库：`./knowledge/shared`（随部署包发布）
- 私有知识库和本机向量索引：`./data/knowledge`（仅保存在本机）
- 登录默认账户：`admin/admin`

说明：

- 本部署目录只使用包内 OpenClaw，不依赖系统 `openclaw` 命令。
- 首次启动不会预置模型、API Key 或通道插件，需要在设置中配置模型并测试通过。
- 启动时会同步 STA-100 Agent 清单；同步失败时 Web 仍会启动，但 Agent 功能不可用。
- 启动后会后台扫描共有资料和用户上传资料，解析后按 Agent 专业方向建立本机向量索引。TXT、Markdown、CSV/TSV、JSON、XML、HTML、RTF、DOCX、XLSX 可直接索引；PDF 会保存但在当前包内标记为待处理，需 OCR 或先转换为文本。
- 可通过绝对路径把知识库放到 TF 卡等容量更大的介质：`STA100_KNOWLEDGE_DATA_DIR=/mnt/tfcard/sta100/knowledge STA100_PRIVATE_DATA_DIR=/mnt/tfcard/sta100/private-files ./scripts/sta100.sh start`。安装 systemd 时请用同样的环境变量执行 `install-autostart`。
- 知识库状态接口：`GET /api/v1/private-files/knowledge/status`；手动同步：`POST /api/v1/private-files/knowledge`。
- `install-autostart` 会安装系统级 systemd 服务：`sta100-openclaw-gateway.service` 和 `sta100-web.service`，重启主机后自动启动。
- systemd 服务会使用执行安装命令的用户运行应用，并自动修正 `data/`、`logs/`、`.run/` 目录权限；如需指定用户，可执行 `sudo STA100_SERVICE_USER=your_user ./scripts/sta100.sh install-autostart`。
- 部署包自带 emoji 字体；执行 `install-autostart` 时会尝试安装到 `/usr/local/share/fonts/sta100` 并刷新字体缓存。前端也带本地图标兜底，目标机缺少系统 emoji 字体时不会空白。
- 部署包自带 LibreOffice 文档转换运行时。CI、PL、合同的 PDF 和 Excel 下载不依赖系统 Office；Word 下载使用系统内置、可编辑的 DOCX 标准模板，并填充同一份业务字段。
- systemd 服务日志可用 `sudo journalctl -u sta100-web.service -n 120 --no-pager` 和 `sudo journalctl -u sta100-openclaw-gateway.service -n 120 --no-pager` 查看；应用自身也会继续写入 `./logs/`。
- Agent 异常排查优先查看：`./scripts/sta100.sh logs sta100-web 160`、`./scripts/sta100.sh logs openclaw-sync 160`、`./scripts/sta100.sh logs openclaw-gateway 160`。
- 打包脚本会从当前部署目录生成归档，并排除 `data/`、`logs/`、`.run/` 测试数据。
EOF
}

copy_runtime() {
  mkdir -p "$RELEASE_DIR/bin" "$RELEASE_DIR/config" "$RELEASE_DIR/assets" "$RELEASE_DIR/openclaw/bin" "$RELEASE_DIR/openclaw/scripts" "$RELEASE_DIR/data" "$RELEASE_DIR/logs" "$RELEASE_DIR/.run" "$RELEASE_DIR/knowledge/shared"
  find "$RELEASE_DIR/knowledge/shared" -mindepth 1 -delete
  if [[ -d "$SHARED_KNOWLEDGE_SOURCE" ]]; then
    echo "copying shared knowledge from $SHARED_KNOWLEDGE_SOURCE..."
    cp -a "$SHARED_KNOWLEDGE_SOURCE/." "$RELEASE_DIR/knowledge/shared/"
  else
    echo "warning: shared knowledge source not found: $SHARED_KNOWLEDGE_SOURCE" >&2
  fi
  if [[ ! -x "$OFFICE_RUNTIME_DIR/bin/soffice" ]]; then
    "$ROOT/scripts/prepare-office-runtime.sh"
  fi
  cp "$STA100_WEB_DIR/config/sta100-agents.json" "$RELEASE_DIR/config/sta100-agents.json"
  if [[ -f "$STA100_WEB_DIR/assets/NotoColorEmoji.ttf" ]]; then
    cp "$STA100_WEB_DIR/assets/NotoColorEmoji.ttf" "$RELEASE_DIR/assets/NotoColorEmoji.ttf"
  fi
  node -e '
    const fs = require("fs");
    const path = process.argv[1];
    const data = JSON.parse(fs.readFileSync(path, "utf8"));
    data.workspace_root = "../openclaw";
    fs.writeFileSync(path, JSON.stringify(data, null, 2) + "\n");
  ' "$RELEASE_DIR/config/sta100-agents.json"
  cp -a "$OPENCLAW_DIR/bin/openclaw" "$RELEASE_DIR/openclaw/bin/openclaw"
  if [[ -d "$OPENCLAW_DIR/app" ]]; then
    mkdir -p "$RELEASE_DIR/openclaw/app"
  cp -a "$OPENCLAW_DIR/app/." "$RELEASE_DIR/openclaw/app/"
  if [[ -d "$OPENCLAW_DIR/extensions" ]]; then
    mkdir -p "$RELEASE_DIR/openclaw/extensions"
    cp -a "$OPENCLAW_DIR/extensions/." "$RELEASE_DIR/openclaw/extensions/"
  fi
  fi
  cp "$OPENCLAW_DIR/scripts/"*.sh "$RELEASE_DIR/openclaw/scripts/" 2>/dev/null || true
  chmod +x "$RELEASE_DIR/openclaw/scripts/"*.sh 2>/dev/null || true
  write_release_sync_script
  if [[ -d "$OPENCLAW_DIR/systemd" ]]; then
    mkdir -p "$RELEASE_DIR/openclaw/systemd"
    cp "$OPENCLAW_DIR/systemd/"*.service "$RELEASE_DIR/openclaw/systemd/" 2>/dev/null || true
  fi
  if [[ -f "$OPENCLAW_DIR/VERSION.json" ]]; then
    cp "$OPENCLAW_DIR/VERSION.json" "$RELEASE_DIR/openclaw/VERSION.json"
  fi
  write_release_launcher
  write_release_readme
}

build_binaries() {
  echo "building sta100-web..."
  mkdir -p "$RELEASE_DIR/bin"
  (cd "$STA100_WEB_DIR" && CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags='-s -w' -o "$RELEASE_DIR/bin/sta100-web" .)
  echo "building openclaw-operator..."
  (cd "$ORCH_DIR" && CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags='-s -w' -o "$RELEASE_DIR/bin/openclaw-operator" ./cmd/openclaw-operator)
}

stop_existing() {
  if [[ "${STA100_DEPLOY_STOP:-1}" == "1" ]] && [[ -x "$RELEASE_DIR/scripts/sta100.sh" ]]; then
    "$RELEASE_DIR/scripts/sta100.sh" stop || true
  fi
}

reset_release_state() {
  if [[ "${STA100_DEPLOY_CLEAN_STATE:-1}" != "1" ]]; then
    return 0
  fi
  echo "cleaning release runtime state..."
  rm -rf "$RELEASE_DIR/data" "$RELEASE_DIR/logs" "$RELEASE_DIR/.run"
  mkdir -p "$RELEASE_DIR/data" "$RELEASE_DIR/logs" "$RELEASE_DIR/.run"
}

start_deployed() {
  if [[ "${STA100_DEPLOY_START:-1}" == "1" ]]; then
    "$RELEASE_DIR/scripts/sta100.sh" start
  fi
}

require_file "$STA100_WEB_DIR/go.mod"
require_file "$OPENCLAW_DIR/bin/openclaw"
require_file "$ORCH_DIR/cmd/openclaw-operator/main.go"
require_file "$STA100_WEB_DIR/config/sta100-agents.json"

stop_existing
reset_release_state
copy_runtime
build_binaries
start_deployed

echo "deployed: $RELEASE_DIR"
