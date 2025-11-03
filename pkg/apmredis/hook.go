package apmredis

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.elastic.co/apm/v2"
)

type hook struct{}

func NewHook() redis.Hook { return hook{} }

func (hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		span, ctx := apm.StartSpan(ctx, "redis.dial", "cache")
		span.Subtype = "redis"
		span.Action = "dial"
		defer span.End()
		return next(ctx, network, addr)
	}
}

func (hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		name := strings.ToLower(cmd.FullName())
		span, ctx := apm.StartSpan(ctx, "redis."+name, "cache")
		span.Subtype = "redis"
		span.Action = name
		start := time.Now()
		err := next(ctx, cmd)
		span.Duration = time.Since(start)
		span.End()
		if err != nil {
			apm.CaptureError(ctx, err).Send()
		}
		return err
	}
}

func (hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		span, ctx := apm.StartSpan(ctx, "redis.pipeline", "cache")
		span.Subtype = "redis"
		span.Action = "pipeline"
		start := time.Now()
		err := next(ctx, cmds)
		span.Duration = time.Since(start)
		span.End()
		if err != nil {
			apm.CaptureError(ctx, err).Send()
		}
		return err
	}
}
