-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS kitchen_tickets (
    id UUID PRIMARY KEY,
    order_id UUID UNIQUE NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    start_cooking_time TIMESTAMPTZ,
    ready_time TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS kitchen_items (
    id BIGSERIAL PRIMARY KEY,
    ticket_id UUID REFERENCES kitchen_tickets(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    comment TEXT
);

CREATE TABLE IF NOT EXISTS kitchen_item_ingredients (
    id BIGSERIAL PRIMARY KEY,
    kitchen_item_id BIGINT REFERENCES kitchen_items(id) ON DELETE CASCADE,
    ingredient_name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kitchen_item_ingredients;
DROP TABLE IF EXISTS kitchen_items;
DROP TABLE IF EXISTS kitchen_tickets;
-- +goose StatementEnd
