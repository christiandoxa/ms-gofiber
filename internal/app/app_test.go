package app

import (
	"context"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"ms-gofiber/internal/app/adapter/controller"
	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/config"
)

func baseConfig(t *testing.T) *config.Config {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	return &config.Config{
		AppHost:            "127.0.0.1",
		AppPort:            18080,
		AppReadTimeout:     1,
		AppWriteTimeout:    1,
		SQLitePath:         filepath.Join(t.TempDir(), "db", "app.db"),
		RedisAddr:          mr.Addr(),
		RedisDB:            0,
		RedisPassword:      "",
		RedisDefaultTTL:    1,
		RedisPingTimeoutMs: 10,
	}
}

func TestBuildNilConfig(t *testing.T) {
	if _, _, err := Build(context.Background(), nil); err == nil {
		t.Fatalf("expected nil config error")
	}
}

func TestBuildSuccess(t *testing.T) {
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
}

func TestBuildInvalidSQLitePath(t *testing.T) {
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

func TestBuildInvalidRedisAddr(t *testing.T) {
	badRedis := baseConfig(t)
	badRedis.RedisAddr = "127.0.0.1:1"
	app3, closer3, err := Build(context.Background(), badRedis)
	if err == nil {
		if closer3 != nil {
			if closeErr := closer3(); closeErr != nil {
				t.Fatalf("close unexpected bad redis app: %v", closeErr)
			}
		}
		t.Fatalf("expected Build error for invalid redis addr, got app=%v", app3)
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

func TestReportTodoCacheError(t *testing.T) {
	reportTodoCacheError(context.Background(), "set", domain.TodoID("1"), errors.New("cache"))
}
