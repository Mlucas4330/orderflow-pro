package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/cache"
	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
	"github.com/mlucas4330/orderflow-pro/internal/messaging"
	"github.com/mlucas4330/orderflow-pro/internal/middleware"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	cfg := config.LoadConfig()

	dbpool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}
	defer dbpool.Close()

	redisClient, err := cache.NewRedisClient(ctx, cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Falha ao conectar com o redis: %v", err)
	}
	defer redisClient.Close()

	healthHandler := handler.NewHealthHandler(dbpool)
	kafkaProducer := messaging.NewKafkaProducer(cfg.KafkaBrokers)
	defer kafkaProducer.Close()

	orderRepository := repository.NewOrderRepository(dbpool, redisClient, kafkaProducer)
	idempotencyRepository := repository.NewIdempotencyRepository(dbpool)
	orderHandler := handler.NewOrderHandler(orderRepository, idempotencyRepository)

	router.GET("/ping", healthHandler.Check)
	apiV1 := router.Group("/api/v1")
	{
		orders := apiV1.Group("/orders")
		{
			orders.POST("/", middleware.AuthMiddleware(), orderHandler.CreateOrder)
			orders.GET("/", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrderById)
			orders.DELETE("/:id", orderHandler.DeleteOrder)
			orders.PATCH("/:id", orderHandler.UpdateOrder)
		}
	}

	err = router.Run()
	if err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
