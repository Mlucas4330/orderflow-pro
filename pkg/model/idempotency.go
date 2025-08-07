package model

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyKey struct {
	Key        uuid.UUID `db:"idempotency_key"`
	UserID     uuid.UUID `db:"user_id"`
	StatusCode int       `db:"response_status_code"`
	Body       []byte    `db:"response_body"`
	CreatedAt  time.Time `db:"created_at"`
}

type IdempotencyResponse struct {
	StatusCode int
	Body       []byte
}
