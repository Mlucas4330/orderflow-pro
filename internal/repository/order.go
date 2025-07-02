package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}

type PostgresOrderRepository struct {
	DB *pgxpool.Pool
}

func NewOrderRepository(pgpool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		DB: pgpool,
	}
}

func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, status, total_amount) 
		VALUES ($1, $2, $3, $4, ...)
	`

	_, err := r.DB.Exec(ctx, query, order.ID, order.CustomerID, order.Status, order.TotalAmount)

	return err
}
