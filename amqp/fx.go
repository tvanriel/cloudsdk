package amqp

import "go.uber.org/fx"

var Module = fx.Module("amqp", fx.Provide(NewAMQPConnection))

