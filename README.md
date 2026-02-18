# Ultralive CRM

Ultralive CRM is a full-stack system for equipment accounting on shoots.

- `server/`: Go API
- `web/`: Nuxt 4 frontend (Nuxt UI + Pinia)
- PostgreSQL
- SQL migrations via `golang-migrate`

## Main CRM Domains

- users and auth (single `users` table, no separate `auths` table in new runtime)
- projects (`projects`)
- drafts/templates (`drafts`)
- equipment and sets (`equipment`, `equipment_sets`, `set_types`)
- dictionaries (`warehouses`, `project_types`)
- linking tables (`equipment_in_project`, `equipment_in_draft`)

Important: API and payload naming follow legacy style (`snake_case`) to stay compatible with the legacy app contract.

## Quick Start

### 1. Generate root `.env` for Docker Compose

```bash
python3 scripts/generate_compose_env.py \
  --acme-email your-email@example.com \
  --web-host home.your-domain.tld \
  --api-host home-server.your-domain.tld \
  --postgres-port 5433
```

Generated `.env` includes `POSTGRES_PORT` (default `5433`).

### 2. Start full stack

```bash
docker compose up -d --build --remove-orphans
```

DB migrations are applied automatically by the `migrate` service before `server` starts.

## Local Development

### Backend

```bash
cd server
cp example.env .env
make docker-db-up
make migrate-up
make air
```

### Frontend

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

## Legacy Data Migration (`LEGACY_APP`)

If you import an existing DB from `LEGACY_APP/ultralive-app-legacy`, run current Go migrations after restore:

```bash
cd server
make migrate-up
```

Migration behavior for compatibility:

- keeps legacy naming in CRM tables/endpoints
- upgrades/patches `users` for new auth model
- if legacy `auths` exists, backfills `users.email` and `users.password` from `auths` using `user_uid/auth_uid`

Recommended safety flow:

1. Backup database before migration.
2. Restore legacy dump.
3. Run `make migrate-up`.
4. Verify users/projects/equipment in UI.

## Verification

- Backend tests: `cd server && make test`
- Frontend build: `cd web && npm run build`
- E2E: `cd web && npm run test:e2e`

## XMPP Release Bot (systemd)

You can run a dedicated XMPP bot as a systemd service that accepts deploy commands (for example from Prosody account `ultralive-release-bot@vyachik-dev.ru`).

Supported bot commands:

- `/help`
- `/ping`
- `/status`
- `/update` (runs deploy workflow: checkout main + fetch/pull + `docker compose up -d --build --remove-orphans`)

### 1. Generate `.env.release-bot`

```bash
python3 scripts/generate_release_bot_env.py \
  --jid ultralive-release-bot@vyachik-dev.ru \
  --password 'BOT_XMPP_PASSWORD' \
  --allowed-sender your-admin@vyachik-dev.ru \
  --repo-path /opt/ultralive-app-v2
```

Notes:

- `--allowed-sender` can be passed multiple times.
- `.env.release-bot` is ignored by git.

### 2. Install and start systemd service

```bash
scripts/setup_release_bot_systemd.sh install
```

### 3. Manage service

```bash
scripts/release_bot_service.sh status
scripts/release_bot_service.sh logs
scripts/release_bot_service.sh restart
scripts/release_bot_service.sh stop
scripts/setup_release_bot_systemd.sh uninstall
```

## Repository Structure

- `server/`: API, services, migrations, tests
- `web/`: pages, stores, Nuxt server proxy routes
- `docs/`: project documentation
- `LEGACY_APP/`: legacy reference app used for contract/domain parity

## Backend Service Layout

HTTP handlers are split by CRM table/domain under `server/service/<name>/service.go`:

- `settype`, `projecttype`, `warehouse`, `equipmentset`
- `equipment`, `project`, `draft`
- `equipmentinproject`, `equipmentindraft`

Shared SQL access remains in `server/service/tracker/store.go`.
