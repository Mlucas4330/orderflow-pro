package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type OrderConfig struct {
	PostgresDSN  string `env:"POSTGRES_DSN,required"`
	RedisAddr    string `env:"REDIS_ADDR,required"`
	RedisDB      int    `env:"REDIS_DB,required"`
	KafkaBrokers string `env:"KAFKA_BROKERS,required"`
	RabbitURL    string `env:"RABBITMQ_URL,required"`
}

func LoadOrderConfig() *OrderConfig {
	cfg := OrderConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar a configuração: %+v", err)
	}
	return &cfg
}
