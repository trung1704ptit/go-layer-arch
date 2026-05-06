# Build Golang RESTful API with Gorm, Gin and Postgres

This project is a layered REST API built with Gin, GORM, and PostgreSQL.

## Setup

### 1) Prerequisites

- Go `1.26.2` or newer
- PostgreSQL running locally or remotely
- Git (optional, but recommended)

### 2) Clone and install dependencies

```bash
git clone <your-repo-url>
cd go-layer
go mod tidy
```

### 3) Configure environment variables

The app reads config from `app.env`.

```bash
cp example.env app.env
```

Update `app.env` values for your environment:

- `POSTGRES_HOST`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_PORT`
- `PORT`
- `CLIENT_ORIGIN`
- `ACCESS_TOKEN_PRIVATE_KEY`
- `ACCESS_TOKEN_PUBLIC_KEY`
- `ACCESS_TOKEN_EXPIRED_IN`
- `ACCESS_TOKEN_MAXAGE`
- `REFRESH_TOKEN_PRIVATE_KEY`
- `REFRESH_TOKEN_PUBLIC_KEY`
- `REFRESH_TOKEN_EXPIRED_IN`
- `REFRESH_TOKEN_MAXAGE`

### 4) Run the API

```bash
go run ./cmd/main.go
```

Server starts on:

- `http://localhost:<PORT>`
- Health check: `GET /api/healthchecker`

### 5) Run with Air (hot reload)

Install Air:

```bash
go install github.com/air-verse/air@latest
```

Run API with auto-reload:

```bash
air
```

This project includes a preconfigured `.air.toml` file:

- builds binary to `tmp/main`
- watches Go and env/template files
- rebuilds and restarts on file changes

### 6) Run SQL migrations

This project uses versioned SQL migrations in `migrations/` with `golang-migrate`.

Run all pending migrations:

```bash
make migrate-up
```

Rollback all applied migrations:

```bash
make migrate-down
```

Create a new migration file pair (sequential: `00001_`, `00002_`, ...):

```bash
make migrate-create name=add_user_phone
```
This creates:

- `migrations/0000X_<name>.up.sql`
- `migrations/0000X_<name>.down.sql`

Then edit both generated files and add your SQL.

## Project Structure

```text
go-layer/
├── cmd/              # App entrypoints
├── internal/         # Application-private layers (controllers/services/repos/...)
├── migrate/          # Migration runner entrypoint
├── migrations/       # Versioned SQL schema migrations (*.up.sql/*.down.sql)
├── cmd/main.go       # Application bootstrap and dependency wiring
├── pkg/shared/       # Shared response helpers
├── pkg/utils/        # Shared utility helpers
├── app.env           # Runtime environment variables
└── example.env       # Environment template
```

## Layering Convention

- `routes` call `controllers`
- `controllers` call `services`
- `services` call `repositories` and utilities
- `repositories` handle database queries
- Keep business logic in `services`, not in controllers
