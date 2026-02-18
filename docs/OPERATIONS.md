# Operations Guide

## Backend Make Targets

Run from `server/`.

- `make build`: compile API binary into `bin/server`
- `make run`: build and run compiled binary
- `make air`: run backend with hot reload
- `make test`: run all Go tests
- `make swagger`: regenerate Swagger/OpenAPI files
- `make docker-up`: build/start compose stack from repository root file
- `make docker-db-up`: start only postgres service
- `make docker-down`: stop compose services

## Migrations

### Apply all migrations

```bash
cd server
make migrate-up
```

### Rollback all migrations

```bash
cd server
make migrate-down
```

### Force migration version

```bash
cd server
make migrate-force
```

Current `force` target pins version to `1` in `cmd/migrate/main.go`.

### Create a new migration template

```bash
cd server
make migration add_some_change
```

This writes SQL files into `server/cmd/migrate/migrations`.

### Regenerate Swagger docs

```bash
cd server
make swagger
```

## Test Commands

### Backend

```bash
cd server
make test
```

### Frontend

```bash
cd web
npm run build
```

## XMPP Release Bot (systemd)

Release bot listens in XMPP and runs deploy on `/update`:

- `git checkout main`
- `git fetch origin`
- `git pull --ff-only origin main`
- `docker compose up -d --build --remove-orphans`

### 1) Generate bot env file

```bash
python3 scripts/generate_release_bot_env.py \
  --jid ultralive-release-bot@vyachik-dev.ru \
  --password 'BOT_XMPP_PASSWORD' \
  --allowed-sender your-admin@vyachik-dev.ru \
  --repo-path /opt/ultralive-app-v2
```

This creates `.env.release-bot` in repository root.

By default, deploy script path is `<repo-path>/release-bot/deploy.sh`.  
If needed, set `DEPLOY_SCRIPT_PATH=/absolute/path/to/release-bot/deploy.sh` in `.env.release-bot`.

### 2) Install and start service

```bash
scripts/setup_release_bot_systemd.sh install
```

### 3) Manage service

```bash
scripts/release_bot_service.sh status
scripts/release_bot_service.sh logs
scripts/release_bot_service.sh restart
scripts/release_bot_service.sh stop
scripts/setup_release_bot_systemd.sh uninstall
```

### 4) XMPP commands

- `/help`
- `/ping`
- `/status`
- `/update`

## Troubleshooting

### Migrations fail with path error

Run migration commands from `server/`, because the runner expects:

- `file://cmd/migrate/migrations`

### `permission denied` on protected route

Check one of:

- missing `Authorization` header
- expired JWT (`expiredAt` claim)
- token signed with different `JWT_SECRET`

### Login succeeds but protected Nuxt API calls fail

The UI must send `Authorization` header from stored token. If token is absent/expired, re-login.
