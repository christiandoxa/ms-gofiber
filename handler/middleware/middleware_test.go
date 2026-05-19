package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/constant/generalkey"
)

type mockValidator struct {
	err error
}

func (m mockValidator) ValidateStruct(any) error {
	return m.err
}

type mockExternalIDService struct {
	err error
}

func (m mockExternalIDService) StoreExternalID(context.Context, string) error {
	return m.err
}

func TestCheckHeader(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(CheckHeader(&model.Service{RequestValidator: mockValidator{}}))
	app.Get("/v1/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	var fiberCtx *fiber.Ctx
	patches := gomonkey.ApplyMethod(fiberCtx, "ReqHeaderParser", func(*fiber.Ctx, any) error {
		return errors.New("parse")
	})
	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(CheckHeader(&model.Service{RequestValidator: mockValidator{}}))
	app.Get("/v1/todos", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	res, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}
	patches.Reset()

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(CheckHeader(&model.Service{
		RequestValidator: mockValidator{err: apperror.New(fiber.StatusBadRequest, "validation failed")},
	}))
	app.Get("/v1/todos", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	req := httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-CLIENT-ID", "client")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(CheckHeader(&model.Service{RequestValidator: mockValidator{}}))
	app.Get("/v1/todos", func(c *fiber.Ctx) error {
		header, ok := c.Locals(generalkey.RequestHeader).(model.Header)
		if !ok || header.ClientID != "client" || header.ExternalID != "external1" {
			t.Fatalf("unexpected header: %+v %v", header, ok)
		}
		return c.SendStatus(fiber.StatusOK)
	})
	req = httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-CLIENT-ID", "client")
	req.Header.Set("X-EXTERNAL-ID", "external1")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}

func TestExternalID(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(ExternalID(&model.Service{ExternalIDService: mockExternalIDService{}}))
	app.Get("/v1/flush/cache", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/flush/cache", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(ExternalID(&model.Service{ExternalIDService: mockExternalIDService{}}))
	app.Get("/v1/todos", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	req := httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-EXTERNAL-ID", "external1")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(ExternalID(&model.Service{ExternalIDService: mockExternalIDService{err: errors.New("duplicate")}}))
	app.Get("/v1/todos", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	req = httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-EXTERNAL-ID", "external1")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}
}
