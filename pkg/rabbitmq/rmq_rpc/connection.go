package rmqrpc

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

// Config -.
type Config struct {
	URL      string
	WaitTime time.Duration
	Attempts int
}

// Connection -.
type Connection struct {
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Delivery   <-chan amqp.Delivery
}

// New -.
func New(consumerExchange string, config Config) *Connection {
	conn := &Connection{
		ConsumerExchange: consumerExchange,
		Config:           config,
	}

	return conn
}

// AttemptConnect -.
func (c *Connection) AttemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		log.Infof("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("rmq_rpc - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - amqp.Dial: %w", err)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - Channel: %w", err)
	}

	err = c.Channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - ExchangeDeclare: %w", err)
	}

	queue, err := c.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - QueueDeclare: %w", err)
	}

	err = c.Channel.QueueBind(
		queue.Name,
		"",
		c.ConsumerExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - QueueBind: %w", err)
	}

	c.Delivery, err = c.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rmq_rpc - connect - Consume: %w", err)
	}

	return nil
}
