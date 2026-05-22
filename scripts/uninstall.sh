#!/usr/bin/env bash
set -euo pipefail

if [[ "$EUID" -ne 0 ]]; then
    echo "Please run as root or with sudo."
    exit 1
fi

echo "Stopping OpsAgent..."
systemctl stop opsagent 2>/dev/null || true
systemctl disable opsagent 2>/dev/null || true

echo "Removing systemd service..."
rm -f /etc/systemd/system/opsagent.service
systemctl daemon-reload

echo "Removing binary..."
rm -f /usr/local/bin/opsagent

echo ""
read -p "Remove config, database and logs? This cannot be undone. [y/N] " confirm
if [[ "$confirm" == "y" || "$confirm" == "Y" ]]; then
    rm -rf /etc/opsagent
    rm -rf /var/lib/opsagent
    rm -rf /var/log/opsagent
    echo "All OpsAgent files removed."
else
    echo "Config, database and logs were kept."
fi

echo "OpsAgent uninstalled."
