-- +goose Up
-- +goose StatementBegin
ALTER TABLE deliveries ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE couriers ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers DROP COLUMN updated_at;
ALTER TABLE deliveries DROP COLUMN updated_at;
-- +goose StatementEnd
