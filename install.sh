#!/usr/bin/env bash

# wpvul Installer
set -e

REPO="spagu/go-wpvul"
APP_NAME="wpvul"
BIN_DIR="/usr/local/bin"

echo "==== Installer for $APP_NAME ===="

# Get OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalize Architecture
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "armv8l" ]; then
    ARCH="arm64"
fi

if [ "$OS" = "darwin" ]; then
    echo "Mac OS / Darwin detected."
elif [ "$OS" = "linux" ]; then
    echo "Linux detected."
elif [ "$OS" = "freebsd" ]; then
    echo "FreeBSD detected."
else
    echo "Unsupported OS: $OS"
    exit 1
fi

BINARY_URL="https://github.com/${REPO}/releases/latest/download/${APP_NAME}-${OS}-${ARCH}"
TMP_DEST="/tmp/${APP_NAME}"

echo ">> Downloading $APP_NAME for $OS ($ARCH) from GitHub Releases..."
echo "   URL: $BINARY_URL"

# We check if curl is installed
if ! command -v curl &> /dev/null; then
    echo "Error: curl could not be found. Please install it first."
    exit 1
fi

# Attempt to download the binary
HTTP_STATUS=$(curl -w "%{http_code}" -sSL -o "$TMP_DEST" "$BINARY_URL")

if [ "$HTTP_STATUS" -ne 200 ] && [ "$HTTP_STATUS" -ne 302 ]; then
    echo -e "\nError: Could not download the release! (HTTP $HTTP_STATUS) \nMake sure the release tag exists on GitHub: https://github.com/${REPO}/releases"
    rm -f "$TMP_DEST"
    exit 1
fi

chmod +x "$TMP_DEST"

echo ">> Installing to $BIN_DIR (might require sudo/root)..."
if [ -w "$BIN_DIR" ]; then
    mv "$TMP_DEST" "$BIN_DIR/$APP_NAME"
else
    sudo mv "$TMP_DEST" "$BIN_DIR/$APP_NAME"
fi

echo "======================================"
echo "✓ Success! $APP_NAME is now installed."
echo "Run '$APP_NAME wp-content/plugins' to start analyzing your folders."
