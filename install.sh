#!/usr/bin/env bash
set -e

REPO="rshdhere/vibecheck"
BIN="vibecheck"

# Get latest tag from GitHub API
TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
if [ -z "$TAG" ]; then
  echo "❌ No releases found for $REPO"
  exit 1
fi

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture
[ "$ARCH" = "x86_64" ] && ARCH="x86_64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

URL="https://github.com/$REPO/releases/download/$TAG/${BIN}_${OS^}_${ARCH}.tar.gz"

echo "⬇️  Downloading $BIN $TAG for $OS/$ARCH..."
curl -fsSL "$URL" -o /tmp/$BIN.tar.gz

echo "📦 Installing to /usr/local/bin..."
sudo tar -xzf /tmp/$BIN.tar.gz -C /usr/local/bin $BIN

echo "✅ Installed!"
echo "Version: $($BIN --version)"
