package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/model"
)

type OrderRepository interface {
	FindOrders(ctx context.Context) []model.Order
	FindOrderById(ctx context.Context) model.Order
	CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem) error
}

type PostgresOrderRepository struct {
	DB *pgxpool.Pool
}

func NewOrderRepository(pgpool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		DB: pgpool,
	}
}

func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem) error {
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

	rows := make([][]any, len(items))
	for i, item := range items {
		rows[i] = []any{item.ID, item.OrderID, item.ProductID, item.Quantity, item.PriceAtTime}
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"order_items"}, itemQuery, pgx.CopyFromRows(rows))
	if err != nil {
		return fmt.Errorf("erro ao fazer bulk insert na tabela order_items: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresOrderRepository) FindOrders(ctx context.Context) {

}

func (r *PostgresOrderRepository) FindOrderById(ctx context.Context) {
	_, err := r.DB.Exec("")

	if(pgx.ErrNoRows){
		
	}
}
