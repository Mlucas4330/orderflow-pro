package cache

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, redisAddr string, redisDB int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao Redis: %w", err)
	}

	return rdb, nil
}
