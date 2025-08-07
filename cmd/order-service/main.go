package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/cache"
	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
	"github.com/mlucas4330/orderflow-pro/internal/messaging/producer"
	"github.com/mlucas4330/orderflow-pro/internal/middleware"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
	pb "github.com/mlucas4330/orderflow-pro/pkg/productpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	cfg := config.LoadOrderConfig()

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

	kafkaProducer := producer.NewKafkaProducer(cfg.KafkaBrokers)
	defer kafkaProducer.Close()

	rabbitProducer := producer.NewRabbitMQProducer(cfg.RabbitURL)
	defer rabbitProducer.Close()

	grpcconn, err := grpc.NewClient(cfg.ProductServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar com o product-service via gRPC: %v", err)
	}
	defer grpcconn.Close()

	orderRepository := repository.NewOrderRepository(dbpool, redisClient, kafkaProducer, rabbitProducer)
	idempotencyRepository := repository.NewIdempotencyRepository(dbpool)
	healthHandler := handler.NewHealthHandler(dbpool)
	productClient := pb.NewProductServiceClient(grpcconn)
	orderHandler := handler.NewOrderHandler(orderRepository, idempotencyRepository, productClient)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecretKey)

	router.GET("/ping", healthHandler.Check)
	apiV1 := router.Group("/api/v1")
	{
		orders := apiV1.Group("/orders")
		{
			orders.POST("/", authMiddleware, orderHandler.CreateOrder)
			orders.GET("/", authMiddleware, orderHandler.GetOrders)
			orders.GET("/:id", authMiddleware, orderHandler.GetOrderById)
			orders.DELETE("/:id", authMiddleware, orderHandler.DeleteOrder)
			orders.PATCH("/:id", authMiddleware, orderHandler.UpdateOrder)
		}
	}

	err = router.Run()
	if err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
