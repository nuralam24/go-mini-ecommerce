-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductWithDetails :one
SELECT p.*,
  c.id as cat_id, c.name as cat_name, c.description as cat_description, c.created_at as cat_created_at, c.updated_at as cat_updated_at,
  b.id as brand_id, b.name as brand_name, b.description as brand_description, b.created_at as brand_created_at, b.updated_at as brand_updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id
WHERE p.id = $1;

-- name: ListProducts :many
SELECT p.*,
  c.id as cat_id, c.name as cat_name, c.description as cat_description, c.created_at as cat_created_at, c.updated_at as cat_updated_at,
  b.id as brand_id, b.name as brand_name, b.description as brand_description, b.created_at as brand_created_at, b.updated_at as brand_updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
JOIN brands b ON p.brand_id = b.id
WHERE (sqlc.narg('category_id')::uuid IS NULL OR p.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('brand_id')::uuid IS NULL OR p.brand_id = sqlc.narg('brand_id'))
ORDER BY p.name;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, image_url, category_id, brand_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = sqlc.narg('description'),
  price = COALESCE(sqlc.narg('price'), price),
  stock = COALESCE(sqlc.narg('stock'), stock),
  image_url = sqlc.narg('image_url'),
  category_id = COALESCE(sqlc.narg('category_id'), category_id),
  brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products SET stock = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
