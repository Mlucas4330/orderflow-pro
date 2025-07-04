package dto

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID   `json:"customer_id" binding:"required,uuid"`
	Items      []OrderItem `json:"items" binding:"required,min=1"`
}

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id" binding:"required,uuid"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
}
