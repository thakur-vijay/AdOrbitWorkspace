# AdOrbit Workspace

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-17-336791?logo=postgresql&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" />
  <img alt="Swagger" src="https://img.shields.io/badge/OpenAPI-Swagger-85EA2D?logo=swagger&logoColor=black" />
  <img alt="CI" src="https://github.com/thakur-vijay/AdOrbitWorkspace/actions/workflows/go-tests.yml/badge.svg" />
</p>

AdOrbit Workspace is a lightweight Go service that manages workspace records for authenticated AdOrbit users. The current implementation is intentionally focused: it exposes a small REST API for creating, listing, fetching, updating, and soft-deleting workspaces, backed by PostgreSQL and protected by JWT and idempotency middleware.

## Table of contents

- [Overview](#overview)
- [Current scope](#current-scope)
- [Architecture](#architecture)
- [Technology stack](#technology-stack)
- [Project structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup and configuration](#setup-and-configuration)
- [Environment variables](#environment-variables)
- [Running locally](#running-locally)
- [API reference](#api-reference)
- [Database schema](#database-schema)
- [Testing and CI](#testing-and-ci)
- [Security considerations](#security-considerations)
- [Deployment status](#deployment-status)
- [License](#license)

## Overview

This repository contains the workspace backend for the AdOrbit domain, not a full multi-service platform. The service follows a layered Go architecture:

- HTTP routing and health checks live in `internal/api`
- request validation and response handling live in `internal/workspace/handler.go`
- business logic is handled in `internal/workspace/service.go`
- persistence is implemented in `internal/workspace/repository.go`
- PostgreSQL connectivity and schema migration are handled in `internal/database`

At startup, `cmd/server/main.go` loads environment variables, validates required configuration, connects to PostgreSQL, runs migrations, and starts the HTTP server. The app exposes both the API and generated Swagger documentation on the same service.

## Current scope

The repository currently implements a single domain workflow around `workspaces`:

- `POST /api/v1/workspaces` creates a workspace for the authenticated user
- `GET /api/v1/workspaces` lists active workspaces belonging to the authenticated user
- `GET /api/v1/workspaces/{id}` fetches a workspace by ID
- `PATCH /api/v1/workspaces/{id}` partially updates a workspace
- `DELETE /api/v1/workspaces/{id}` soft-deletes a workspace by setting `status = 'deleted'`

Supported workspace types are:

- `advertiser`
- `publisher`

Supported statuses are:

- `active`
- `deleted`

The request layer validates input payloads, enforces bearer-token authentication, and requires an `Idempotency-Key` header for create, update, and delete operations.

## Architecture

```mermaid
flowchart LR
    Client[Client / API consumer] -->|HTTP + Bearer JWT| API[API router\ninternal/api]
    API --> Auth[JWT middleware]
    API --> Idem[Idempotency middleware]
    Auth --> Handler[Workspace handler]
    Idem --> Handler
    Handler --> Service[Workspace service]
    Service --> Repo[Workspace repository]
    Repo --> DB[(PostgreSQL\nworkspaces table)]
    API --> Health[GET /health]
    API --> Swagger[Swagger UI\n/swagger/]
```

Key implementation details verified from the repository:

- `internal/app/app.go` creates the database connection and runs migrations during bootstrap
- `internal/app/dependencies.go` wires JWT middleware, idempotency middleware, the repo, and the handler
- `internal/workspace/routes.go` registers all workspace endpoints and scopes idempotency keys by user ID
- `internal/database/migrate.go` applies both workspace schema migrations and idempotency repository migrations
- `internal/api/api.go` serves the health endpoint and Swagger UI, and mounts versioned routes under `/api/v1`

## Technology stack

| Category | Technology | Verified in repo | Purpose |
| --- | --- | --- | --- |
| Language | Go 1.27.1 | `go.mod`, `Dockerfile`, `Dockerfile.dev` | Primary backend implementation |
| HTTP server | Go `net/http` | `internal/api/api.go`, `internal/app/server.go` | Routing, handlers, lifecycle management |
| Database | PostgreSQL 17 | `docker-compose.yml`, migration SQL | Persistent workspace storage |
| Driver | `pgx/v5` | `go.mod`, `internal/database/database.go` | PostgreSQL connectivity |
| Migrations | `golang-migrate/migrate` | `internal/database/migrate.go` | Schema and migration execution |
| Auth | JWT validation middleware | `internal/app/dependencies.go`, `internal/config/config.go` | Verifies bearer tokens |
| Idempotency | `github.com/thakur-vijay/adorbit-go/idempotency` | `internal/app/dependencies.go`, `internal/workspace/routes.go` | Protects retries for write operations |
| API docs | Swagger / OpenAPI | `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go` | Generated API documentation |
| Live reload | Air | `.air.toml`, `Dockerfile.dev` | Watches and rebuilds during development |
| Environment loading | `godotenv` | `cmd/server/main.go` | Loads `.env` at startup |
| Containerization | Docker + Compose | `Dockerfile`, `Dockerfile.dev`, `docker-compose.yml` | Local orchestration of app + Postgres |
| CI | GitHub Actions | `.github/workflows/go-tests.yml` | Automated Go test execution |

## Project structure

```text
.
├── .air.toml
├── .env
├── .github/
│   └── workflows/
│       └── go-tests.yml
├── .gitignore
├── Dockerfile
├── Dockerfile.dev
├── README.md
├── README_PROMPT.md
├── Start prompt.txt
├── WorkspaceTesting.md
├── cmd/
│   └── server/
│       ├── main.go
│       └── main_test.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal/
│   ├── api/
│   │   ├── api.go
│   │   └── api_test.go
│   ├── app/
│   │   ├── app.go
│   │   ├── app_test.go
│   │   ├── dependencies.go
│   │   ├── dependencies_test.go
│   │   └── server.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── database/
│   │   ├── database.go
│   │   ├── migrate.go
│   │   └── migrations/
│   │       ├── 000001_create_workspace.down.sql
│   │       └── 000001_create_workspace.up.sql
│   └── workspace/
│       ├── handler.go
│       ├── handler_test.go
│       ├── model.go
│       ├── repository.go
│       ├── routes.go
│       ├── service.go
│       └── test_helpers_test.go
├── reset_docker.sh
├── temp.txt
└── tmp/
```

Important directories:

- `cmd/server` contains the application entrypoint and Swagger setup
- `internal/api` mounts the HTTP router and health endpoints
- `internal/workspace` contains the workspace model, service, repository, handler, and route registration
- `internal/database` handles PostgreSQL setup and migration execution
- `docs/` stores generated Swagger output used by the `/swagger/` route

## Prerequisites

To run the service locally, you need:

- Go 1.27.1 or a compatible Go toolchain
- PostgreSQL 17 or a compatible database instance
- Docker and Docker Compose for the included local stack
- A valid JWT public key for request authentication
- Two 64-character hexadecimal secret keys for idempotency scope and response encryption

This repository does not include a license file, so no explicit license is declared in the codebase.

## Setup and configuration

1. Clone the repository:

```bash
git clone https://github.com/thakur-vijay/AdOrbitWorkspace.git
cd AdOrbitWorkspace
```

2. Create a local `.env` file or export the required environment variables before startup.

3. Start the PostgreSQL + app stack with Docker Compose:

```bash
docker compose up --build
```

4. To run the service directly with Go instead of Docker:

```bash
go run ./cmd/server
```

The app automatically loads `.env` via `godotenv.Load(".env")` in `cmd/server/main.go`.

## Environment variables

The application reads the following variables from the environment at runtime. Use placeholders for secrets rather than real credentials.

| Variable | Required | Default / usage |
| --- | --- | --- |
| `APP_PORT` | No | HTTP port for the server. Defaults to `8080`; the config normalizes it to `:8080` when unset. |
| `WORKSPACE_DATABASE_URL` | Yes | PostgreSQL connection string used by the app and migration runner. |
| `JWT_PUBLIC_KEY` | Yes | Public RSA key used to validate incoming bearer tokens. Escaped `\n` sequences are converted to real newlines. |
| `IDEMPOTENCY_SCOPE_KEY` | Yes | 64-character hexadecimal string used for idempotency scope hashing. |
| `IDEMPOTENCY_ENCRYPTION_KEY` | Yes | 64-character hexadecimal string used for idempotency response encryption. |
| `SWAGGER_HOST` | No | Overrides the host value rendered in the generated Swagger UI. |
| `WORKSPACE_POSTGRES_USER` | Compose-only | PostgreSQL user for the local Compose database. |
| `WORKSPACE_POSTGRES_PASSWORD` | Compose-only | PostgreSQL password for the local Compose database. |
| `WORKSPACE_POSTGRES_DB` | Compose-only | PostgreSQL database name used in Compose. |
| `WORKSPACE_POSTGRES_HOST_PORT` | Compose-only | Host port exposed for PostgreSQL in Compose. |
| `WORKSPACE_HOST_PORT` | Compose-only | Host port exposed for the app in Compose. |

The checked-in `.env` file in this repository is a local development configuration and should be treated as configuration data, not production secrets.

## Running locally

### Docker Compose

The repository defines a `postgres` service and a `workspace` service in `docker-compose.yml`.

```bash
docker compose up --build
```

This starts:

- PostgreSQL on the configured host port
- the Go server on the configured application port
- the repo mounted into `/app` for development with `air`-based reloading

### Direct Go execution

```bash
export WORKSPACE_DATABASE_URL="postgres://user:password@localhost:5432/adorbit_workspace?sslmode=disable"
export JWT_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
export IDEMPOTENCY_SCOPE_KEY="<64 hex chars>"
export IDEMPOTENCY_ENCRYPTION_KEY="<64 hex chars>"

go run ./cmd/server
```

The app binds to the configured `APP_PORT`, and Swagger defaults to `http://localhost:PORT` when `SWAGGER_HOST` is not set. In the provided local config, the app listens on `:8082`.

### Swagger UI

Once the service is running, Swagger is available at the application host:

- `http://localhost:8082/swagger/` for the default local configuration
- `http://localhost:<APP_PORT>/swagger/` when a different port is configured

## API reference

The service is versioned under `/api/v1` and currently exposes the following endpoints.

| Method | Endpoint | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | None | Returns `{"status":"ok"}` |
| `GET` | `/swagger/` | None | Serves Swagger UI |
| `POST` | `/api/v1/workspaces` | Bearer JWT + `Idempotency-Key` | Create a workspace |
| `GET` | `/api/v1/workspaces` | Bearer JWT | Fetch active workspaces for the authenticated user |
| `GET` | `/api/v1/workspaces/{id}` | Bearer JWT | Fetch a single workspace by ID |
| `PATCH` | `/api/v1/workspaces/{id}` | Bearer JWT + `Idempotency-Key` | Partially update a workspace |
| `DELETE` | `/api/v1/workspaces/{id}` | Bearer JWT + `Idempotency-Key` | Soft-delete a workspace |

### Create workspace payload

```json
{
  "name": "Acme Media",
  "type": "advertiser",
  "country": "US",
  "timezone": "America/New_York",
  "currency": "USD",
  "industry": "Marketing",
  "website": "https://example.com",
  "description": "Regional media buying workspace"
}
```

Validation rules verified in the code:

- `name` is required and limited to 150 characters
- `type` must be `advertiser` or `publisher`
- `country` must be exactly 2 characters
- `timezone` must not exceed 100 characters
- `currency` must be exactly 3 characters
- `industry` is optional and capped at 100 characters
- `website` is optional, but must be a valid URL when provided
- `PATCH` requests require at least one field to be present

### Idempotency behavior

Write operations are wrapped by the idempotency middleware in `internal/workspace/routes.go`:

- `POST /workspaces` requires `Idempotency-Key`
- `PATCH /workspaces/{id}` requires `Idempotency-Key`
- `DELETE /workspaces/{id}` requires `Idempotency-Key`

The scope is derived from the authenticated user ID, so retries for the same user and operation are tracked and protected.

## Database schema

The database schema is created by the migration in `internal/database/migrations/000001_create_workspace.up.sql`.

```sql
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(180) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('advertiser', 'publisher')),
    country CHAR(2) NOT NULL,
    timezone VARCHAR(100) NOT NULL,
    currency CHAR(3) NOT NULL,
    industry VARCHAR(100),
    website TEXT,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Indexes created by the migration:

- `idx_workspaces_created_by`
- `idx_workspaces_type`
- `idx_workspaces_status`

Important behaviors:

- Workspace records are soft-deleted by updating `status` to `deleted`
- `slug` is unique and generated from the workspace name using the `slug` library
- `Fetch` filters to `created_by = userID` and `status = 'active'`

## Testing and CI

The repository includes Go tests under the application packages and a GitHub Actions workflow for automated validation.

### Local test commands

```bash
go test ./...
go test -race -count=1 ./...
```

### CI workflow

The workflow in `.github/workflows/go-tests.yml` runs on push, pull request, and manual dispatch. It sets up a PostgreSQL 17 service and executes:

- `go test -count=1 ./...`
- `go test -race -count=1 ./...`
- `go test -count=10 ./internal/workspace/...`
- `go test -cover -count=1 ./...`
- coverage report generation with `go tool cover`

## Security considerations

The repository implements a few concrete security controls that are visible in code:

- JWT validation using a configured public key (`JWT_PUBLIC_KEY`)
- Authentication required on all workspace routes except `/health` and `/swagger/`
- Idempotency enforcement for state-changing requests using request headers and persisted server-side state
- Input validation using the response decoder and DTO validation tags
- PostgreSQL constraints for allowed workspace types and valid status values

This README does not claim a complete security posture beyond the behavior that is actually implemented in the repository. Sensitive keys and credentials are expected to be supplied through environment configuration, not committed to source control.

## Deployment status

The repository contains Docker and Compose configuration for local development, and a production build pipeline is generated through `Dockerfile`:

```bash
docker build -t adorbit-workspace .
```

The project does not define a deployment manifest for a cloud platform, orchestration system, or infrastructure-as-code stack in this repository. There is no verified deployment configuration beyond local container-based execution.

## License

No license file is present in the repository root or project tree, so no explicit license is declared for this codebase at the moment.
