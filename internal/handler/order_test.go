package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mlucas4330/orderflow-pro/internal/dto"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
	"github.com/mlucas4330/orderflow-pro/internal/middleware"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
	pb "github.com/mlucas4330/orderflow-pro/pkg/productpb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var jwtSecretKey = "jwt-secret-key"

func generateTestToken(t *testing.T, userID uuid.UUID) string {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	require.NoError(t, err)
	return tokenString
}

func TestCreateOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockOrderRepo := new(repository.MockOrderRepository)
	mockIdemRepo := new(repository.MockIdempotencyRepository)
	mockProductClient := new(repository.MockProductServiceClient)

	mockIdemRepo.On("GetResponse", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockProductClient.On("GetProductDetails", mock.Anything, mock.Anything).Return(&pb.GetProductDetailsResponse{Price: "19.99"}, nil)
	mockOrderRepo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdemRepo.On("SaveResponse", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	orderHandler := handler.NewOrderHandler(mockOrderRepo, mockIdemRepo, mockProductClient)
	authMiddleware := middleware.NewAuthMiddleware(jwtSecretKey)

	router := gin.New()
	router.POST("/api/v1/orders", authMiddleware, orderHandler.CreateOrder)

	userID := uuid.New()
	createDTO := dto.CreateOrderRequest{
		CustomerID: userID,
		Items:      []dto.OrderItem{{ProductID: uuid.New(), Quantity: 1}},
	}
	body, _ := json.Marshal(createDTO)
	token := generateTestToken(t, userID)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", uuid.New().String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	mockOrderRepo.AssertExpectations(t)
	mockIdemRepo.AssertExpectations(t)
	mockProductClient.AssertExpectations(t)
}
