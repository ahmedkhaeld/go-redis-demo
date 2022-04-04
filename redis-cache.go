package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// RedisCache holds a pointer to redis client
type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache() *RedisCache {
	var opts *redis.Options

	if os.Getenv("LOCAL") == "true" {
		redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))
		opts = &redis.Options{
			Addr:     redisAddress,
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	} else {
		builtOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}
		opts = builtOpts
	}

	rdb := redis.NewClient(opts)

	return &RedisCache{
		cache: rdb,
	}
}
