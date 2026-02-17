-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    order_id UUID UNIQUE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    method SMALLINT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    transaction_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
