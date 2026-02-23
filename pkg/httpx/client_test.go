package httpx

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/christiandoxa/welog/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type errCloseBody struct {
	io.Reader
}

func (b errCloseBody) Close() error { return errors.New("close error") }

func TestHeaderToInterface(t *testing.T) {
	h := http.Header{}
	h.Set("A", "1")
	h["B"] = []string{"2", "3"}
	m := headerToInterface(h)
	if m["A"] != "1" {
		t.Fatalf("unexpected single header: %+v", m)
	}
	if _, ok := m["B"].([]string); !ok {
		t.Fatalf("expected multi-value slice for header B: %+v", m)
	}
}

func TestDoBranches(t *testing.T) {
	origReq := newHTTPRequest
	origWrap := wrapHTTPClient
	origLog := logFiberClient
	t.Cleanup(func() {
		newHTTPRequest = origReq
		wrapHTTPClient = origWrap
		logFiberClient = origLog
	})
	logFiberClient = func(*fiber.Ctx, model.TargetRequest, model.TargetResponse) {}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})

	// new request error
	newHTTPRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("new request error")
	}
	app.Get("/reqerr", func(c *fiber.Ctx) error {
		_, _, _, err := Do(c, "GET", "http://example.com", "", nil, nil, 0)
		if err == nil {
			t.Fatalf("expected request creation error")
		}
		return c.SendStatus(200)
	})
	if _, err := app.Test(httptest.NewRequest("GET", "/reqerr", nil)); err != nil {
		t.Fatalf("fiber test failed: %v", err)
	}

	// client do error path
	newHTTPRequest = origReq
	wrapHTTPClient = func(_ *http.Client, _ ...apmhttp.ClientOption) *http.Client {
		return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("transport error")
		})}
	}
	app.Get("/doerr", func(c *fiber.Ctx) error {
		_, _, _, err := Do(c, "GET", "http://example.com", "", nil, nil, time.Second)
		if err == nil {
			t.Fatalf("expected do error")
		}
		return c.SendStatus(200)
	})
	if _, err := app.Test(httptest.NewRequest("GET", "/doerr", nil)); err != nil {
		t.Fatalf("fiber test failed: %v", err)
	}

	// success path with close-error branch
	wrapHTTPClient = func(_ *http.Client, _ ...apmhttp.ClientOption) *http.Client {
		return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"X-Test": []string{"ok"}},
				Body:       errCloseBody{Reader: bytes.NewBufferString(`{"ok":true}`)},
			}, nil
		})}
	}
	app.Get("/ok", func(c *fiber.Ctx) error {
		status, body, hdr, err := Do(c, "GET", "http://example.com", "application/json", map[string]string{"A": "B"}, []byte(`{"in":1}`), time.Second)
		if err != nil {
			t.Fatalf("unexpected success error: %v", err)
		}
		if status != 200 || string(body) != `{"ok":true}` || hdr.Get("X-Test") != "ok" {
			t.Fatalf("unexpected do response status=%d body=%s hdr=%v", status, string(body), hdr)
		}
		return c.SendStatus(200)
	})
	if _, err := app.Test(httptest.NewRequest("GET", "/ok", nil)); err != nil {
		t.Fatalf("fiber test failed: %v", err)
	}

	// real server success path
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-From", "server")
		_, _ = w.Write([]byte("pong"))
	}))
	defer srv.Close()
	wrapHTTPClient = origWrap
	app.Get("/real", func(c *fiber.Ctx) error {
		status, body, hdr, err := Do(c, "GET", srv.URL, "", nil, nil, 0)
		if err != nil {
			t.Fatalf("real do error: %v", err)
		}
		if status != 200 || string(body) != "pong" || hdr.Get("X-From") != "server" {
			t.Fatalf("unexpected real do response")
		}
		return c.SendStatus(200)
	})
	if _, err := app.Test(httptest.NewRequest("GET", "/real", nil)); err != nil {
		t.Fatalf("fiber test failed: %v", err)
	}
}
