# AdOrbit Workspace

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-17-336791?logo=postgresql&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" />
  <img alt="Swagger" src="https://img.shields.io/badge/OpenAPI-Swagger-85EA2D?logo=swagger&logoColor=black" />
</p>

AdOrbit Workspace is the backend service responsible for creating and storing workspace records for AdOrbit users. The service is implemented in Go with a small HTTP API, JWT-based authentication, idempotency protection, PostgreSQL persistence, and generated Swagger documentation.

## Table of contents

- [Overview](#overview)
- [Key capabilities](#key-capabilities)
- [Architecture](#architecture)
- [Technology stack](#technology-stack)
- [Project structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Environment variables](#environment-variables)
- [Running the application](#running-the-application)
- [API reference](#api-reference)
- [Database schema](#database-schema)
- [Testing](#testing)
- [Development workflow](#development-workflow)
- [Security considerations](#security-considerations)
- [Current scope and limitations](#current-scope-and-limitations)
- [License](#license)

## Overview

This repository contains a compact Go service for the AdOrbit workspace domain. Its current implementation centers on a single user-facing workflow: creating a workspace for an authenticated user.

The project is organized around a layered pattern:

- a `net/http` API layer for routing and health checks,
- a workspace service layer for business logic,
- a repository layer for PostgreSQL access,
- a database migration layer for schema setup,
- middleware for JWT authorization and idempotent request handling.

The app bootstraps itself from `cmd/server/main.go`, loads runtime configuration from environment variables, connects to PostgreSQL, runs database migrations, and then serves both the API and Swagger documentation.

## Key capabilities

- Create workspaces for an authenticated user with a required bearer token.
- Enforce JWT validation using a configured public key.
- Protect repeated create requests with idempotency middleware and per-user scopes.
- Persist workspace records in PostgreSQL using `pgx` and SQL migrations.
- Generate and serve Swagger/OpenAPI docs from the `docs/` directory.
- Run locally with Docker Compose or with a direct Go command.

## Architecture

The current repository is a small, single-service application rather than a multi-service platform.

```mermaid
flowchart LR
    Client[Client / API consumer] -->|HTTP + Bearer JWT| API[API Router\ninternal/api]
    API --> JWT[JWT middleware]
    API --> ID[Idempotency middleware]
    JWT --> WS[Workspace handler]
    ID --> WS
    WS --> SVC[Workspace service]
    SVC --> REPO[Workspace repository]
    REPO --> PG[(PostgreSQL\nworkspaces table)]
    API --> Health[GET /health]
    API --> Swagger[Swagger UI /swagger/]
```

Important implementation details from the codebase:

- The application creates a single top-level HTTP server in `internal/app/server.go`.
- `internal/app/dependencies.go` wires together the JWT middleware, idempotency middleware, repository, and handler.
- `internal/workspace/routes.go` registers the workspace route and wraps it with authentication and idempotency checks.
- `internal/database/migrate.go` applies both the application migrations and the idempotency schema migrations.

## Technology stack

| Category | Technology | Evidence in repo | Purpose |
| --- | --- | --- | --- |
| Language | Go 1.27.1 | `go.mod`, `Dockerfile.dev` | Primary application backend |
| Web server | Go `net/http` standard library | `internal/api/api.go`, `internal/app/server.go` | HTTP routing and server lifecycle |
| Database | PostgreSQL 17 | `docker-compose.yml`, migration SQL | Persistent workspace storage |
| DB driver | `pgx` | `go.mod`, `internal/database/database.go` | PostgreSQL connectivity |
| Migrations | `golang-migrate/migrate` | `internal/database/migrate.go` | Schema and idempotency migration management |
| Auth | JWT verification middleware | `internal/app/dependencies.go`, `internal/config/config.go` | Bearer token validation |
| API contracts | Swagger / OpenAPI | `docs/swagger.yaml`, `docs/swagger.json`, `.air.toml` | API documentation and UI |
| Dependency management | Go modules | `go.mod` | Versioned dependencies |
| Containerization | Docker + Compose | `Dockerfile.dev`, `docker-compose.yml` | Local development and service orchestration |
| Live reload | Air | `.air.toml` | Continuous rebuild during development |
| Environment config | `godotenv` | `cmd/server/main.go` | Loading `.env` files at startup |

## Project structure

```text
.
├── .air.toml
├── .env
├── .gitignore
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
└── tmp/
```

Notable directories:

- `cmd/server` contains the application entrypoint.
- `internal/api` exposes the HTTP router and health endpoint.
- `internal/workspace` implements the workspace domain object, service, handler, and route registration.
- `internal/database` handles PostgreSQL connectivity and migration execution.
- `docs/` is generated Swagger output used by `/swagger/`.

## Prerequisites

To run this service locally, you need:

- Go 1.27.1 or compatible tooling
- PostgreSQL 17 (or a compatible running instance)
- Docker and Docker Compose for the included local stack
- A valid JWT public key for verifying bearer tokens
- A 64-character hex idempotency scope key and 64-character hex encryption key

The repository does not include a license file, so no explicit software license is declared in the repo.

## Getting started

1. Clone the repository:

```bash
git clone https://github.com/thakur-vijay/AdOrbitWorkspace.git
cd AdOrbitWorkspace
```

2. Configure the required environment variables in a `.env` file or exported shell variables.

3. Start the local database and app with Docker Compose:

```bash
docker compose up --build
```

4. Or run the Go server directly:

```bash
go run ./cmd/server
```

When running natively, the application expects the database and JWT configuration to already be available in the environment.

## Environment variables

The code reads the following runtime settings from environment variables. Use placeholders for secret values rather than actual credentials.

| Variable | Required | Purpose |
| --- | --- | --- |
| `APP_PORT` | No | HTTP port for the Go server. Defaults to `8080` and is normalized with a leading `:` when needed. |
| `WORKSPACE_DATABASE_URL` | Yes | PostgreSQL connection string used by the app and migration runner. |
| `JWT_PUBLIC_KEY` | Yes | Public key used to validate incoming bearer JWT tokens. Newlines are normalized from escaped `\n` sequences. |
| `IDEMPOTENCY_SCOPE_KEY` | Yes | 64-character hexadecimal key used for idempotency scope hashing. |
| `IDEMPOTENCY_ENCRYPTION_KEY` | Yes | 64-character hexadecimal key used to encrypt idempotency responses. |
| `SWAGGER_HOST` | No | Overrides the Swagger host used in generated docs if set. |
| `WORKSPACE_POSTGRES_USER` | Yes for Compose stack | PostgreSQL database user used by Docker Compose. |
| `WORKSPACE_POSTGRES_PASSWORD` | Yes for Compose stack | PostgreSQL password used by Docker Compose. |
| `WORKSPACE_POSTGRES_DB` | Yes for Compose stack | PostgreSQL database name used by Docker Compose. |
| `WORKSPACE_POSTGRES_HOST_PORT` | Yes for Compose stack | Host port mapped to PostgreSQL container port `5432`. |
| `WORKSPACE_HOST_PORT` | Yes for Compose stack | Host port mapped to the application container port. |

The repository contains a local `.env` file with example values for development, but that file should not be treated as a source of truth for production secrets.

## Running the application

### Docker Compose

The repository defines a `workspace` service and a `postgres` service in `docker-compose.yml`.

```bash
docker compose up --build
```

This starts:

- PostgreSQL on the configured host port,
- the Go application on the configured `WORKSPACE_HOST_PORT`,
- the app with the project mounted into `/app` for development.

### Local Go process

```bash
export WORKSPACE_DATABASE_URL="postgres://user:password@localhost:5432/adorbit_workspace?sslmode=disable"
export JWT_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
export IDEMPOTENCY_SCOPE_KEY="<64 hex chars>"
export IDEMPOTENCY_ENCRYPTION_KEY="<64 hex chars>"

go run ./cmd/server
```

The application listens on the configured `APP_PORT` with a normalized `:port` format.

## API reference

The service currently exposes a small API surface.

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | No | Returns a JSON status object `{"status":"ok"}`. |
| `GET` | `/swagger/` | No | Serves the generated Swagger UI for the API. |
| `POST` | `/api/v1/workspace` | Yes (`Bearer <JWT>`) | Creates a workspace for the authenticated user. |

The workspace creation payload matches `CreateWorkspaceDTO` and supports the following fields:

- `name` (required)
- `type` (`advertiser` or `publisher`)
- `country` (required)
- `timezone` (required)
- `currency` (required)
- `industry` (optional)
- `website` (optional)
- `description` (optional)

The route is protected by both JWT middleware and idempotency middleware. The idempotency scope is keyed per authenticated user, which prevents duplicate create requests from being processed repeatedly.

## Database schema

The database migration file `internal/database/migrations/000001_create_workspace.up.sql` creates a single `workspaces` table with the following structure:

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `UUID` | Primary key, generated with `gen_random_uuid()` |
| `name` | `VARCHAR(150)` | Not null |
| `slug` | `VARCHAR(180)` | Unique, generated from workspace name |
| `type` | `VARCHAR(20)` | Restricted to `advertiser` or `publisher` |
| `country` | `CHAR(2)` | Not null |
| `timezone` | `VARCHAR(100)` | Not null |
| `currency` | `CHAR(3)` | Not null |
| `industry` | `VARCHAR(100)` | Optional |
| `website` | `TEXT` | Optional |
| `description` | `TEXT` | Optional |
| `status` | `VARCHAR(20)` | Defaults to `active`; allowed values are `active` and `deleted` |
| `created_by` | `UUID` | User who created the record |
| `created_at` | `TIMESTAMPTZ` | Defaults to `NOW()` |
| `updated_at` | `TIMESTAMPTZ` | Defaults to `NOW()` |

The migration also creates the following indexes:

- `idx_workspaces_created_by`
- `idx_workspaces_type`
- `idx_workspaces_status`

## Testing

The repository contains a small but meaningful set of Go tests. The following command runs the current suite:

```bash
go test ./...
```

The current tests verify:

- config loading and validation,
- basic app creation failure behavior,
- graceful nil handling in app shutdown,
- the health endpoint response,
- default and custom Swagger host configuration.

## Development workflow

The repository includes a local development setup via Air and Docker.

### Air

`.air.toml` configures a Go build process that:

- compiles the server binary into `./tmp/main`,
- runs `swag init` before rebuilds to refresh Swagger docs,
- launches the binary with Delve for remote debugging,
- watches project source for automatic rebuilds.

### Typical workflow

```bash
# run all tests
go test ./...

# regenerate Swagger docs
swag init -g cmd/server/main.go -o docs --parseInternal --parseDependency

# start the server with Air
air -c .air.toml
```

## Security considerations

The implementation includes a few explicit security controls, based on the repository code:

- JWT validation is performed using a configured public key (`JWT_PUBLIC_KEY`).
- Authentication is enforced on the workspace route with the JWT middleware.
- Idempotency middleware is applied to workspace creation to reduce duplicate processing.
- Database access uses PostgreSQL and stores records with constrained enum-like values.

This service is not documented as providing a full identity platform, RBAC model, or external secret manager. The repository does not include additional security infrastructure beyond the code above.

## Current scope and limitations

This repository is intentionally narrow in scope. As implemented today, it provides:

- one domain: workspace creation,
- one database table: `workspaces`,
- one primary route: `POST /api/v1/workspace`,
- health and Swagger endpoints,
- JWT and idempotency enforcement.

There are no repository files indicating broader multi-tenant systems, billing flows, user management, or deployment automation beyond the local Docker Compose setup.

## License

There is no license file in the repository root. The project currently does not declare a license, so no license information is documented here beyond that fact.
