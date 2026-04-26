package httpx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.elastic.co/apm/module/apmhttp/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type errCloseBody struct {
	io.Reader
}

func (b errCloseBody) Close() error { return errors.New("close error") }

type noopLogger struct{}

func (noopLogger) Log(context.Context, RequestLog, ResponseLog) {}

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
	t.Cleanup(func() {
		newHTTPRequest = origReq
		wrapHTTPClient = origWrap
	})
	// new request error
	newHTTPRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("new request error")
	}
	_, err := Do(context.Background(), Request{
		Method: "GET",
		URL:    "http://example.com",
	}, noopLogger{})
	if err == nil {
		t.Fatalf("expected request creation error")
	}

	// client do error path
	newHTTPRequest = origReq
	wrapHTTPClient = func(_ *http.Client, _ ...apmhttp.ClientOption) *http.Client {
		return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("transport error")
		})}
	}
	_, err = Do(context.Background(), Request{
		Method:  "GET",
		URL:     "http://example.com",
		Timeout: time.Second,
	}, noopLogger{})
	if err == nil {
		t.Fatalf("expected do error")
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
	res, err := Do(context.Background(), Request{
		Method:      "GET",
		URL:         "http://example.com",
		ContentType: "application/json",
		Header:      map[string]string{"A": "B"},
		Body:        []byte(`{"in":1}`),
		Timeout:     time.Second,
	}, noopLogger{})
	if err != nil {
		t.Fatalf("unexpected success error: %v", err)
	}
	if res.StatusCode != 200 || string(res.Body) != `{"ok":true}` || res.Header.Get("X-Test") != "ok" {
		t.Fatalf("unexpected do response status=%d body=%s hdr=%v", res.StatusCode, string(res.Body), res.Header)
	}

	// real server success path
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-From", "server")
		_, _ = w.Write([]byte("pong"))
	}))
	defer srv.Close()
	wrapHTTPClient = origWrap
	res, err = Do(context.Background(), Request{
		Method: "GET",
		URL:    srv.URL,
	}, nil)
	if err != nil {
		t.Fatalf("real do error: %v", err)
	}
	if res.StatusCode != 200 || string(res.Body) != "pong" || res.Header.Get("X-From") != "server" {
		t.Fatalf("unexpected real do response")
	}
}
