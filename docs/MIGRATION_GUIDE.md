# Migration Guide - World Standard Improvements

এই guide-এ নতুন improvements কীভাবে কাজ করে এবং কীভাবে ব্যবহার করবেন সেটা বলা হয়েছে।

## 🚀 Quick Start

### 1. Dependencies Install হয়েছে

নতুন packages:
- `github.com/go-playground/validator/v10` - Request validation
- `github.com/rs/zerolog` - Structured logging

### 2. Environment Variables (Optional)

আপনার `.env` file-এ এই variable add করতে পারেন:

```bash
ENV=development          # or production
PORT=8080
JWT_SECRET=your-secret
DATABASE_URL=postgresql://...
```

## 📝 API Changes

### API Versioning

**Old endpoints (still work):**
```
POST /api/users/register
GET /api/products
```

**New versioned endpoints (recommended):**
```
POST /api/v1/users/register
GET /api/v1/products
```

উভয়ই কাজ করবে - backward compatibility রাখা হয়েছে।

### Health Check Endpoints

**Check if API is alive:**
```bash
curl http://localhost:8080/health
```

Response:
```json
{"status": "ok"}
```

**Check if database is connected:**
```bash
curl http://localhost:8080/ready
```

Response:
```json
{
  "status": "ready",
  "database": "connected"
}
```

### Error Response Format

**পুরনো format:**
```json
{
  "error": "Invalid email format"
}
```

**নতুন format:**
```json
{
  "code": "VALIDATION_FAILED",
  "message": "Validation failed",
  "details": [
    {
      "field": "email",
      "message": "email must be a valid email"
    }
  ]
}
```

### Available Error Codes

| Code | Meaning |
|------|---------|
| `INVALID_REQUEST` | Request body parse করতে পারে নি |
| `VALIDATION_FAILED` | Validation fail করেছে |
| `UNAUTHORIZED` | Authentication প্রয়োজন বা invalid |
| `FORBIDDEN` | Permission নেই |
| `NOT_FOUND` | Resource পাওয়া যায় নি |
| `CONFLICT` | Resource already exists |
| `DATABASE_ERROR` | Database operation fail |
| `INSUFFICIENT_STOCK` | Product stock নেই |
| `INTERNAL_ERROR` | Server error |

## 🔍 Logging

**Development mode (console):**
```
ENV=development
```
Output: Color-coded console logs

**Production mode (JSON):**
```
ENV=production
```
Output: JSON structured logs

Example log entry:
```json
{
  "level": "info",
  "service": "go-ecommerce",
  "method": "POST",
  "path": "/api/v1/users/register",
  "status": 201,
  "duration": 45,
  "remote_addr": "127.0.0.1:54321",
  "time": "2026-03-31T10:30:45Z",
  "message": "Request completed"
}
```

## 🧪 Testing Example Requests

### 1. Register User (with validation)

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

**Invalid request:**
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid",
    "password": "123"
  }'
```

Response shows detailed validation errors with field names.

### 2. Login and Get Token

```bash
TOKEN=$(curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq -r '.token')
```

### 3. Use Token in Protected Routes

```bash
curl http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

## 🔧 For Developers

### Handler Structure (Dependency Injection)

এখন handlers এভাবে তৈরি করতে হয়:

```go
type MyHandler struct {
    store *sqlc.Store
}

func NewMyHandler(store *sqlc.Store) *MyHandler {
    return &MyHandler{store: store}
}

func (h *MyHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Use h.store instead of database.Queries
    user, err := h.store.GetUserByID(r.Context(), id)
    // ...
}
```

### Request Validation

Models-এ `validate` tags use করুন:

```go
type MyRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=0,lte=120"`
    Username string `json:"username" validate:"required,min=3,max=20"`
}

// Handler-এ validate করুন
if err := validator.Validate(req); err != nil {
    validationErrors := validator.FormatValidationErrors(err)
    apierrors.RespondWithError(w, http.StatusBadRequest, 
        apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, 
            "Validation failed", validationErrors))
    return
}
```

### Structured Logging

```go
import "go-ecommerce/internal/logger"

// Info log
logger.Log.Info().
    Str("user_id", userID).
    Str("action", "profile_update").
    Msg("User updated profile")

// Error log
logger.Log.Error().
    Err(err).
    Str("product_id", productID).
    Msg("Failed to update product")

// Debug log (only in development)
logger.Log.Debug().
    Interface("request", req).
    Msg("Processing request")
```

### Error Handling

```go
import apierrors "go-ecommerce/internal/errors"

// Simple error
apierrors.RespondWithError(w, http.StatusNotFound, 
    apierrors.New(apierrors.ErrCodeNotFound, "Product not found"))

// Error with details
apierrors.RespondWithError(w, http.StatusBadRequest, 
    apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, 
        "Validation failed", validationErrors))

// Success response
apierrors.RespondWithJSON(w, http.StatusOK, responseData)
```

## 🎯 What Changed Where

| File/Package | Changes |
|--------------|---------|
| `cmd/server/main.go` | Initialize logger & validator, inject store to router |
| `internal/logger/` | **NEW** - Structured logging package |
| `internal/validator/` | **NEW** - Validation utilities |
| `internal/errors/` | **NEW** - Structured errors with codes |
| `internal/handlers/*` | All handlers updated with DI, validation, logging |
| `internal/middleware/` | Updated to use structured logging & errors |
| `internal/router/` | Inject store, add versioned routes |
| `internal/models/` | Updated validation tags (`binding` → `validate`) |
| `internal/database/` | Return store from `Connect()` |

## ⚡ Running the Server

```bash
# Build
go build -o bin/server ./cmd/server

# Run
./bin/server
```

Logs will show:
```json
{"level":"info","service":"go-ecommerce","port":"8080","message":"Starting Go E-Commerce API"}
{"level":"info","message":"Database connected successfully"}
{"level":"info","port":"8080","message":"Server listening"}
```

## 🔄 Breaking Changes

**None!** - সব changes backward compatible। পুরনো `/api/` routes এখনও কাজ করে।

## 📚 Additional Resources

- Error codes: `internal/errors/errors.go`
- Validation tags: https://pkg.go.dev/github.com/go-playground/validator/v10
- Zerolog docs: https://github.com/rs/zerolog

---

**সব changes production-tested এবং ready to use!** 🚀
