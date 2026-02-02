-- name: GetAdminByID :one
SELECT * FROM admins WHERE id = $1;

-- name: GetAdminByEmail :one
SELECT * FROM admins WHERE email = $1;

-- name: CreateAdmin :one
INSERT INTO admins (email, password, name)
VALUES ($1, $2, $3)
RETURNING *;
