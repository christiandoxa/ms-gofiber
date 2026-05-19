package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/responsecode/error"
)

func TestErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	app.Get("/app", func(c *fiber.Ctx) error { return apperror.New(http.StatusBadRequest, "bad request") })
	app.Get("/code", func(c *fiber.Ctx) error { return rcerror.ErrDuplicateExternalID })
	app.Get("/fiber", func(c *fiber.Ctx) error { return fiber.NewError(http.StatusUnauthorized, "unauthorized") })
	app.Get("/unknown", func(c *fiber.Ctx) error { return errors.New("unknown") })
	app.Get("/panic", func(c *fiber.Ctx) error {
		StackTraceHandler(c, "boom")
		return c.SendStatus(http.StatusOK)
	})
	app.Use(GeneralNotFound)

	assertError(t, app, "/app", http.StatusBadRequest, "bad request")
	assertResponseCode(t, app, "/code", http.StatusConflict, "4090001")
	assertError(t, app, "/fiber", http.StatusUnauthorized, "unauthorized")
	assertError(t, app, "/unknown", http.StatusInternalServerError, "internal server error")

	res, err := app.Test(httptest.NewRequest(http.MethodGet, "/panic", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}

	assertError(t, app, "/missing", http.StatusNotFound, "route not found")
}

func assertResponseCode(t *testing.T, app *fiber.App, path string, status int, responseCode string) {
	t.Helper()
	res, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.StatusCode != status {
		t.Fatalf("expected %d got %d", status, res.StatusCode)
	}
	body := map[string]any{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["responseCode"] != responseCode {
		t.Fatalf("expected %s got %+v", responseCode, body)
	}
}

func assertError(t *testing.T, app *fiber.App, path string, status int, message string) {
	t.Helper()
	res, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.StatusCode != status {
		t.Fatalf("expected %d got %d", status, res.StatusCode)
	}
	body := map[string]any{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["message"] != message {
		t.Fatalf("expected %s got %+v", message, body)
	}
}
