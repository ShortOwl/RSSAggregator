-- name: CreatePost :one

INSERT INTO posts(id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = sqlc.arg('user_id')
AND (sqlc.narg('feed_id')::uuid IS NULL or posts.feed_id = sqlc.narg('feed_id'))
AND (
  sqlc.arg('search')::text = ''
  OR
  to_tsvector('english',posts.title || ' ' || COALESCE(posts.description,''))
  @@ plainto_tsquery('english',sqlc.arg('search'))
)
AND (
  NOT sqlc.arg('unread_only')::boolean
  OR NOT EXISTS (
    SELECT 1 FROM read_posts
    WHERE read_posts.user_id = sqlc.arg('user_id')
      AND read_posts.post_id = posts.id
  )
)
AND (
  sqlc.narg('cursor_published_at')::timestamp IS NULL
  OR (posts.published_at,posts.id) < (sqlc.narg('cursor_published_at'),sqlc.narg('cursor_id')::uuid)
)
ORDER BY posts.published_at DESC, posts.id DESC
LIMIT sqlc.arg('limit');
