#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATE_DIR="${OPENCLAW_STATE_DIR:-$HOME/.openclaw}"
BACKUP_DIR="$ROOT/backups"
STAMP="$(date +%Y%m%d-%H%M%S)"

mkdir -p "$BACKUP_DIR"
if [[ ! -d "$STATE_DIR" ]]; then
  echo "No OpenClaw state directory found at $STATE_DIR"
  exit 0
fi

ARCHIVE="$BACKUP_DIR/openclaw-state-$STAMP.tar.gz"
tar -C "$(dirname "$STATE_DIR")" -czf "$ARCHIVE" "$(basename "$STATE_DIR")"
sha256sum "$ARCHIVE" > "$ARCHIVE.sha256"
chmod 600 "$ARCHIVE" "$ARCHIVE.sha256"
echo "$ARCHIVE"
