-- +goose Up
-- +goose StatementBegin
INSERT INTO products (id, name, description, category, base_price, is_available) VALUES
('019c53e1-4f90-7373-8d00-6472a7dc310c', 'Margherita', 'Classic pizza with tomato sauce and mozzarella', 0, 10.99, TRUE),
('019c53e1-4f90-7373-8d00-6472a7dc310d', 'Pepperoni', 'Classic pizza with pepperoni and mozzarella', 0, 12.99, TRUE),
('019c53e1-4f90-7373-8d00-6472a7dc310e', 'Cola', 'Refreshing drink', 4, 2.50, TRUE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE id IN ('019c53e1-4f90-7373-8d00-6472a7dc310c', '019c53e1-4f90-7373-8d00-6472a7dc310d', '019c53e1-4f90-7373-8d00-6472a7dc310e');
-- +goose StatementEnd
