-- name: GetBrandByID :one
SELECT * FROM brands WHERE id = $1;

-- name: GetBrandByName :one
SELECT * FROM brands WHERE name = $1;

-- name: ListBrands :many
SELECT * FROM brands ORDER BY name;

-- name: CreateBrand :one
INSERT INTO brands (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateBrand :one
UPDATE brands
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = sqlc.narg('description'),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = $1;
