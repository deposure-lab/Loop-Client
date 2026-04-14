#!/bin/bash
BASE_URL="https://github.com/deposure-lab/Loop-Client/releases/download/main"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/agglabs"
TEMP_BIN="/tmp/aggloop_tmp"

log() {
    echo "[INFO] $1"
}

error() {
    echo "[ERROR] $1" >&2
    exit 1
}

log "Initializing system environment detection..."

OS="$(uname -s)"
ARCH="$(uname -m)"
BINARY_NAME=""

case "${OS}" in
    Linux*)     OS_TYPE="linux";;
    Darwin*)    OS_TYPE="macos";;
    MINGW*|MSYS*|CYGWIN*) OS_TYPE="windows";;
    *)          error "Unsupported operating system: ${OS}";;
esac

case "${ARCH}" in
    x86_64|amd64) ARCH_TYPE="amd64";;
    arm64|aarch64) ARCH_TYPE="arm64";;
    *)            error "Unsupported architecture: ${ARCH}";;
esac

if [ "$OS_TYPE" == "windows" ]; then
    BINARY_NAME="aggloop-windows-${ARCH_TYPE}.exe"
else
    BINARY_NAME="aggloop-${OS_TYPE}-${ARCH_TYPE}"
fi

DOWNLOAD_URL="${BASE_URL}/${BINARY_NAME}"

log "Detected platform: ${OS_TYPE}"
log "Detected architecture: ${ARCH_TYPE}"
log "Fetching binary from: ${DOWNLOAD_URL}"

curl -L -s -o "${TEMP_BIN}" "${DOWNLOAD_URL}"
if [ $? -ne 0 ]; then
    error "Failed to download the binary from the repository."
fi

log "Installing binary to ${INSTALL_DIR}/aggloop..."

if [ ! -d "${CONFIG_DIR}" ]; then
    sudo mkdir -p "${CONFIG_DIR}"
    sudo chmod 755 "${CONFIG_DIR}"
fi

sudo mv "${TEMP_BIN}" "${INSTALL_DIR}/aggloop"
sudo chmod +x "${INSTALL_DIR}/aggloop"

if command -v aggloop >/dev/null 2>&1; then
    echo "--------------------------------------------------"
    log "Installation completed successfully."
    log "The 'aggloop' command is now available globally."
    echo "--------------------------------------------------"
else
    log "Installation successful, but ${INSTALL_DIR} is not in the system PATH."
fi