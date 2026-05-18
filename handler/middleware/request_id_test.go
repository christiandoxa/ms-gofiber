package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/constant/generalkey"
)

func TestRequestID(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		if c.Locals(generalkey.RequestID) != "rid" {
			t.Fatalf("unexpected request id locals")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set(generalkey.RequestIDHeader, "rid")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.Header.Get(generalkey.RequestIDHeader) != "rid" {
		t.Fatalf("unexpected response request id")
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}

	app = fiber.New()
	app.Use(RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		if c.Locals(generalkey.RequestID) == "" {
			t.Fatalf("expected generated request id")
		}
		return c.SendStatus(fiber.StatusOK)
	})
	res, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.Header.Get(generalkey.RequestIDHeader) == "" {
		t.Fatalf("expected generated response request id")
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}
}
