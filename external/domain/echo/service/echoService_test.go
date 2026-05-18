package service

import (
	"context"
	"errors"
	"testing"

	"ms-gofiber/pkg/client"
)

type mockEchoRepository struct {
	response *client.Response
	err      error
}

func (m mockEchoRepository) Fetch(context.Context, string) (*client.Response, error) {
	return m.response, m.err
}

func TestEchoService(t *testing.T) {
	service := New(mockEchoRepository{response: &client.Response{Body: []byte("ok")}})
	response, err := service.Echo(context.Background(), "target")
	if err != nil {
		t.Fatalf("echo failed: %v", err)
	}
	if response != "ok" {
		t.Fatalf("unexpected response: %s", response)
	}

	errExpected := errors.New("fetch")
	service = New(mockEchoRepository{err: errExpected})
	if _, err := service.Echo(context.Background(), "target"); !errors.Is(err, errExpected) {
		t.Fatalf("expected fetch error, got %v", err)
	}
}
