-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: GetOrderItemWithProduct :one
SELECT oi.*, p.id as product_id, p.name as product_name, p.description as product_description, p.price as product_price, p.stock as product_stock, p.image_url as product_image_url,
  p.category_id as product_category_id, p.brand_id as product_brand_id, p.created_at as product_created_at, p.updated_at as product_updated_at
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.id = $1;
