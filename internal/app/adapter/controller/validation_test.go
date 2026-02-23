package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	mw "ms-gofiber/internal/middleware"
	"ms-gofiber/pkg/apperror"
)

func setupValidationApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})
	return app
}

func TestValidationControllerPrepareExample(t *testing.T) {
	// invalid body
	app := setupValidationApp()
	h := NewValidation(mockValidator{})
	app.Post("/v", h.PrepareExample)
	res, err := app.Test(httptest.NewRequest("POST", "/v", strings.NewReader("{")))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// validator error
	app = setupValidationApp()
	h = NewValidation(mockValidator{err: apperror.New(apperror.ErrValidation, "invalid")})
	app.Post("/v", h.PrepareExample)
	req := httptest.NewRequest("POST", "/v", strings.NewReader(`{"terminalType":"APP"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// success
	app = setupValidationApp()
	h = NewValidation(mockValidator{})
	app.Post("/v", h.PrepareExample)
	req = httptest.NewRequest("POST", "/v", strings.NewReader(`{"terminalType":"APP","osType":"ANDROID","osVersion":"14","grantType":"AUTHORIZATION_CODE","paymentMethodType":"DANA","scope":["SEND_OTP"],"transactionTime":"2026-02-23T10:30:00Z","merchantName":"Demo Merchant"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}
