#!/bin/bash
BINARY="/usr/local/bin/aggloop"
CONFIG_DIR="/etc/agglabs"

log() {
    echo "[INFO] $1"
}

log "Starting uninstallation procedure..."

if command -v aggloop >/dev/null 2>&1; then
    log "Scanning for active background tunnels..."
    if [ -d "$CONFIG_DIR" ]; then
        for f in "$CONFIG_DIR"/*.pid; do
            if [ -f "$f" ]; then
                APP_NAME=$(basename "$f" .pid)
                log "Terminating application: $APP_NAME"
                sudo aggloop stop "$APP_NAME" >/dev/null 2>&1
            fi
        done
    fi
fi

if [ -f "$BINARY" ]; then
    sudo rm "$BINARY"
    log "Removed binary: $BINARY"
else
    log "Binary not found in $BINARY. Skipping."
fi

if [ -d "$CONFIG_DIR" ]; then
    echo "Do you wish to remove the configuration directory and logs ($CONFIG_DIR)? [y/N]"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        sudo rm -rf "$CONFIG_DIR"
        log "Configuration directory removed."
    else
        log "Configuration directory preserved."
    fi
fi

log "Uninstallation procedure finalized."