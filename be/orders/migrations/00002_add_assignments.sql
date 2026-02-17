-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN IF NOT EXISTS chef_id UUID;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS courier_id UUID;

CREATE INDEX IF NOT EXISTS idx_orders_chef_id ON orders(chef_id);
CREATE INDEX IF NOT EXISTS idx_orders_courier_id ON orders(courier_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN IF EXISTS chef_id;
ALTER TABLE orders DROP COLUMN IF EXISTS courier_id;
-- +goose StatementEnd
