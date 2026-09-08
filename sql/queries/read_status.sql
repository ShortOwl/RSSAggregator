-- name: MarkPostAsRead :one
INSERT INTO read_posts (user_id, post_id)
SELECT sqlc.arg('user_id'), posts.id
FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE posts.id = sqlc.arg('post_id')
  AND feed_follows.user_id = sqlc.arg('user_id')
ON CONFLICT (user_id, post_id)
DO UPDATE SET post_id = EXCLUDED.post_id
RETURNING read_posts.*;

-- name: MarkPostAsUnread :exec
DELETE FROM read_posts
WHERE user_id = $1 AND post_id = $2;
