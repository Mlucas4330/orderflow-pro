package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID            uuid.UUID       `db:"id"`
	Name          string          `db:"name"`
	SKU           string          `db:"sku"`
	Price         decimal.Decimal `db:"price"`
	StockQuantity int             `db:"stock_quantity"`
	IsActive      bool            `db:"is_active"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}
