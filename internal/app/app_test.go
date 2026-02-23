package app

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"ms-gofiber/internal/config"
)

func baseConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		AppHost:         "127.0.0.1",
		AppPort:         18080,
		AppReadTimeout:  1,
		AppWriteTimeout: 1,
		SQLitePath:      filepath.Join(t.TempDir(), "db", "app.db"),
		RedisAddr:       "127.0.0.1:1",
		RedisDB:         0,
		RedisPassword:   "",
		RedisDefaultTTL: 1,
	}
}

func TestBuildSuccessAndError(t *testing.T) {
	cfg := baseConfig(t)
	app, closer, err := Build(cfg)
	if err != nil {
		t.Fatalf("Build success expected, got err: %v", err)
	}
	if closer == nil {
		t.Fatalf("expected closer not nil")
	}
	defer closer()

	res, err := app.Test(httptest.NewRequest("GET", "/v1/health", nil))
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	bad := baseConfig(t)
	bad.SQLitePath = "/proc/1/forbidden/db.sqlite"
	app2, closer2, err := Build(bad)
	if err == nil {
		if closer2 != nil {
			closer2()
		}
		t.Fatalf("expected Build error for invalid sqlite path, got app=%v", app2)
	}
}
