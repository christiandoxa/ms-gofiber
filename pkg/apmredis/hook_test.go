package apmredis

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestHookWrappers(t *testing.T) {
	h := NewHook()
	ctx := context.Background()

	dial := h.DialHook(func(ctx context.Context, network, addr string) (net.Conn, error) {
		c1, c2 := net.Pipe()
		_ = c2.Close()
		return c1, nil
	})
	conn, err := dial(ctx, "tcp", "x")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	_ = conn.Close()

	processOK := h.ProcessHook(func(context.Context, redis.Cmder) error { return nil })
	if err := processOK(ctx, redis.NewStringCmd(ctx, "GET", "k")); err != nil {
		t.Fatalf("process ok error: %v", err)
	}

	processErr := h.ProcessHook(func(context.Context, redis.Cmder) error { return errors.New("x") })
	if err := processErr(ctx, redis.NewStringCmd(ctx, "GET", "k")); err == nil {
		t.Fatalf("expected process error")
	}

	pipeOK := h.ProcessPipelineHook(func(context.Context, []redis.Cmder) error { return nil })
	if err := pipeOK(ctx, []redis.Cmder{redis.NewStringCmd(ctx, "GET", "k")}); err != nil {
		t.Fatalf("pipeline ok error: %v", err)
	}

	pipeErr := h.ProcessPipelineHook(func(context.Context, []redis.Cmder) error { return errors.New("x") })
	if err := pipeErr(ctx, []redis.Cmder{redis.NewStringCmd(ctx, "GET", "k")}); err == nil {
		t.Fatalf("expected pipeline error")
	}
}
