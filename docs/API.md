# API Reference

Base URL: `http://localhost:8000/api/v1`

Swagger UI (requires auth): `http://localhost:8000/swagger/index.html`

## Authentication

Public endpoints:

- `POST /register`
- `POST /login`

Protected endpoints accept one of:

- `Authorization: Bearer <token>`
- `Authorization: <token>`
- `task_tracker_token` cookie

## Payload Naming

CRM endpoints keep legacy naming conventions:

- route names: `snake_case` (`/equipment_set`, `/equipment_in_project`, ...)
- JSON fields: `snake_case` (`project_id`, `equipment_set_name`, ...)

## User/Auth Endpoints

- `POST /register`
- `POST /login`
- `GET /profile`
- `PUT /profile`
- `PUT /profile/password`
- `GET /users`
- `GET /users/{id}`
- `GET /users/search/{name}`
- `GET /users/lookup`

## CRM Dictionary Endpoints

### Set Types

- `GET /set_types/`
- `GET /set_types/{id}`
- `POST /set_types/`
- `PUT /set_types/{id}`
- `DELETE /set_types/{id}`

### Project Types

- `GET /project_types/`
- `GET /project_types/{id}`
- `POST /project_types/`
- `PUT /project_types/{id}`
- `DELETE /project_types/{id}`

### Warehouses

- `GET /warehouse/`
- `POST /warehouse/`
- `PUT /warehouse/{id}`
- `DELETE /warehouse/{id}`

## CRM Entity Endpoints

### Equipment Sets

- `GET /equipment_set/`
- `GET /equipment_set/search/{id}`
- `GET /equipment_set/maintenance`
- `GET /equipment_set/storage`
- `POST /equipment_set/`
- `PUT /equipment_set/{id}`
- `DELETE /equipment_set/{id}`

### Equipment

- `GET /equipment/`
- `GET /equipment/set/{id}`
- `GET /equipment/search/{id}`
- `POST /equipment/`
- `PUT /equipment/{id}`
- `DELETE /equipment/{id}` (returns `204 No Content`)

### Projects

- `GET /projects/`
- `GET /projects/archived`
- `GET /projects/search/{id}`
- `POST /projects/`
- `PUT /projects/{id}`
- `DELETE /projects/{id}`

### Drafts

- `GET /drafts/`
- `GET /drafts/search/{id}`
- `POST /drafts/`
- `PUT /drafts/{id}`
- `DELETE /drafts/{id}`

## Linking Endpoints

### Equipment in Project

- `GET /equipment_in_project/{id}`
- `POST /equipment_in_project/add`
- `PUT /equipment_in_project/del`
- `POST /equipment_in_project/add_set`
- `PUT /equipment_in_project/del_set`
- `POST /equipment_in_project/equipment_in_set`
- `POST /equipment_in_project/conflicting`
- `DELETE /equipment_in_project/reset/{id}`
- `POST /equipment_in_project/add_draft`
- `POST /equipment_in_project/conflicting_projects`

### Equipment in Draft

- `GET /equipment_in_draft/{id}`
- `POST /equipment_in_draft/add`
- `PUT /equipment_in_draft/del`
- `POST /equipment_in_draft/add_set`
- `PUT /equipment_in_draft/del_set`
- `POST /equipment_in_draft/equipment_in_set`

## Example Requests

### Register

```http
POST /api/v1/register
Content-Type: application/json

{
  "firstName": "Alice",
  "lastName": "Smith",
  "email": "alice@example.com",
  "password": "secret123"
}
```

### Login

```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "alice@example.com",
  "password": "secret123"
}
```

Success (`200`):

```json
{
  "token": "<jwt>"
}
```

### Add Equipment to Project

```http
POST /api/v1/equipment_in_project/add
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "project_id": 10,
  "equipment_id": 25
}
```

## Error Shape

Errors are returned as JSON. Typical statuses:

- `400` invalid payload/validation/invalid reference
- `403` unauthorized or permission denied
- `404` entity not found
- `500` unexpected server/database error
