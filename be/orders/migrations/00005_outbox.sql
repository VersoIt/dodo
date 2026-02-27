-- +goose Up
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE outbox_events;
