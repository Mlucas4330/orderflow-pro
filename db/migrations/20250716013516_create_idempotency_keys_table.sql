-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  idempotency_keys (
    idempotency_key UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    response_status_code INT NOT NULL,
    response_body BYTEA,
    created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE idempotency_keys;

-- +goose StatementEnd