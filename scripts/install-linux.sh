#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-hduhelp/hdu-cli}"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_PATH="${CONFIG_PATH:-$HOME/.hdu-cli.yaml}"
SYSTEMD_DIR="${SYSTEMD_DIR:-$HOME/.config/systemd/user}"
DOWNLOAD_PREFIX="${DOWNLOAD_PREFIX:-}"
BIN_NAME="hdu-cli"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)
      echo "unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

release_url() {
  local arch
  local asset
  arch="$(detect_arch)"
  if [[ "$VERSION" == "latest" ]]; then
    asset="https://github.com/${REPO}/releases/latest/download/${BIN_NAME}_Linux_${arch}.tar.gz"
  else
    asset="https://github.com/${REPO}/releases/download/${VERSION}/${BIN_NAME}_Linux_${arch}.tar.gz"
  fi
  echo "${DOWNLOAD_PREFIX}${asset}"
}

main() {
  need_cmd curl
  need_cmd tar
  need_cmd systemctl

  mkdir -p "$INSTALL_DIR" "$SYSTEMD_DIR"

  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT

  echo "Downloading ${BIN_NAME}..."
  curl -fsSL "$(release_url)" -o "$tmp/${BIN_NAME}.tar.gz"

  tar -xzf "$tmp/${BIN_NAME}.tar.gz" -C "$tmp"
  install -m 0755 "$tmp/${BIN_NAME}" "$INSTALL_DIR/${BIN_NAME}"

  if [[ ! -f "$CONFIG_PATH" ]]; then
    echo "Writing default config to $CONFIG_PATH"
    "$INSTALL_DIR/${BIN_NAME}" config init --config "$CONFIG_PATH"
  fi

  echo "Installing user service..."
  "$INSTALL_DIR/${BIN_NAME}" service install --config "$CONFIG_PATH"
  systemctl --user daemon-reload
  systemctl --user enable --now hdu-cli.service

  cat <<EOF

Installed successfully.

Next:
  1. Edit your account:
     $INSTALL_DIR/${BIN_NAME} config set --config "$CONFIG_PATH" --key username --value YOUR_STUDENT_ID
  2. Save your password quickly:
     $INSTALL_DIR/${BIN_NAME} config init --config "$CONFIG_PATH" --username YOUR_STUDENT_ID --password YOUR_PASSWORD
  3. Check status:
     $INSTALL_DIR/${BIN_NAME} service status
     $INSTALL_DIR/${BIN_NAME} net status --config "$CONFIG_PATH"
EOF
}

main "$@"
