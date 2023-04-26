package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	instance *redis.Client
}

func NewRedis() (*Redis, error) {
	return &Redis{
		instance: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}, nil
}

func (r *Redis) Set(ctx context.Context, key string) {
	err := r.instance.Set(ctx, "key", true, 0).Err()
	if err != nil {
		panic(err)
	}
}
