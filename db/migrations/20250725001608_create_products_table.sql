-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    sku VARCHAR(255) UNIQUE NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

INSERT INTO
  products (id, name, sku, price, stock_quantity)
VALUES
  (
    '00000000-0000-0000-0000-000000000001',
    'Café Especial Grão Mágico 250g',
    'CF-GM-250G',
    55.00,
    100
  ),
  (
    '00000000-0000-0000-0000-000000000002',
    'Blend da Casa Intenso 500g',
    'CF-BCI-500G',
    89.90,
    100
  );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE products;

-- +goose StatementEnd