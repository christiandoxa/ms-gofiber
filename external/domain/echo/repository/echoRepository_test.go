package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoRepositoryFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	response, err := New().Fetch(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("fetch failed: %v", err)
	}
	if string(response.Body) != "ok" {
		t.Fatalf("unexpected response: %s", string(response.Body))
	}
}
