-- name: GetUserByID :one
SELECT
    id,
    username,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT
    id,
    username,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    password_hash
) VALUES (
    $1, $2
)
RETURNING 
    id,
    username,
    password_hash,
    created_at,
    updated_at;

-- name: UpdateUser :one
UPDATE users
SET 
    password_hash = $2,
    updated_at = now()
WHERE id = $1
RETURNING 
    id,
    username,
    password_hash,
    created_at,
    updated_at;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;