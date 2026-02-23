package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Todo is redis implementation of todo cache
type Todo struct {
	client *redis.Client
}

// NewTodo is constructor of redis todo cache
func NewTodo(client *redis.Client) *Todo {
	return &Todo{client: client}
}

func (c *Todo) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *Todo) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Todo) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
