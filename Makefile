.PHONY: help install migrate run build watch clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install Go dependencies
	go mod download
	go mod tidy

migrate: ## Apply SQL migrations (set DATABASE_URL or use .env)
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; \
	if [ -z "$$DATABASE_URL" ]; then echo "Set DATABASE_URL or create .env"; exit 1; fi; \
	psql "$$DATABASE_URL" -f db/migrations/001_schema.sql

run: ## Run the server
	go run cmd/server/main.go

AIR := $(shell go env GOPATH)/bin/air
watch: ## Run server with live reload (air). Install: go install github.com/air-verse/air@latest
	@test -f $(AIR) || (echo "Install air: go install github.com/air-verse/air@latest"; exit 1)
	$(AIR)

build: ## Build the application
	go build -o bin/server cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/
