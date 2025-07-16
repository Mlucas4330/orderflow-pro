package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/model"
)

type IdempotencyRepository interface {
	GetResponse(ctx context.Context, key uuid.UUID, userID uuid.UUID) (*model.IdempotencyResponse, error)
	SaveResponse(ctx context.Context, key uuid.UUID, userID uuid.UUID, response *model.IdempotencyResponse) error
}

type PostgresIdempotencyRepository struct {
	DB *pgxpool.Pool
}

func NewIdempotencyRepository(dbpool *pgxpool.Pool) *PostgresIdempotencyRepository {
	return &PostgresIdempotencyRepository{DB: dbpool}
}

func (r *PostgresIdempotencyRepository) GetResponse(ctx context.Context, key uuid.UUID, userID uuid.UUID) (*model.IdempotencyResponse, error) {
	query := `
		SELECT response_status_code, response_body 
		FROM idempotency_keys 
		WHERE idempotency_key = $1 AND user_id = $2
	`

	var res model.IdempotencyResponse

	err := r.DB.QueryRow(ctx, query, key, userID).Scan(&res.StatusCode, &res.Body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar chave de idempotência: %w", err)
	}

	return &res, nil
}

func (r *PostgresIdempotencyRepository) SaveResponse(ctx context.Context, key uuid.UUID, userID uuid.UUID, response *model.IdempotencyResponse) error {
	query := `
		INSERT INTO idempotency_keys (idempotency_key, user_id, response_status_code, response_body)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.DB.Exec(ctx, query, key, userID, response.StatusCode, response.Body)
	if err != nil {
		return fmt.Errorf("erro ao salvar chave de idempotência: %w", err)
	}

	return nil
}
