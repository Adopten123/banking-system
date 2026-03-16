package rabbitmq

import (
	"github.com/Adopten123/banking-system/service-community/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SSOConsumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	profileService domain.ProfileService
}

func NewSSOConsumer(connURL string, service domain.ProfileService) (*SSOConsumer, error) {
	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"sso.events",
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

	q, err := ch.QueueDeclare(
		"community.user.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"user.created",
		"sso.events",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &SSOConsumer{
		conn:           conn,
		channel:        ch,
		profileService: service,
	}, nil
}