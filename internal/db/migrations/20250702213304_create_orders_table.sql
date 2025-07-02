-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- customer_id UUID NOT NULL REFERENCES customers(id),

  status TEXT NOT NULL CHECK (status IN (
    'pending', 'paid', 'shipped', 'delivered', 'cancelled', 'refunded'
  )),

  total_amount NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
  shipping_cost NUMERIC(10, 2) DEFAULT 0 CHECK (shipping_cost >= 0),
  discount_amount NUMERIC(10, 2) DEFAULT 0 CHECK (discount_amount >= 0),

  currency CHAR(3) NOT NULL DEFAULT 'BRL',

  payment_method TEXT NOT NULL CHECK (payment_method IN (
    'credit_card', 'pix', 'boleto', 'paypal', 'cash', 'bank_transfer'
  )),

  -- shipping_address_id UUID REFERENCES addresses(id),
  -- billing_address_id UUID REFERENCES addresses(id),

  items_count INTEGER NOT NULL CHECK (items_count >= 0),

  notes TEXT,
  is_test BOOLEAN NOT NULL DEFAULT false,

  metadata JSONB,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
