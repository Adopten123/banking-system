package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type transferPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewTransferPublisher(connURL string) (domain.TransferPublisher, error) {
	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"community.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &transferPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *transferPublisher) PublishTransferCreated(ctx context.Context, event domain.TransferCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer event: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		"community.events",
		"transfer.created",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish transfer event: %w", err)
	}

	log.Printf("Published TransferCreated event for MessageID: %d", event.MessageID)
	return nil
}
