package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/events"
	"github.com/mlucas4330/orderflow-pro/internal/messaging/producer"
	"github.com/mlucas4330/orderflow-pro/pkg/model"
	redis "github.com/redis/go-redis/v9"
)

type OrderRepository interface {
	FindOrders(ctx context.Context) ([]model.Order, error)
	FindOrderById(ctx context.Context, id uuid.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, order *model.Order, orderItems []model.OrderItem) error
	UpdateOrder(ctx context.Context, id uuid.UUID, status model.Status) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}

type PostgresOrderRepository struct {
	DB               *pgxpool.Pool
	Redis            *redis.Client
	KafkaProducer    *producer.KafkaProducer
	RabbitMQProducer *producer.RabbitMQProducer
}

func NewOrderRepository(pgpool *pgxpool.Pool, redis *redis.Client, kafkaProducer *producer.KafkaProducer, rabbitProducer *producer.RabbitMQProducer) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		DB:               pgpool,
		Redis:            redis,
		KafkaProducer:    kafkaProducer,
		RabbitMQProducer: rabbitProducer,
	}
}

func (r *PostgresOrderRepository) FindOrders(ctx context.Context) ([]model.Order, error) {
	key := "orders:list:all"

	result, err := r.Redis.Get(ctx, key).Result()

	if err == nil {
		var orders []model.Order
		err := json.Unmarshal([]byte(result), &orders)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do json: %w", err)
		}

		return orders, nil
	}

	if err != redis.Nil {
		log.Printf("Erro ao buscar do Redis, mas não é um cache miss: %v", err)
	}

	query := `SELECT id, customer_id, status, total, currency, created_at, updated_at FROM orders`
	orderRows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos: %w", err)
	}
	defer orderRows.Close()

	var orders []model.Order
	for orderRows.Next() {
		var order model.Order
		err := orderRows.Scan(
			&order.ID, &order.CustomerID, &order.Status, &order.Total,
			&order.Currency, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}
		orders = append(orders, order)
	}
	if err = orderRows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos pedidos: %w", err)
	}

	if len(orders) == 0 {
		return []model.Order{}, nil
	}

	orderIDs := make([]uuid.UUID, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price_at_time FROM order_items WHERE order_id = ANY($1)`
	itemRows, err := r.DB.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar itens dos pedidos: %w", err)
	}
	defer itemRows.Close()

	itemsByOrderID := make(map[uuid.UUID][]model.OrderItem)
	for itemRows.Next() {
		var orderItem model.OrderItem
		if err := itemRows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.PriceAtTime); err != nil {
			return nil, fmt.Errorf("erro ao escanear item: %w", err)
		}
		itemsByOrderID[orderItem.OrderID] = append(itemsByOrderID[orderItem.OrderID], orderItem)
	}
	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos itens: %w", err)
	}

	for i, order := range orders {
		if orderItems, ok := itemsByOrderID[order.ID]; ok {
			orders[i].OrderItems = orderItems
		} else {
			orders[i].OrderItems = []model.OrderItem{}
		}
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		return nil, fmt.Errorf("erro ao transformar pedido em json: %w", err)
	}

	if err := r.Redis.Set(ctx, key, jsonData, 30*time.Second).Err(); err != nil {
		log.Printf("Falha ao salvar pedidos no cache do Redis: %v", err)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) FindOrderById(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	key := fmt.Sprintf("order:%s", id.String())

	result, err := r.Redis.Get(ctx, key).Result()

	if err == nil {
		var order model.Order
		err := json.Unmarshal([]byte(result), &order)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do json: %w", err)
		}

		return &order, nil
	}

	if err != redis.Nil {
		log.Printf("Erro ao buscar do Redis, mas não é um cache miss: %v", err)
	}

	orderQuery := `
		SELECT id, customer_id, status, total, currency, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	var order model.Order
	err = r.DB.QueryRow(ctx, orderQuery, id).Scan(
		&order.ID, &order.CustomerID, &order.Status, &order.Total,
		&order.Currency, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("erro ao buscar o pedido: %w", err)
	}

	itemsQuery := `
		SELECT id, product_id, quantity, price_at_time
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := r.DB.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar os itens do pedido: %w", err)
	}
	defer rows.Close()

	var orderItems []model.OrderItem
	for rows.Next() {
		var orderItem model.OrderItem
		if err := rows.Scan(&orderItem.ID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.PriceAtTime); err != nil {
			return nil, fmt.Errorf("erro ao escanear item do pedido: %w", err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a leitura dos itens do pedido: %w", err)
	}

	order.OrderItems = orderItems

	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("erro ao transformar pedido em json: %w", err)
	}

	if err := r.Redis.Set(ctx, key, jsonData, 10*time.Minute).Err(); err != nil {
		log.Printf("Falha ao salvar pedido %s no cache do Redis: %v", id.String(), err)
	}

	return &order, nil
}

func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *model.Order, orderItems []model.OrderItem) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	orderQuery := `
		INSERT INTO orders (id, customer_id, status, total, currency, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, orderQuery, order.ID, order.CustomerID, order.Status, order.Total, order.Currency, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao inserir na tabela orders: %w", err)
	}

	itemQuery := []string{"id", "order_id", "product_id", "quantity", "price_at_time"}
	rows := make([][]any, len(orderItems))
	for i, orderItem := range orderItems {
		rows[i] = []any{orderItem.ID, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.PriceAtTime}
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"order_items"}, itemQuery, pgx.CopyFromRows(rows))
	if err != nil {
		return fmt.Errorf("erro ao fazer bulk insert na tabela order_items: %w", err)
	}

	eventItems := make([]events.OrderItemCreated, len(orderItems))
	for i, item := range orderItems {
		eventItems[i] = events.OrderItemCreated{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	event := events.OrderCreatedEvent{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Total:      order.Total,
		Items:      eventItems,
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao comitar transação: %w", err)
	}

	go func() {
		err := r.KafkaProducer.PublishOrderCreated(context.Background(), event)
		if err != nil {
			log.Printf("ERRO ao publicar evento OrderCreated no Kafka: %v", err)
		}
	}()

	return nil
}

func (r *PostgresOrderRepository) UpdateOrder(ctx context.Context, id uuid.UUID, status model.Status) error {
	query := `
		UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3 
	`

	tag, err := r.DB.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar a tabela orders: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *PostgresOrderRepository) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM orders WHERE id = $1
	`

	tag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir ordem: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
