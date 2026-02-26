#!/bin/sh
# TinyClaw installer — curl -fsSL https://raw.githubusercontent.com/tinyland-inc/tinyclaw/main/scripts/install.sh | sh
set -eu

REPO="tinyland-inc/tinyclaw"
INSTALL_DIR="${TINYCLAW_INSTALL_DIR:-$HOME/.local/bin}"

info()  { printf '  \033[1;34m→\033[0m %s\n' "$1"; }
ok()    { printf '  \033[1;32m✓\033[0m %s\n' "$1"; }
fail()  { printf '  \033[1;31m✗\033[0m %s\n' "$1" >&2; exit 1; }

detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "Linux"  ;;
    Darwin*) echo "Darwin" ;;
    *)       fail "Unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)   echo "x86_64" ;;
    aarch64|arm64)   echo "arm64"  ;;
    *)               fail "Unsupported architecture: $(uname -m)" ;;
  esac
}

latest_version() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
      grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" |
      grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/'
  else
    fail "Neither curl nor wget found"
  fi
}

download() {
  url="$1"; dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$dest" "$url"
  else
    wget -qO "$dest" "$url"
  fi
}

main() {
  printf '\n  \033[1mTinyClaw Installer\033[0m\n\n'

  OS="$(detect_os)"
  ARCH="$(detect_arch)"
  info "Detected ${OS}/${ARCH}"

  VERSION="$(latest_version)"
  [ -z "$VERSION" ] && fail "Could not determine latest version"
  info "Latest version: v${VERSION}"

  TARBALL="tinyclaw_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"
  info "Downloading ${URL}"

  TMPDIR="$(mktemp -d)"
  trap 'rm -rf "$TMPDIR"' EXIT

  download "$URL" "${TMPDIR}/${TARBALL}"
  tar -xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

  mkdir -p "$INSTALL_DIR"
  mv "${TMPDIR}/tinyclaw" "${INSTALL_DIR}/tinyclaw"
  chmod +x "${INSTALL_DIR}/tinyclaw"
  ok "Installed tinyclaw to ${INSTALL_DIR}/tinyclaw"

  if "${INSTALL_DIR}/tinyclaw" version >/dev/null 2>&1; then
    ok "Verified: $("${INSTALL_DIR}/tinyclaw" version 2>&1 | head -1)"
  else
    fail "Binary verification failed"
  fi

  case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) ;;
    *) printf '\n  \033[1;33m!\033[0m Add %s to your PATH\n' "$INSTALL_DIR" ;;
  esac

  printf '\n'
}

main "$@"
