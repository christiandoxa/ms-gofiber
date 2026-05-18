package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing content type header")
	}
	if res.Header.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing frame header")
	}
	if res.Header.Get("X-XSS-Protection") != "1; mode=block" {
		t.Fatalf("missing xss header")
	}
}
