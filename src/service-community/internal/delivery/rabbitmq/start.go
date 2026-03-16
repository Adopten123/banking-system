package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

func (c *SSOConsumer) Start(ctx context.Context) {
	msgs, err := c.channel.Consume(
		"community.user.created",
		"community_sso_consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		log.Println("RabbitMQ Consumer started. Waiting for SSO events...")
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled. Stopping RabbitMQ consumer...")
				c.channel.Close()
				c.conn.Close()
				return
			case msg := <-msgs:
				var event domain.UserCreatedEvent
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Printf("Error decoding RabbitMQ message: %v", err)
					msg.Nack(false, false)
					continue
				}

				log.Printf("Received UserCreated event for %s", event.Username)

				err = c.profileService.CreateProfileFromEvent(ctx, event)
				if err != nil {
					log.Printf("Failed to process event: %v", err)
					msg.Nack(false, false)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}
