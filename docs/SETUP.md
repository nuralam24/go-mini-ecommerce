# Setup Guide

## 1. Go & PostgreSQL

- **Go**: 1.22 বা তার উপরের version install করুন
- **PostgreSQL**: 12+ (local অথবা Neon/cloud)

## 2. Environment

```bash
cp .env.example .env
```

`.env` এ set করুন:

- `DATABASE_URL`: PostgreSQL connection string (e.g. `postgresql://user:pass@localhost:5432/ecommerce?sslmode=disable`)
- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Strong secret for JWT signing

## 3. Dependencies

```bash
go mod download
go mod tidy
```

## 4. Database

Database create করুন (যদি নেই):

```bash
psql postgres -c "CREATE DATABASE ecommerce;"
```

Schema apply করুন:

```bash
psql "$DATABASE_URL" -f db/migrations/001_schema.sql
```

অথবা Makefile:

```bash
make migrate
```

## 5. Run

```bash
go run cmd/server/main.go
```

অথবা `./scripts/run.sh` অথবা `make run`।

## Optional: sqlc

Code generation from SQL queries (optional):

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```

Generated code goes to `internal/database/sqlc/`. The app works without this (hand-written store is included).

## Troubleshooting

- **Connection refused**: Ensure PostgreSQL is running and `DATABASE_URL` is correct.
- **Migration errors**: If tables already exist, use a fresh database or drop objects before re-running the migration.
- **Missing go.sum**: Run `go mod tidy` (with network).
