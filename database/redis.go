package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"nymphicus-service/config"
)

func NewRedisClient(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.RedisAddr,
		Password: c.Redis.Password,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return rdb, nil
}
