-- +goose Up
ALTER TABLE notifications ADD COLUMN status VARCHAR(20) DEFAULT 'queued';
ALTER TABLE notifications ADD COLUMN error_msg TEXT;

-- +goose Down
ALTER TABLE notifications DROP COLUMN error_msg;
ALTER TABLE notifications DROP COLUMN status;
