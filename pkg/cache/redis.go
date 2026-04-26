package cache

import (
	"github.com/redis/go-redis/v9"

	"ms-gofiber/pkg/apmredis"
)

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

func NewRedis(opts RedisOptions) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	r.AddHook(apmredis.NewHook())
	return r
}
