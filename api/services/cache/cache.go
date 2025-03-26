package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(addr, password string, db int) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &RedisService{Client: client}, nil
}

func (r *RedisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisService) SetJSON(ctx context.Context, key string, value interface{}, expiration int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, data, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisService) Close() error {
	return r.Client.Close()
}
