package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type Config struct {
	PostgresDSN  string `env:"POSTGRES_DSN,required"`
	RedisAddr    string `env:"REDIS_ADDR,required"`
	RedisDB      int    `env:"REDIS_DB,required"`
	KafkaBrokers string `end:"KAFKA_BROKERS,required"`
}

func LoadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar a configuração: %+v", err)
	}
	return &cfg
}
