package controller

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/christiandoxa/welog/pkg/model"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/httpx"
)

func TestFiberHTTPLogger(t *testing.T) {
	req := httpx.RequestLog{
		URL:         "http://service.local",
		Method:      "GET",
		ContentType: "application/json",
		Header:      map[string]interface{}{"X-Test": "1"},
		Body:        []byte(`{"request":true}`),
		Timestamp:   time.Now().UTC(),
	}
	res := httpx.ResponseLog{
		Header:  map[string]interface{}{"X-Result": "ok"},
		Body:    []byte(`{"ok":true}`),
		Status:  200,
		Latency: time.Millisecond,
	}

	fiberHTTPLogger{}.Log(context.Background(), req, res)

	called := false
	patches := gomonkey.ApplyGlobalVar(&logFiberClient, func(_ *fiber.Ctx, gotReq model.TargetRequest, gotRes model.TargetResponse) {
		called = true
		if gotReq.URL != req.URL || gotReq.Method != req.Method || gotReq.ContentType != req.ContentType {
			t.Fatalf("unexpected request log: %+v", gotReq)
		}
		if gotRes.Status != res.Status || gotRes.Latency != res.Latency {
			t.Fatalf("unexpected response log: %+v", gotRes)
		}
	})
	defer patches.Reset()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fiberHTTPLogger{ctx: c}.Log(context.Background(), req, res)
		return c.SendStatus(fiber.StatusOK)
	})
	response, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if !called {
		t.Fatalf("expected logFiberClient call")
	}
}

func TestHTTPLogConverters(t *testing.T) {
	now := time.Now().UTC()
	req := toTargetRequest(httpx.RequestLog{
		URL:         "http://example.com",
		Method:      "POST",
		ContentType: "application/json",
		Header:      map[string]interface{}{"A": "B"},
		Body:        []byte("body"),
		Timestamp:   now,
	})
	if req.URL != "http://example.com" || req.Method != "POST" || !req.Timestamp.Equal(now) {
		t.Fatalf("unexpected target request: %+v", req)
	}

	res := toTargetResponse(httpx.ResponseLog{
		Header:  map[string]interface{}{"C": "D"},
		Body:    []byte("response"),
		Status:  201,
		Latency: time.Second,
	})
	if res.Status != 201 || res.Latency != time.Second {
		t.Fatalf("unexpected target response: %+v", res)
	}
}
