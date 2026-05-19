package service

import (
	"context"
	"errors"
	"testing"
)

type mockExternalIDRepository struct {
	err error
}

func (m mockExternalIDRepository) StoreExternalID(context.Context, string) error {
	return m.err
}

func TestExternalIDService(t *testing.T) {
	service := New(mockExternalIDRepository{})
	if err := service.StoreExternalID(context.Background(), "id1"); err != nil {
		t.Fatalf("store external id: %v", err)
	}

	expectedErr := errors.New("repo")
	service = New(mockExternalIDRepository{err: expectedErr})
	if err := service.StoreExternalID(context.Background(), "id1"); !errors.Is(err, expectedErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
