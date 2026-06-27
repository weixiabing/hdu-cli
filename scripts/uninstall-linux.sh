#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
SYSTEMD_DIR="${SYSTEMD_DIR:-$HOME/.config/systemd/user}"
BIN_NAME="hdu-cli"

if command -v systemctl >/dev/null 2>&1; then
  systemctl --user disable --now hdu-cli.service >/dev/null 2>&1 || true
  systemctl --user daemon-reload >/dev/null 2>&1 || true
fi

rm -f "$SYSTEMD_DIR/hdu-cli.service"
rm -f "$INSTALL_DIR/$BIN_NAME"

echo "Removed user service and binary."
