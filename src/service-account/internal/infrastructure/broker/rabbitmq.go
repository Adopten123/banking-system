package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"account_events", // Name of exchange
		"fanout",         // Type
		true,             // Durable (If RabbitMQ crashes,our exchange will not be deleted)
		false,            // Auto-deleted
		false,            // Internal
		false,            // No-wait
		nil,              // Arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return &RabbitMQPublisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *RabbitMQPublisher) Close() error {
	if err := p.ch.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
