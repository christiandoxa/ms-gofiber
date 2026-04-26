package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"ms-gofiber/internal/app/domain"
)

var marshalTodo = json.Marshal

type Todo struct {
	client *redis.Client
}

func NewTodo(client *redis.Client) *Todo {
	return &Todo{client: client}
}

func (c *Todo) GetTodo(ctx context.Context, id domain.TodoID) (*domain.Todo, bool, error) {
	value, err := c.client.Get(ctx, todoKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var todo domain.Todo
	if err := json.Unmarshal(value, &todo); err != nil {
		return nil, false, err
	}
	return &todo, true, nil
}

func (c *Todo) SetTodo(ctx context.Context, todo *domain.Todo, ttl time.Duration) error {
	value, err := marshalTodo(todo)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, todoKey(todo.ID), value, ttl).Err()
}

func (c *Todo) DeleteTodo(ctx context.Context, id domain.TodoID) error {
	return c.client.Del(ctx, todoKey(id)).Err()
}

func todoKey(id domain.TodoID) string {
	return fmt.Sprintf("todo:%s", id)
}
