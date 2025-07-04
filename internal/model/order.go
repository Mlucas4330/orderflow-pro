package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

type Order struct {
	ID         uuid.UUID       `db:"id"`
	CustomerID uuid.UUID       `db:"customer_id"`
	Status     Status          `db:"status"`
	Total      decimal.Decimal `db:"total"`
	Currency   string          `db:"currency"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}
