# Go E-Commerce API

একটি mini e-commerce API যা Go, PostgreSQL, এবং sqlc ব্যবহার করে তৈরি করা হয়েছে। এটি world standard best practices অনুসরণ করে এবং clean code architecture implement করে।

## 🚀 Features

- **User Management**: User registration, login, profile management
- **Admin Management**: Admin registration and authentication
- **Category Management**: CRUD operations for product categories
- **Brand Management**: CRUD operations for product brands
- **Product Management**: Complete product catalog with filtering
- **Order Management**: Order creation, tracking, and status updates
- **JWT Authentication**: Secure token-based authentication
- **Swagger Documentation**: Complete API documentation
- **Clean Architecture**: Well-organized code structure following best practices

## 📋 Prerequisites

নিচের software গুলো install করা থাকতে হবে:

- **Go** 1.22 বা তার উপরের version
- **PostgreSQL** 12 বা তার উপরের version
- **Make** (optional, but recommended)

## 🛠️ Installation & Setup

### 1. Repository Clone করুন

```bash
git clone <repository-url>
cd go-ecommerce
```

### 2. Environment Variables Setup করুন

`.env.example` file থেকে `.env` file তৈরি করুন:

```bash
cp .env.example .env
```

`.env` file এ আপনার database credentials এবং JWT secret set করুন:

**Neon Database (Recommended for Cloud):**
```env
DATABASE_URL="postgresql://neondb_owner:YOUR_PASSWORD@ep-aged-surf-a19a8yxg-pooler.ap-southeast-1.aws.neon.tech/go_commerce?sslmode=require&channel_binding=require"
PORT=8080
ENV=development
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

**Local PostgreSQL:**
```env
DATABASE_URL="postgresql://username:password@localhost:5432/ecommerce?sslmode=disable"
PORT=8080
ENV=development
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

**Note**: Neon connection string এ `sslmode=require` এবং `channel_binding=require` parameters প্রয়োজনীয়।

### 3. Dependencies Install করুন

```bash
go mod download
go mod tidy
```

### 4. Database Create ও Migration Run করুন

প্রথমে PostgreSQL database create করুন:

```bash
psql postgres
CREATE DATABASE ecommerce;
\q
```

তারপর schema apply করুন (SQL migration):

```bash
psql -d ecommerce -f db/migrations/001_schema.sql
```

অথবা আপনার connection string দিয়ে:

```bash
psql "$DATABASE_URL" -f db/migrations/001_schema.sql
```

### 5. Server Run করুন

```bash
go run cmd/server/main.go
```

অথবা script ব্যবহার করুন:

```bash
./scripts/run.sh
```

অথবা Makefile:

```bash
make run
```

Server `http://localhost:8080` এ start হবে।

## 📚 API Documentation

Server run করার পর Swagger documentation এ access করতে পারেন:

**Swagger UI**: http://localhost:8080/swagger/index.html

## 🏗️ Project Structure

```
go-ecommerce/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── db/
│   ├── migrations/          # SQL schema migrations
│   └── queries/             # sqlc query files (optional: run sqlc generate)
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and sqlc store
│   │   └── sqlc/            # sqlc models and store (hand-written or generated)
│   ├── handlers/            # HTTP handlers for all endpoints
│   ├── middleware/          # HTTP middleware (auth, CORS, logging)
│   ├── models/              # Data models and DTOs
│   ├── router/              # Route registration
│   └── utils/               # Utility functions (JWT, password hashing)
├── sqlc.yaml                # sqlc config (optional: for sqlc generate)
├── .env.example             # Environment variables template
├── go.mod
├── Makefile
└── README.md
```

## 🔐 Authentication

API endpoints দুই ধরনের:

### Public Endpoints
- User registration and login
- Admin registration and login
- Get categories, brands, products (read-only)

### Protected Endpoints
Protected endpoints এর জন্য `Authorization` header এ Bearer token পাঠাতে হবে:

```
Authorization: Bearer <your-jwt-token>
```

### Role-Based Access
- **User Role**: Can create orders, view own orders, update profile
- **Admin Role**: Can manage categories, brands, products, and update order status

## 📝 API Endpoints

### User Endpoints
- `POST /api/users/register` - Register new user
- `POST /api/users/login` - User login
- `GET /api/users/profile` - Get user profile (Protected)
- `PUT /api/users/profile` - Update user profile (Protected)

### Admin Endpoints
- `POST /api/admin/register` - Register new admin
- `POST /api/admin/login` - Admin login

### Category Endpoints
- `GET /api/categories` - Get all categories
- `GET /api/categories/{id}` - Get category by ID
- `POST /api/categories` - Create category (Admin only)
- `PUT /api/categories/{id}` - Update category (Admin only)
- `DELETE /api/categories/{id}` - Delete category (Admin only)

### Brand Endpoints
- `GET /api/brands` - Get all brands
- `GET /api/brands/{id}` - Get brand by ID
- `POST /api/brands` - Create brand (Admin only)
- `PUT /api/brands/{id}` - Update brand (Admin only)
- `DELETE /api/brands/{id}` - Delete brand (Admin only)

### Product Endpoints
- `GET /api/products` - Get all products (supports ?category=id&brand=id filters)
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Create product (Admin only)
- `PUT /api/products/{id}` - Update product (Admin only)
- `DELETE /api/products/{id}` - Delete product (Admin only)

### Order Endpoints
- `POST /api/orders` - Create new order (User only)
- `GET /api/orders` - Get all orders (User sees own, Admin sees all)
- `GET /api/orders/{id}` - Get order by ID
- `PUT /api/orders/{id}/status` - Update order status (Admin only)

## 🧪 Testing the API

### 1. Create an Admin

```bash
curl -X POST http://localhost:8080/api/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123",
    "name": "Admin User"
  }'
```

### 2. Login as Admin

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

Response থেকে `token` save করুন।

### 3. Create a Category

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-admin-token>" \
  -d '{
    "name": "Electronics",
    "description": "Electronic products"
  }'
```

### 4. Create a Brand

```bash
curl -X POST http://localhost:8080/api/brands \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-admin-token>" \
  -d '{
    "name": "Samsung",
    "description": "Samsung brand"
  }'
```

### 5. Create a Product

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-admin-token>" \
  -d '{
    "name": "Samsung Galaxy S21",
    "description": "Latest smartphone",
    "price": 899.99,
    "stock": 50,
    "category_id": "<category-id>",
    "brand_id": "<brand-id>"
  }'
```

### 6. Register a User

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "user123",
    "name": "Test User"
  }'
```

### 7. Create an Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-user-token>" \
  -d '{
    "items": [
      {
        "product_id": "<product-id>",
        "quantity": 2
      }
    ]
  }'
```

## 🛠️ Development Commands

Makefile commands:

```bash
make help        # Show all available commands
make install     # Install Go dependencies
make migrate     # Apply SQL migrations (see Makefile)
make run         # Run the server
make build       # Build the application
make clean       # Clean build artifacts
```

Optional: sqlc code generation (if sqlc installed):

```bash
sqlc generate   # Regenerate internal/database/sqlc from db/queries
```

## 🔧 Configuration

Configuration file: `internal/config/config.go`

Environment variables:
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development/production)
- `JWT_SECRET`: Secret key for JWT token signing

## 📦 Dependencies

Main dependencies:
- **github.com/golang-jwt/jwt/v5**: JWT token handling
- **github.com/jackc/pgx/v5**: PostgreSQL driver
- **github.com/joho/godotenv**: Environment variable management
- **github.com/swaggo/http-swagger**: Swagger documentation
- **golang.org/x/crypto**: Password hashing (bcrypt)

## 🏛️ Architecture

এই project clean architecture principles follow করে:

1. **Separation of Concerns**: Handlers, models, database logic আলাদা
2. **Dependency Injection**: Database store এবং utilities properly injected
3. **Middleware Pattern**: Authentication, CORS, logging middleware
4. **Error Handling**: Consistent error responses
5. **Type Safety**: Strong typing throughout
6. **Code Organization**: Logical package structure

## 🔒 Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- Role-based access control (RBAC)
- CORS middleware for cross-origin requests
- Input validation

## 📄 License

MIT License

## 🤝 Contributing

Contributions welcome! Please follow the existing code style and architecture patterns.

## 📞 Support

For issues and questions, please open an issue in the repository.

## 📚 Additional Documentation

- **docs/SETUP.md**: Detailed setup and troubleshooting
- **docs/QUICKSTART.md**: Quick start guide
- **docs/PROJECT_STRUCTURE.md**: Architecture and project structure
- **docs/FILE_GUIDE.md**: File-by-file guide and request-response lifecycle
- **docs/GO_LIFECYCLE.md**: Go lifecycle (compilation, runtime, goroutines)

---

**Happy Coding! 🚀**
