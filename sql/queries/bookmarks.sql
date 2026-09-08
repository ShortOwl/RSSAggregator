-- name: CreateBookmark :one
INSERT INTO bookmarks (user_id, post_id)
SELECT
  sqlc.arg('user_id'),
  posts.id
FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE posts.id = sqlc.arg('post_id')
  AND feed_follows.user_id = sqlc.arg('user_id')
ON CONFLICT (user_id, post_id)
DO UPDATE SET post_id = EXCLUDED.post_id
RETURNING bookmarks.*;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND post_id = $2;

-- name: GetBookmarkedPosts :many
SELECT posts.* FROM posts
JOIN bookmarks ON posts.id = bookmarks.post_id
WHERE bookmarks.user_id = sqlc.arg('user_id')
AND (
  sqlc.narg('cursor_published_at')::timestamp IS NULL
  OR (posts.published_at, posts.id) < (
    sqlc.narg('cursor_published_at'),
    sqlc.narg('cursor_id')::uuid
  )
)
ORDER BY posts.published_at DESC, posts.id DESC
LIMIT sqlc.arg('limit');
