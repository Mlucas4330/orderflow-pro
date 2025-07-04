package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mlucas4330/orderflow-pro/internal/dto"
	"github.com/mlucas4330/orderflow-pro/internal/model"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
	"github.com/shopspring/decimal"
)

type OrderHandler struct {
	Repo repository.OrderRepository
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido: " + err.Error()})
		return
	}

	orderID := uuid.New()

	var orderItems []model.OrderItem

	total := decimal.NewFromInt(0)

	for _, itemDTO := range req.Items {
		priceAtTime := decimal.NewFromFloat(19.99)

		orderItem := model.OrderItem{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   itemDTO.ProductID,
			Quantity:    itemDTO.Quantity,
			PriceAtTime: priceAtTime,
		}

		orderItems = append(orderItems, orderItem)

		itemTotal := priceAtTime.Mul(decimal.NewFromInt(int64(itemDTO.Quantity)))
		total = total.Add(itemTotal)
	}

	order := &model.Order{
		ID:         orderID,
		CustomerID: req.CustomerID,
		Status:     model.StatusPending,
		Total:      total,
		Currency:   "BRL",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	ctx := c.Request.Context()
	err := h.Repo.CreateOrder(ctx, order, orderItems)
	if err != nil {
		log.Printf("Erro ao criar pedido no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao processar o pedido"})
		return
	}

	c.JSON(http.StatusCreated, order)
}
