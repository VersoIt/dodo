-- +goose Up
CREATE TYPE role_type AS ENUM ('client', 'courier', 'support', 'system');

CREATE TABLE chat_messages (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    role role_type NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_read BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_chat_messages_order_id ON chat_messages(order_id);
CREATE INDEX idx_chat_messages_order_created ON chat_messages(order_id, created_at DESC);

-- +goose Down
DROP TABLE chat_messages;
DROP TYPE role_type;
