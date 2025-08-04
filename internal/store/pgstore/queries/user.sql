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

-- name: CheckEmailExistsExcludingID :one
SELECT EXISTS(
  SELECT 1
  FROM users
  WHERE email = $1 AND id != $2
);

-- name: UpdateUser :one
UPDATE users
SET username = $2,
    email = $3,
    bio = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, username, email, bio, created_at, updated_at;

-- name: ChangePassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;