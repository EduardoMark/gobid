-- name: CreateUser :one
INSERT INTO users (
  username, email,
  password_hash, bio
) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateUser :exec
UPDATE users
SET username = $2,
    email = $3,
    password_hash = $4,
    updated_at = now()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;