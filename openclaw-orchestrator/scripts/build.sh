#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$ROOT/bin"
CGO_ENABLED=0 GOOS="${GOOS:-linux}" GOARCH="${GOARCH:-arm64}" \
  go build -trimpath -ldflags='-s -w' \
  -o "$ROOT/bin/openclaw-operator" "$ROOT/cmd/openclaw-operator"
echo "$ROOT/bin/openclaw-operator"
