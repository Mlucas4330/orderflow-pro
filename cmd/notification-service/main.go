package main

import (
	"encoding/json"
	"log"

	"github.com/mlucas4330/orderflow-pro/internal/config"
	"github.com/mlucas4330/orderflow-pro/internal/messaging/consumer"
	"github.com/mlucas4330/orderflow-pro/pkg/messaging"
)

func main() {
	cfg := config.LoadNotificationConfig()

	rabbitConsumer := consumer.NewRabbitMQConsumer(cfg.RabbitURL)
	defer rabbitConsumer.Close()

	msgs, err := rabbitConsumer.Consume("email_notifications")
	if err != nil {
		log.Fatalf("Falha ao consumir fila RabbitMQ: %s", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var payload messaging.NotificationPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Erro ao desserializar mensagem: %s", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("TAREFA RECEBIDA: Enviando e-mail de confirmação para o cliente %s sobre o pedido %s.", payload.CustomerID, payload.OrderID)
			d.Ack(false)
		}
	}()

	log.Printf("Serviço de notificação iniciado. Aguardando tarefas...")
	<-forever
}
