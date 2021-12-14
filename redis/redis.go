package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type DB struct {
	redisClient *redis.Client
}

func NewDB(addr string) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &DB{
		redisClient: client,
	}, nil
}