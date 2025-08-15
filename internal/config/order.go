package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type OrderConfig struct {
	PostgresUser       string `env:"POSTGRES_USER,required"`
	PostgresPass       string `env:"POSTGRES_PASS,required"`
	PostgresHost       string `env:"POSTGRES_HOST,required"`
	PostgresDb         string `env:"POSTGRES_DB,required"`
	RedisAddr          string `env:"REDIS_ADDR,required"`
	RedisDB            int    `env:"REDIS_DB,required"`
	KafkaBrokers       string `env:"KAFKA_BROKERS,required"`
	ProductServiceAddr string `env:"PRODUCT_SERVICE_ADDR,required"`
	JWTSecretKey       string `env:"JWT_SECRET_KEY,required"`
	RabbitmqUser       string `env:"RABBITMQ_USER,required"`
	RabbitmqPass       string `env:"RABBITMQ_PASS,required"`
	RabbitmqHost       string `env:"RABBITMQ_HOST,required"`
}

func LoadOrderConfig() *OrderConfig {
	cfg := OrderConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar as configuração do pedido: %+v", err)
	}
	return &cfg
}
