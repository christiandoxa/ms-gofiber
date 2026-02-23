package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	mw "ms-gofiber/internal/middleware"
)

func setupSelfApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})
	return app
}

func TestInternalEcho(t *testing.T) {
	app := setupSelfApp()
	h := NewInternal()
	app.Get("/echo", h.Echo)

	res, err := app.Test(httptest.NewRequest("GET", "/echo?msg=hello", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if body["echo"] != "hello" {
		t.Fatalf("unexpected echo response: %+v", body)
	}
}

func TestInternalSelfCallBranches(t *testing.T) {
	orig := httpDo
	t.Cleanup(func() { httpDo = orig })

	// error path
	httpDo = func(*fiber.Ctx, string, string, string, map[string]string, []byte, time.Duration) (int, []byte, http.Header, error) {
		return 0, nil, nil, errors.New("boom")
	}
	app := setupSelfApp()
	h := NewInternal()
	app.Get("/self", h.SelfCall)
	res, err := app.Test(httptest.NewRequest("GET", "/self", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}

	// success path
	httpDo = func(*fiber.Ctx, string, string, string, map[string]string, []byte, time.Duration) (int, []byte, http.Header, error) {
		return 200, []byte(`{"echo":"ok"}`), http.Header{}, nil
	}
	res, err = app.Test(httptest.NewRequest("GET", "/self", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if body["code"] != "OK" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
