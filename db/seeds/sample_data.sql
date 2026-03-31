-- Insert sample data for performance testing
-- This creates realistic test data to demonstrate index performance

-- Insert categories if not exist
INSERT INTO categories (name, description) 
VALUES 
    ('Electronics', 'Electronic devices and gadgets'),
    ('Clothing', 'Fashion and apparel'),
    ('Books', 'Books and literature')
ON CONFLICT (name) DO NOTHING;

-- Insert brands if not exist
INSERT INTO brands (name, description) 
VALUES 
    ('Samsung', 'Korean electronics brand'),
    ('Apple', 'American tech company'),
    ('Nike', 'Sportswear brand')
ON CONFLICT (name) DO NOTHING;

-- Insert 100 sample products
INSERT INTO products (name, description, price, stock, category_id, brand_id)
SELECT 
    'Product ' || i,
    'Description for product ' || i,
    (50 + (random() * 950))::numeric(10,2),
    (random() * 100)::int,
    (SELECT id FROM categories ORDER BY random() LIMIT 1),
    (SELECT id FROM brands ORDER BY random() LIMIT 1)
FROM generate_series(1, 100) AS i;

-- Print summary
SELECT 
    (SELECT COUNT(*) FROM categories) as categories,
    (SELECT COUNT(*) FROM brands) as brands,
    (SELECT COUNT(*) FROM products) as products;
