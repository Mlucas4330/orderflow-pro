package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	env "github.com/joho/godotenv"

	"github.com/mlucas4330/orderflow-pro/internal/cache"
	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/messaging"
	"github.com/mlucas4330/orderflow-pro/internal/model"
	redis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*PostgresOrderRepository, *pgxpool.Pool, *redis.Client) {
	_ = env.Load("../../.env.test")

	cfg := config.LoadOrderConfig()

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	require.NoError(t, err, "Falha ao conectar ao banco de dados de teste")

	redisClient, err := cache.NewRedisClient(ctx, cfg.RedisAddr, cfg.RedisDB)
	require.NoError(t, err, "Falha ao conectar ao Redis de teste")

	kafkaProducer := messaging.NewKafkaProducer(cfg.KafkaBrokers)
	defer kafkaProducer.Close()

	rabbitProducer := messaging.NewRabbitMQProducer(cfg.RabbitURL)
	defer rabbitProducer.Close()

	repo := NewOrderRepository(dbpool, redisClient, kafkaProducer, rabbitProducer)

	return repo, dbpool, redisClient
}

func cleanup(t *testing.T, dbpool *pgxpool.Pool, redisClient *redis.Client) {
	_, err := dbpool.Exec(context.Background(), "TRUNCATE TABLE order_items, orders RESTART IDENTITY")
	require.NoError(t, err)

	err = redisClient.FlushDB(context.Background()).Err()
	require.NoError(t, err)
}

func TestCreateAndFindOrder(t *testing.T) {
	repo, dbpool, redisClient := setupTest(t)
	t.Cleanup(func() {
		cleanup(t, dbpool, redisClient)
		dbpool.Close()
		redisClient.Close()
	})

	ctx := context.Background()

	orderID := uuid.New()
	customerID := uuid.New()
	productID := uuid.New()

	order := &model.Order{
		ID:         orderID,
		CustomerID: customerID,
		Status:     model.StatusPending,
		Total:      decimal.NewFromFloat(99.95),
		Currency:   "BRL",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	items := []model.OrderItem{
		{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   productID,
			Quantity:    5,
			PriceAtTime: decimal.NewFromFloat(19.99),
		},
	}

	err := repo.CreateOrder(ctx, order, items)

	require.NoError(t, err, "CreateOrder não deveria retornar um erro")

	foundOrder, err := repo.FindOrderById(ctx, orderID)

	require.NoError(t, err, "FindOrderByID não deveria retornar um erro")
	require.NotNil(t, foundOrder, "O pedido encontrado não deveria ser nulo")

	require.Equal(t, order.ID, foundOrder.ID, "O ID do pedido não bate")
	require.Equal(t, order.Status, foundOrder.Status, "O Status do pedido não bate")
	require.True(t, order.Total.Equal(foundOrder.Total), "O Total do pedido não bate")

	require.Len(t, foundOrder.OrderItems, 1, "Deveria haver 1 item no pedido")
	require.Equal(t, items[0].ProductID, foundOrder.OrderItems[0].ProductID, "O ProductID do item não bate")
}

func TestFindOrderCache(t *testing.T) {
	repo, dbpool, redisClient := setupTest(t)
	t.Cleanup(func() {
		cleanup(t, dbpool, redisClient)
		dbpool.Close()
		redisClient.Close()
	})
	ctx := context.Background()

	orderID := uuid.New()
	customerID := uuid.New()
	productID := uuid.New()
	order := &model.Order{
		ID:         orderID,
		CustomerID: customerID,
		Status:     model.StatusPending,
		Total:      decimal.NewFromFloat(99.95),
		Currency:   "BRL",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	items := []model.OrderItem{
		{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   productID,
			Quantity:    5,
			PriceAtTime: decimal.NewFromFloat(19.99),
		},
	}
	err := repo.CreateOrder(ctx, order, items)
	require.NoError(t, err)

	_, err = repo.FindOrderById(ctx, orderID)
	require.NoError(t, err)

	keyExists, err := redisClient.Exists(ctx, "order:"+orderID.String()).Result()
	require.NoError(t, err)
	require.Equal(t, int64(1), keyExists, "A chave do pedido deveria existir no cache após o primeiro 'find'")

	_, err = dbpool.Exec(ctx, "TRUNCATE TABLE order_items, orders RESTART IDENTITY")
	require.NoError(t, err)

	cachedOrder, err := repo.FindOrderById(ctx, orderID)
	require.NoError(t, err, "A busca no cache não deveria dar erro, mesmo com o banco limpo")
	require.NotNil(t, cachedOrder, "Deveria encontrar o pedido no cache")
	require.Equal(t, orderID, cachedOrder.ID, "O ID do pedido do cache está incorreto")
}
