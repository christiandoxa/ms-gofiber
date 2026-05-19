package repository

import (
	"context"
	"errors"
	"testing"

	"ms-gofiber/pkg/infrastructure/cache"
)

func TestCacheRepositoryFlush(t *testing.T) {
	cacheClient := cache.New()
	if err := cacheClient.Store(context.Background(), cache.Data{Key: "key", Value: "value"}); err != nil {
		t.Fatalf("store cache: %v", err)
	}
	repository := New(cacheClient)
	if err := repository.Flush(context.Background()); err != nil {
		t.Fatalf("flush cache: %v", err)
	}
	if cacheClient.Len() != 0 {
		t.Fatalf("expected empty cache")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := repository.Flush(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context error, got %v", err)
	}
}
