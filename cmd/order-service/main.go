package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mlucas4330/orderflow-pro/internal/db"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	dsn := os.Getenv("POSTGRES_DSN")

	dbpool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}
	defer dbpool.Close()

	healthHandler := handler.NewHealthHandler(dbpool)
	orderRepository := repository.NewOrderRepository(dbpool)
	orderHandler := handler.NewOrderHandler(orderRepository)

	router.GET("/ping", healthHandler.Check)
	apiV1 := router.Group("/api/v1")
	{
		orders := apiV1.Group("/orders")
		{
			orders.POST("/", orderHandler.Create)
			orders.GET("/", orderHandler.Create)
			orders.GET("/:id", orderHandler.Create)
			orders.DELETE("/:id", orderHandler.Create)
			orders.PUT("/:id", orderHandler.Create)
		}
	}

	router.Run()
}
