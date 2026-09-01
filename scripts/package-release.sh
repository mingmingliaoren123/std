#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_DIR="${RELEASE_DIR:-$ROOT/release/sta100-release}"
ARCHIVE_NAME="${ARCHIVE_NAME:-sta100-release-linux-arm64-$(date +%Y%m%d-%H%M%S).tar.gz}"
ARCHIVE_PATH="${ARCHIVE_PATH:-$ROOT/release/$ARCHIVE_NAME}"
STAGE="$(mktemp -d)"

cleanup() {
  rm -rf "$STAGE"
}
trap cleanup EXIT

require_file() {
  [[ -e "$1" ]] || { echo "missing required path: $1" >&2; exit 1; }
}

require_file "$RELEASE_DIR/bin/sta100-web"
require_file "$RELEASE_DIR/bin/openclaw-operator"
require_file "$RELEASE_DIR/openclaw/bin/openclaw"
require_file "$RELEASE_DIR/config/sta100-agents.json"
require_file "$RELEASE_DIR/scripts/sta100.sh"
require_file "$RELEASE_DIR/office/bin/soffice"

mkdir -p "$(dirname "$ARCHIVE_PATH")"
rm -f "$ARCHIVE_PATH"

PACKAGE_ROOT="$STAGE/sta100-release"
mkdir -p "$PACKAGE_ROOT"

for item in bin config assets office openclaw scripts README-DEPLOY.md; do
  if [[ -e "$RELEASE_DIR/$item" ]]; then
    cp -a "$RELEASE_DIR/$item" "$PACKAGE_ROOT/"
  fi
done

mkdir -p "$PACKAGE_ROOT/data" "$PACKAGE_ROOT/logs" "$PACKAGE_ROOT/.run"

find "$PACKAGE_ROOT" -type f \( -name '*.log' -o -name '*.pid' -o -name '*.db' -o -name '*.db-shm' -o -name '*.db-wal' \) -delete
for runtime_dir in downloads runtime backups; do
  rm -rf "$PACKAGE_ROOT/openclaw/$runtime_dir"
done

(cd "$STAGE" && tar -czf "$ARCHIVE_PATH" sta100-release)

echo "$ARCHIVE_PATH"
