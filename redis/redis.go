package redis

import (
	"context"

	gredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type NewRedisClientParams struct {
	fx.In

	Configuration Configuration
}

func NewRedisClient(p NewRedisClientParams) (*RedisClient, error) {
	conn := gredis.NewClient(&gredis.Options{
		Addr:       p.Configuration.Address,
		Password:   p.Configuration.Password,
		DB:         p.Configuration.DatabaseIndex,
		MaxRetries: 0,
	})
	err := conn.Ping(context.Background()).Err()

	return &RedisClient{
		conn: conn,
	}, err

}

type RedisClient struct {
	conn *gredis.Client
}

func (r *RedisClient) Conn() *gredis.Client {
	return r.conn
}
