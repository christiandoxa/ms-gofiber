package server

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewServer(t *testing.T) {
	app := NewServer()

	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}

	req := httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-CLIENT-ID", "client")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}

	req = httptest.NewRequest(fiber.MethodGet, "/missing", nil)
	req.Header.Set("X-CLIENT-ID", "client")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected 404 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}
}
