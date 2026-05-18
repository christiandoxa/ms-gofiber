package config

import (
	"fmt"
	"time"

	"ms-gofiber/pkg/tool"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppHost         string
	AppPort         int
	AppReadTimeout  int
	AppWriteTimeout int

	SQLitePath string

	RedisAddr          string
	RedisDB            int
	RedisPassword      string
	RedisDefaultTTL    int
	RedisPingTimeoutMs int
}

func Load() (*Config, error) {
	if err := tool.LoadDotenvIfExists(".env"); err != nil {
		return nil, err
	}

	appPort, err := tool.IntFromEnv("APP_PORT", 8080)
	if err != nil {
		return nil, err
	}
	readTimeout, err := tool.IntFromEnv("APP_READ_TIMEOUT_SEC", 10)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := tool.IntFromEnv("APP_WRITE_TIMEOUT_SEC", 10)
	if err != nil {
		return nil, err
	}
	redisDB, err := tool.IntFromEnv("REDIS_DB", 0)
	if err != nil {
		return nil, err
	}
	redisTTL, err := tool.IntFromEnv("REDIS_DEFAULT_TTL_SEC", 60)
	if err != nil {
		return nil, err
	}
	redisPingTimeout, err := tool.IntFromEnv("REDIS_PING_TIMEOUT_MS", 5000)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppName:         tool.StringFromEnv("APP_NAME", "ms-gofiber"),
		AppEnv:          tool.StringFromEnv("APP_ENV", "local"),
		AppHost:         tool.StringFromEnv("APP_HOST", "0.0.0.0"),
		AppPort:         appPort,
		AppReadTimeout:  readTimeout,
		AppWriteTimeout: writeTimeout,

		SQLitePath: tool.StringFromEnv("SQLITE_PATH", "data/ms-gofiber.db"),

		RedisAddr:          tool.StringFromEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:            redisDB,
		RedisPassword:      tool.StringFromEnv("REDIS_PASSWORD", ""),
		RedisDefaultTTL:    redisTTL,
		RedisPingTimeoutMs: redisPingTimeout,
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if c.AppHost == "" {
		return fmt.Errorf("APP_HOST is required")
	}
	if c.AppPort <= 0 || c.AppPort > 65535 {
		return fmt.Errorf("APP_PORT must be between 1 and 65535")
	}
	if c.AppReadTimeout <= 0 {
		return fmt.Errorf("APP_READ_TIMEOUT_SEC must be positive")
	}
	if c.AppWriteTimeout <= 0 {
		return fmt.Errorf("APP_WRITE_TIMEOUT_SEC must be positive")
	}
	if c.SQLitePath == "" {
		return fmt.Errorf("SQLITE_PATH is required")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}
	if c.RedisDB < 0 {
		return fmt.Errorf("REDIS_DB must be zero or positive")
	}
	if c.RedisDefaultTTL <= 0 {
		return fmt.Errorf("REDIS_DEFAULT_TTL_SEC must be positive")
	}
	if c.RedisPingTimeoutMs <= 0 {
		return fmt.Errorf("REDIS_PING_TIMEOUT_MS must be positive")
	}
	return nil
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.AppHost, c.AppPort)
}

func (c *Config) ReadTimeout() time.Duration {
	return time.Duration(c.AppReadTimeout) * time.Second
}

func (c *Config) WriteTimeout() time.Duration {
	return time.Duration(c.AppWriteTimeout) * time.Second
}

func (c *Config) RedisDefaultTTLDuration() time.Duration {
	return time.Duration(c.RedisDefaultTTL) * time.Second
}

func (c *Config) RedisPingTimeout() time.Duration {
	return time.Duration(c.RedisPingTimeoutMs) * time.Millisecond
}
