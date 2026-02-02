#!/bin/bash
# Project run script - dependencies, migrations info, এবং server চালায়
# Run: ./scripts/run.sh

set -e
cd "$(dirname "$0")/.."

echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo "📤 Database: Apply migrations manually if needed (e.g. psql -f db/migrations/001_schema.sql)"
if [ -f .env ]; then set -a; source .env 2>/dev/null; set +a; fi

echo "🚀 Starting server..."
echo "   Swagger: http://localhost:8080/swagger/index.html"
go run cmd/server/main.go
