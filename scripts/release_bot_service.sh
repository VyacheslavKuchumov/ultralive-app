#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="${SERVICE_NAME:-ultralive-release-bot}"

usage() {
  cat <<USAGE
Usage: scripts/release_bot_service.sh <start|stop|restart|status|logs>

Environment overrides:
  SERVICE_NAME   systemd service name (default: ultralive-release-bot)
USAGE
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

if ! command -v sudo >/dev/null 2>&1; then
  echo "sudo is required" >&2
  exit 1
fi

case "$1" in
  start)
    sudo systemctl start "${SERVICE_NAME}"
    ;;
  stop)
    sudo systemctl stop "${SERVICE_NAME}"
    ;;
  restart)
    sudo systemctl restart "${SERVICE_NAME}"
    ;;
  status)
    sudo systemctl status "${SERVICE_NAME}" --no-pager
    ;;
  logs)
    sudo journalctl -u "${SERVICE_NAME}" -f
    ;;
  *)
    usage
    exit 1
    ;;
esac
