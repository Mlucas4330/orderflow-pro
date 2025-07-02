package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mlucas4330/orderflow-pro/internal/model"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
)

type OrderHandler struct {
	Repo repository.OrderRepository
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	

	h.Repo.CreateOrder(ctx, &model.Order{})
}
