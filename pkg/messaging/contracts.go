package messaging

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderCreatedEvent struct {
	OrderID    uuid.UUID       `json:"order_id"`
	CustomerID uuid.UUID       `json:"customer_id"`
	Total      decimal.Decimal `json:"total"`
	Items      []OrderItem   `json:"items"`
}

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type NotificationPayload struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Message    string `json:"message"`
}