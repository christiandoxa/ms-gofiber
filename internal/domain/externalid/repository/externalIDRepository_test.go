package repository

import (
	"context"
	"errors"
	"testing"

	"ms-gofiber/pkg/infrastructure/cache"
	"ms-gofiber/pkg/responsecode/error"
)

func TestExternalIDRepository(t *testing.T) {
	repository := New(cache.New())
	if err := repository.StoreExternalID(context.Background(), "id1"); err != nil {
		t.Fatalf("store external id: %v", err)
	}
	if err := repository.StoreExternalID(context.Background(), "id1"); !errors.Is(err, rcerror.ErrDuplicateExternalID) {
		t.Fatalf("expected duplicate external id, got %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := repository.StoreExternalID(ctx, "id2"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context error, got %v", err)
	}
}
