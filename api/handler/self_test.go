package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/christiandoxa/welog"
	"github.com/gofiber/fiber/v2"

	mw "ms-gofiber/api/middleware"
	"ms-gofiber/pkg/httpx"
)

func setupSelfApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(welog.NewFiber(fiber.Config{ErrorHandler: mw.ErrorHandler()}))
	return app
}

func TestInternalEcho(t *testing.T) {
	app := setupSelfApp()
	h := &Internal{}
	app.Get("/echo", h.Echo)

	res, err := app.Test(httptest.NewRequest("GET", "/echo?msg=hello", nil))
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
	if body["echo"] != "hello" {
		t.Fatalf("unexpected echo response: %+v", body)
	}
}

func TestInternalSelfCallBranches(t *testing.T) {
	patches := gomonkey.ApplyFunc(httpx.Do, func(context.Context, httpx.Request, httpx.Logger) (*httpx.Response, error) {
		return nil, errors.New("boom")
	})
	app := setupSelfApp()
	h := &Internal{}
	app.Get("/self", h.SelfCall)
	assertStatus(t, app, httptest.NewRequest("GET", "/self", nil), 500)
	patches.Reset()

	patches = gomonkey.ApplyFunc(httpx.Do, func(context.Context, httpx.Request, httpx.Logger) (*httpx.Response, error) {
		return &httpx.Response{StatusCode: 200, Body: []byte(`{`)}, nil
	})
	assertStatus(t, app, httptest.NewRequest("GET", "/self", nil), 500)
	patches.Reset()

	patches = gomonkey.ApplyFunc(httpx.Do, func(context.Context, httpx.Request, httpx.Logger) (*httpx.Response, error) {
		return &httpx.Response{StatusCode: 200, Body: []byte(`{"echo":"ok"}`)}, nil
	})
	defer patches.Reset()

	res, err := app.Test(httptest.NewRequest("GET", "/self", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if body["code"] != "OK" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
