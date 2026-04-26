package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"ms-gofiber/internal/app/domain"
)

func TestTodoCacheCRUD(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis run error: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cache := NewTodo(client)
	ctx := context.Background()

	todo := &domain.Todo{ID: "1", Title: "cached"}
	if err := cache.SetTodo(ctx, todo, time.Minute); err != nil {
		t.Fatalf("set error: %v", err)
	}
	got, found, err := cache.GetTodo(ctx, "1")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if !found || got.ID != "1" || got.Title != "cached" {
		t.Fatalf("unexpected value: found=%v todo=%+v", found, got)
	}
	if err := cache.DeleteTodo(ctx, "1"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	_, found, err = cache.GetTodo(ctx, "1")
	if err != nil || found {
		t.Fatalf("expected cache miss after delete, found=%v err=%v", found, err)
	}

	if err := mr.Set(todoKey("bad"), "{bad"); err != nil {
		t.Fatalf("set invalid payload: %v", err)
	}
	if _, _, err := cache.GetTodo(ctx, "bad"); err == nil {
		t.Fatalf("expected invalid cache payload error")
	}
}

func TestTodoCacheErrorBranches(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	cache := NewTodo(client)
	ctx := context.Background()

	if err := client.Close(); err != nil {
		t.Fatalf("close redis client: %v", err)
	}
	if _, _, err := cache.GetTodo(ctx, "closed"); err == nil {
		t.Fatalf("expected closed client get error")
	}

	expected := errors.New("marshal error")
	patches := gomonkey.ApplyGlobalVar(&marshalTodo, func(any) ([]byte, error) {
		return nil, expected
	})
	defer patches.Reset()

	if err := cache.SetTodo(ctx, &domain.Todo{ID: "broken"}, time.Minute); !errors.Is(err, expected) {
		t.Fatalf("expected marshal error, got %v", err)
	}
}
