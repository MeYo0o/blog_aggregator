-- name: CreateFeed :one
INSERT INTO feeds(id, name, url, user_id, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetFeeds :many
SELECT *
FROM feeds;
-- name: GetFeedByID :one
SELECT *
FROM feeds
WHERE id = $1;
-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;