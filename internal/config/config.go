package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppHost         string
	AppPort         int
	AppReadTimeout  int
	AppWriteTimeout int

	SQLitePath string

	RedisAddr       string
	RedisDB         int
	RedisPassword   string
	RedisDefaultTTL int
}

func Load() (*Config, error) {
	// Load .env bila ada; jika tidak ada, biarkan saja (pakai env OS)
	_ = godotenv.Load()

	cfg := &Config{
		AppName:         getenv("APP_NAME", "ms-gofiber"),
		AppEnv:          getenv("APP_ENV", "local"),
		AppHost:         getenv("APP_HOST", "0.0.0.0"),
		AppPort:         getint("APP_PORT", 8080),
		AppReadTimeout:  getint("APP_READ_TIMEOUT_SEC", 10),
		AppWriteTimeout: getint("APP_WRITE_TIMEOUT_SEC", 10),

		SQLitePath: getenv("SQLITE_PATH", "data/ms-gofiber.db"),

		RedisAddr:       getenv("REDIS_ADDR", "localhost:6379"),
		RedisDB:         getint("REDIS_DB", 0),
		RedisPassword:   getenv("REDIS_PASSWORD", ""),
		RedisDefaultTTL: getint("REDIS_DEFAULT_TTL_SEC", 60),
	}
	return cfg, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func getint(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.AppHost, c.AppPort)
}
func (c *Config) ReadTimeout() time.Duration  { return time.Duration(c.AppReadTimeout) * time.Second }
func (c *Config) WriteTimeout() time.Duration { return time.Duration(c.AppWriteTimeout) * time.Second }
