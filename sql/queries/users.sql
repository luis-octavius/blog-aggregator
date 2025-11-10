-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1, 
  $2, 
  $3, 
  $4
)
RETURNING *;

-- name: GetUser :one 
SELECT * FROM users 
WHERE name = $1 LIMIT 1;

-- name: DeleteUsers :exec 
DELETE FROM users;

-- name: GetUsers :many 
SELECT * FROM users;

-- name: CreateFeed :one 
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
  $1,
  $2, 
  $3,
  $4, 
  $5,
  $6
)
RETURNING *;
