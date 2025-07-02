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
	router.POST("/v1/orders", orderHandler.Create)

	router.Run()
}
