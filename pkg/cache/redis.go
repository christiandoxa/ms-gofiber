package cache

import (
	"github.com/redis/go-redis/v9"

	"ms-gofiber/internal/config"
	"ms-gofiber/pkg/apmredis"
)

func NewRedis(cfg *config.Config) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	// APM: trace semua command redis dengan context dari handler
	r.AddHook(apmredis.NewHook())
	return r
}
