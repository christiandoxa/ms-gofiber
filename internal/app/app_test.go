package app

import (
	"context"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alicebob/miniredis/v2"

	apivalidation "ms-gofiber/api/validation"
	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/config"
	appvalidator "ms-gofiber/internal/validator"
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

func TestBuildSuccess(t *testing.T) {
	cfg := baseConfig(t)
	runtime, err := Build(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Build success expected, got err: %v", err)
	}
	if runtime == nil {
		t.Fatalf("expected runtime not nil")
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatalf("close app: %v", err)
		}
	}()

	res, err := runtime.Test(httptest.NewRequest("GET", "/v1/health", nil))
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
	runtime, err := Build(context.Background(), bad)
	if err == nil {
		t.Fatalf("expected Build error for invalid sqlite path, got runtime=%v", runtime)
	}
}

func TestBuildInvalidRedisAddr(t *testing.T) {
	badRedis := baseConfig(t)
	badRedis.RedisAddr = "127.0.0.1:1"
	runtime, err := Build(context.Background(), badRedis)
	if err == nil {
		t.Fatalf("expected Build error for invalid redis addr, got runtime=%v", runtime)
	}
}

func TestBuildValidatorError(t *testing.T) {
	patches := gomonkey.ApplyFunc(apivalidation.NewStructValidator, func() (*appvalidator.StructValidator, error) {
		return nil, errors.New("validator")
	})
	defer patches.Reset()

	runtime, err := Build(context.Background(), baseConfig(t))
	if err == nil || !strings.Contains(err.Error(), "validator") {
		t.Fatalf("expected validator error, got %v", err)
	}
	if runtime != nil {
		t.Fatalf("expected nil runtime on validator error")
	}
}

func TestReportTodoCacheError(t *testing.T) {
	reportTodoCacheError(context.Background(), "set", domain.TodoID("1"), errors.New("cache"))
}
