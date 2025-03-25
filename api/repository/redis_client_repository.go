package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}
// NewRedisClient cria uma nova instância do cliente Redis.
func NewRedisClient(addr, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("Fala ao conectar ao Redis: " + err.Error())
	}
	return &RedisClient{Client: client}
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}