#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.release-bot.yml"
ENV_FILE="${ROOT_DIR}/.env.release-bot"

usage() {
  cat <<USAGE
Usage: scripts/release_bot_stack.sh <command>

Commands:
  up        Build and start release bot container
  down      Stop and remove release bot container
  restart   Restart release bot container
  logs      Show release bot logs (follow)
  ps        Show release bot container status
USAGE
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "Missing ${ENV_FILE}. Generate it first:" >&2
  echo "  python3 scripts/generate_release_bot_env.py --help" >&2
  exit 1
fi

compose=(docker compose --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}")

case "$1" in
  up)
    "${compose[@]}" up -d --build --remove-orphans
    ;;
  down)
    "${compose[@]}" down
    ;;
  restart)
    "${compose[@]}" down
    "${compose[@]}" up -d --build --remove-orphans
    ;;
  logs)
    "${compose[@]}" logs -f release-bot
    ;;
  ps)
    "${compose[@]}" ps
    ;;
  *)
    usage
    exit 1
    ;;
esac
