package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/mlucas4330/orderflow-pro/internal/dto"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
	"github.com/mlucas4330/orderflow-pro/pkg/model"
	pb "github.com/mlucas4330/orderflow-pro/pkg/productpb"

	"github.com/shopspring/decimal"
)

type OrderHandler struct {
	OrderRepo       repository.OrderRepository
	IdempotencyRepo repository.IdempotencyRepository
	ProductClient   pb.ProductServiceClient
}

func NewOrderHandler(orderRepo repository.OrderRepository, idempotencyRepo repository.IdempotencyRepository, productClient pb.ProductServiceClient) *OrderHandler {
	return &OrderHandler{OrderRepo: orderRepo, IdempotencyRepo: idempotencyRepo, ProductClient: productClient}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	orders, err := h.OrderRepo.FindOrders(ctx)
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

	order, err := h.OrderRepo.FindOrderById(ctx, id)
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
	ctx := c.Request.Context()

	idempotencyKeyStr := c.GetHeader("Idempotency-Key")
	idempotencyKey, err := uuid.Parse(idempotencyKeyStr)

	if err == nil {
		mockCustomerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		if savedResponse, err := h.IdempotencyRepo.GetResponse(ctx, idempotencyKey, mockCustomerID); err == nil && savedResponse != nil {
			log.Printf("HIT de idempotência para a chave: %s", idempotencyKeyStr)
			c.Data(savedResponse.StatusCode, "application/json; charset=utf-8", savedResponse.Body)
			return
		}
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido: " + err.Error()})
		return
	}

	orderID := uuid.New()
	var orderItems []model.OrderItem
	total := decimal.NewFromInt(0)

	for _, itemDTO := range req.Items {
		log.Printf("Buscando detalhes do produto %s via gRPC...", itemDTO.ProductID)
		productDetails, err := h.ProductClient.GetProductDetails(ctx, &pb.GetProductDetailsRequest{
			ProductId: itemDTO.ProductID.String(),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "produto inválido: " + itemDTO.ProductID.String()})
			return
		}
		priceAtTime, err := decimal.NewFromString(productDetails.GetPrice())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "preço inválido retornado pelo serviço de produto"})
			return
		}

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
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		OrderItems: orderItems,
	}

	if err := h.OrderRepo.CreateOrder(ctx, order, orderItems); err != nil {
		log.Printf("Erro ao criar pedido no repositório: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno ao processar o pedido"})
		return
	}

	if idempotencyKey != uuid.Nil {
		responseBody, _ := json.Marshal(order)
		responseToSave := &model.IdempotencyResponse{
			StatusCode: http.StatusCreated,
			Body:       responseBody,
		}

		if err := h.IdempotencyRepo.SaveResponse(ctx, idempotencyKey, req.CustomerID, responseToSave); err != nil {
			log.Printf("AVISO CRÍTICO: Falha ao salvar a resposta de idempotência: %v", err)
		}
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
	err = h.OrderRepo.UpdateOrder(ctx, id, model.Status(req.Status))

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
	err = h.OrderRepo.DeleteOrder(ctx, id)

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
