# Database Joins & Indexing - বিস্তারিত ব্যাখ্যা

আপনার Go E-Commerce API-তে কোথায় কোন **JOIN** হচ্ছে এবং কী কী **INDEX** আছে।

---

## 📊 Part 1: Joins (কোথায় কী Join হচ্ছে)

### 1. Product Queries (সবচেয়ে গুরুত্বপূর্ণ!)

#### Query 1: `GetProductWithDetails` - Single Product with Join

**Location:** `internal/database/sqlc/store.go:271-287`

**SQL Query:**
```sql
SELECT 
  -- Product fields
  p.id, p.name, p.description, p.price, p.stock, 
  p.image_url, p.category_id, p.brand_id, 
  p.created_at, p.updated_at,
  
  -- Category fields (joined)
  c.id as cat_id, 
  c.name as cat_name, 
  c.description as cat_description, 
  c.created_at as cat_created_at, 
  c.updated_at as cat_updated_at,
  
  -- Brand fields (joined)
  b.id as brand_id, 
  b.name as brand_name, 
  b.description as brand_description, 
  b.created_at as brand_created_at, 
  b.updated_at as brand_updated_at

FROM products p 
JOIN categories c ON p.category_id = c.id    -- ← JOIN 1
JOIN brands b ON p.brand_id = b.id           -- ← JOIN 2
WHERE p.id = $1
```

**Join Type:** INNER JOIN (দুইটা)

**কী হচ্ছে:**
1. `products` table থেকে data নিচ্ছে
2. `category_id` দিয়ে `categories` table-এর সাথে join করছে
3. `brand_id` দিয়ে `brands` table-এর সাথে join করছে

**Result:** একটা product + তার category + তার brand একসাথে আসছে

**Visual Representation:**

```
products table          categories table        brands table
┌─────────────┐        ┌──────────────┐        ┌──────────────┐
│ id          │        │ id           │        │ id           │
│ name        │  ┌────►│ name         │  ┌────►│ name         │
│ price       │  │     │ description  │  │     │ description  │
│ category_id │──┘     └──────────────┘  │     └──────────────┘
│ brand_id    │────────────────────────┘
└─────────────┘

Result: Product + Category + Brand (একসাথে!)
```

---

#### Query 2: `ListProducts` - Multiple Products with Join

**Location:** `internal/database/sqlc/store.go:290-311`

**SQL Query:**
```sql
SELECT 
  p.id, p.name, p.description, p.price, p.stock, 
  p.image_url, p.category_id, p.brand_id, 
  p.created_at, p.updated_at,
  
  c.id as cat_id, c.name as cat_name, c.description as cat_description,
  c.created_at as cat_created_at, c.updated_at as cat_updated_at,
  
  b.id as brand_id, b.name as brand_name, b.description as brand_description,
  b.created_at as brand_created_at, b.updated_at as brand_updated_at

FROM products p 
JOIN categories c ON p.category_id = c.id    -- ← JOIN 1
JOIN brands b ON p.brand_id = b.id           -- ← JOIN 2

WHERE 
  ($1::uuid IS NULL OR p.category_id = $1)   -- Optional filter
  AND ($2::uuid IS NULL OR p.brand_id = $2)  -- Optional filter
  
ORDER BY p.name
```

**Join Type:** INNER JOIN (দুইটা)

**Filter Support:**
- `?category=uuid` → শুধু ঐ category-র products
- `?brand=uuid` → শুধু ঐ brand-র products
- দুইটা একসাথে → specific category AND brand

**Example Response Structure:**

```json
[
  {
    "id": "product-uuid",
    "name": "Samsung Galaxy S21",
    "price": 899.99,
    "category": {                    // ← Join থেকে এসেছে
      "id": "electronics-uuid",
      "name": "Electronics"
    },
    "brand": {                       // ← Join থেকে এসেছে
      "id": "samsung-uuid",
      "name": "Samsung"
    }
  }
]
```

---

### 2. Order Queries (N+1 Pattern - No Direct Join)

**Location:** `internal/handlers/order_handler.go`

**Pattern:** Multiple separate queries (not JOINed in SQL)

```go
// 1. Get order
order, _ := h.store.GetOrderByID(r.Context(), id)

// 2. Get user separately
user, _ := h.store.GetUserByID(r.Context(), order.UserID)

// 3. Get order items separately
items, _ := h.store.ListOrderItemsByOrderID(r.Context(), order.ID)

// 4. Get product details for each item (loop করে)
for _, item := range items {
    prod, _ := h.store.GetProductWithDetails(r.Context(), item.ProductID)
    itemProducts[item.ProductID] = prod
}

// 5. Build complete response
resp := models.ToOrderResponse(order, user, items, itemProducts)
```

**কেন JOIN নেই?**
- Order response complex (nested items with products)
- Application-level joining করা easier to maintain
- Small dataset-এ performance impact minimal

**এটা কি bad practice?**
- ❌ No - small/medium projects-এর জন্য acceptable
- ✅ Readable and maintainable
- ⚠️ Large scale-এ optimize করা যেতে পারে (single query with JSON aggregation)

---

## 🔍 Part 2: Indexing (কী কী Index আছে)

### Automatic Indexes (PostgreSQL automatically creates)

#### 1. PRIMARY KEY Indexes

প্রতিটা table-এ `id UUID PRIMARY KEY` আছে → **Automatic B-tree index**

```sql
-- Automatically created indexes:
admins_pkey        ON admins(id)         -- ✅
users_pkey         ON users(id)          -- ✅
categories_pkey    ON categories(id)     -- ✅
brands_pkey        ON brands(id)         -- ✅
products_pkey      ON products(id)       -- ✅
orders_pkey        ON orders(id)         -- ✅
order_items_pkey   ON order_items(id)    -- ✅
```

**Type:** B-tree (default)  
**Used for:** WHERE id = $1 queries (very fast!)

---

#### 2. UNIQUE Constraint Indexes

```sql
-- Schema থেকে:
email TEXT NOT NULL UNIQUE     -- admins, users
name TEXT NOT NULL UNIQUE      -- categories, brands
```

**Automatically created indexes:**
```sql
admins_email_key       ON admins(email)      -- ✅
users_email_key        ON users(email)       -- ✅
categories_name_key    ON categories(name)   -- ✅
brands_name_key        ON brands(name)       -- ✅
```

**Type:** Unique B-tree  
**Used for:** 
- Login queries: `WHERE email = $1` (very fast!)
- Duplicate checking: `GetCategoryByName`, `GetBrandByName`

---

#### 3. FOREIGN KEY Indexes

**Schema:**
```sql
products:
  category_id UUID REFERENCES categories(id)   -- FK 1
  brand_id UUID REFERENCES brands(id)          -- FK 2

orders:
  user_id UUID REFERENCES users(id)            -- FK 3

order_items:
  order_id UUID REFERENCES orders(id)          -- FK 4
  product_id UUID REFERENCES products(id)      -- FK 5
```

**PostgreSQL automatically creates indexes on:**
```sql
-- Target columns (referenced side)
categories(id)    -- ✅ Via PRIMARY KEY
brands(id)        -- ✅ Via PRIMARY KEY
users(id)         -- ✅ Via PRIMARY KEY
orders(id)        -- ✅ Via PRIMARY KEY
products(id)      -- ✅ Via PRIMARY KEY
```

**Note:** PostgreSQL does NOT auto-index the **referencing side** (source columns)!

---

### ⚠️ Missing Indexes (Add করলে better হবে)

Current schema-তে এই columns-এ index নেই:

```sql
products.category_id   -- ⚠️ Used in JOINs and filters
products.brand_id      -- ⚠️ Used in JOINs and filters
orders.user_id         -- ⚠️ Used in ListOrdersByUserID
order_items.order_id   -- ⚠️ Used in ListOrderItemsByOrderID
order_items.product_id -- ⚠️ Used in lookups
```

**Performance Impact:**
- Small dataset (< 10,000 products): না থাকলেও problem নেই
- Large dataset (> 100,000 products): Sequential scan slow হবে

---

## 🎯 Part 3: Join Analysis (কোথায় কী হচ্ছে)

### Summary Table

| Query Function | Tables Involved | Join Type | Join Columns | Index Used? |
|----------------|-----------------|-----------|--------------|-------------|
| `GetProductWithDetails` | products + categories + brands | INNER JOIN x2 | category_id, brand_id | ⚠️ No FK index |
| `ListProducts` | products + categories + brands | INNER JOIN x2 | category_id, brand_id | ⚠️ No FK index |
| `ListOrdersByUserID` | orders (only) | No join | user_id | ⚠️ No index |
| `ListOrderItemsByOrderID` | order_items (only) | No join | order_id | ⚠️ No index |

---

## 📈 Performance Analysis

### Current Query Performance

#### Fast Queries (Index ব্যবহার করছে):

```sql
-- ✅ VERY FAST (uses PRIMARY KEY index)
SELECT * FROM products WHERE id = $1;

-- ✅ FAST (uses UNIQUE index on email)
SELECT * FROM users WHERE email = $1;

-- ✅ FAST (uses UNIQUE index on name)
SELECT * FROM categories WHERE name = $1;
```

#### Slow Queries (Index নেই):

```sql
-- ⚠️ SLOWER (no index on category_id)
SELECT * FROM products WHERE category_id = $1;

-- ⚠️ SLOWER (no index on brand_id)  
SELECT * FROM products WHERE brand_id = $1;

-- ⚠️ SLOWER (no index on user_id)
SELECT * FROM orders WHERE user_id = $1;

-- ⚠️ SLOW when products table is large (JOIN without FK index)
SELECT p.*, c.*, b.* 
FROM products p 
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id;
```

---

## 🛠️ Recommended Indexes (Add করুন performance-এর জন্য)

### Migration File তৈরি করুন

Create: `db/migrations/002_add_indexes.sql`

```sql
-- Performance indexes for foreign keys
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Composite index for common filter patterns
CREATE INDEX idx_products_category_brand ON products(category_id, brand_id);

-- Index for sorting by created_at (used in order queries)
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- Optional: Full-text search index (if you want product name search)
CREATE INDEX idx_products_name_trgm ON products USING gin(name gin_trgm_ops);
-- Requires: CREATE EXTENSION pg_trgm;
```

### Apply Migration

```bash
psql $DATABASE_URL -f db/migrations/002_add_indexes.sql
```

---

## 📖 Join Flow Visualization

### Product List Query Flow

```
Client Request: GET /api/v1/products
           │
           ▼
    ProductHandler.GetAll()
           │
           ▼
    store.ListProducts(ctx, nil, nil)
           │
           ▼
    SQL Query:
    ┌─────────────────────────────────────────┐
    │  SELECT p.*, c.*, b.*                   │
    │  FROM products p                        │
    │  JOIN categories c ON p.category_id = c.id  ← JOIN 1
    │  JOIN brands b ON p.brand_id = b.id         ← JOIN 2
    │  ORDER BY p.name                        │
    └─────────────────────────────────────────┘
           │
           ▼
    Database returns: ProductWithDetails[]
           │
           ▼
    models.ToProductResponseFromDetails()
           │
           ▼
    JSON Response:
    [
      {
        "id": "product-1",
        "name": "Samsung Galaxy S21",
        "price": 899.99,
        "category": {          ← Category data from JOIN
          "id": "cat-1",
          "name": "Electronics"
        },
        "brand": {             ← Brand data from JOIN
          "id": "brand-1",
          "name": "Samsung"
        }
      }
    ]
```

---

### Product Join Diagram

```
Request: GET /api/v1/products?category=electronics-uuid

Database Tables:
┌─────────────────┐
│   products      │
├─────────────────┤
│ id: prod-1      │
│ name: "Galaxy"  │◄─┐
│ price: 899.99   │  │
│ category_id: c1 │──┼─────► ┌──────────────┐
│ brand_id: b1    │──┼─────► │  categories  │
└─────────────────┘  │       ├──────────────┤
                     │       │ id: c1       │
                     │       │ name: "Elec" │
                     │       └──────────────┘
                     │
                     └─────► ┌──────────────┐
                             │   brands     │
                             ├──────────────┤
                             │ id: b1       │
                             │ name: "Sam"  │
                             └──────────────┘

Result: All three tables merged into one object!
```

---

## 🔍 Part 2: Index Analysis

### Current Indexes (Automatic)

#### Primary Key Indexes (B-tree)

```sql
Table: admins
Index: admins_pkey
Column: id
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast (O(log n))

Table: users
Index: users_pkey  
Column: id
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast

Table: categories
Index: categories_pkey
Column: id  
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast

Table: brands
Index: brands_pkey
Column: id
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast

Table: products
Index: products_pkey
Column: id
Type: B-tree (unique)  
Usage: WHERE id = $1  ✅ Very fast

Table: orders
Index: orders_pkey
Column: id
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast

Table: order_items
Index: order_items_pkey
Column: id
Type: B-tree (unique)
Usage: WHERE id = $1  ✅ Very fast
```

---

#### UNIQUE Constraint Indexes

```sql
Table: admins
Index: admins_email_key
Column: email
Type: B-tree (unique)
Usage: WHERE email = $1  ✅ Fast (for login)
Query: GetAdminByEmail()

Table: users  
Index: users_email_key
Column: email
Type: B-tree (unique)
Usage: WHERE email = $1  ✅ Fast (for login)
Query: GetUserByEmail()

Table: categories
Index: categories_name_key
Column: name
Type: B-tree (unique)
Usage: WHERE name = $1  ✅ Fast
Query: GetCategoryByName()

Table: brands
Index: brands_name_key  
Column: name
Type: B-tree (unique)
Usage: WHERE name = $1  ✅ Fast
Query: GetBrandByName()
```

---

### Missing Indexes (Should Add!)

#### Foreign Key Columns

```sql
-- ⚠️ NO INDEX (Sequential scan হবে!)
products.category_id   -- Used in: JOIN, WHERE category_id = $1
products.brand_id      -- Used in: JOIN, WHERE brand_id = $1
orders.user_id         -- Used in: WHERE user_id = $1
order_items.order_id   -- Used in: WHERE order_id = $1
order_items.product_id -- Used in: lookups
```

**Impact:**
- Small data: OK (< 1000 records)
- Medium data: Noticeable (1000-10000)
- Large data: Very slow (> 10000)

---

## 📊 Query Performance Comparison

### Example: Get Products by Category

#### Without Index (Current):

```sql
EXPLAIN ANALYZE 
SELECT * FROM products WHERE category_id = 'some-uuid';
```

**Result:**
```
Seq Scan on products  (cost=0.00..35.50 rows=10 width=...)
  Filter: (category_id = 'some-uuid')
  Rows Removed by Filter: 990
Planning Time: 0.123 ms
Execution Time: 2.456 ms   ← SLOW for 1000 rows
```

**কী হচ্ছে:** 
- পুরো table scan করছে (1000টা row)
- প্রতিটা row check করছে category_id match করে কিনা
- 990টা row discard করছে
- Only 10টা row return করছে

#### With Index (After adding):

```sql
CREATE INDEX idx_products_category_id ON products(category_id);

EXPLAIN ANALYZE 
SELECT * FROM products WHERE category_id = 'some-uuid';
```

**Result:**
```
Index Scan using idx_products_category_id  (cost=0.15..12.25 rows=10 width=...)
  Index Cond: (category_id = 'some-uuid')
Planning Time: 0.089 ms
Execution Time: 0.234 ms   ← 10x FASTER!
```

**কী হচ্ছে:**
- Directly index-এ গিয়ে matching rows খুঁজছে
- Only 10টা row read করছে
- 990টা unnecessary row skip করছে

---

## 🎯 Join Performance with Index

### Product JOIN Query Performance

#### Current (Without FK Index):

```sql
SELECT p.*, c.*, b.*
FROM products p 
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id;
```

**Execution Plan:**
```
Hash Join  (cost=...)
  -> Seq Scan on products     ← Full table scan
  -> Hash
    -> Seq Scan on categories ← Full table scan
  -> Hash Join
    -> Seq Scan on brands     ← Full table scan
```

**Time:** ~5-10ms (1000 products)

#### After Adding Indexes:

```sql
-- Add these:
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
```

**Execution Plan:**
```
Nested Loop  (cost=...)
  -> Seq Scan on products
  -> Index Scan on categories using categories_pkey   ← Uses PK index
  -> Index Scan on brands using brands_pkey           ← Uses PK index
```

**Time:** ~1-2ms (1000 products) - **5x faster!**

---

## 🛠️ How to Add Missing Indexes

### Create Migration File

**File:** `db/migrations/002_add_performance_indexes.sql`

```sql
-- =====================================================
-- Performance Indexes for Foreign Keys and Queries
-- =====================================================

-- Products foreign key indexes (for JOINs)
CREATE INDEX IF NOT EXISTS idx_products_category_id 
  ON products(category_id);

CREATE INDEX IF NOT EXISTS idx_products_brand_id 
  ON products(brand_id);

-- Composite index for common filter pattern
CREATE INDEX IF NOT EXISTS idx_products_category_brand 
  ON products(category_id, brand_id);

-- Orders foreign key index
CREATE INDEX IF NOT EXISTS idx_orders_user_id 
  ON orders(user_id);

-- Order items foreign keys
CREATE INDEX IF NOT EXISTS idx_order_items_order_id 
  ON order_items(order_id);

CREATE INDEX IF NOT EXISTS idx_order_items_product_id 
  ON order_items(product_id);

-- Sorting indexes (for ORDER BY queries)
CREATE INDEX IF NOT EXISTS idx_orders_created_at 
  ON orders(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_products_name 
  ON products(name);

-- Optional: Text search index (if you add search feature)
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX idx_products_name_search 
--   ON products USING gin(name gin_trgm_ops);

-- =====================================================
-- Index Statistics
-- =====================================================

COMMENT ON INDEX idx_products_category_id IS 'Speeds up product JOINs and category filters';
COMMENT ON INDEX idx_products_brand_id IS 'Speeds up product JOINs and brand filters';
COMMENT ON INDEX idx_orders_user_id IS 'Speeds up user order lookups';
```

### Apply Migration

```bash
psql $DATABASE_URL -f db/migrations/002_add_performance_indexes.sql
```

### Verify Indexes Created

```bash
psql $DATABASE_URL -c "\d products"
```

**Output দেখবেন:**
```
Indexes:
    "products_pkey" PRIMARY KEY, btree (id)
    "idx_products_category_id" btree (category_id)     ← New!
    "idx_products_brand_id" btree (brand_id)           ← New!
    "idx_products_category_brand" btree (category_id, brand_id) ← New!
```

---

## 🔬 Deep Dive: Join Execution

### Example: Get Product with Details

**SQL:**
```sql
SELECT p.*, c.*, b.*
FROM products p 
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id
WHERE p.id = 'some-uuid';
```

**Execution Steps:**

```
1. Index Scan on products (products_pkey)
   → Find product row by ID          [Fast: Uses PK index]
   → Returns: category_id = 'c1', brand_id = 'b1'

2. Index Scan on categories (categories_pkey)  
   → Find category row where id = 'c1'   [Fast: Uses PK index]
   
3. Index Scan on brands (brands_pkey)
   → Find brand row where id = 'b1'      [Fast: Uses PK index]
   
4. Merge results
   → Return single row with all columns
```

**Result:** Very fast! (< 1ms)

**কেন fast?** 
- WHERE p.id = $1 uses PRIMARY KEY index ✅
- JOIN conditions use PRIMARY KEY indexes ✅
- Only 3 index lookups (no table scans)

---

### Example: List Products by Category

**SQL:**
```sql
SELECT p.*, c.*, b.*
FROM products p 
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id
WHERE p.category_id = 'electronics-uuid'
ORDER BY p.name;
```

**Execution Steps (WITHOUT index on category_id):**

```
1. Seq Scan on products               [SLOW: Full table scan]
   → Read ALL 10,000 products
   → Filter: category_id = 'electronics-uuid'
   → Keep: 500 matching rows
   → Discard: 9,500 non-matching rows
   
2. For each of 500 products:
   → Index Scan on categories (PK)    [Fast x500]
   → Index Scan on brands (PK)        [Fast x500]
   
3. Sort by name                        [Medium: 500 rows]

Time: ~50-100ms for 10,000 products
```

**Execution Steps (WITH index on category_id):**

```
1. Index Scan on idx_products_category_id  [FAST!]
   → Directly find 500 matching products
   → Skip 9,500 non-matching rows
   
2. For each of 500 products:
   → Index Scan on categories (PK)    [Fast x500]
   → Index Scan on brands (PK)        [Fast x500]
   
3. Sort by name                        [Medium: 500 rows]

Time: ~5-10ms for 10,000 products  ← 10x FASTER!
```

---

## 🎯 Index Strategy Summary

### Already Have (Automatic):

| Index Type | Column(s) | Usage | Performance |
|------------|-----------|-------|-------------|
| PRIMARY KEY | id (all tables) | WHERE id = $1 | ✅ Excellent |
| UNIQUE | email (users, admins) | Login queries | ✅ Excellent |
| UNIQUE | name (categories, brands) | Duplicate check | ✅ Good |

### Should Add (Recommended):

| Index | Column(s) | Usage | Expected Improvement |
|-------|-----------|-------|----------------------|
| idx_products_category_id | category_id | JOINs, filters | 5-10x faster |
| idx_products_brand_id | brand_id | JOINs, filters | 5-10x faster |
| idx_orders_user_id | user_id | User orders | 10-20x faster |
| idx_order_items_order_id | order_id | Order details | 10x faster |

---

## 📈 When to Add Indexes

### Small Dataset (< 1,000 records)
- ❌ Not urgent
- Current performance: OK
- Index overhead: Not worth it

### Medium Dataset (1,000 - 50,000)
- ⚠️ Recommended
- Noticeable performance gain
- Small overhead

### Large Dataset (> 50,000)
- ✅ **MUST ADD**
- Major performance impact
- Essential for production

---

## 🎓 Index Best Practices

### When PostgreSQL Uses Indexes:

```sql
-- ✅ USES INDEX
WHERE category_id = $1
WHERE user_id = $1
WHERE email = $1

-- ❌ DOESN'T USE INDEX (function on column)
WHERE LOWER(email) = $1
WHERE created_at::date = $1

-- ✅ USES INDEX (function index needed)
CREATE INDEX idx_email_lower ON users(LOWER(email));
WHERE LOWER(email) = $1  -- Now uses index!
```

### Index Maintenance:

```sql
-- See index usage stats
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

-- Unused indexes (consider dropping)
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes  
WHERE schemaname = 'public' AND idx_scan = 0;
```

---

## 🚀 Performance Optimization Roadmap

### Phase 1: Essential Indexes (Do Now)
```sql
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

### Phase 2: Query Optimization (If Needed)
```sql
CREATE INDEX idx_products_category_brand ON products(category_id, brand_id);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
```

### Phase 3: Advanced Features (Optional)
```sql
-- Full-text search
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_products_name_search ON products USING gin(name gin_trgm_ops);

-- Partial index (only active products)
CREATE INDEX idx_products_in_stock ON products(id) WHERE stock > 0;
```

---

## 📝 Summary

### Joins Currently Used:

1. **Product queries** - 2 INNER JOINs
   - `products JOIN categories` (via category_id)
   - `products JOIN brands` (via brand_id)

2. **Order queries** - Application-level joins
   - Multiple separate queries
   - Combined in Go code

### Indexes Currently Have:

- ✅ **7 PRIMARY KEY indexes** (id columns)
- ✅ **4 UNIQUE indexes** (email, name columns)
- ⚠️ **0 FOREIGN KEY indexes** (missing!)

### Performance Impact:

- **Small dataset (current):** ✅ Good enough
- **Large dataset (future):** ⚠️ Need FK indexes

### Next Steps:

1. Add FK indexes (recommended)
2. Monitor query performance
3. Add composite indexes if needed

---

## 🔧 Quick Commands

```bash
# Check what indexes exist
psql $DATABASE_URL -c "\di"

# Check specific table indexes
psql $DATABASE_URL -c "\d products"

# See index usage statistics
psql $DATABASE_URL -c "SELECT tablename, indexname, idx_scan FROM pg_stat_user_indexes;"

# Add recommended indexes
psql $DATABASE_URL -f db/migrations/002_add_indexes.sql
```

---

**📖 Summary:**
- ✅ JOINs properly implemented in product queries
- ✅ Basic indexes (PK, UNIQUE) working well
- ⚠️ FK indexes missing (add করলে better performance)
- 🎯 Current setup good for small/medium data
