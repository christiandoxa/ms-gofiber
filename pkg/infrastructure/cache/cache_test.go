package cache

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	instance = nil
	once = sync.Once{}

	first := Connect()
	second := Connect()
	if first == nil || first != second {
		t.Fatalf("expected singleton cache")
	}
}

func TestCacheStoreGetDeleteFlush(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cache := New()
	cache.now = func() time.Time { return now }

	if err := cache.Store(ctx, Data{Key: "key", Value: "value"}); err != nil {
		t.Fatalf("store cache: %v", err)
	}
	value, err := cache.Get(ctx, "key")
	if err != nil || value != "value" {
		t.Fatalf("unexpected cache value: %v %v", value, err)
	}
	if err := cache.Store(ctx, Data{Key: "key", Value: "next"}); !errors.Is(err, ErrCacheKeyExists) {
		t.Fatalf("expected duplicate key error, got %v", err)
	}
	if err := cache.Store(ctx, Data{Key: "key", Value: "next", Override: true}); err != nil {
		t.Fatalf("override cache: %v", err)
	}
	value, err = cache.Get(ctx, "key")
	if err != nil || value != "next" {
		t.Fatalf("unexpected override value: %v %v", value, err)
	}
	if cache.Len() != 1 {
		t.Fatalf("expected cache len 1")
	}
	if err := cache.Delete(ctx, "key"); err != nil {
		t.Fatalf("delete cache: %v", err)
	}
	if _, err := cache.Get(ctx, "key"); !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("expected cache miss, got %v", err)
	}

	if err := cache.Store(ctx, Data{Key: "one", Value: 1}); err != nil {
		t.Fatalf("store one: %v", err)
	}
	if err := cache.Store(ctx, Data{Key: "two", Value: 2}); err != nil {
		t.Fatalf("store two: %v", err)
	}
	if err := cache.Flush(ctx); err != nil {
		t.Fatalf("flush cache: %v", err)
	}
	if cache.Len() != 0 {
		t.Fatalf("expected empty cache")
	}
}

func TestCacheExpiration(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cache := New()
	cache.now = func() time.Time { return now }

	if err := cache.Store(ctx, Data{Key: "key", Value: "value", Duration: time.Second}); err != nil {
		t.Fatalf("store cache: %v", err)
	}
	now = now.Add(time.Second)
	if _, err := cache.Get(ctx, "key"); !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("expected expired cache miss, got %v", err)
	}

	if err := cache.Store(ctx, Data{Key: "key", Value: "value", Duration: time.Second}); err != nil {
		t.Fatalf("store cache again: %v", err)
	}
	now = now.Add(time.Second)
	if cache.Len() != 0 {
		t.Fatalf("expected expired entry pruned")
	}

	if err := cache.Store(ctx, Data{Key: "key", Value: "value", Duration: time.Second}); err != nil {
		t.Fatalf("store cache third: %v", err)
	}
	now = now.Add(time.Second)
	if err := cache.Store(ctx, Data{Key: "key", Value: "new"}); err != nil {
		t.Fatalf("expected expired key to be replaceable: %v", err)
	}
}

func TestCacheValidationAndContext(t *testing.T) {
	cache := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := cache.Store(ctx, Data{Key: "key", Value: "value"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected store context error, got %v", err)
	}
	if _, err := cache.Get(ctx, "key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected get context error, got %v", err)
	}
	if err := cache.Delete(ctx, "key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected delete context error, got %v", err)
	}
	if err := cache.Flush(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected flush context error, got %v", err)
	}

	if err := cache.Store(context.Background(), Data{}); !errors.Is(err, ErrCacheKeyEmpty) {
		t.Fatalf("expected empty store key error, got %v", err)
	}
	if _, err := cache.Get(context.Background(), ""); !errors.Is(err, ErrCacheKeyEmpty) {
		t.Fatalf("expected empty get key error, got %v", err)
	}
	if err := cache.Delete(context.Background(), ""); !errors.Is(err, ErrCacheKeyEmpty) {
		t.Fatalf("expected empty delete key error, got %v", err)
	}
	if err := cache.Flush(context.Background()); err != nil {
		t.Fatalf("expected flush to pass: %v", err)
	}
}
