-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategoryByName :one
SELECT * FROM categories WHERE name = $1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = sqlc.narg('description'),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
