-- +goose Up
CREATE INDEX posts_search_idx ON posts USING GIN (
    to_tsvector('english', title || ' ' || COALESCE(description, ''))
);

-- +goose Down
DROP INDEX posts_search_idx;
