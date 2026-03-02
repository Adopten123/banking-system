package main

import (
	"log"

	"github.com/Adopten123/banking-system/service-notification/internal/config"
	"github.com/Adopten123/banking-system/service-notification/internal/infrastructure/broker"
	"github.com/Adopten123/banking-system/service-notification/internal/server"
	"github.com/Adopten123/banking-system/service-notification/internal/service"
)

func main() {
	cfg := config.MustLoad()

	log.Printf("Starting notification service in %s env", cfg.Env)

	conn, ch, err := broker.NewRabbitMQConn(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}

	notificationSvc := service.NewNotificationService()

	consumer := broker.NewNotificationConsumer(
		ch,
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.Queue,
		notificationSvc,
	)

	if err := consumer.Setup(); err != nil {
		log.Fatalf("Failed to setup RabbitMQ topology: %v", err)
	}

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	log.Println("Notification service is successfully running!")

	server.GracefulShutdown(func() {
		log.Println("Closing RabbitMQ connections...")

		if err := ch.Close(); err != nil {
			log.Printf("[ERROR] Failed to close RabbitMQ channel: %v", err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("[ERROR] Failed to close RabbitMQ connection: %v", err)
		}

		log.Println("RabbitMQ connections closed properly")
	})
}
