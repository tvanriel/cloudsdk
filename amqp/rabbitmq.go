package amqp

import (
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Connection struct {
	config Configuration
	Amqp   *amqp.Connection
	Log    *zap.Logger
}

type Channel struct {
	*amqp.Channel
	closed int32
}

const delay = 3

func NewAMQPConnection(config Configuration, l *zap.Logger) (*Connection, error) {
	conn, err := amqp.Dial(config.Dsn())
	if err != nil {
		return nil, err
	}
	return &Connection{
		config: config,
		Amqp:   conn,
		Log:    l,
	}, nil
}

func (a *Connection) Reconnect() error {
	cfg := amqp.Config{Properties: amqp.NewConnectionProperties()}
	cfg.Properties.SetClientConnectionName(a.config.ConsumerName)
	var err error
	a.Amqp, err = amqp.DialConfig(a.config.Dsn(), cfg)
	if err != nil {
		return err
	}

	go func() {
		for {
			reason, ok := <-a.Amqp.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				break
			}
			a.Log.Warn("Reconnecting to AMQP queue", zap.String("reason", reason.Reason))

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(a.config.Dsn())
				if err == nil {
					a.Amqp = conn
					break
				}
				a.Log.Warn("Could not reconnect to AMQP queue", zap.Error(err))
			}
		}
	}()

	return nil
}

func (a *Connection) Channel() (*Channel, error) {
	ch, err := a.Amqp.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			a.Log.Warn("channel closed",
				zap.String("reason", reason.Reason),
			)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(delay * time.Second)

				ch, err := a.Amqp.Channel()
				if err == nil {
					channel.Channel = ch
					break
				}
			}
		}
	}()

	return channel, nil
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return (atomic.LoadInt32(&ch.closed) == 1)
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume wrap amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				time.Sleep(delay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
