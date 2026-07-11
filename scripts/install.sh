#!/bin/sh
# OmniConfig — cross-platform installation script
# Usage: curl -fsSL https://omnicofig.sh/install | sh
# Or:   wget -qO- https://omnicofig.sh/install | sh

set -eu

BINARY="omniconfig"
VERSION="${VERSION:-latest}"
REPO="${REPO:-omnicofig/cli}"
PREFIX="${PREFIX:-/usr/local/bin}"

# --- Detect OS and architecture ---
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$OS" in
        linux|darwin) ;;
        cygwin*|mingw*|msys*) OS="windows" ;;
        *) echo "Unsupported OS: $OS"; exit 1 ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l|armv8l) ARCH="arm64" ;;
        *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    # Detect if running under Termux on Android
    if [ -n "${TERMUX_VERSION:-}" ] || [ -d "/data/data/com.termux" ]; then
        OS="android"
    fi

    echo "${OS}-${ARCH}"
}

# --- Determine install directory ---
detect_prefix() {
    if [ -n "${PREFIX:-}" ] && [ -d "$(dirname "$PREFIX")" ]; then
        echo "$PREFIX"
        return
    fi

    # Prefer /usr/local/bin, fallback to ~/.local/bin
    if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
    elif [ -d "$HOME/.local/bin" ] && [ -w "$HOME/.local/bin" ]; then
        echo "$HOME/.local/bin"
    elif [ -d "$HOME/bin" ] && [ -w "$HOME/bin" ]; then
        echo "$HOME/bin"
    else
        mkdir -p "$HOME/.local/bin"
        echo "$HOME/.local/bin"
    fi
}

# --- Main installation ---
main() {
    PLATFORM=$(detect_platform)
    INSTALL_DIR=$(detect_prefix)

    echo "==> OmniConfig Installer"
    echo "    Platform: ${PLATFORM}"
    echo "    Target:   ${INSTALL_DIR}/${BINARY}"

    # Determine download URL
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}-${PLATFORM}"
    else
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}-${PLATFORM}"
    fi

    # Add .exe for Windows
    case "$PLATFORM" in
        windows-*) DOWNLOAD_URL="${DOWNLOAD_URL}.exe" ;;
    esac

    echo "    Download: ${DOWNLOAD_URL}"

    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "${DOWNLOAD_URL}" -o "/tmp/${BINARY}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${DOWNLOAD_URL}" -O "/tmp/${BINARY}"
    else
        echo "Error: need curl or wget to download"
        exit 1
    fi

    # Verify download
    if [ ! -s "/tmp/${BINARY}" ]; then
        echo "Error: download failed or file is empty"
        exit 1
    fi

    # Make executable and install
    chmod +x "/tmp/${BINARY}"
    mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"

    echo ""
    echo "==> OmniConfig installed successfully!"
    echo "    ${INSTALL_DIR}/${BINARY}"
    echo ""
    echo "    Run '${BINARY} --help' to get started"
    echo "    Run '${BINARY} detect' to detect your OS and config format"
}

main "$@"