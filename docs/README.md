# Project Documentation

This folder contains documentation for the Ultralive CRM project.

## Contents

- `docs/SETUP.md`: local and Docker setup
- `docs/ARCHITECTURE.md`: system design and service wiring
- `docs/API.md`: backend API contract and routes
- `docs/OPERATIONS.md`: migrations, tests, and troubleshooting
- `docs/CONTRIBUTING.md`: branch strategy, commits, and pull requests

## Quick Start

Run full stack with Docker:

```bash
docker compose up -d --build --remove-orphans
```

Open:

- Frontend: `http://localhost:3000`
- Backend API base: `http://localhost:8000/api/v1`
