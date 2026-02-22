-- +goose Up
ALTER TABLE chat_messages ADD COLUMN sender_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE chat_messages DROP COLUMN sender_name;
