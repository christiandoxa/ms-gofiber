package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
	"ms-gofiber/pkg/apperror"
)

type mockEchoService struct {
	body string
	err  error
}

func (m mockEchoService) Echo(context.Context, string) (string, error) {
	return m.body, m.err
}

func TestEcho(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Get("/echo", Echo(&model.Service{EchoService: mockEchoService{body: "ok"}}))

	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/echo", nil), http.StatusBadRequest)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/echo?target=http://example.test", nil), http.StatusOK)

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Get("/echo", Echo(&model.Service{EchoService: mockEchoService{err: apperror.New(http.StatusInternalServerError, "failed")}}))
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/echo?target=http://example.test", nil), http.StatusInternalServerError)

	app = fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Get("/echo", Echo(&model.Service{EchoService: mockEchoService{err: errors.New("failed")}}))
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/echo?target=http://example.test", nil), http.StatusInternalServerError)
}

func assertStatus(t *testing.T, app *fiber.App, req *http.Request, expected int) {
	t.Helper()
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.StatusCode != expected {
		t.Fatalf("expected %d got %d", expected, res.StatusCode)
	}
}
