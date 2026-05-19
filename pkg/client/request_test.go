package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("X-Test") != "ok" {
			t.Fatalf("unexpected request")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	response := map[string]bool{}
	status, err := RequestJSON(context.Background(), map[string]string{"name": "test"}, &response, server.URL, map[string]string{"X-Test": "ok"})
	if err != nil {
		t.Fatalf("request json: %v", err)
	}
	if status != http.StatusOK || !response["ok"] {
		t.Fatalf("unexpected response: %d %+v", status, response)
	}
}

func TestRequestJSONErrors(t *testing.T) {
	if _, err := RequestJSON(context.Background(), make(chan int), nil, "http://example.test", nil); err == nil {
		t.Fatalf("expected marshal error")
	}

	if _, err := RequestJSON(context.Background(), map[string]string{}, nil, string([]byte{0x7f}), nil); err == nil {
		t.Fatalf("expected request error")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte(`{`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	response := map[string]bool{}
	if _, err := RequestJSON(context.Background(), map[string]string{}, &response, server.URL, nil); err == nil {
		t.Fatalf("expected decode error")
	}

	if _, err := RequestJSON(context.Background(), map[string]string{}, nil, server.URL, nil); err != nil {
		t.Fatalf("nil response should skip decode: %v", err)
	}
}
