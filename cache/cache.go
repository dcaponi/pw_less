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
}

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(c RedisCacheConfig) (RedisCache, error) {
	cache := RedisCache{
		Client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.DB,
		}),
	}
	_, err := cache.Client.Ping(ctx).Result()
	if err != nil {
		return cache, err
	}
	log.Println("established connection to redis!")
	return cache, nil
}

func (c RedisCache) Get(key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c RedisCache) Set(key string, val interface{}) error {
	return c.Client.Set(ctx, key, val, 0).Err()
}

func (c RedisCache) Setex(key string, val interface{}, ttl time.Duration) error {
	return c.Client.Set(ctx, key, val, ttl).Err()
}

func (c RedisCache) Del(key string) error {
	return c.Client.Del(ctx, key).Err()
}
