package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TodoCache struct {
	client *redis.Client
}

func NewTodoCache(client *redis.Client) *TodoCache {
	return &TodoCache{client: client}
}

func (c *TodoCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *TodoCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *TodoCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
