# AdOrbit Workspace

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-17-336791?logo=postgresql&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" />
  <img alt="Swagger" src="https://img.shields.io/badge/OpenAPI-Swagger-85EA2D?logo=swagger&logoColor=black" />
</p>

AdOrbit Workspace is the Go backend service that creates and persists workspace records for AdOrbit users. The codebase implements a small HTTP API, validates bearer JWT tokens, enforces request idempotency for workspace creation, stores data in PostgreSQL, and serves generated Swagger documentation.

## Table of contents

- [Overview](#overview)
- [Current capabilities](#current-capabilities)
- [Architecture](#architecture)
- [Technology stack](#technology-stack)
- [Repository layout](#repository-layout)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Environment variables](#environment-variables)
- [Running the service](#running-the-service)
- [API reference](#api-reference)
- [Database schema](#database-schema)
- [Testing](#testing)
- [Security considerations](#security-considerations)
- [License](#license)

## Overview

This repository is a compact, single-service backend for the AdOrbit workspace domain. The current implementation centers on one primary workflow: creating an authenticated user workspace and persisting it to PostgreSQL.

The runtime boot sequence is straightforward and verifiable from the code:

1. `cmd/server/main.go` loads environment variables through `godotenv` and reads configuration from `internal/config`.
2. `internal/app.New` opens a PostgreSQL connection, runs migrations, and initializes the HTTP router.
3. `internal/api.NewRouter` registers the health endpoint, Swagger UI, and the versioned workspace route.
4. `internal/workspace` handles request validation, JWT-based auth, and the `workspaces` insert workflow.

This is not a multi-service application; it is a focused backend service with a small API surface and one persistence model.

## Current capabilities

- Create a workspace for an authenticated user via `POST /api/v1/workspace`.
- Require a bearer JWT to access the workspace route.
- Require an `Idempotency-Key` header for duplicate-safe workspace creation requests.
- Store workspace records in a PostgreSQL `workspaces` table.
- Generate and serve Swagger docs from the `docs/` directory.
- Run locally with either Docker Compose or a native Go process.

## Architecture

```mermaid
flowchart LR
    Client[Client / API consumer] -->|HTTP + Bearer JWT| API[API router\ninternal/api]
    API --> JWT[JWT middleware]
    API --> ID[Idempotency middleware]
    JWT --> Handler[Workspace handler]
    ID --> Handler
    Handler --> Service[Workspace service]
    Service --> Repo[Workspace repository]
    Repo --> PG[(PostgreSQL\nworkspaces table)]
    API --> Health[GET /health]
    API --> Swagger[Swagger UI /swagger/]
```

Important implementation details confirmed in the repository:

- `internal/app/server.go` starts a `net/http` server with timeouts for reads, writes, and shutdown.
- `internal/app/dependencies.go` wires the JWT middleware, idempotency middleware, repository, and workspace handler together.
- `internal/workspace/routes.go` registers the workspace endpoint and wraps it with authentication and idempotency checks.
- `internal/database/migrate.go` applies application migrations from `internal/database/migrations` and also runs the idempotency schema migrations from the shared Go library.

## Technology stack

| Category | Technology | Evidence in repo | Purpose |
| --- | --- | --- | --- |
| Language | Go 1.27.1 | `go.mod`, `Dockerfile`, `Dockerfile.dev` | Backend implementation |
| HTTP server | Go `net/http` | `internal/api/api.go`, `internal/app/server.go` | Routing and HTTP lifecycle |
| Database | PostgreSQL | `docker-compose.yml`, `internal/database/migrations/*.sql` | Persistent workspace storage |
| Database driver | `pgx` | `go.mod`, `internal/database/database.go` | PostgreSQL connectivity |
| Migrations | `golang-migrate/migrate` | `internal/database/migrate.go` | Schema management |
| Authentication | JWT validation middleware | `internal/config/config.go`, `internal/app/dependencies.go` | Bearer token validation |
| Idempotency | `adorbit-go/idempotency` | `internal/app/dependencies.go`, `internal/workspace/routes.go` | Duplicate-safe request handling |
| API docs | Swagger / OpenAPI | `docs/swagger.yaml`, `docs/swagger.json`, `.air.toml` | API documentation UI |
| Environment config | `godotenv` | `cmd/server/main.go` | `.env` loading at startup |
| Containerization | Docker + Compose | `Dockerfile`, `Dockerfile.dev`, `docker-compose.yml` | Local and packaged runtime |
| Dev tooling | Air + Delve + Swag | `.air.toml`, `Dockerfile.dev` | Live reload, debugging, docs generation |

## Repository layout

```text
.
├── .air.toml
├── .env
├── .gitignore
├── Dockerfile
├── Dockerfile.dev
├── README.md
├── README_PROMPT.md
├── Start prompt.txt
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
│       ├── model.go
│       ├── repository.go
│       ├── routes.go
│       └── service.go
├── reset_docker.sh
└── tmp/
```

Notable areas:

- `cmd/server` is the application entrypoint.
- `internal/api` defines the router and the health route.
- `internal/workspace` contains the domain model, service logic, handler, and route registration.
- `internal/database` handles PostgreSQL connection and schema migration execution.
- `docs/` contains generated Swagger output used by the `/swagger/` UI.

## Prerequisites

To run the project locally, you need:

- Go 1.27.1 or a compatible toolchain
- Docker and Docker Compose for the included local stack
- PostgreSQL access for the application runtime and migrations
- A valid JWT public key in the environment
- A 64-character hex value for `IDEMPOTENCY_SCOPE_KEY`
- A 64-character hex value for `IDEMPOTENCY_ENCRYPTION_KEY`

The repository does not contain a root `LICENSE` file, so the repository does not declare a license in the current state.

## Setup

1. Clone the repository:

```bash
git clone https://github.com/thakur-vijay/AdOrbitWorkspace.git
cd AdOrbitWorkspace
```

2. Create or update the required environment variables. The project loads values from a `.env` file and from the process environment.

3. Start the database and app with Docker Compose:

```bash
docker compose up --build
```

4. Or run the server directly:

```bash
go run ./cmd/server
```

The repo includes both a production Dockerfile and a development Dockerfile. The development setup uses Air and Delve as configured in `.air.toml`.

## Environment variables

The application reads the following variables directly in `internal/config/config.go`.

| Variable | Required | Notes |
| --- | --- | --- |
| `APP_PORT` | No | HTTP port. Defaults to `8080`; the loader normalizes it to `:port` format when needed. |
| `WORKSPACE_DATABASE_URL` | Yes | PostgreSQL connection string used by the app and migration process. |
| `JWT_PUBLIC_KEY` | Yes | Public key used to validate incoming bearer JWT tokens. Literal `\n` sequences are converted to real newlines. |
| `IDEMPOTENCY_SCOPE_KEY` | Yes | 64-character hexadecimal key used for hashing idempotency scopes. |
| `IDEMPOTENCY_ENCRYPTION_KEY` | Yes | 64-character hexadecimal key used for encrypting idempotency responses. |
| `SWAGGER_HOST` | No | Optional override used when generating Swagger host metadata. |

The Compose stack also relies on these variables:

| Variable | Required | Notes |
| --- | --- | --- |
| `WORKSPACE_POSTGRES_USER` | Yes | PostgreSQL user for the local Compose database. |
| `WORKSPACE_POSTGRES_PASSWORD` | Yes | PostgreSQL password for the local Compose database. |
| `WORKSPACE_POSTGRES_DB` | Yes | PostgreSQL database name. |
| `WORKSPACE_POSTGRES_HOST_PORT` | Yes | Host port mapped to the PostgreSQL container. |
| `WORKSPACE_HOST_PORT` | Yes | Host port mapped to the application container. |

Use placeholders instead of real secrets when documenting or sharing configuration.

## Running the service

### Docker Compose

The repository defines a `workspace` service and a `postgres` service in `docker-compose.yml`.

```bash
docker compose up --build
```

This starts:

- PostgreSQL on the configured host port,
- the Go application on the configured `WORKSPACE_HOST_PORT`,
- the project mounted in `/app` for development.

### Local Go process

```bash
export WORKSPACE_DATABASE_URL="postgres://user:password@localhost:5432/adorbit_workspace?sslmode=disable"
export JWT_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
export IDEMPOTENCY_SCOPE_KEY="<64 hexadecimal characters>"
export IDEMPOTENCY_ENCRYPTION_KEY="<64 hexadecimal characters>"

go run ./cmd/server
```

The server listens on `APP_PORT` using a normalized `:port` address, for example `:8080` or `:8082`.

## API reference

The service exposes a very small API surface.

| Method | Endpoint | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | No | Returns a JSON payload: `{"status":"ok"}` |
| `GET` | `/swagger/` | No | Serves the generated Swagger UI |
| `POST` | `/api/v1/workspace` | Bearer JWT + `Idempotency-Key` header | Creates a workspace for the authenticated user |

The workspace creation payload matches `CreateWorkspaceDTO` in `internal/workspace/model.go`.

### POST /api/v1/workspace

Request body example:

```json
{
  "name": "Northwind Media",
  "type": "advertiser",
  "country": "US",
  "timezone": "America/New_York",
  "currency": "USD",
  "industry": "Retail",
  "website": "https://example.com",
  "description": "Advertising workspace for the Northwind account."
}
```

Required fields:

- `name` — required string, max 150 chars
- `type` — required enum: `advertiser` or `publisher`
- `country` — required 2-character ISO country code
- `timezone` — required string, max 100 chars
- `currency` — required 3-character currency code

Optional fields:

- `industry` — optional string, max 100 chars
- `website` — optional URL
- `description` — optional string

The route is protected by JWT middleware and idempotency middleware. The idempotency scope is derived from the authenticated user, so duplicate create requests are isolated per user.

## Database schema

The initial migration is defined in `internal/database/migrations/000001_create_workspace.up.sql` and creates a single table named `workspaces`.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `UUID` | Primary key generated via `gen_random_uuid()` |
| `name` | `VARCHAR(150)` | Not null |
| `slug` | `VARCHAR(180)` | Unique, generated from the workspace name |
| `type` | `VARCHAR(20)` | Restricted to `advertiser` or `publisher` |
| `country` | `CHAR(2)` | Not null |
| `timezone` | `VARCHAR(100)` | Not null |
| `currency` | `CHAR(3)` | Not null |
| `industry` | `VARCHAR(100)` | Optional |
| `website` | `TEXT` | Optional |
| `description` | `TEXT` | Optional |
| `status` | `VARCHAR(20)` | Defaults to `active`; allowed values are `active` and `deleted` |
| `created_by` | `UUID` | Authenticated user ID |
| `created_at` | `TIMESTAMPTZ` | Defaults to `NOW()` |
| `updated_at` | `TIMESTAMPTZ` | Defaults to `NOW()` |

Indexes created by the migration:

- `idx_workspaces_created_by`
- `idx_workspaces_type`
- `idx_workspaces_status`

## Testing

The repository includes Go tests for config loading, app lifecycle, health checks, and Swagger host setup. The project test suite is run with:

```bash
go test ./...
```

Current repository state: this command currently fails in `cmd/server/main_test.go` because the Swagger host assertion expects `localhost:8081` but the implementation builds `localhost8081` without the colon. The test suite is therefore present but not fully green in the current state of the repo.

## Security considerations

The codebase includes a few explicit security controls that are directly implemented:

- JWT validation with `JWT_PUBLIC_KEY`.
- Route-level authentication on the workspace endpoint via the JWT middleware.
- Idempotency enforcement for workspace creation to reduce duplicate processing.
- PostgreSQL table constraints such as `type` and `status` checks in the migration.

This repository does not show a broader identity, RBAC, rate-limit, or secret-management layer beyond those checks.

## License

There is no license file in the repository root or the current project tree, so the project does not currently declare a license in the repository state.
