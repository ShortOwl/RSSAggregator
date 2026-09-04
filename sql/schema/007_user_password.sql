-- +goose Up
ALTER TABLE users
ADD COLUMN email TEXT UNIQUE NOT NULL,
ADD COLUMN password_hash TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN password_hash,
DROP COLUMN email;
