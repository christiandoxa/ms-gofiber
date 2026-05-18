package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRequestID(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"header": c.GetRespHeader("X-Request-ID"),
			"local":  c.Locals("request_id"),
		})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "custom-id")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if body["header"] != "custom-id" || body["local"] != "custom-id" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestRequestIDGeneratedWhenMissing(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"header": c.GetRespHeader("X-Request-ID"),
			"local":  c.Locals("request_id"),
		})
	})

	res, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	header, ok := body["header"].(string)
	if !ok {
		t.Fatalf("expected string header: %+v", body)
	}
	local, ok := body["local"].(string)
	if !ok {
		t.Fatalf("expected string local: %+v", body)
	}
	if header == "" || local == "" {
		t.Fatalf("expected generated request id: %+v", body)
	}
	if header != local {
		t.Fatalf("header and local must match: %+v", body)
	}
}
