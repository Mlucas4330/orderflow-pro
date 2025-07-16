-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    customer_id UUID NOT NULL,
    status TEXT NOT NULL CHECK (
      status IN (
        'pending',
        'paid',
        'shipped',
        'delivered',
        'cancelled',
        'refunded'
      )
    ),
    total NUMERIC(10, 2) NOT NULL CHECK (total >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'BRL',
    created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;

-- +goose StatementEnd