# Project Structure

```
go-ecommerce/
├── cmd/server/
│   └── main.go              # Entry point, server start
├── db/
│   ├── migrations/          # SQL schema (001_schema.sql)
│   └── queries/             # sqlc query files (optional)
├── internal/
│   ├── config/              # Config loading (env)
│   ├── database/            # DB connection, sqlc store
│   │   └── sqlc/            # Models + Store (hand-written or sqlc generated)
│   ├── handlers/            # HTTP handlers (admin, user, category, brand, product, order)
│   ├── middleware/          # Auth, CORS, logging
│   ├── models/              # DTOs and request/response types
│   ├── router/              # Route registration
│   └── utils/               # JSON, JWT, password
├── docs/                    # Documentation
├── scripts/
│   └── run.sh               # Run script
├── sqlc.yaml                # sqlc config (optional)
├── go.mod
├── Makefile
└── README.md
```

## Data flow

- **main.go** → loads config, connects DB, starts router
- **Router** → registers routes, applies middleware
- **Handlers** → use `database.Queries` (sqlc Store) for DB access
- **Models** → DTOs and conversion from sqlc types to API responses

## Database layer

- **database/database.go**: Opens PostgreSQL with pgx driver, creates `Queries` (Store).
- **database/sqlc/**: Store with methods (GetAdminByEmail, CreateUser, ListProducts, etc.) and models (Admin, User, Category, Brand, Product, Order, OrderItem). Can be replaced by `sqlc generate` output if you use sqlc.

## Key files

- `db/migrations/001_schema.sql`: Full schema (admins, users, categories, brands, products, orders, order_items).
- `internal/database/sqlc/store.go`: All DB operations used by handlers.
- `internal/models/models.go`: API request/response types and To*Response helpers.
