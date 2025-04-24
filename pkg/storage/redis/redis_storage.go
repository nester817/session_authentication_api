package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientRedis struct {
	client *redis.Client
}

func NewClientRedis(addr string) *ClientRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &ClientRedis{
		client: client,
	}
}

func (cr *ClientRedis) Get(ctx context.Context, key string) (string, error) {
	value, err := cr.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, err
}

func (cr *ClientRedis) Set(ctx context.Context, key, value string) error {
	err := cr.client.Set(ctx, key, value, 24*time.Hour).Err()
	return err
}
