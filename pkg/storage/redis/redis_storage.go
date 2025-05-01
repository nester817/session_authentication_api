package storageRedis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DbRedis struct {
	client *redis.Client
}

func NewDbRedis(addr string) *DbRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &DbRedis{
		client: client,
	}
}

func (db *DbRedis) Get(ctx context.Context, key string) (string, error) {
	return db.client.Get(ctx, key).Result()
}

func (db *DbRedis) Set(ctx context.Context, key, value string) error {
	return db.client.Set(ctx, key, value, 24*24*time.Hour).Err()
}

func (db *DbRedis) Delete(ctx context.Context, key string) error {
	return db.client.Del(ctx, key).Err()
}
