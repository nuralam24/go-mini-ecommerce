# কিভাবে Code Explore করবেন এবং Run করবেন

এই guide এ step-by-step code explore এবং run করার instructions আছে।

## 📁 Code Structure Overview

```
go-ecommerce/
├── cmd/server/main.go          # 🚀 Entry point - এখান থেকে সব শুরু
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection
│   ├── handlers/               # API endpoints handlers
│   ├── middleware/             # Auth, CORS, Logging
│   ├── models/                 # Data models
│   ├── router/                 # Route definitions
│   └── utils/                  # Helper functions
└── db/
    └── migrations/             # SQL schema
```

---

## 🔍 Step 1: Code Explore করুন

### 1.1 Entry Point দেখুন

**File**: `cmd/server/main.go`

```go
func main() {
    // 1. Configuration load
    cfg := config.Load()
    
    // 2. JWT initialize
    utils.InitJWT(cfg.JWTSecret)
    
    // 3. Database connect
    database.Connect()
    
    // 4. Router setup
    r := router.NewRouter()
    r.RegisterRoutes()
    
    // 5. Middleware apply
    handler := middleware.CORS(middleware.Logging(r))
    
    // 6. Server start
    srv := &http.Server{...}
    srv.ListenAndServe()
}
```

**কি করছে:**
- Configuration load করছে
- Database connect করছে
- Routes register করছে
- Server start করছে

### 1.2 Configuration দেখুন

**File**: `internal/config/config.go`

```go
type Config struct {
    DatabaseURL string  // Neon database URL
    Port        string  // Server port (8080)
    Env         string  // development/production
    JWTSecret   string  // JWT token secret
}
```

**কি করছে:**
- `.env` file থেকে variables read করছে
- Default values set করছে

### 1.3 Database Connection দেখুন

**File**: `internal/database/database.go`

```go
var Queries *sqlc.Store  // Database store (sqlc-style)

func Connect() error {
    DB, _ = sql.Open("pgx", connStr)
    Queries = sqlc.NewStore(DB)
    // PostgreSQL এর সাথে connect
}
```

**কি করছে:**
- PostgreSQL connection open করছে (pgx driver)
- Store (Queries) create করছে যেটা handlers use করে

### 1.4 Routes দেখুন

**File**: `internal/router/router.go`

```go
func RegisterRoutes() {
    // Public routes
    r.mux.HandleFunc("POST /api/users/register", ...)
    r.mux.HandleFunc("POST /api/users/login", ...)
    
    // Protected routes
    r.mux.HandleFunc("GET /api/users/profile", authMiddleware(...))
    
    // Admin routes
    r.mux.HandleFunc("POST /api/categories", authMiddleware(adminMiddleware(...)))
}
```

**Available Routes:**
- `POST /api/users/register` - User registration
- `POST /api/users/login` - User login
- `POST /api/admin/register` - Admin registration
- `POST /api/admin/login` - Admin login
- `GET /api/categories` - Get all categories
- `GET /api/products` - Get all products
- `POST /api/orders` - Create order (protected)
- ইত্যাদি...

### 1.5 Handlers দেখুন

**File**: `internal/handlers/user_handler.go`

```go
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Decode request
    var req models.CreateUserRequest
    utils.DecodeJSON(r, &req)
    
    // 2. Hash password
    hashedPassword := utils.HashPassword(req.Password)
    
    // 3. Create user in database
    user := database.Client.User.CreateOne(...)
    
    // 4. Send response
    utils.RespondWithJSON(w, 201, user)
}
```

**কি করছে:**
- Request body decode করছে
- Password hash করছে
- Database এ save করছে
- Response send করছে

---

## 🚀 Step 2: Project Run করুন

### 2.1 Prerequisites Check

```bash
# Go version check (1.22+ required)
go version

# Output দেখুন:
# go version go1.25.4 darwin/amd64 ✅
```

### 2.2 Environment Setup

**`.env` file check করুন:**

```bash
cat .env
```

**Expected content:**
```env
DATABASE_URL="postgresql://neondb_owner:password@ep-aged-surf-a19a8yxg-pooler.ap-southeast-1.aws.neon.tech/go_commerce?sslmode=require&channel_binding=require"
PORT=8080
ENV=development
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 2.3 Dependencies Install

```bash
# Go dependencies download
go mod download

# Verify
go mod verify
```

**Output:**
```
all modules verified ✅
```

### 2.4 Database Schema Apply

```bash
# SQL migration run করুন (একবার)
psql "$DATABASE_URL" -f db/migrations/001_schema.sql
# অথবা: make migrate
```

**কি করছে:**
- Tables create করছে (admins, users, categories, brands, products, orders, order_items)
- Enum ও triggers create করছে

### 2.5 Server Run

**Method 1: Direct Run**
```bash
go run cmd/server/main.go
```

**Method 2: Using Makefile**
```bash
make run
```

**Expected output:**
```
2024/01/22 17:30:00 Database connected successfully
2024/01/22 17:30:00 Server starting on port 8080
2024/01/22 17:30:00 Swagger UI available at http://localhost:8080/swagger/index.html
```

### 2.6 Verify Server Running

**Browser এ open করুন:**
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Base**: http://localhost:8080

**Terminal এ test করুন:**
```bash
# Health check (if you have a health endpoint)
curl http://localhost:8080/api/categories

# Expected: [] (empty array, no categories yet)
```

---

## 🧪 Step 3: API Test করুন

### 3.1 Admin Account Create করুন

```bash
curl -X POST http://localhost:8080/api/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123",
    "name": "Admin User"
  }'
```

**Expected Response:**
```json
{
  "message": "Admin created successfully"
}
```

### 3.2 Admin Login করুন

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123"
  }'
```

**Expected Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@test.com",
    "name": "Admin User"
  }
}
```

**Token save করুন** - পরবর্তী requests এ use করবেন!

### 3.3 Category Create করুন (Admin)

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Electronics",
    "description": "Electronic products"
  }'
```

### 3.4 Products দেখুন

```bash
# Get all products (public endpoint)
curl http://localhost:8080/api/products
```

---

## 🔧 Step 4: Development Workflow

### 4.1 Code Change করার পর

```bash
# 1. Code edit করুন
# 2. Server restart করুন (Ctrl+C, then run again)
# 3. Test করুন
```

### 4.2 Database Schema Change করার পর

```bash
# 1. db/migrations/ এ নতুন .sql file add করুন অথবা existing edit করুন
# 2. Migration apply করুন
psql "$DATABASE_URL" -f db/migrations/002_add_field.sql
# 3. Optional: sqlc generate (যদি sqlc use করেন)
sqlc generate
# 4. Server restart করুন
```

### 4.3 Common Commands

```bash
# Dependencies update
go mod tidy

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server

# Clean build artifacts
make clean

# Full setup
make setup
```

---

## 🐛 Troubleshooting

### Problem 1: Database Connection Failed

**Error:**
```
Failed to connect to database: connection refused
```

**Solution:**
```bash
# 1. .env file check করুন
cat .env

# 2. DATABASE_URL verify করুন
# 3. Neon database running আছে কিনা check করুন
# 4. Connection string এ password correct আছে কিনা check করুন
```

### Problem 2: Database / Import Error

**Error:** `missing go.sum` or package not found

**Solution:**
```bash
go mod tidy
go build ./...
```

### Problem 3: Port Already in Use

**Error:**
```
bind: address already in use
```

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
PORT=8081
```

### Problem 4: Migration Failed

**Error:** Schema already exists or syntax error

**Solution:**
```bash
# Fresh DB: drop and recreate database, then:
psql "$DATABASE_URL" -f db/migrations/001_schema.sql
```

---

## 📊 Code Flow Summary

```
1. main() function starts
   ↓
2. config.Load() - .env file read
   ↓
3. database.Connect() - Neon database connect
   ↓
4. router.RegisterRoutes() - Routes setup
   ↓
5. middleware.CORS() + Logging() - Middleware chain
   ↓
6. http.Server.ListenAndServe() - Server starts
   ↓
7. Request arrives
   ↓
8. Middleware → Router → Handler → Database → Response
```

---

## 🎯 Quick Start Checklist

- [ ] Go installed (1.22+)
- [ ] `.env` file created with DATABASE_URL, JWT_SECRET
- [ ] `go mod tidy` completed
- [ ] Schema applied: `psql "$DATABASE_URL" -f db/migrations/001_schema.sql`
- [ ] Server running on port 8080
- [ ] Swagger UI accessible
- [ ] Admin account created
- [ ] API endpoints tested

---

## 📚 Next Steps

1. **Swagger UI Explore করুন**: http://localhost:8080/swagger/index.html
2. **API Endpoints Test করুন**: Swagger UI থেকে directly test করতে পারেন
3. **Code Explore করুন**: Handlers, models, middleware files দেখুন
4. **Features Add করুন**: নতুন endpoints add করুন

---

**Happy Coding! 🚀**
