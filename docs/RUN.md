# Run Guide

## Quick Run

1. Environment: `.env` এ `DATABASE_URL`, `PORT`, `JWT_SECRET` set করুন
2. Database: schema apply করুন (একবার)
   ```bash
   psql "$DATABASE_URL" -f db/migrations/001_schema.sql
   ```
3. Server চালান:
   ```bash
   go run cmd/server/main.go
   ```
   অথবা `./scripts/run.sh` অথবা `make run`

## Step-by-Step

### 1. Dependencies

```bash
go mod download
go mod tidy
```

### 2. Database schema apply

```bash
# .env load করুন (optional)
export $(grep -v '^#' .env | xargs)
# Migration run করুন
psql "$DATABASE_URL" -f db/migrations/001_schema.sql
```

অথবা Makefile:

```bash
make migrate
```

### 3. Server start

```bash
go run cmd/server/main.go
```

Swagger: http://localhost:8080/swagger/index.html

## Troubleshooting

- **Database connection error**: Check `DATABASE_URL` in `.env`
- **Schema already exists**: If tables exist, migration may error; use a fresh DB or drop tables first
- **Port in use**: Change `PORT` in `.env` (default 8080)
