-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    delivery_city VARCHAR(100),
    delivery_street VARCHAR(255),
    delivery_house VARCHAR(20),
    delivery_apartment VARCHAR(20),
    delivery_floor VARCHAR(10),
    delivery_comment TEXT,
    
    delivery_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    promo_code VARCHAR(50),
    final_price DECIMAL(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    base_price DECIMAL(12,2) NOT NULL,
    size_multiplier DOUBLE PRECISION NOT NULL DEFAULT 1.0
);

CREATE TABLE IF NOT EXISTS order_item_toppings (
    id BIGSERIAL PRIMARY KEY,
    order_item_id UUID REFERENCES order_items(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(12,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_item_toppings;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
