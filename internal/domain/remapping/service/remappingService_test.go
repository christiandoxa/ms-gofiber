package service

import (
	"context"
	"errors"
	"testing"

	remappingmodel "ms-gofiber/internal/domain/remapping/model"
)

type mockRemappingRepository struct {
	responseCode *remappingmodel.ResponseCode
	err          error
}

func (m mockRemappingRepository) Remapping(context.Context, string, string, string) (*remappingmodel.ResponseCode, error) {
	return m.responseCode, m.err
}

func TestRemappingService(t *testing.T) {
	service := New(mockRemappingRepository{responseCode: &remappingmodel.ResponseCode{
		ResponseCode:    "4000001",
		ResponseMessage: "invalid request",
		StatusCode:      400,
	}})

	responseCode, err := service.Remapping(context.Background(), "/v1/test", "E001", "failed")
	if err != nil {
		t.Fatalf("remapping failed: %v", err)
	}
	if responseCode.ResponseCode != "4000001" || responseCode.StatusCode != 400 {
		t.Fatalf("unexpected response code: %+v", responseCode)
	}

	expectedErr := errors.New("repo")
	service = New(mockRemappingRepository{err: expectedErr})
	responseCode, err = service.Remapping(context.Background(), "/v1/test", "E001", "failed")
	if !errors.Is(err, expectedErr) || responseCode != nil {
		t.Fatalf("expected repo error, got %v", err)
	}
}
