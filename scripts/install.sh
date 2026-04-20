#!/usr/bin/env bash
# DevPulse one-line installer
# Usage: curl -fsSL https://raw.githubusercontent.com/devpulse-cli/devpulse/main/scripts/install.sh | bash

set -euo pipefail

REPO="devpulse-cli/devpulse"
BINARY="pulse"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
CYAN='\033[0;36m'; GREEN='\033[0;32m'; RED='\033[0;31m'; BOLD='\033[1m'; R='\033[0m'

echo -e "${CYAN}${BOLD}"
cat <<'EOF'
  ██████╗ ███████╗██╗   ██╗██████╗ ██╗   ██╗██╗     ███████╗███████╗
  ██╔══██╗██╔════╝██║   ██║██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
  ██║  ██║█████╗  ██║   ██║██████╔╝██║   ██║██║     ███████╗█████╗
  ██║  ██║██╔══╝  ╚██╗ ██╔╝██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝
  ██████╔╝███████╗ ╚████╔╝ ██║     ╚██████╔╝███████╗███████║███████╗
  ╚═════╝ ╚══════╝  ╚═══╝  ╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝
EOF
echo -e "${R}"
echo -e "${CYAN}installing DevPulse...${R}"
echo ""

# ── Detect OS + arch ────────────────────────────────────────────────
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)         ARCH="amd64"  ;;
    aarch64|arm64)  ARCH="arm64"  ;;
    *)  echo -e "${RED}unsupported arch: $ARCH${R}"; exit 1 ;;
esac

# ── Fetch latest release tag from GitHub API ────────────────────────
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}could not fetch latest version — check your internet connection${R}"
    exit 1
fi

# ── Download the binary ─────────────────────────────────────────────
FILENAME="${BINARY}_${OS}_${ARCH}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

echo -e "  downloading ${CYAN}${BINARY} ${VERSION}${R} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "/tmp/${BINARY}"
chmod +x "/tmp/${BINARY}"

# ── Install ─────────────────────────────────────────────────────────
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
    echo -e "  needs sudo to write to ${INSTALL_DIR}..."
    sudo mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

echo ""
echo -e "  ${GREEN}✓ DevPulse installed!${R}"
echo ""
echo -e "  run ${CYAN}pulse init${R} to start tracking ur grind 🔥"
echo ""
