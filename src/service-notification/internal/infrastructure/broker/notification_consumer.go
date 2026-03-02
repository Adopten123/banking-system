package broker

import (
	"log"

	"github.com/Adopten123/banking-system/service-notification/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationConsumer struct {
	channel      *amqp.Channel
	exchange     string
	queue        string
	notifService *service.NotificationService
}

func NewNotificationConsumer(
	ch *amqp.Channel,
	exchange, queue string,
	notifSvc *service.NotificationService,
) *NotificationConsumer {
	return &NotificationConsumer{
		channel:      ch,
		exchange:     exchange,
		queue:        queue,
		notifService: notifSvc,
	}
}

// Setup - method for config exchange, queue and binding
func (c *NotificationConsumer) Setup() error {
	// 1. Declare exchange (fanout)
	err := c.channel.ExchangeDeclare(
		c.exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	// 2. Declare queue
	q, err := c.channel.QueueDeclare(
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	// 3. Queue binding
	err = c.channel.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		c.exchange, // exchange
		false,
		nil,
	)

	return err
}

// Start - starting async reading messages
func (c *NotificationConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer tag
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			err := c.notifService.HandleMessage(d.Type, d.Body)

			if err != nil {
				log.Printf("[ERROR] Error processing message type %s: %v", d.Type, err)
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages in queue: %s", c.queue)
	return nil
}
