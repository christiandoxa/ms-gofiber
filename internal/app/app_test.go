package app

import (
	"context"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"ms-gofiber/internal/app/adapter/controller"
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
	if _, _, err := Build(context.Background(), nil); err == nil {
		t.Fatalf("expected nil config error")
	}

	cfg := baseConfig(t)
	app, closer, err := Build(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Build success expected, got err: %v", err)
	}
	if closer == nil {
		t.Fatalf("expected closer not nil")
	}
	defer func() {
		if err := closer(); err != nil {
			t.Fatalf("close app: %v", err)
		}
	}()

	res, err := app.Test(httptest.NewRequest("GET", "/v1/health", nil))
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	bad := baseConfig(t)
	bad.SQLitePath = "/proc/1/forbidden/db.sqlite"
	app2, closer2, err := Build(context.Background(), bad)
	if err == nil {
		if closer2 != nil {
			if closeErr := closer2(); closeErr != nil {
				t.Fatalf("close unexpected bad app: %v", closeErr)
			}
		}
		t.Fatalf("expected Build error for invalid sqlite path, got app=%v", app2)
	}
}

func TestBuildValidatorError(t *testing.T) {
	orig := newValidator
	t.Cleanup(func() { newValidator = orig })

	newValidator = func() (controller.RequestValidator, error) {
		return nil, errors.New("validator")
	}

	_, closer, err := Build(context.Background(), baseConfig(t))
	if err == nil || !strings.Contains(err.Error(), "validator") {
		t.Fatalf("expected validator error, got %v", err)
	}
	if closer != nil {
		t.Fatalf("expected nil closer on validator error")
	}
}
