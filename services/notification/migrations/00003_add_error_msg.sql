-- +goose Up
ALTER TABLE notifications ADD COLUMN error_msg TEXT;

-- +goose Down
ALTER TABLE notifications DROP COLUMN error_msg;
