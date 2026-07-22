#!/usr/bin/env bash
# game-dev-cli 一键安装脚本
# 用法: curl -fsSL https://github.com/8liang/game-dev-cli/releases/latest/download/install.sh | bash
set -euo pipefail

VERSION="${GAME_DEV_CLI_VERSION:-latest}"
INSTALL_DIR="${GAME_DEV_CLI_DIR:-$HOME/.game-dev-cli/bin}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ 不支持的架构: $ARCH (仅 amd64/arm64)"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "❌ 不支持的平台: $OS (仅 linux/darwin)"; exit 1 ;;
esac

TARBALL="game-dev-cli-${OS}-${ARCH}.tar.gz"

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/8liang/game-dev-cli/releases/latest/download/${TARBALL}"
else
  URL="https://github.com/8liang/game-dev-cli/releases/download/${VERSION}/${TARBALL}"
fi

mkdir -p "$INSTALL_DIR"
echo "⬇️  下载 $URL"

# try curl, fallback to wget
if command -v curl &>/dev/null; then
  curl -fsSL "$URL" -o "${INSTALL_DIR}/${TARBALL}"
elif command -v wget &>/dev/null; then
  wget -q "$URL" -O "${INSTALL_DIR}/${TARBALL}"
else
  echo "❌ 需要 curl 或 wget"; exit 1
fi

tar xzf "${INSTALL_DIR}/${TARBALL}" -C "$INSTALL_DIR"
mv "${INSTALL_DIR}/game-dev-cli-${OS}-${ARCH}" "${INSTALL_DIR}/game-dev-cli" 2>/dev/null || true
chmod +x "${INSTALL_DIR}/game-dev-cli"
rm -f "${INSTALL_DIR}/${TARBALL}"

echo "✅ game-dev-cli 已安装到: ${INSTALL_DIR}/game-dev-cli"
echo "   添加到 PATH:  export PATH=\"${INSTALL_DIR}:\$PATH\""
