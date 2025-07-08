package dto

import (
	"github.com/google/uuid"
	"github.com/mlucas4330/orderflow-pro/internal/model"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID   `json:"customer_id" binding:"required,uuid"`
	Items      []OrderItem `json:"items" binding:"required,min=1"`
}

type UpdateOrderRequest struct {
	Status model.Status `json:"status" binding:"required"`
}

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id" binding:"required,uuid"`
	Quantity  int       `json:"quantity" binding:"required,gt=0"`
}
