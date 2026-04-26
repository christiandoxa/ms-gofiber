package cache

import (
	"testing"
)

func TestNewRedis(t *testing.T) {
	c := NewRedis(RedisOptions{Addr: "127.0.0.1:1", DB: 0, Password: ""})
	if c == nil {
		t.Fatalf("expected redis client")
	}
	_ = c.Close()
}
