package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/events"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()
	
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(cfg.KafkaBrokers, ","),
		Topic:          "orders",
		GroupID:        "inventory-service",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,
	})
	defer reader.Close()

	log.Println("Serviço de inventário iniciado. A ouvir por eventos de 'order.created'...")

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Erro ao buscar mensagem do Kafka: %v", err)
			continue
		}

		var event events.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Erro ao desserializar evento OrderCreated: %v", err)
			continue
		}

		log.Printf("Evento OrderCreated recebido! Pedido ID: %s. A dar baixa no stock...", event.OrderID)

		if err := reader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("Erro ao fazer commit da mensagem: %v", err)
		}
	}
}