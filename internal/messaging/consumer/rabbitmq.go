package consumer

import (
	"fmt"
	"log"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *rabbitmq.Connection
	channel *rabbitmq.Channel
}

func NewRabbitMQConsumer(rabbitURL string) *RabbitMQConsumer {
	conn, err := rabbitmq.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir um canal no RabbitMQ: %v", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: channel,
	}
}

func (c *RabbitMQConsumer) Consume(queueName string) (<-chan rabbitmq.Delivery, error) {
	_, err := c.channel.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	msgs, err := c.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao registrar um consumidor: %w", err)
	}

	return msgs, nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("Conexão do produtor RabbitMQ fechada.")
}
