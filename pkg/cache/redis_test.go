package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestNewRedis(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	defer mr.Close()

	c, err := NewRedis(context.Background(), RedisOptions{Addr: mr.Addr(), DB: 0, Password: ""})
	if err != nil {
		t.Fatalf("new redis client: %v", err)
	}
	if c == nil {
		t.Fatalf("expected redis client")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close redis client: %v", err)
	}

	_, err = NewRedis(context.Background(), RedisOptions{
		Addr:        "bad address",
		PingTimeout: time.Millisecond,
	})
	if err == nil {
		t.Fatalf("expected redis ping error")
	}
}
