package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
)

type mockCacheService struct {
	err error
}

func (m mockCacheService) Flush(context.Context) error {
	return m.err
}

func TestFlush(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Get("/flush", Flush(&model.Service{CacheService: mockCacheService{}}))
	assertFlushStatus(t, app, http.StatusOK)

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Get("/flush", Flush(&model.Service{CacheService: mockCacheService{err: errors.New("flush")}}))
	assertFlushStatus(t, app, http.StatusInternalServerError)
}

func assertFlushStatus(t *testing.T, app *fiber.App, expected int) {
	t.Helper()
	res, err := app.Test(httptest.NewRequest(http.MethodGet, "/flush", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.StatusCode != expected {
		t.Fatalf("expected %d got %d", expected, res.StatusCode)
	}
}
