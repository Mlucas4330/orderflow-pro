package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

type NotificationConfig struct {
	RabbitURL string `env:"RABBITMQ_URL,required"`
}

func LoadNotificationConfig() *NotificationConfig {
	cfg := NotificationConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Não foi possível carregar a configuração: %+v", err)
	}
	return &cfg
}
