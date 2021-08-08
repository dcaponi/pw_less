package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Cache interface {
	Get(key string) (string, error)
	Set(key string, val interface{}) error
	Setex(key string, val interface{}, ttl time.Duration) error
	Del(key string) error
}

type RedisCacheConfig struct {
	Addr, Password string
	DB             int
	Client         *redis.Client
}

func NewRedisCache(c RedisCacheConfig) (Cache, error) {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	_, err := c.Client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	log.Println("established connection to redis!")
	return c, nil
}

func (c RedisCacheConfig) Get(key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c RedisCacheConfig) Set(key string, val interface{}) error {
	return c.Client.Set(ctx, key, val, 0).Err()
}

func (c RedisCacheConfig) Setex(key string, val interface{}, ttl time.Duration) error {
	return c.Client.Set(ctx, key, val, ttl).Err()
}

func (c RedisCacheConfig) Del(key string) error {
	return c.Client.Del(ctx, key).Err()
}
