#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICE_NAME="${SERVICE_NAME:-ultralive-release-bot}"
ENV_FILE="${ROOT_DIR}/.env.release-bot"
VENV_DIR="${ROOT_DIR}/release-bot/.venv"
UNIT_PATH="/etc/systemd/system/${SERVICE_NAME}.service"
BOT_USER="${BOT_USER:-${SUDO_USER:-$(id -un)}}"
BOT_GROUP="${BOT_GROUP:-$(id -gn "${BOT_USER}")}"

usage() {
  cat <<USAGE
Usage: scripts/setup_release_bot_systemd.sh <install|uninstall>

Environment overrides:
  SERVICE_NAME   systemd service name (default: ultralive-release-bot)
  BOT_USER       user to run the bot under (default: current user)
  BOT_GROUP      group to run the bot under (default: BOT_USER primary group)
USAGE
}

require_sudo() {
  if ! command -v sudo >/dev/null 2>&1; then
    echo "sudo is required to manage systemd unit files" >&2
    exit 1
  fi
}

install_service() {
  if [[ ! -f "${ENV_FILE}" ]]; then
    echo "Missing ${ENV_FILE}. Generate it first:" >&2
    echo "  python3 scripts/generate_release_bot_env.py --help" >&2
    exit 1
  fi

  python3 -m venv "${VENV_DIR}"
  "${VENV_DIR}/bin/pip" install --upgrade pip
  "${VENV_DIR}/bin/pip" install -r "${ROOT_DIR}/release-bot/requirements.txt"

  tmp_unit="$(mktemp)"
  cat > "${tmp_unit}" <<UNIT
[Unit]
Description=Ultralive XMPP Release Bot
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${BOT_USER}
Group=${BOT_GROUP}
WorkingDirectory=${ROOT_DIR}
EnvironmentFile=${ENV_FILE}
ExecStart=${VENV_DIR}/bin/python ${ROOT_DIR}/release-bot/bot.py
Restart=always
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=false
ReadWritePaths=${ROOT_DIR}

[Install]
WantedBy=multi-user.target
UNIT

  require_sudo
  sudo install -m 0644 "${tmp_unit}" "${UNIT_PATH}"
  rm -f "${tmp_unit}"

  sudo systemctl daemon-reload
  sudo systemctl enable --now "${SERVICE_NAME}"
  sudo systemctl status "${SERVICE_NAME}" --no-pager

  echo "Service installed: ${SERVICE_NAME}"
  echo "View logs: sudo journalctl -u ${SERVICE_NAME} -f"
}

uninstall_service() {
  require_sudo
  sudo systemctl disable --now "${SERVICE_NAME}" >/dev/null 2>&1 || true
  sudo rm -f "${UNIT_PATH}"
  sudo systemctl daemon-reload
  echo "Service removed: ${SERVICE_NAME}"
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

case "$1" in
  install)
    install_service
    ;;
  uninstall)
    uninstall_service
    ;;
  *)
    usage
    exit 1
    ;;
esac
