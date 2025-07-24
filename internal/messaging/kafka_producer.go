package messaging

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/mlucas4330/orderflow-pro/internal/events"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(kafkaBrokers string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(kafkaBrokers, ",")...),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) PublishOrderCreated(ctx context.Context, event events.OrderCreatedEvent) error {
	msgValue, err := json.Marshal(event)
	if err != nil {
		log.Printf("Erro ao serializar evento OrderCreated: %v", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.OrderID.String()),
		Value: msgValue,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Erro ao publicar mensagem no Kafka: %v", err)
		return err
	}

	log.Printf("Evento OrderCreated publicado para o pedido: %s", event.OrderID)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
