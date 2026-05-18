package handler

import (
	"testing"

	"github.com/christiandoxa/welog"
	"github.com/gofiber/fiber/v2"

	mw "ms-gofiber/api/middleware"
	"ms-gofiber/pkg/apperror"
)

func setupValidationApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(welog.NewFiber(fiber.Config{ErrorHandler: mw.ErrorHandler()}))
	return app
}

func TestValidationControllerPrepareExample(t *testing.T) {
	// invalid body
	app := setupValidationApp()
	h := NewValidation(mockValidator{})
	app.Post("/v", h.PrepareExample)
	assertStatus(t, app, jsonRequest("POST", "/v", "{"), 400)

	// validator error
	app = setupValidationApp()
	h = NewValidation(mockValidator{err: apperror.New(apperror.ErrValidation, "invalid")})
	app.Post("/v", h.PrepareExample)
	req := jsonRequest("POST", "/v", `{"terminalType":"APP"}`)
	assertStatus(t, app, req, 400)

	// success
	app = setupValidationApp()
	h = NewValidation(mockValidator{})
	app.Post("/v", h.PrepareExample)
	req = jsonRequest("POST", "/v", `{"terminalType":"APP","osType":"ANDROID","osVersion":"14","grantType":"AUTHORIZATION_CODE","paymentMethodType":"DANA","scope":["SEND_OTP"],"transactionTime":"2026-02-23T10:30:00Z","merchantName":"Demo Merchant"}`)
	assertStatus(t, app, req, 200)
}
