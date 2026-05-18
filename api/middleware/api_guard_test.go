package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alicebob/miniredis/v2"
	"github.com/christiandoxa/welog"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"ms-gofiber/pkg/apperror"
)

type mockReqValidator struct {
	err error
}

func (m mockReqValidator) ValidateStruct(any) error { return m.err }

func withLogger(app *fiber.App) {
	app.Use(welog.NewFiber(fiber.Config{ErrorHandler: ErrorHandler()}))
}

func TestDefaultSkipHelpers(t *testing.T) {
	s := skippedPaths()
	if !shouldSkipPath("/v1/health", s) {
		t.Fatalf("health should be skipped")
	}
	if shouldSkipPath("/v1/other", s) {
		t.Fatalf("other should not be skipped")
	}
}

func TestHeaderGuardBranches(t *testing.T) {
	// parse error branch
	var fiberCtx *fiber.Ctx
	patches := gomonkey.ApplyMethod(fiberCtx, "ReqHeaderParser", func(*fiber.Ctx, interface{}) error {
		return errors.New("parse")
	})
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app)
	app.Use(HeaderGuard(mockReqValidator{}))
	app.Get("/x", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	res, err := app.Test(httptest.NewRequest("GET", "/x", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
	patches.Reset()

	// validation error branch
	app2 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app2)
	app2.Use(HeaderGuard(mockReqValidator{err: apperror.New(apperror.ErrValidation, "invalid")}))
	app2.Get("/x", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("X-PARTNER-ID", "A1")
	req.Header.Set("CHANNEL-ID", "B1")
	req.Header.Set("X-EXTERNAL-ID", "123")
	res, err = app2.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 for validation error, got %d", res.StatusCode)
	}

	// success branch + set locals
	app3 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app3)
	app3.Use(HeaderGuard(mockReqValidator{}))
	app3.Get("/v1/todos", func(c *fiber.Ctx) error {
		header, ok := c.Locals("request_header").(RequestHeader)
		return c.JSON(fiber.Map{"ok": ok, "partner": header.XPartnerID})
	})
	req = httptest.NewRequest("GET", "/v1/todos", nil)
	req.Header.Set("X-PARTNER-ID", "A1")
	req.Header.Set("CHANNEL-ID", "B1")
	req.Header.Set("X-EXTERNAL-ID", "123")
	res, err = app3.Test(req)
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
	if body["ok"] != true || body["partner"] != "A1" {
		t.Fatalf("expected request_header locals present: %+v", body)
	}

	// skip path branch
	app4 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app4)
	app4.Use(HeaderGuard(mockReqValidator{err: apperror.New(apperror.ErrValidation, "invalid")}))
	app4.Get("/v1/health", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	res, err = app4.Test(httptest.NewRequest("GET", "/v1/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected skipped path success, got %d", res.StatusCode)
	}
}

func TestExternalIDGuardBranches(t *testing.T) {
	// missing header
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app)
	app.Use(ExternalIDGuard(nil, 0))
	app.Get("/x", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	res, err := app.Test(httptest.NewRequest("GET", "/x", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// redis nil branch
	app2 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app2)
	app2.Use(ExternalIDGuard(nil, 0))
	app2.Get("/x", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("X-EXTERNAL-ID", "100")
	res, err = app2.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	// redis error branch
	badRedis := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:1",
		MaxRetries: 0,
		Dialer: func(context.Context, string, string) (net.Conn, error) {
			return nil, errors.New("dial error")
		},
	})
	app3 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app3)
	app3.Use(ExternalIDGuard(badRedis, time.Second))
	app3.Get("/x", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req = httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("X-EXTERNAL-ID", "100")
	res, err = app3.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}

	// success + duplicate
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis error: %v", err)
	}
	defer mr.Close()
	goodRedis := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	app4 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app4)
	app4.Use(ExternalIDGuard(goodRedis, time.Second))
	app4.Get("/v1/todos", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req = httptest.NewRequest("GET", "/v1/todos", nil)
	req.Header.Set("X-EXTERNAL-ID", "101")
	res, err = app4.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected first 200 got %d", res.StatusCode)
	}

	req = httptest.NewRequest("GET", "/v1/todos", nil)
	req.Header.Set("X-EXTERNAL-ID", "101")
	res, err = app4.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 409 {
		t.Fatalf("expected duplicate 409 got %d", res.StatusCode)
	}

	// skip path
	app5 := fiber.New(fiber.Config{ErrorHandler: ErrorHandler()})
	withLogger(app5)
	app5.Use(ExternalIDGuard(goodRedis, time.Second))
	app5.Get("/v1/health", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	res, err = app5.Test(httptest.NewRequest("GET", "/v1/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected health skip 200 got %d", res.StatusCode)
	}
}

func TestLogMiddlewareErrorNoLogger(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		logMiddlewareError(c, errors.New("x"))
		return c.SendStatus(200)
	})
	if _, err := app.Test(httptest.NewRequest("GET", "/", nil)); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
