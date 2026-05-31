-- +goose Up
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS error_msg TEXT;

-- +goose Down
ALTER TABLE notifications DROP COLUMN IF EXISTS error_msg;
