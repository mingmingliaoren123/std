#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OFFICE_DIR="${STA100_OFFICE_RUNTIME_DIR:-$ROOT/release/sta100-release/office}"
WORK_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

case "$OFFICE_DIR" in
  "$ROOT"/release/*/office) ;;
  *) echo "STA100_OFFICE_RUNTIME_DIR must be a release/*/office directory" >&2; exit 2 ;;
esac

command -v apt-get >/dev/null 2>&1 || { echo "apt-get is required to assemble the ARM64 Office runtime" >&2; exit 1; }
command -v dpkg-deb >/dev/null 2>&1 || { echo "dpkg-deb is required to assemble the Office runtime" >&2; exit 1; }
LDCONFIG="$(command -v ldconfig 2>/dev/null || true)"
if [[ -z "$LDCONFIG" ]] && [[ -x /sbin/ldconfig ]]; then
  LDCONFIG=/sbin/ldconfig
fi
[[ -n "$LDCONFIG" ]] || { echo "ldconfig is required to assemble the Office runtime" >&2; exit 1; }

if [[ "$(dpkg --print-architecture 2>/dev/null || true)" != "arm64" ]]; then
  echo "the bundled Office runtime is currently built for Ubuntu ARM64" >&2
  exit 1
fi

if [[ -x "$OFFICE_DIR/bin/soffice" ]] && [[ "${STA100_OFFICE_FORCE_REBUILD:-0}" != "1" ]]; then
  echo "Office runtime already exists: $OFFICE_DIR"
  exit 0
fi

mkdir -p "$WORK_DIR/debs" "$WORK_DIR/root"
apt-get -s --no-install-recommends install libreoffice-core libreoffice-calc libreoffice-writer | awk '/^Inst / {print $2}' > "$WORK_DIR/packages.txt"
[[ -s "$WORK_DIR/packages.txt" ]] || { echo "LibreOffice package resolution returned no packages" >&2; exit 1; }

(
  cd "$WORK_DIR/debs"
  xargs -r apt-get download < "$WORK_DIR/packages.txt"
)

for deb in "$WORK_DIR"/debs/*.deb; do
  dpkg-deb -x "$deb" "$WORK_DIR/root"
done

LIB_DIR="$WORK_DIR/root/usr/lib/aarch64-linux-gnu"
PROGRAM_DIR="$WORK_DIR/root/usr/lib/libreoffice/program"
[[ -x "$PROGRAM_DIR/soffice" ]] || { echo "LibreOffice soffice binary was not assembled" >&2; exit 1; }

# Debian package configuration normally creates these library symlinks. The
# portable runtime is unpacked only, so create them without touching ld.so.conf.
"$LDCONFIG" -n "$LIB_DIR"

# Keep the runtime relocatable after tar extraction rather than retaining the
# build machine's absolute /usr/lib/libreoffice and /etc/libreoffice paths.
sed -i 's#^FHS_CONFIG_FILE=.*#FHS_CONFIG_FILE=${ORIGIN}/../../../../etc/libreoffice/sofficerc#' "$PROGRAM_DIR/sofficerc"
sed -i 's#^BRAND_BASE_DIR=.*#BRAND_BASE_DIR=${ORIGIN}/..#' "$PROGRAM_DIR/fundamentalrc"

rm -rf "$OFFICE_DIR"
mkdir -p "$OFFICE_DIR/bin"
cp -a "$WORK_DIR/root/." "$OFFICE_DIR/root/"

cat > "$OFFICE_DIR/bin/soffice" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

OFFICE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OFFICE_ROOT="$OFFICE_DIR/root"
PROGRAM_DIR="$OFFICE_ROOT/usr/lib/libreoffice/program"
LIB_DIR="$OFFICE_ROOT/usr/lib/aarch64-linux-gnu"
PROFILE_DIR="${STA100_OFFICE_PROFILE_DIR:-$OFFICE_DIR/../data/office-profile}"
OFFICE_HOME="${STA100_OFFICE_HOME:-$OFFICE_DIR/../data/office-home}"

[[ -x "$PROGRAM_DIR/soffice" ]] || { echo "bundled LibreOffice runtime is incomplete" >&2; exit 1; }
mkdir -p "$PROFILE_DIR" "$OFFICE_HOME"

export LD_LIBRARY_PATH="$PROGRAM_DIR:$LIB_DIR${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
export HOME="$OFFICE_HOME"
export SAL_USE_VCLPLUGIN=gen
export SAL_DISABLE_OPENGL=true

exec "$PROGRAM_DIR/soffice" "-env:UserInstallation=file://$PROFILE_DIR" "$@"
EOF
chmod 0755 "$OFFICE_DIR/bin/soffice"

echo "Office runtime assembled: $OFFICE_DIR"
