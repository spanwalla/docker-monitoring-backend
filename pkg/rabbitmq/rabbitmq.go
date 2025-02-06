package rabbitmq

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultConnAttempts = 10
	defaultConnDelay    = 5 * time.Second
)

type RabbitMQ struct {
	url          string
	conn         *amqp.Connection
	connAttempts int
	connDelay    time.Duration
}

func New(url string, opts ...Option) (*RabbitMQ, error) {
	r := &RabbitMQ{
		url:          url,
		connAttempts: defaultConnAttempts,
		connDelay:    defaultConnDelay,
	}

	for _, opt := range opts {
		opt(r)
	}

	var err error
	for r.connAttempts > 0 {
		r.conn, err = amqp.Dial(r.url)
		if err == nil {
			break
		}

		log.Infof("RabbitMQ is trying to connect, attempts left: %d", r.connAttempts)
		time.Sleep(r.connDelay)
		r.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQ) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQ) Channel() (*amqp.Channel, error) {
	if r.conn != nil {
		return r.conn.Channel()
	}
	return nil, errors.New("RabbitMQ.Channel: connection is nil")
}
