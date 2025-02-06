package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewRabbitMQPublisher(channel *amqp.Channel, queueName string) (*RabbitMQPublisher, error) {
	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewRabbitMQPublisher - QueueDeclare: %v", err)
	}
	return &RabbitMQPublisher{Channel: channel, Queue: q}, nil
}

func (p *RabbitMQPublisher) Publish(body []byte) error {
	return p.Channel.Publish(
		"",
		p.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
