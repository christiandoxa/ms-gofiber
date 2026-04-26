package config

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadWithDefaultsAndOverrides(t *testing.T) {
	t.Setenv("APP_NAME", "app1")
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_HOST", "127.0.0.1")
	t.Setenv("APP_PORT", "9091")
	t.Setenv("APP_READ_TIMEOUT_SEC", "12")
	t.Setenv("APP_WRITE_TIMEOUT_SEC", "13")
	t.Setenv("SQLITE_PATH", "tmp/test.db")
	t.Setenv("REDIS_ADDR", "127.0.0.1:6380")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("REDIS_PASSWORD", "p")
	t.Setenv("REDIS_DEFAULT_TTL_SEC", "77")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if cfg.AppName != "app1" || cfg.AppEnv != "test" || cfg.AppHost != "127.0.0.1" || cfg.AppPort != 9091 {
		t.Fatalf("unexpected app config: %+v", cfg)
	}
	if cfg.SQLitePath != "tmp/test.db" {
		t.Fatalf("unexpected sqlite path: %s", cfg.SQLitePath)
	}
	if cfg.RedisAddr != "127.0.0.1:6380" || cfg.RedisDB != 2 || cfg.RedisPassword != "p" || cfg.RedisDefaultTTL != 77 {
		t.Fatalf("unexpected redis config: %+v", cfg)
	}
	if cfg.ListenAddr() != "127.0.0.1:9091" {
		t.Fatalf("unexpected listen addr: %s", cfg.ListenAddr())
	}
	if cfg.ReadTimeout() != 12*time.Second || cfg.WriteTimeout() != 13*time.Second {
		t.Fatalf("unexpected timeout read=%v write=%v", cfg.ReadTimeout(), cfg.WriteTimeout())
	}
	if cfg.RedisDefaultTTLDuration() != 77*time.Second {
		t.Fatalf("unexpected redis ttl: %v", cfg.RedisDefaultTTLDuration())
	}
}

func TestHelpersFallback(t *testing.T) {
	if err := os.Unsetenv("X_NOPE"); err != nil {
		t.Fatalf("unset env: %v", err)
	}
	if got := getenv("X_NOPE", "d"); got != "d" {
		t.Fatalf("expected default, got %s", got)
	}

	got, err := getenvInt("X_NOPE", 123)
	if err != nil || got != 123 {
		t.Fatalf("expected fallback int, got %d err=%v", got, err)
	}

	t.Setenv("X_BAD_INT", "nan")
	if _, err := getenvInt("X_BAD_INT", 123); err == nil {
		t.Fatalf("expected invalid integer error")
	}
}

func TestLoadOptionalDotenv(t *testing.T) {
	origStat := statEnvFile
	origLoad := loadEnvFile
	t.Cleanup(func() {
		statEnvFile = origStat
		loadEnvFile = origLoad
	})

	statEnvFile = func() (os.FileInfo, error) { return nil, os.ErrNotExist }
	if err := loadOptionalDotenv(); err != nil {
		t.Fatalf("missing .env should be optional: %v", err)
	}

	statEnvFile = func() (os.FileInfo, error) { return nil, errors.New("stat") }
	if err := loadOptionalDotenv(); err == nil || !strings.Contains(err.Error(), "stat .env") {
		t.Fatalf("expected stat error, got %v", err)
	}

	statEnvFile = func() (os.FileInfo, error) { return nil, nil }
	loadEnvFile = func(filenames ...string) error { return errors.New("load") }
	if err := loadOptionalDotenv(); err == nil || !strings.Contains(err.Error(), "load .env") {
		t.Fatalf("expected load error, got %v", err)
	}

	loadEnvFile = func(filenames ...string) error { return nil }
	if err := loadOptionalDotenv(); err != nil {
		t.Fatalf("expected load success, got %v", err)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	origStat := statEnvFile
	t.Cleanup(func() { statEnvFile = origStat })
	statEnvFile = func() (os.FileInfo, error) { return nil, errors.New("stat") }
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "stat .env") {
		t.Fatalf("expected dotenv stat error, got %v", err)
	}
	statEnvFile = origStat

	cases := []struct {
		key   string
		value string
	}{
		{"APP_PORT", "nan"},
		{"APP_READ_TIMEOUT_SEC", "nan"},
		{"APP_WRITE_TIMEOUT_SEC", "nan"},
		{"REDIS_DB", "nan"},
		{"REDIS_DEFAULT_TTL_SEC", "nan"},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)
			if _, err := Load(); err == nil || !strings.Contains(err.Error(), tc.key) {
				t.Fatalf("expected %s error, got %v", tc.key, err)
			}
		})
	}

	t.Setenv("APP_PORT", "0")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "APP_PORT") {
		t.Fatalf("expected APP_PORT range error, got %v", err)
	}
}

func TestValidateBranches(t *testing.T) {
	base := Config{
		AppName:         "app",
		AppHost:         "127.0.0.1",
		AppPort:         8080,
		AppReadTimeout:  1,
		AppWriteTimeout: 1,
		SQLitePath:      "data/app.db",
		RedisAddr:       "127.0.0.1:6379",
		RedisDB:         0,
		RedisDefaultTTL: 1,
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}

	cases := []struct {
		name   string
		mutate func(*Config)
		want   string
	}{
		{"app name", func(c *Config) { c.AppName = "" }, "APP_NAME"},
		{"app host", func(c *Config) { c.AppHost = "" }, "APP_HOST"},
		{"app port high", func(c *Config) { c.AppPort = 65536 }, "APP_PORT"},
		{"read timeout", func(c *Config) { c.AppReadTimeout = 0 }, "APP_READ_TIMEOUT_SEC"},
		{"write timeout", func(c *Config) { c.AppWriteTimeout = 0 }, "APP_WRITE_TIMEOUT_SEC"},
		{"sqlite path", func(c *Config) { c.SQLitePath = "" }, "SQLITE_PATH"},
		{"redis addr", func(c *Config) { c.RedisAddr = "" }, "REDIS_ADDR"},
		{"redis db", func(c *Config) { c.RedisDB = -1 }, "REDIS_DB"},
		{"redis ttl", func(c *Config) { c.RedisDefaultTTL = 0 }, "REDIS_DEFAULT_TTL_SEC"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := base
			tc.mutate(&cfg)
			if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %s validation error, got %v", tc.want, err)
			}
		})
	}
}
