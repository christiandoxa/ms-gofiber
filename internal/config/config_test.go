package config

import (
	"os"
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
}

func TestHelpersFallback(t *testing.T) {
	_ = os.Unsetenv("X_NOPE")
	if got := getenv("X_NOPE", "d"); got != "d" {
		t.Fatalf("expected default, got %s", got)
	}

	t.Setenv("X_BAD_INT", "nan")
	if got := getint("X_BAD_INT", 123); got != 123 {
		t.Fatalf("expected fallback int, got %d", got)
	}
}
