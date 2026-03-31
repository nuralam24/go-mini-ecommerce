# 🚀 Run Guide - Step by Step

আপনার world-standard Go E-Commerce API run করার complete guide।

---

## ⚡ Quick Start (3 Steps)

### Step 1: Database Setup

আপনার কাছে দুইটা option আছে:

#### Option A: Local PostgreSQL (Recommended for Development)

```bash
# PostgreSQL install (if not installed)
# macOS:
brew install postgresql@15
brew services start postgresql@15

# Create database
createdb go_commerce

# Update .env file
DATABASE_URL='postgresql://postgres@localhost:5432/go_commerce?sslmode=disable'
```

#### Option B: Cloud Database (Neon/Supabase/etc)

আপনার Neon database quota শেষ হয়ে গেছে। নতুন database create করুন অথবা upgrade করুন।

```bash
# .env-এ নতুন DATABASE_URL set করুন
DATABASE_URL='postgresql://user:pass@host:5432/dbname?sslmode=require'
```

### Step 2: Run Migrations

```bash
make migrate
```

**Expected output:**
```
CREATE TABLE
CREATE TYPE
CREATE TABLE
...
Migration completed ✓
```

### Step 3: Run Server

```bash
make run
```

**Expected output:**
```json
{"level":"info","service":"go-ecommerce","port":"8080","message":"Starting Go E-Commerce API"}
{"level":"info","message":"Database connected successfully"}
{"level":"info","port":"8080","message":"Server listening"}
{"level":"info","swagger_url":"http://localhost:8080/swagger/index.html","message":"Swagger UI available"}
```

✅ **Server running on http://localhost:8080**

---

## 🔧 Detailed Setup (First Time)

### 1. Prerequisites Check

```bash
# Go version (need 1.21+)
go version

# PostgreSQL check
psql --version

# Make check
make help
```

### 2. Install Dependencies

```bash
make install
```

### 3. Environment Variables

Check your `.env` file:

```bash
cat .env
```

Should have:
```bash
DATABASE_URL='postgresql://...'
PORT=8080
ENV=development
JWT_SECRET='your-secret-key'
```

### 4. Database Setup

#### Local PostgreSQL:

```bash
# Start PostgreSQL
brew services start postgresql@15

# Create database
createdb go_commerce

# Update .env
DATABASE_URL='postgresql://postgres@localhost:5432/go_commerce?sslmode=disable'

# Run migrations
make migrate
```

#### Docker PostgreSQL (Alternative):

```bash
# Run PostgreSQL in Docker
docker run --name postgres-ecommerce \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=go_commerce \
  -p 5432:5432 \
  -d postgres:15

# Update .env
DATABASE_URL='postgresql://postgres:postgres@localhost:5432/go_commerce?sslmode=disable'

# Run migrations
make migrate
```

### 5. Run Server

```bash
make run
```

---

## 🧪 Testing the API

### Test 1: Health Checks

```bash
# Liveness
curl http://localhost:8080/health

# Readiness (checks DB)
curl http://localhost:8080/ready
```

### Test 2: Register User (New!)

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123",
    "name": "John Doe"
  }'
```

**Success Response:**
```json
{
  "id": "...",
  "email": "john@example.com",
  "name": "John Doe",
  "created_at": "2026-03-31T10:00:00Z",
  "updated_at": "2026-03-31T10:00:00Z"
}
```

### Test 3: Login

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Success Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "...",
    "email": "john@example.com",
    "name": "John Doe"
  }
}
```

### Test 4: Get Profile (Protected Route)

```bash
# Save token from login
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Test 5: Validation Error (See New Error Format!)

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "123"
  }'
```

**New Structured Error Response:**
```json
{
  "code": "VALIDATION_FAILED",
  "message": "Validation failed",
  "details": [
    {
      "field": "email",
      "message": "email must be a valid email"
    },
    {
      "field": "password",
      "message": "password must be at least 6 characters"
    },
    {
      "field": "name",
      "message": "name is required"
    }
  ]
}
```

---

## 📱 Swagger UI (Best for Testing)

1. Server run করুন: `make run`
2. Browser-এ open করুন: http://localhost:8080/swagger/index.html
3. All endpoints test করতে পারবেন interactive UI দিয়ে!

**Swagger UI Features:**
- ✅ Try out all endpoints
- ✅ See request/response schemas
- ✅ Authentication with JWT token
- ✅ See validation rules

---

## 🎯 Common Use Cases

### Development (কোড লিখছেন):

```bash
# Terminal 1: Live reload
make watch

# কোড change করলে auto-restart হবে
```

### Testing (API test করছেন):

```bash
# Terminal 1: Run server
make run

# Terminal 2: Use curl or Swagger
curl http://localhost:8080/api/v1/...
# OR
open http://localhost:8080/swagger/index.html
```

### Production (Deploy করছেন):

```bash
# Build
make build

# Run binary
ENV=production \
DATABASE_URL=$PROD_DB_URL \
JWT_SECRET=$PROD_SECRET \
./bin/server
```

---

## 🐳 Docker Compose (Easiest for Local Development)

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_commerce
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/go_commerce?sslmode=disable
      ENV: development
      JWT_SECRET: dev-secret
      PORT: 8080
    depends_on:
      - postgres

volumes:
  postgres_data:
```

**Run everything:**
```bash
docker-compose up
```

---

## 🎨 Viewing Logs

### Development Mode:
Beautiful color-coded console output!

```
2026-03-31T16:00:06+06:00 INF Starting Go E-Commerce API port=8080
2026-03-31T16:00:07+06:00 INF Database connected successfully
2026-03-31T16:00:07+06:00 INF Server listening port=8080
```

### Production Mode:
JSON structured logs for log aggregation:

```json
{"level":"info","service":"go-ecommerce","port":"8080","time":"2026-03-31T10:00:06Z","message":"Starting Go E-Commerce API"}
{"level":"info","message":"Database connected successfully","time":"2026-03-31T10:00:07Z"}
```

---

## 🛑 Stopping the Server

```bash
# Press Ctrl+C in terminal
# Server will gracefully shutdown

# Expected output:
{"level":"info","message":"Shutting down server..."}
{"level":"info","message":"Server exited gracefully"}
```

---

## ✅ Quick Checklist

Before running:
- [ ] PostgreSQL running
- [ ] `.env` file configured
- [ ] Dependencies installed (`make install`)
- [ ] Migrations applied (`make migrate`)

To run:
- [ ] `make run` (simple)
- [ ] `make watch` (development)
- [ ] `make build && ./bin/server` (production)

To test:
- [ ] `curl http://localhost:8080/health`
- [ ] Open http://localhost:8080/swagger/index.html
- [ ] Try API endpoints

---

## 🎉 All Commands Summary

```bash
# Setup (one time)
make install                 # Install dependencies
make migrate                 # Setup database

# Development
make run                     # Run server
make watch                   # Run with live reload
make build                   # Build binary

# Testing
curl localhost:8080/health   # Quick test
# OR
open localhost:8080/swagger/index.html  # Interactive testing

# Cleanup
make clean                   # Remove build files
```

---

**সবচেয়ে সহজ:** 

```bash
make run
```

Then open: **http://localhost:8080/swagger/index.html** 🎉

Server running হলে console-এ beautiful structured logs দেখবেন! ✨
