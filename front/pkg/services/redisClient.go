package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// InitRedis initializes the Redis client
func InitRedis() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	defer client.Close()

	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (r *RedisClient) Set(key, value string) error {
	err := r.client.Set(r.ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

type redisClientKey struct{}

func WithRedisClient(ctx context.Context, client *RedisClient) context.Context {
	return context.WithValue(ctx, redisClientKey{}, client)
}
