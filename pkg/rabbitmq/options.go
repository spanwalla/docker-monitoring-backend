package rabbitmq

import "time"

type Option func(*RabbitMQ)

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(r *RabbitMQ) {
		r.connAttempts = attempts
	}
}

// ConnDelay -.
func ConnDelay(delay time.Duration) Option {
	return func(r *RabbitMQ) {
		r.connDelay = delay
	}
}
