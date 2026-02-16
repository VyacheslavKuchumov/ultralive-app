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
  --api-host home-server.your-domain.tld
```

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
