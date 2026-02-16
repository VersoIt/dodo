-- +goose Up
-- +goose StatementBegin
ALTER TABLE kitchen_tickets ADD COLUMN order_number VARCHAR(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kitchen_tickets DROP COLUMN order_number;
-- +goose StatementEnd
