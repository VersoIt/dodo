-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category SMALLINT NOT NULL,
    base_price DECIMAL(12,2) NOT NULL,
    image_url TEXT,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ingredients (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cost DECIMAL(12,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    in_stock BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS product_ingredients (
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    ingredient_id UUID REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity DOUBLE PRECISION NOT NULL,
    is_removable BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (product_id, ingredient_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
