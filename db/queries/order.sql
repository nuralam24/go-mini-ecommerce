-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrdersByUserID :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListOrdersAll :many
SELECT * FROM orders ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders SET status = $2, updated_at = now() WHERE id = $1 RETURNING *;
