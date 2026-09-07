-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name,url,user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
WHERE (
    sqlc.arg('search')::text = ''
    OR name ILIKE '%' || sqlc.arg('search') || '%'
    OR url ILIKE '%' || sqlc.arg('search') || '%'
  )
  AND (
    sqlc.narg('cursor_created_at')::timestamp IS NULL
    OR (created_at, id) < (
      sqlc.narg('cursor_created_at'),
      sqlc.narg('cursor_id')::uuid
    )
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit');


-- name: GetNextFeedsToFetch :many
SELECT * from feeds ORDER BY  last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;
