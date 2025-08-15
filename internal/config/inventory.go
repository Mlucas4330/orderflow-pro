package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type InventoryConfig struct {
	PostgresUser string `env:"POSTGRES_USER,required"`
	PostgresPass string `env:"POSTGRES_PASS,required"`
	PostgresHost string `env:"POSTGRES_HOST,required"`
	PostgresDb   string `env:"POSTGRES_DB,required"`
	KafkaBrokers string `env:"KAFKA_BROKERS,required"`
}

func LoadInventoryConfig() *InventoryConfig {
	cfg := InventoryConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar a configuração: %+v", err)
	}
	return &cfg
}
