# File Guide & Request-Response Lifecycle

এই document এ project এর প্রতিটি file এর কাজ এবং Go request-response lifecycle explain করা হয়েছে।

## 📁 File Structure & Purpose

### Root Level Files

#### `go.mod`
**কাজ**: Go module definition file
- Project এর dependencies define করে
- Module name এবং Go version specify করে
- `go get` command দিয়ে dependencies install করা যায়

#### `go.sum`
**কাজ**: Dependencies এর checksums store করে
- Auto-generated file
- Security এর জন্য important
- Git এ commit করতে হবে

#### `.env.example`
**কাজ**: Environment variables এর template
- Database URL, PORT, JWT_SECRET এর example values
- Copy করে `.env` file তৈরি করতে হবে
- Git এ commit করা safe (sensitive data নেই)

#### `.env`
**কাজ**: Actual environment variables
- **Git এ commit করবেন না!** (.gitignore এ আছে)
- Database credentials, secrets store করে
- Runtime এ application এই values use করে

#### `.gitignore`
**কাজ**: Git ignore rules
- `.env`, `bin/`, generated files ignore করে
- Sensitive data protect করে

#### `Makefile`
**কাজ**: Build automation commands
- `make run`, `make build` এর মতো commands define করে
- Development workflow সহজ করে

#### `package.json`
**কাজ**: Optional project metadata

---

### `cmd/server/main.go`
**কাজ**: Application entry point
- Server start করে
- Database connect করে
- Router initialize করে
- Middleware chain setup করে
- Graceful shutdown handle করে

**Key Functions**:
- `main()`: Entry point
- Configuration load
- Database connection
- Server startup
- Signal handling (SIGINT, SIGTERM)

---

### `internal/config/config.go`
**কাজ**: Configuration management
- Environment variables load করে
- `.env` file read করে
- Configuration struct provide করে
- Default values set করে

**Key Functions**:
- `Load()`: Configuration load করে return করে
- `getEnv()`: Environment variable read করে

**Returns**:
```go
Config {
    DatabaseURL: "postgresql://..."
    Port: "8080"
    Env: "development"
    JWTSecret: "secret-key"
}
```

---

### `internal/database/database.go`
**কাজ**: Database connection management
- PostgreSQL connection open করে (pgx driver)
- Store (Queries) create করে
- Graceful disconnect handle করে

**Key Functions**:
- `Connect()`: Database connect করে, Queries set করে
- `Disconnect()`: Database disconnect করে

**Global Variables**:
- `DB`: *sql.DB
- `Queries`: sqlc Store (সব handlers use করে)

---

### `internal/models/models.go`
**কাজ**: Data models এবং DTOs (Data Transfer Objects)
- Request models (API থেকে data receive করার জন্য)
- Response models (API থেকে data send করার জন্য)
- Conversion functions (sqlc models → Response models)

**Key Types**:
- `UserResponse`, `CreateUserRequest`, `LoginRequest`
- `CategoryResponse`, `CreateCategoryRequest`
- `ProductResponse`, `CreateProductRequest`
- `OrderResponse`, `CreateOrderRequest`
- ইত্যাদি...

**Key Functions**:
- `ToUserResponse()`: sqlc User → UserResponse
- `ToCategoryResponse()`: sqlc Category → CategoryResponse
- `ToProductResponse()` / `ToProductResponseFromDetails()`: sqlc Product → ProductResponse
- `ToOrderResponse()`: sqlc Order + user + items → OrderResponse

---

### `internal/utils/` Directory

#### `jwt.go`
**কাজ**: JWT token operations
- Token generate করে
- Token validate করে
- Claims extract করে

**Key Functions**:
- `InitJWT(secret)`: JWT secret initialize করে
- `GenerateToken(userID, email, role)`: JWT token create করে
- `ValidateToken(token)`: Token verify করে, claims return করে

#### `password.go`
**কাজ**: Password hashing operations
- Password hash করে (bcrypt)
- Password verify করে

**Key Functions**:
- `HashPassword(password)`: Password hash করে
- `CheckPasswordHash(password, hash)`: Password match করে কিনা check করে

#### `json.go`
**কাজ**: JSON response helpers
- JSON response send করে
- JSON decode করে
- Error response send করে

**Key Functions**:
- `RespondWithJSON()`: JSON response send করে
- `RespondWithError()`: Error response send করে
- `DecodeJSON()`: Request body decode করে

---

### `internal/middleware/` Directory

#### `cors.go`
**কাজ**: CORS (Cross-Origin Resource Sharing) handling
- Browser থেকে cross-origin requests allow করে
- CORS headers set করে
- OPTIONS request handle করে

**Key Function**:
- `CORS(next http.Handler)`: CORS middleware

#### `logging.go`
**কাজ**: Request logging
- প্রতিটি request log করে
- Method, path, status code, duration log করে
- Debugging এর জন্য helpful

**Key Function**:
- `Logging(next http.Handler)`: Logging middleware

#### `auth.go`
**কাজ**: JWT authentication
- Authorization header check করে
- JWT token validate করে
- User info context এ add করে
- Admin role check করে

**Key Functions**:
- `AuthMiddleware()`: JWT validation middleware
- `AdminMiddleware()`: Admin role check middleware
- `GetUserID()`, `GetUserRole()`: Context থেকে user info extract করে

---

### `internal/router/router.go`
**কাজ**: Route registration
- সব routes define করে
- Middleware apply করে
- Public vs Protected routes separate করে
- Admin-only routes protect করে

**Key Functions**:
- `NewRouter()`: Router instance create করে
- `RegisterRoutes()`: সব routes register করে
- `ServeHTTP()`: HTTP request handle করে

**Route Types**:
- **Public Routes**: No authentication required
- **Protected Routes**: Authentication required
- **Admin Routes**: Admin role required

---

### `internal/handlers/` Directory

#### `user_handler.go`
**কাজ**: User-related endpoints
- User registration
- User login
- Profile get/update

**Endpoints**:
- `POST /api/users/register` - Public
- `POST /api/users/login` - Public
- `GET /api/users/profile` - Protected
- `PUT /api/users/profile` - Protected

#### `admin_handler.go`
**কাজ**: Admin authentication
- Admin registration
- Admin login

**Endpoints**:
- `POST /api/admin/register` - Public
- `POST /api/admin/login` - Public

#### `category_handler.go`
**কাজ**: Category CRUD operations
- Category create, read, update, delete

**Endpoints**:
- `GET /api/categories` - Public
- `GET /api/categories/{id}` - Public
- `POST /api/categories` - Admin only
- `PUT /api/categories/{id}` - Admin only
- `DELETE /api/categories/{id}` - Admin only

#### `brand_handler.go`
**কাজ**: Brand CRUD operations
- Brand create, read, update, delete

**Endpoints**:
- `GET /api/brands` - Public
- `GET /api/brands/{id}` - Public
- `POST /api/brands` - Admin only
- `PUT /api/brands/{id}` - Admin only
- `DELETE /api/brands/{id}` - Admin only

#### `product_handler.go`
**কাজ**: Product CRUD operations
- Product create, read, update, delete
- Product filtering (by category, brand)

**Endpoints**:
- `GET /api/products` - Public (supports ?category=id&brand=id)
- `GET /api/products/{id}` - Public
- `POST /api/products` - Admin only
- `PUT /api/products/{id}` - Admin only
- `DELETE /api/products/{id}` - Admin only

#### `order_handler.go`
**কাজ**: Order management
- Order create
- Order list (user sees own, admin sees all)
- Order status update (admin only)

**Endpoints**:
- `POST /api/orders` - User only
- `GET /api/orders` - Protected
- `GET /api/orders/{id}` - Protected
- `PUT /api/orders/{id}/status` - Admin only

---

### `db/migrations/001_schema.sql`
**কাজ**: Database schema definition
- Tables ও enum define করে (admins, users, categories, brands, products, orders, order_items)
- Optional: sqlc generate দিয়ে code generate করা যায়

**Models**:
- `Admin`, `User`, `Category`, `Brand`, `Product`, `Order`, `OrderItem`
- `OrderStatus` enum

---

## 🔄 Go Request-Response Lifecycle

### Complete Flow Diagram

```
1. HTTP Request arrives
   ↓
2. Server receives request (main.go)
   ↓
3. CORS Middleware (cors.go)
   - Adds CORS headers
   - Handles OPTIONS requests
   ↓
4. Logging Middleware (logging.go)
   - Logs request start
   ↓
5. Router (router.go)
   - Matches URL pattern
   - Finds handler
   ↓
6. Auth Middleware (auth.go) [if protected route]
   - Extracts Authorization header
   - Validates JWT token
   - Adds user info to context
   ↓
7. Admin Middleware (auth.go) [if admin route]
   - Checks user role == "admin"
   ↓
8. Handler Function (handlers/*.go)
   - Decodes request body (json.go)
   - Validates input
   - Calls database operations (database.go)
   - Processes business logic
   - Converts sqlc models to DTOs (models.go)
   ↓
9. Database Operations (database.go + sqlc Store)
   - Executes queries
   - Returns data
   ↓
10. Handler prepares response
    - Creates response DTO (models.go)
    - Sets status code
    ↓
11. JSON Response (json.go)
    - Encodes response to JSON
    - Sends HTTP response
    ↓
12. Logging Middleware (logging.go)
    - Logs response (status, duration)
    ↓
13. CORS Middleware (cors.go)
    - Adds CORS headers to response
    ↓
14. Response sent to client
```

---

## 📝 Detailed Lifecycle Example

### Example: Creating a Product (Admin)

#### Step 1: Request Arrives
```http
POST /api/products HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Samsung Galaxy S21",
  "price": 899.99,
  "stock": 50,
  "category_id": "uuid-here",
  "brand_id": "uuid-here"
}
```

#### Step 2: Server Receives (main.go)
```go
// Server is listening on :8080
srv.ListenAndServe()
// Request arrives at server
```

#### Step 3: CORS Middleware (cors.go)
```go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // ... more headers
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r) // Pass to next middleware
    })
}
```

#### Step 4: Logging Middleware (logging.go)
```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // ... logging setup
        
        next.ServeHTTP(w, r) // Pass to router
        
        duration := time.Since(start)
        log.Printf("%s %s %d %v", r.Method, r.RequestURI, status, duration)
    })
}
```

#### Step 5: Router (router.go)
```go
// Route matching
r.mux.HandleFunc("POST /api/products", 
    authMiddleware(
        adminMiddleware(
            http.HandlerFunc(productHandler.Create)
        )
    ).ServeHTTP
)
```

#### Step 6: Auth Middleware (auth.go)
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from header
        authHeader := r.Header.Get("Authorization")
        token := strings.Split(authHeader, " ")[1]
        
        // Validate token
        claims, err := utils.ValidateToken(token)
        if err != nil {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
            return
        }
        
        // Add to context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
        
        next.ServeHTTP(w, r.WithContext(ctx)) // Pass to admin middleware
    })
}
```

#### Step 7: Admin Middleware (auth.go)
```go
func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role := r.Context().Value(UserRoleKey).(string)
        
        if role != "admin" {
            utils.RespondWithError(w, http.StatusForbidden, "Admin access required")
            return
        }
        
        next.ServeHTTP(w, r) // Pass to handler
    })
}
```

#### Step 8: Handler Function (product_handler.go)
```go
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    // 8.1: Decode request body
    var req models.CreateProductRequest
    utils.DecodeJSON(r, &req) // json.go
    
    // 8.2: Validate input
    if req.Name == "" {
        utils.RespondWithError(w, http.StatusBadRequest, "Name required")
        return
    }
    
    // 8.3: Verify category exists
    category, err := database.Client.Category.FindUnique(
        db.Category.ID.Equals(req.CategoryID),
    ).Exec(r.Context())
    
    // 8.4: Verify brand exists
    brand, err := database.Client.Brand.FindUnique(
        db.Brand.ID.Equals(req.BrandID),
    ).Exec(r.Context())
    
    // 8.5: Create product
    product, err := database.Client.Product.CreateOne(
        db.Product.Name.Set(req.Name),
        db.Product.Price.Set(req.Price),
        // ... more fields
    ).Exec(r.Context())
    
    // 8.6: Convert to response model
    response := models.ToProductResponse(product) // models.go
    
    // 8.7: Send response
    utils.RespondWithJSON(w, http.StatusCreated, response) // json.go
}
```

#### Step 9: Database Operations (database.go + sqlc Store)
```go
// Store executes SQL query
// Returns sqlc model
product, _ := database.Queries.CreateProduct(ctx, name, desc, price, stock, imageURL, categoryID, brandID)
// product is *sqlc.Product
```

#### Step 10-11: Response Preparation & Sending (json.go)
```go
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(payload)
    // Response sent!
}
```

#### Step 12: Response Sent
```http
HTTP/1.1 201 Created
Content-Type: application/json
Access-Control-Allow-Origin: *

{
  "id": "uuid",
  "name": "Samsung Galaxy S21",
  "price": 899.99,
  "stock": 50,
  "created_at": "2024-01-01T00:00:00Z",
  ...
}
```

---

## 🔑 Key Concepts

### 1. Middleware Chain
Middleware গুলো chain আকারে execute হয়:
```
Request → CORS → Logging → Auth → Admin → Handler → Response
```

### 2. Context
Context দিয়ে request এর সাথে data pass করা যায়:
- User ID
- User Role
- Request metadata

### 3. Handler Pattern
প্রতিটি handler:
- `http.HandlerFunc` type
- `(w http.ResponseWriter, r *http.Request)` signature
- Request process করে, response send করে

### 4. Error Handling
Consistent error responses:
```go
utils.RespondWithError(w, statusCode, "Error message")
```

### 5. Response Models
sqlc models → DTOs conversion:
```go
user, _ := database.Queries.GetUserByID(ctx, userID)
response := models.ToUserResponse(user)
```

---

## 🎯 Summary

### File Responsibilities

| File | Responsibility |
|------|---------------|
| `main.go` | Server startup, initialization |
| `config.go` | Configuration management |
| `database.go` | Database connection |
| `models.go` | Data structures, DTOs |
| `handlers/*.go` | Business logic, request handling |
| `middleware/*.go` | Cross-cutting concerns |
| `router.go` | Route registration |
| `utils/*.go` | Helper functions |
| `db/migrations/*.sql` | Database schema |

### Request Flow Summary

1. **Receive** → Server receives HTTP request
2. **Middleware** → CORS, Logging, Auth, Admin checks
3. **Route** → Router matches URL to handler
4. **Handler** → Business logic execution
5. **Database** → Data operations via sqlc Store (Queries)
6. **Response** → JSON response sent back
7. **Log** → Request logged

এই structure follow করে code maintainable, testable, এবং scalable হয়! 🚀
