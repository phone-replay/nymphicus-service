package database

import (
	"github.com/go-redis/redis/v8"
	"nymphicus-service/config"
)

func NewRedisClient(c *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.RedisAddr,
		Password: c.Redis.Password,
	})
	return rdb
}
