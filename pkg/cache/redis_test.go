package cache

import (
	"testing"

	"ms-gofiber/internal/config"
)

func TestNewRedis(t *testing.T) {
	cfg := &config.Config{RedisAddr: "127.0.0.1:1", RedisDB: 0, RedisPassword: ""}
	c := NewRedis(cfg)
	if c == nil {
		t.Fatalf("expected redis client")
	}
	_ = c.Close()
}
