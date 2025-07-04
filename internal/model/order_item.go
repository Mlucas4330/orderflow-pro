package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
    ID          uuid.UUID       `db:"id"`
    OrderID     uuid.UUID       `db:"order_id"`
    ProductID   uuid.UUID       `db:"product_id"`
    Quantity    int             `db:"quantity"`
    PriceAtTime decimal.Decimal `db:"price_at_time"`
}