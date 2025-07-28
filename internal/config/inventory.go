package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type InventoryConfig struct {
	PostgresDSN  string `env:"POSTGRES_DSN,required"`
	KafkaBrokers string `env:"KAFKA_BROKERS,required"`
}

func LoadInventoryConfig() *InventoryConfig {
	cfg := InventoryConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar a configuração: %+v", err)
	}
	return &cfg
}
