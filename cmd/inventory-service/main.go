package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/events"
	"github.com/mlucas4330/orderflow-pro/internal/repository"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadInventoryConfig()

	postgresDsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresDb)

	dbpool, err := pgxpool.New(ctx, postgresDsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}

	inventoryRepo := repository.NewInventoryRepository(dbpool)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(cfg.KafkaBrokers, ","),
		Topic:       "orders",
		GroupID:     "inventory-service",
		Logger:      kafka.LoggerFunc(log.Printf),
		ErrorLogger: kafka.LoggerFunc(log.Printf),
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

		for _, item := range event.Items {
			const maxAttempts = 3
			var err error

			for attempt := 1; attempt <= maxAttempts; attempt++ {
				log.Printf("Tentativa %d de atualizar estoque para o produto %s: -%d", attempt, item.ProductID, item.Quantity)

				err = inventoryRepo.DecrementStock(ctx, item.ProductID, item.Quantity)

				if err == nil {
					log.Printf("Estoque para o produto %s atualizado com sucesso.", item.ProductID)
					break
				}

				log.Printf("AVISO: Falha na tentativa %d para o produto %s: %v", attempt, item.ProductID, err)

				if attempt < maxAttempts {
					time.Sleep(time.Duration(attempt) * time.Second)
				}
			}

			if err != nil {
				log.Printf("ERRO FINAL: Todas as %d tentativas falharam para o produto %s. Erro: %v", maxAttempts, item.ProductID, err)
			}
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Erro ao fazer commit da mensagem: %v", err)
		}
	}
}
