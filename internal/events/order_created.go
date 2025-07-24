package events

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderCreatedEvent struct {
	OrderID    uuid.UUID          `json:"order_id"`
	CustomerID uuid.UUID          `json:"customer_id"`
	Total      decimal.Decimal    `json:"total"`
	Items      []OrderItemCreated `json:"items"`
}

type OrderItemCreated struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
