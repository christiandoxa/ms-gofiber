package tool

import (
	"errors"
	"io"
	"math/big"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestEnvHelpers(t *testing.T) {
	t.Setenv("STRING_KEY", "value")
	t.Setenv("INT_KEY", "42")
	t.Setenv("BAD_INT_KEY", "bad")
	t.Setenv("DURATION_KEY", "2s")
	t.Setenv("BAD_DURATION_KEY", "bad")

	if StringFromEnv("STRING_KEY", "fallback") != "value" {
		t.Fatalf("unexpected string value")
	}
	if StringFromEnv("MISSING_STRING_KEY", "fallback") != "fallback" {
		t.Fatalf("unexpected fallback string")
	}
	if IntFromEnv("INT_KEY", 1) != 42 {
		t.Fatalf("unexpected int value")
	}
	if IntFromEnv("BAD_INT_KEY", 1) != 1 {
		t.Fatalf("unexpected bad int fallback")
	}
	if IntFromEnv("MISSING_INT_KEY", 1) != 1 {
		t.Fatalf("unexpected missing int fallback")
	}
	if DurationFromEnv("DURATION_KEY", time.Second) != 2*time.Second {
		t.Fatalf("unexpected duration value")
	}
	if DurationFromEnv("BAD_DURATION_KEY", time.Second) != time.Second {
		t.Fatalf("unexpected bad duration fallback")
	}
	if DurationFromEnv("MISSING_DURATION_KEY", time.Second) != time.Second {
		t.Fatalf("unexpected missing duration fallback")
	}
}

func TestHeaderToMap(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		headers := HeaderToMap(c, "X-Client-ID", "X-Missing")
		if headers["X-Client-ID"] != "client" {
			t.Fatalf("unexpected headers: %+v", headers)
		}
		if _, ok := headers["X-Missing"]; ok {
			t.Fatalf("missing header should be skipped")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-Client-ID", "client")
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var requestHeader fasthttp.RequestHeader
	requestHeader.Set("X-Test", "request")
	headers := HeaderToMap(&requestHeader)
	if headers["X-Test"] != "request" {
		t.Fatalf("unexpected request header map: %+v", headers)
	}

	var responseHeader fasthttp.ResponseHeader
	responseHeader.Set("X-Test", "response")
	headers = HeaderToMap(&responseHeader)
	if headers["X-Test"] != "response" {
		t.Fatalf("unexpected response header map: %+v", headers)
	}

	if len(HeaderToMap(struct{}{})) != 0 {
		t.Fatalf("unsupported header should return empty map")
	}
}

func TestNowUTC(t *testing.T) {
	if NowUTC().Location() != time.UTC {
		t.Fatalf("expected UTC time")
	}
}

func TestExpiration(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	if EndOfDayExpiration(now) <= 0 {
		t.Fatalf("expected positive expiration")
	}

	endOfDay := time.Date(2026, 5, 18, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
	if EndOfDayExpiration(endOfDay) <= 0 {
		t.Fatalf("expected next day expiration")
	}

	if GetExpiration() <= 0 {
		t.Fatalf("expected current expiration")
	}
}

func TestGenerateRequestReference(t *testing.T) {
	originalRandomInt := randomInt
	t.Cleanup(func() {
		randomInt = originalRandomInt
	})

	randomInt = func(io.Reader, *big.Int) (*big.Int, error) {
		return big.NewInt(42), nil
	}
	if got := GenerateRequestReference(); got != "000000000042" {
		t.Fatalf("unexpected request reference: %s", got)
	}
	if got := GenerateRequestRefnum(); got != "000000000042" {
		t.Fatalf("unexpected request refnum: %s", got)
	}

	expectedErr := errors.New("random")
	randomInt = func(io.Reader, *big.Int) (*big.Int, error) {
		return nil, expectedErr
	}
	if got := GenerateRequestReference(); got != "" {
		t.Fatalf("expected empty request reference, got %s", got)
	}
}
