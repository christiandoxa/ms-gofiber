package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "ok" {
			t.Fatalf("unexpected header")
		}
		w.WriteHeader(http.StatusAccepted)
		if _, err := w.Write([]byte("accepted")); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	response, err := Do(context.Background(), Request{
		URL:    server.URL,
		Header: map[string]string{"X-Test": "ok"},
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if response.StatusCode != http.StatusAccepted || string(response.Body) != "accepted" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestDoBranches(t *testing.T) {
	if _, err := Do(context.Background(), Request{URL: string([]byte{0x7f})}); err == nil {
		t.Fatalf("expected invalid url error")
	}
	if _, err := Do(context.Background(), Request{URL: "http://127.0.0.1:1", Timeout: time.Millisecond}); err == nil {
		t.Fatalf("expected request error")
	}
}

func TestDoReadBodyError(t *testing.T) {
	originalClient := http.DefaultClient
	t.Cleanup(func() {
		http.DefaultClient = originalClient
	})
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{},
				Body:       errReader{},
			}, nil
		}),
	}
	if _, err := Do(context.Background(), Request{URL: "http://example.test"}); err == nil {
		t.Fatalf("expected body read error")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read")
}

func (errReader) Close() error {
	return nil
}

var _ io.ReadCloser = errReader{}
