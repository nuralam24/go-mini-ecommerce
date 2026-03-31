-- =====================================================
-- Performance Indexes for Foreign Keys and Queries
-- Created: 2026-03-31
-- Purpose: Optimize JOINs and WHERE clauses on foreign keys
-- =====================================================

-- Products foreign key indexes (speeds up JOINs and filters)
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_brand_id ON products(brand_id);

-- Composite index for common filter pattern (category + brand together)
CREATE INDEX IF NOT EXISTS idx_products_category_brand ON products(category_id, brand_id);

-- Orders foreign key index (speeds up user order lookups)
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

-- Order items foreign keys (speeds up order details and product lookups)
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

-- Sorting indexes (for ORDER BY queries)
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);

-- =====================================================
-- Index Comments (for documentation)
-- =====================================================

COMMENT ON INDEX idx_products_category_id IS 'Speeds up product-category JOINs and WHERE category_id = $1 queries';
COMMENT ON INDEX idx_products_brand_id IS 'Speeds up product-brand JOINs and WHERE brand_id = $1 queries';
COMMENT ON INDEX idx_products_category_brand IS 'Optimizes queries filtering by both category AND brand';
COMMENT ON INDEX idx_orders_user_id IS 'Speeds up user order history queries';
COMMENT ON INDEX idx_order_items_order_id IS 'Speeds up order details fetching';
COMMENT ON INDEX idx_order_items_product_id IS 'Speeds up product sales tracking';
COMMENT ON INDEX idx_orders_created_at IS 'Optimizes ORDER BY created_at DESC (recent orders first)';
COMMENT ON INDEX idx_products_name IS 'Optimizes ORDER BY name queries';

-- =====================================================
-- Verification
-- =====================================================
-- After running this migration, verify with:
--   psql $DATABASE_URL -c "\d products"
--   psql $DATABASE_URL -c "SELECT tablename, indexname FROM pg_indexes WHERE schemaname = 'public' ORDER BY tablename, indexname;"
