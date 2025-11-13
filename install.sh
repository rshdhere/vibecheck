#!/usr/bin/env bash

# --- FORCE BASH EVEN IF USER RUNS VIA ZSH ---
# macOS uses zsh by default and will break bash features
if [ -n "$ZSH_VERSION" ]; then
  exec bash "$0" "$@"
fi

set -e

REPO="rshdhere/vibecheck"
BIN="vibecheck"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to find all existing installations
find_existing_installations() {
  local locations=(
    "/usr/local/bin/$BIN"
    "/usr/bin/$BIN"
    "$HOME/go/bin/$BIN"
    "$HOME/.local/bin/$BIN"
    "$HOME/bin/$BIN"
  )

  if [ -n "$GOPATH" ]; then
    locations+=("$GOPATH/bin/$BIN")
  fi

  local found=()
  for loc in "${locations[@]}"; do
    if [ -f "$loc" ] && [ "$loc" != "$INSTALL_DIR/$BIN" ]; then
      found+=("$loc")
    fi
  done

  echo "${found[@]}"
}

get_version() {
  local binary=$1
  $binary --version 2>/dev/null | head -1 || echo "unknown"
}

# --- FIXED FOR MACOS (NO GNU GREP NEEDED) ---
echo -e "${BLUE}🔍 Checking for latest release...${NC}"
TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest |
  grep '"tag_name"' |
  sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
  echo -e "${RED}❌ No releases found for $REPO${NC}"
  exit 1
fi
# --------------------------------------------

existing=($(find_existing_installations))
if [ ${#existing[@]} -gt 0 ]; then
  echo -e "${YELLOW}⚠️  Found existing installation(s):${NC}"
  for loc in "${existing[@]}"; do
    ver=$(get_version "$loc")
    echo -e "   ${loc} ${BLUE}(version: ${ver})${NC}"
  done
  echo ""
  echo -e "${YELLOW}🧹 Cleaning up old installations to avoid PATH conflicts...${NC}"
  for loc in "${existing[@]}"; do
    if [ -w "$loc" ]; then
      rm -f "$loc"
      echo -e "   ${GREEN}✓${NC} Removed $loc"
    elif [ -w "$(dirname "$loc")" ]; then
      rm -f "$loc"
      echo -e "   ${GREEN}✓${NC} Removed $loc"
    else
      sudo rm -f "$loc" 2>/dev/null && echo -e "   ${GREEN}✓${NC} Removed $loc" ||
        echo -e "   ${YELLOW}⚠${NC}  Couldn't remove $loc (please remove manually)"
    fi
  done
  echo ""
fi

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

[ "$ARCH" = "x86_64" ] && ARCH="x86_64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"
[ "$ARCH" = "i386" ] && ARCH="i386"
[ "$ARCH" = "i686" ] && ARCH="i386"

# --- FIXED: BASH-ONLY UPPERCASE BROKE ON MAC ZSH ---
OS_UC=$(echo "$OS" | tr '[:lower:]' '[:upper:]')
URL="https://github.com/$REPO/releases/download/$TAG/${BIN}_${OS_UC}_${ARCH}.tar.gz"
# -----------------------------------------------

echo -e "${BLUE}⬇️  Downloading $BIN $TAG for $OS/$ARCH...${NC}"
curl -fsSL "$URL" -o /tmp/$BIN.tar.gz

echo -e "${BLUE}📦 Installing to $INSTALL_DIR...${NC}"
sudo tar -xzf /tmp/$BIN.tar.gz -C $INSTALL_DIR $BIN
sudo chmod +x $INSTALL_DIR/$BIN
rm -f /tmp/$BIN.tar.gz

echo -e "${GREEN}✅ Successfully installed!${NC}"
echo ""

INSTALLED_VERSION=$($INSTALL_DIR/$BIN --version 2>&1 | head -1)
echo -e "📌 Installed version: ${GREEN}${INSTALLED_VERSION}${NC}"
echo -e "📍 Location: ${BLUE}$INSTALL_DIR/$BIN${NC}"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo ""
  echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
  echo -e "   Add this to your ~/.bashrc or ~/.zshrc:"
  echo -e "   ${BLUE}export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
fi

echo ""
echo -e "🚀 Run ${GREEN}vibecheck --help${NC} to get started!"
