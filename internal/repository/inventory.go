package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	DecrementStock(ctx context.Context, productId uuid.UUID, quantity int) error
}

type PostgresInventoryRepository struct {
	DB *pgxpool.Pool
}

func NewInventoryRepository(dbpool *pgxpool.Pool) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{DB: dbpool}
}

func (r *PostgresInventoryRepository) DecrementStock(ctx context.Context, productId uuid.UUID, quantity int) error {
	query := `
		UPDATE products SET stock_quantity = stock_quantity - $1, updated_at = $2 WHERE id = $3 
	`

	tag, err := r.DB.Exec(ctx, query, quantity, time.Now(), productId)
	if err != nil {
		return fmt.Errorf("erro ao atualizar a tabela products: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
