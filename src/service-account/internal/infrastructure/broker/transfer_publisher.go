package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishTransferCreated serializes the event to JSON and sends it to RabbitMQ
func (p *RabbitMQPublisher) publishMessage(ctx context.Context, event domain.DomainEvent) error {
	eventType := event.EventName()

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	err = p.ch.PublishWithContext(ctx,
		"account_events", // Exchange
		"",               // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Type:         eventType,
			Body:         body,
			DeliveryMode: amqp.Persistent, // message is saved to disk
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	return nil
}

func (p *RabbitMQPublisher) PublishTransferCreated(ctx context.Context, event domain.TransferCreatedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishAccountCreated(ctx context.Context, event domain.AccountCreatedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishAccountStatusChanged(ctx context.Context, event domain.AccountStatusChangedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishDepositCompleted(ctx context.Context, event domain.DepositCompletedEvent) error {
	return p.publishMessage(ctx, event)
}

func (p *RabbitMQPublisher) PublishWithdrawalCompleted(ctx context.Context, event domain.WithdrawalCompletedEvent) error {
	return p.publishMessage(ctx, event)
}
