package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	res, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing X-Content-Type-Options")
	}
	if res.Header.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing X-Frame-Options")
	}
	if res.Header.Get("X-XSS-Protection") != "1; mode=block" {
		t.Fatalf("missing X-XSS-Protection")
	}
}
