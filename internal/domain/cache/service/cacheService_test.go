package service

import (
	"context"
	"errors"
	"testing"
)

type mockCacheRepository struct {
	err error
}

func (m mockCacheRepository) Flush(context.Context) error {
	return m.err
}

func TestCacheServiceFlush(t *testing.T) {
	service := New(mockCacheRepository{})
	if err := service.Flush(context.Background()); err != nil {
		t.Fatalf("flush cache: %v", err)
	}

	expectedErr := errors.New("repo")
	service = New(mockCacheRepository{err: expectedErr})
	if err := service.Flush(context.Background()); !errors.Is(err, expectedErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
