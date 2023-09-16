package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewAMQPConnection(config Configuration) (*amqp.Connection, error) {
	cfg := amqp.Config{Properties: amqp.NewConnectionProperties()}
	cfg.Properties.SetClientConnectionName(config.ConsumerName)

	return amqp.DialConfig(config.Dsn(), cfg)
}
