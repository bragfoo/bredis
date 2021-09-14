package bredis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type realRedisBRedis struct {
	rdb *redis.Client
}

func (r *realRedisBRedis) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	ctx := context.Background()

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		return "", err
	}
	return val, nil
}

func (r *realRedisBRedis) Set(key string, val string) error {
	if key == "" {
		return ErrEmptyKey
	}
	ctx := context.Background()
	err := r.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRealRedisBRedis() BRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &realRedisBRedis{
		rdb: rdb,
	}
}
