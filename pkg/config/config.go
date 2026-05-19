package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"

	"ms-gofiber/pkg/constant/envkey"
	"ms-gofiber/pkg/tool"
)

type Config struct {
	AppPort      string
	DatabasePath string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		AppPort:      tool.StringFromEnv(envkey.AppPort, "8080"),
		DatabasePath: tool.StringFromEnv(envkey.DatabasePath, ":memory:"),
		ReadTimeout:  tool.DurationFromEnv(envkey.AppReadTimeout, 10*time.Second),
		WriteTimeout: tool.DurationFromEnv(envkey.AppWriteTimeout, 10*time.Second),
	}
	setDefaultEnv(envkey.AppPort, cfg.AppPort)
	return cfg
}

func setDefaultEnv(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
