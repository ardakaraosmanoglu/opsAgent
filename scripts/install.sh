#!/usr/bin/env bash
set -euo pipefail

APP_NAME="opsagent"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/opsagent"
DATA_DIR="/var/lib/opsagent"
LOG_DIR="/var/log/opsagent"
SERVICE_FILE="/etc/systemd/system/opsagent.service"
# Get latest release tag dynamically
LATEST_TAG=$(curl -fsSL https://api.github.com/repos/ardakaraosmanoglu/opsAgent/releases/latest | grep '"tag_name"' | sed 's/.*"tag_name": "\([^"]*\)".*/\1/')
BINARY_URL_BASE="${OPSAGENT_URL:-https://github.com/ardakaraosmanoglu/opsAgent/releases/download/${LATEST_TAG}}"

if [[ "$EUID" -ne 0 ]]; then
    echo "Please run as root or with sudo."
    exit 1
fi

# Check if already installed
if systemctl is-enabled opsagent &>/dev/null; then
    echo "OpsAgent Updater"
    UPDATE_MODE=true
else
    echo "OpsAgent Installer"
    UPDATE_MODE=false
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

if [[ "$OS" != "linux" ]]; then
    echo "Unsupported OS: $OS"
    exit 1
fi

echo "[1/9] Detected OS: $OS / $ARCH"

if [[ "$UPDATE_MODE" == "true" ]]; then
    echo "[1/3] Stopping existing service..."
    systemctl stop opsagent || true
else
    mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
fi

if [[ -n "$BINARY_URL_BASE" ]]; then
    echo "[$([ "$UPDATE_MODE" == "true" ] && echo "2" || echo "2")/3] Downloading OpsAgent binary..."
    curl -fsSL -L "$BINARY_URL_BASE/opsagent-linux-$ARCH" -o "$INSTALL_DIR/$APP_NAME"
else
    echo "[$([ "$UPDATE_MODE" == "true" ] && echo "2" || echo "2")/3] Installing OpsAgent from local build..."
    if [[ ! -f "./opsagent" ]]; then
        echo "Error: No binary found at ./opsagent and no OPSAGENT_URL set."
        echo "Build the binary first: go build -o opsagent ./cmd/opsagent"
        exit 1
    fi
    cp ./opsagent "$INSTALL_DIR/$APP_NAME"
fi

chmod 755 "$INSTALL_DIR/$APP_NAME"
echo "[$([ "$UPDATE_MODE" == "true" ] && echo "3" || echo "3")/3] Binary installed to $INSTALL_DIR/$APP_NAME"

if [[ "$UPDATE_MODE" == "true" ]]; then
    echo "Starting OpsAgent..."
    systemctl start opsagent
    sleep 2
    if systemctl is-active --quiet opsagent; then
        echo ""
        echo "OpsAgent updated successfully."
        echo ""
        echo "Service:"
        echo "active"
        echo ""
        echo "Dashboard:"
        echo "http://$(hostname -I | awk '{print $1}'):8787"
    else
        echo "Service failed to start. Check logs with: journalctl -u opsagent -n 50"
        exit 1
    fi
else
    echo "[4/9] Creating directories and setting permissions..."
    chmod 700 "$CONFIG_DIR" 2>/dev/null || mkdir -p "$CONFIG_DIR" && chmod 700 "$CONFIG_DIR"
    chmod 700 "$DATA_DIR" 2>/dev/null || mkdir -p "$DATA_DIR" && chmod 700 "$DATA_DIR"
    chmod 700 "$LOG_DIR" 2>/dev/null || mkdir -p "$LOG_DIR" && chmod 700 "$LOG_DIR"

    echo "[5/9] Creating config..."
    if [[ ! -f "$CONFIG_DIR/config.yaml" ]]; then
        SESSION_SECRET="$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | xxd -p)"

        cat > "$CONFIG_DIR/config.yaml" <<EOF
app:
  name: "OpsAgent"
  environment: "production"

server:
  bind_address: "0.0.0.0"
  port: 8787

database:
  path: "/var/lib/opsagent/opsagent.db"

logging:
  level: "info"
  file: "/var/log/opsagent/opsagent.log"

security:
  require_approval_for_write_actions: true
  command_timeout_seconds: 120
  max_output_size_kb: 1024

ai:
  enabled: false
  provider: ""
  api_key: ""
  model: ""

monitoring:
  interval_seconds: 30
  disk_warning_threshold: 85
  disk_critical_threshold: 95
  memory_warning_threshold: 90
  cpu_warning_threshold: 90
  cpu_warning_duration_seconds: 300

dashboard:
  session_secret: "$SESSION_SECRET"
EOF
    fi
    chmod 600 "$CONFIG_DIR/config.yaml"
    chown root:root "$CONFIG_DIR/config.yaml"

    echo "[6/9] Creating systemd service..."
    cat > "$SERVICE_FILE" <<'EOF'
[Unit]
Description=OpsAgent - Local Linux AI Ops Assistant
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/opsagent serve --config /etc/opsagent/config.yaml
Restart=always
RestartSec=5
User=root
WorkingDirectory=/var/lib/opsagent

[Install]
WantedBy=multi-user.target
EOF

    echo "[7/9] Reloading systemd and enabling service..."
    systemctl daemon-reload
    systemctl enable opsagent

    echo "[8/9] Starting OpsAgent..."
    systemctl restart opsagent

    sleep 2

    echo "[9/9] Verifying service status..."
    if systemctl is-active --quiet opsagent; then
        echo ""
        echo "OpsAgent installed successfully."
        echo ""
        echo "Service:"
        echo "active"
        echo ""
        echo "Dashboard:"
        echo "http://$(hostname -I | awk '{print $1}'):8787"
        echo ""
        echo "Access dashboard directly:"
        echo ""
        echo "  http://YOUR_SERVER_IP:8787"
        echo ""
        echo "Your server data stays on this machine."
        echo "No write operation will run without approval."
    else
        echo "Service failed to start. Check logs with: journalctl -u opsagent -n 50"
        exit 1
    fi
fi
