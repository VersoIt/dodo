-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS couriers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    current_lat DOUBLE PRECISION,
    current_lng DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_id UUID PRIMARY KEY,
    courier_id UUID REFERENCES couriers(id) ON DELETE SET NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pickup_time TIMESTAMPTZ,
    delivery_time TIMESTAMPTZ,
    current_lat DOUBLE PRECISION,
    current_lng DOUBLE PRECISION
);

CREATE INDEX IF NOT EXISTS idx_deliveries_courier_id ON deliveries(courier_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS couriers;
-- +goose StatementEnd
