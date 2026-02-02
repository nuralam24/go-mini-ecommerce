# Quick Start

## Prerequisites

- Go 1.22+
- PostgreSQL 12+
- (Optional) sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Setup

```bash
# Clone & enter project
cd go-ecommerce

# Copy env
cp .env.example .env
# Edit .env: DATABASE_URL, JWT_SECRET

# Dependencies
go mod tidy

# Create DB (if needed)
psql postgres -c "CREATE DATABASE ecommerce;"

# Apply schema
psql "$DATABASE_URL" -f db/migrations/001_schema.sql

# Run server
go run cmd/server/main.go
```

## Optional: sqlc

If you use sqlc to regenerate code from `db/queries`:

```bash
sqlc generate
```

Output goes to `internal/database/sqlc/`. The project currently uses hand-written store in that package; you can replace it with sqlc-generated code if desired.

## API

- Swagger: http://localhost:8080/swagger/index.html
- Base URL: http://localhost:8080
