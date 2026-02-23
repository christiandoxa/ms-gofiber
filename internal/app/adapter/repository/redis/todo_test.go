package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
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

	if err := cache.Set(ctx, "k", []byte("v"), time.Minute); err != nil {
		t.Fatalf("set error: %v", err)
	}
	v, err := cache.Get(ctx, "k")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if string(v) != "v" {
		t.Fatalf("unexpected value: %s", v)
	}
	if err := cache.Delete(ctx, "k"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
}
