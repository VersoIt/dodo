-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    channel VARCHAR(20) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMPTZ,
    is_sent BOOLEAN NOT NULL DEFAULT FALSE,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
