# AGENTS.md

## Purpose

Project-specific context for coding agents in the `ultralive-app` repository.

## Project Summary

Ultralive CRM is a full-stack app for equipment accounting during shoots:

- Go backend API (`server/`)
- Nuxt 4 frontend (`web/`)
- PostgreSQL
- SQL migrations with `golang-migrate`
- Docker Compose orchestration at repository root

## Repository Layout

- `server/`: Go API, business logic, migrations, tests
- `web/`: Nuxt UI app, Pinia stores, Nuxt server API proxy routes
- `LEGACY_APP/`: legacy app used as domain/contract reference
- `docker-compose.yml`: postgres + server + web + traefik
- `docs/`: project documentation

## Naming Contract (Critical)

Keep legacy API/domain naming:

- Route names like `/projects`, `/equipment_set`, `/warehouse`, `/equipment_in_project`
- Payload/response field names in `snake_case`
- Dictionary/entity names as in legacy app

Do not silently migrate to task-tracker naming (`goals/tasks`) in new code.

## Core Runtime Flow

1. Frontend calls Nuxt server routes under `web/server/api/...` (primary proxy: `web/server/api/backend/[...path].ts`).
2. Nuxt route handlers proxy requests to Go API `/api/v1/...` via `callBackend`.
3. Go API validates JWT, runs table-specific service handlers, and executes CRM operations in PostgreSQL.

## Backend Service Layout

- HTTP services are split by domain/table in `server/service/<name>/service.go`.
- Shared SQL repository is `server/service/tracker/store.go`.
- Shared auth/validation/error HTTP helpers are in `server/service/crmhttp/helpers.go`.

## Auth Model

Use the new-project auth approach:

- single `users` table for identity + credentials
- no split `auths + users` model in runtime
- protected routes accept:
  - `Authorization: Bearer <token>`
  - `task_tracker_token` cookie

## Commands Agents Should Use

### Full stack (recommended)

```bash
docker compose up -d --build --remove-orphans
```

`docker compose` runs migrations automatically via the `migrate` service before `server`.

### Standalone Docker testing

```bash
python3 scripts/dev_docker_stack.py up
python3 scripts/dev_docker_stack.py down --remove-network
```

`dev_docker_stack.py up` runs DB migrations before starting backend/frontend containers.

### Backend local development

```bash
cd server
cp example.env .env
make docker-db-up
make migrate-up
make air
```

### Frontend local development

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

## Testing and Verification

Before finishing backend-impacting changes:

```bash
cd server
make test
```

Before finishing frontend-impacting changes:

```bash
cd web
npm run build
```

For flow/auth/integration changes (requires backend + frontend):

```bash
cd web
npm run test:e2e
```

## Migration Notes

- Migrations live in `server/cmd/migrate/migrations`.
- File names may contain old labels (`products`, `orders`) but current schema is CRM.
- Run migrations with:

```bash
cd server
make migrate-up
```

### Legacy DB Import Compatibility

When migrating data from `LEGACY_APP/ultralive-app-legacy`:

1. Restore legacy DB dump.
2. Run current migrations.
3. Verify users/projects/equipment in UI.

Compatibility behavior in migration:

- patches `users` schema for the new single-table auth model
- backfills `users.email/password` from legacy `auths` when present (`user_uid/auth_uid` mapping)

Always make a DB backup before running migrations on imported legacy data.

## Environment Notes

- Backend env: `server/.env`
- Frontend env: `web/.env` (`BACKEND_URL`)
- Default local ports:
  - Web: `3000`
  - API: `8000`
  - Postgres host port: `5433`

## Change Discipline

- Keep backend contracts and frontend store/proxy calls in sync.
- If endpoint shapes change, update in same PR:
  - `server/service/...`
  - `web/server/api/...`
  - `web/app/stores/...`
  - `docs/API.md`
- Prefer focused edits and keep tests/build green.

## Git Workflow Reminder

- Do not work directly on `main`.
- Sync first:
  - `git fetch origin`
  - `git switch main`
  - `git pull --ff-only origin main`
- Create a feature branch for each task.
- Run relevant verification commands before commit.
- Open PR to `main`; do not merge directly to `main`.

## Commit Message Guidelines

- Use imperative, focused subject.
- Prefer `type: short summary` (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`).
- Keep subject concise (~50-72 chars).

## Pull Request Guidelines

PR description should include:

- what changed
- why it changed
- how it was tested (commands)

If contracts changed, include backend + frontend + docs updates in the same PR.
