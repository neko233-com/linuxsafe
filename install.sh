#!/bin/bash
set -e

REPO="neko233-com/linuxsafe"
BINARY="linuxsafe"
INSTALL_DIR="/usr/local/bin"

echo "Installing linuxsafe..."

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux)
            PLATFORM="linux"
            ;;
        darwin)
            PLATFORM="darwin"
            ;;
        *)
            echo "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l|armhf)
            ARCH="arm"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    echo "${PLATFORM}_${ARCH}"
}

get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4
}

download_and_install() {
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        echo "Failed to get latest version"
        exit 1
    fi
    
    PLATFORM=$(detect_platform)
    FILENAME="${BINARY}_${VERSION}_${PLATFORM}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"
    
    echo "Downloading ${VERSION} for ${PLATFORM}..."
    
    TMPDIR=$(mktemp -d)
    curl -sL "$URL" -o "${TMPDIR}/${FILENAME}"
    
    echo "Extracting..."
    tar -xzf "${TMPDIR}/${FILENAME}" -C "${TMPDIR}"
    
    echo "Installing to ${INSTALL_DIR}..."
    sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/"
    sudo chmod +x "${INSTALL_DIR}/${BINARY}"
    
    rm -rf "${TMPDIR}"
    
    echo "Installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
}

verify_installation() {
    if command -v ${BINARY} &> /dev/null; then
        echo ""
        echo "✓ Installation successful!"
        echo ""
        echo "Run '${BINARY} --help' to get started"
        echo "Run '${BINARY} scan /' to scan your system"
    else
        echo ""
        echo "Installation completed but ${BINARY} not found in PATH"
        echo "You may need to add ${INSTALL_DIR} to your PATH"
    fi
}

download_and_install
verify_installation
