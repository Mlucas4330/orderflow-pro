package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func New(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
