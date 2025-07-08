package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mlucas4330/orderflow-pro/internal/cache"
	"github.com/mlucas4330/orderflow-pro/internal/db"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	dsn := os.Getenv("POSTGRES_DSN")
	redisAddr := os.Getenv("REDIS_ADDR")

	dbpool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}
	defer dbpool.Close()

	redisClient, err := cache.NewRedisClient(ctx, redisAddr)
	if err != nil {
		log.Fatalf("Falha ao conectar com o redis: %v", err)
	}
	defer redisClient.Close()

	healthHandler := handler.NewHealthHandler(dbpool)
	orderRepository := repository.NewOrderRepository(dbpool, redisClient)
	orderHandler := handler.NewOrderHandler(orderRepository)

	router.GET("/ping", healthHandler.Check)
	apiV1 := router.Group("/api/v1")
	{
		orders := apiV1.Group("/orders")
		{
			orders.POST("/", orderHandler.CreateOrder)
			orders.GET("/", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrderById)
			orders.DELETE("/:id", orderHandler.DeleteOrder)
			orders.PUT("/:id", orderHandler.UpdateOrder)
		}
	}

	router.Run()
}
