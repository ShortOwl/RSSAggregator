-- +goose Up
CREATE INDEX IF NOT EXISTS posts_published_at_idx ON posts (published_at);

-- +goose Down
DROP INDEX IF EXISTS posts_published_at_idx;
