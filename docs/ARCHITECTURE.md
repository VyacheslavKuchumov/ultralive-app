# Architecture

## High-Level Components

- `web` (Nuxt 4 + Nuxt UI): CRM UI, auth state, server-side proxy routes
- `server` (Go + chi): JWT auth, CRM HTTP services, SQL store layer
- `postgres` (PostgreSQL 16): CRM data storage

## Backend Structure (`server/`)

- `cmd/main.go`: API entrypoint
- `cmd/server/server.go`: HTTP router and service wiring
- `cmd/migrate/migrations/`: SQL migrations
- `service/user/`: auth and user profile handlers
- `service/<table>/`: one HTTP service per CRM table/domain
- `service/crmhttp/`: shared HTTP helpers (auth check, payload validation, error mapping)
- `service/tracker/store.go`: shared SQL store implementation for CRM entities
- `types/`: payloads and response structs (legacy-compatible `snake_case` fields)

### Service-per-table routing

`server/cmd/server/server.go` wires separate service packages for:

- `set_types`
- `project_types`
- `warehouse`
- `equipment_set`
- `equipment`
- `projects`
- `drafts`
- `equipment_in_project`
- `equipment_in_draft`

## Request Flow

### Public auth flow

1. Frontend calls Nuxt route (`/api/auth/login` or `/api/auth/register`).
2. Nuxt route proxies request to backend `/api/v1/login` or `/api/v1/register`.
3. Backend validates payload, creates/validates user, returns JSON.
4. On login backend sets `task_tracker_token` cookie and returns JWT token.

### Protected CRM flow

1. Frontend sends JWT via `Authorization` header (or cookie).
2. Backend middleware validates JWT and injects user ID into request context.
3. CRM service handler validates path/body and calls store methods.
4. Store executes SQL against legacy-compatible CRM schema.

## Data Model (CRM)

Core tables used by API:

- `users`
- `set_types`
- `project_types`
- `warehouses`
- `equipment_sets`
- `equipment`
- `projects`
- `drafts`
- `equipment_in_project`
- `equipment_in_draft`

Naming is intentionally aligned with legacy app contract to keep migration compatibility.
