#!/usr/bin/env sh
set -eu

REPO_PATH="${REPO_PATH:-.}"
DEPLOY_REMOTE="${DEPLOY_REMOTE:-origin}"
DEPLOY_BRANCH="${DEPLOY_BRANCH:-main}"
DEPLOY_COMPOSE_FILE="${DEPLOY_COMPOSE_FILE:-}"

cd "${REPO_PATH}"

git checkout "${DEPLOY_BRANCH}"
git fetch "${DEPLOY_REMOTE}"
git pull --ff-only "${DEPLOY_REMOTE}" "${DEPLOY_BRANCH}"

if [ -n "${DEPLOY_COMPOSE_FILE}" ]; then
  if docker compose version >/dev/null 2>&1; then
    docker compose -f "${DEPLOY_COMPOSE_FILE}" up -d --build --remove-orphans
  else
    docker-compose -f "${DEPLOY_COMPOSE_FILE}" up -d --build --remove-orphans
  fi
else
  if docker compose version >/dev/null 2>&1; then
    docker compose up -d --build --remove-orphans
  else
    docker-compose up -d --build --remove-orphans
  fi
fi
