-- +goose Up
-- +goose StatementBegin
ALTER TABLE deliveries ADD COLUMN order_number VARCHAR(50);
ALTER TABLE deliveries ADD COLUMN city VARCHAR(255);
ALTER TABLE deliveries ADD COLUMN street VARCHAR(255);
ALTER TABLE deliveries ADD COLUMN house VARCHAR(50);
ALTER TABLE deliveries ADD COLUMN apartment VARCHAR(50);

CREATE TABLE IF NOT EXISTS delivery_items (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID REFERENCES deliveries(order_id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS delivery_items;
ALTER TABLE deliveries DROP COLUMN order_number;
ALTER TABLE deliveries DROP COLUMN city;
ALTER TABLE deliveries DROP COLUMN street;
ALTER TABLE deliveries DROP COLUMN house;
ALTER TABLE deliveries DROP COLUMN apartment;
-- +goose StatementEnd
