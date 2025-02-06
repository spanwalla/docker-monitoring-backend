package broker

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQConsumer struct {
	Channel     *amqp.Channel
	Queue       amqp.Queue
	ProcessFunc func(ctx context.Context, deliveryBody []byte) error
	ConsumerTag string
}

func NewRabbitMQConsumer(channel *amqp.Channel, queueName, consumerTag string, processFunc func(ctx context.Context, deliveryBody []byte) error) (*RabbitMQConsumer, error) {
	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewRabbitMQConsumer - QueueDeclare: %v", err)
	}
	return &RabbitMQConsumer{
		Channel:     channel,
		Queue:       q,
		ProcessFunc: processFunc,
		ConsumerTag: consumerTag,
	}, nil
}

func (c *RabbitMQConsumer) Start() error {
	messages, err := c.Channel.Consume(
		c.Queue.Name,
		c.ConsumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("RabbitMQConsumer.Start - Consume: %v", err)
	}

	go func() {
		for msg := range messages {
			if err = c.ProcessFunc(context.Background(), msg.Body); err != nil {
				log.Errorf("RabbitMQConsumer.Start - ProcessFunc: %v", err)
				err = msg.Nack(false, true)
				if err != nil {
					log.Errorf("RabbitMQConsumer.Start - Nack: %v", err)
					return
				}
				continue
			}
			err = msg.Ack(false)
			if err != nil {
				log.Errorf("RabbitMQConsumer.Start - Ack: %v", err)
				return
			}
		}
	}()

	log.Infof("Consumer %s waiting for messages from queue %s", c.ConsumerTag, c.Queue.Name)
	return nil
}

func (c *RabbitMQConsumer) Stop() error {
	log.Infof("Consumer %s stopping...", c.ConsumerTag)
	return c.Channel.Cancel(c.ConsumerTag, false)
}
