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
go run main.go
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

## Project Structure

```text
go-layer/
├── controllers/      # HTTP handlers (request/response only)
├── services/         # Business logic layer
├── routes/           # Route registration and middleware wiring
├── middleware/       # Cross-cutting HTTP middleware (auth, etc.)
├── models/           # GORM models and request/response structs
├── initializers/     # Config loading and DB connection bootstrap
├── utils/            # Shared helpers (tokens, password hashing)
├── migrate/          # Database migration entrypoint
├── main.go           # Application bootstrap and dependency wiring
├── app.env           # Runtime environment variables
└── example.env       # Environment template
```

## Layering Convention

- `routes` call `controllers`
- `controllers` call `services`
- `services` work with database/models/utilities
- Keep business logic in `services`, not in controllers
