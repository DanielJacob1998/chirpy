-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE users
ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- +goose Down
-- Note: be cautious with downgrading extensions; usually, you don't need to drop them
ALTER TABLE users
ALTER COLUMN id DROP DEFAULT;
