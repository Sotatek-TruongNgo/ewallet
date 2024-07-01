package redis

import (
	"context"
	"fmt"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

const RedisNil = redisv8.Nil

type Redis interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	MSet(ctx context.Context, entries []Entry, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys []string) ([]interface{}, error)
	Del(ctx context.Context, key string) error
	Close()
}

type Entry struct {
	Key   string
	Value interface{}
}

type redis struct {
	client *redisv8.Client
}

func NewRedis(
	host string,
	password string,
	writeTimeout time.Duration,
	readTimeout time.Duration,
) (Redis, error) {
	redisClient := redisv8.NewClient(&redisv8.Options{
		Addr:         host,
		Password:     password,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("ping redis error: %w", err)
	}

	return &redis{client: redisClient}, nil
}

func (r *redis) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	_, err := r.client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		return err
	}
	return nil
}

// MSet set multiple key value pairs with ttl.
// Implemented with pipeline, since redis doesn't support MSet with ttl
// This method is not an atomic operation.
func (r *redis) MSet(
	ctx context.Context,
	entries []Entry,
	ttl time.Duration,
) error {
	pipe := r.client.Pipeline()
	for _, e := range entries {
		pipe.Set(ctx, e.Key, e.Value, ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) Get(
	ctx context.Context,
	key string,
) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// MGet get multiple keys. the result should be casting by caller.
func (r *redis) MGet(
	ctx context.Context,
	keys []string,
) ([]interface{}, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	return r.client.MGet(ctx, keys...).Result()
}

func (r *redis) Del(
	ctx context.Context,
	key string,
) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) Close() {
	r.client.Close()
}
