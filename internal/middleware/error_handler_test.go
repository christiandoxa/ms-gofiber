package middleware

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"ms-gofiber/pkg/apperror"
)

func setupErrorApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})
	return app
}

func TestErrorHandlerBranches(t *testing.T) {
	app := setupErrorApp()
	app.Get("/aerr500", func(c *fiber.Ctx) error { return apperror.New(apperror.ErrInternal, "boom") })
	app.Get("/aerr400", func(c *fiber.Ctx) error { return apperror.New(apperror.ErrBadRequest, "bad") })
	app.Get("/fiber500", func(c *fiber.Ctx) error { return fiber.NewError(fiber.StatusInternalServerError, "f500") })
	app.Get("/fiber400", func(c *fiber.Ctx) error { return fiber.NewError(fiber.StatusBadRequest, "f400") })
	app.Get("/unknown", func(c *fiber.Ctx) error { return errors.New("x") })

	cases := []struct {
		path   string
		status int
		code   string
	}{
		{"/aerr500", 500, string(apperror.ErrInternal)},
		{"/aerr400", 400, string(apperror.ErrBadRequest)},
		{"/fiber500", 500, string(apperror.ErrBadRequest)},
		{"/fiber400", 400, string(apperror.ErrBadRequest)},
		{"/unknown", 500, string(apperror.ErrInternal)},
	}

	for _, tc := range cases {
		res, err := app.Test(httptest.NewRequest("GET", tc.path, nil))
		if err != nil {
			t.Fatalf("request %s failed: %v", tc.path, err)
		}
		defer res.Body.Close()
		if res.StatusCode != tc.status {
			t.Fatalf("path %s expected status %d got %d", tc.path, tc.status, res.StatusCode)
		}
		var body map[string]any
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("decode %s failed: %v", tc.path, err)
		}
		if body["code"] != tc.code {
			t.Fatalf("path %s expected code %s got %v", tc.path, tc.code, body["code"])
		}
	}
}
