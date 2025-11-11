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
INSERT INTO feeds (name, url, user_id, created_at, updated_at)
VALUES (
  $1,
  $2, 
  $3,
  $4, 
  $5
)
RETURNING *;

-- name: GetFeeds :many 
SELECT feeds.name, feeds.url, users.name 
FROM feeds 
INNER JOIN users 
ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one 
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
  VALUES (
    $1,
    $2,
    $3,
    $4
  ) RETURNING *
)
  SELECT 
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
  FROM inserted_feed_follow 
  INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id 
  INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedByUrl :one 
SELECT * FROM feeds 
WHERE url = $1 LIMIT 1;
