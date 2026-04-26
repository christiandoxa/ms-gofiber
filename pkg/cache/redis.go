package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"ms-gofiber/pkg/apmredis"
)

const defaultRedisPingTimeout = 5 * time.Second

type RedisOptions struct {
	Addr        string
	Password    string
	DB          int
	PingTimeout time.Duration
}

func NewRedis(ctx context.Context, opts RedisOptions) (*redis.Client, error) {
	pingTimeout := opts.PingTimeout
	if pingTimeout <= 0 {
		pingTimeout = defaultRedisPingTimeout
	}

	r := redis.NewClient(&redis.Options{
		Addr:         opts.Addr,
		Password:     opts.Password,
		DB:           opts.DB,
		DialTimeout:  pingTimeout,
		ReadTimeout:  pingTimeout,
		WriteTimeout: pingTimeout,
		MaxRetries:   0,
	})

	r.AddHook(apmredis.NewHook())

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := r.Ping(pingCtx).Err(); err != nil {
		return nil, errors.Join(err, r.Close())
	}
	return r, nil
}
