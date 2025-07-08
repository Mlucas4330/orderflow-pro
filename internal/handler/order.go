package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (h *OrderHandler) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	orders, err := h.Repo.FindOrders(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pedidos não encontrados"})
			return
		}

		log.Printf("Erro ao buscar pedidos no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao buscar o pedido"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pedido inválido"})
		return
	}

	order, err := h.Repo.FindOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pedido não encontrado"})
			return
		}

		log.Printf("Erro ao buscar pedido por ID no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao buscar o pedido"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
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

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pedido inválido"})
		return
	}

	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	err = h.Repo.UpdateOrder(ctx, id, req.Status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pedido não encontrado para atualização"})
			return
		}
		log.Printf("Erro ao atualizar pedido no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao processar o pedido"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pedido inválido"})
		return
	}

	ctx := c.Request.Context()
	err = h.Repo.DeleteOrder(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pedido não encontrado para exclusão"})
			return
		}
		log.Printf("Erro ao excluir pedido no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao processar o pedido"})
		return
	}

	c.Status(http.StatusNoContent)
}
