package producer

import (
	"context"
	"fmt"
	"log"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn    *rabbitmq.Connection
	channel *rabbitmq.Channel
}

func NewRabbitMQProducer(rabbitURL string) *RabbitMQProducer {
	conn, err := rabbitmq.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir um canal no RabbitMQ: %v", err)
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: channel,
	}
}

func (p *RabbitMQProducer) Publish(ctx context.Context, queueName string, body []byte) error {
	_, err := p.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		rabbitmq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("falha ao publicar mensagem: %w", err)
	}

	log.Printf("Mensagem publicada na fila: %s", queueName)
	return nil
}

func (p *RabbitMQProducer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	log.Println("Conexão do produtor RabbitMQ fechada.")
}
