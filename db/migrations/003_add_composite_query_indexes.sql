-- Composite indexes for filtered+sorted hot paths
-- Created: 2026-03-31

CREATE INDEX IF NOT EXISTS idx_orders_user_created_at
ON orders(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_products_category_name
ON products(category_id, name);

CREATE INDEX IF NOT EXISTS idx_products_brand_name
ON products(brand_id, name);

COMMENT ON INDEX idx_orders_user_created_at IS 'Optimizes WHERE user_id = ? ORDER BY created_at DESC';
COMMENT ON INDEX idx_products_category_name IS 'Optimizes WHERE category_id = ? ORDER BY name';
COMMENT ON INDEX idx_products_brand_name IS 'Optimizes WHERE brand_id = ? ORDER BY name';
