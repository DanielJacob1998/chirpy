-- +goose Up
ALTER TABLE users
ALTER COLUMN created_at SET DEFAULT NOW();

-- +goose Down
ALTER TABLE users
ALTER COLUMN created_at DROP DEFAULT;
