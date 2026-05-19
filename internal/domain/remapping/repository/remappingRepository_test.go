package repository

import (
	"context"
	"errors"
	"testing"

	remappingmodel "ms-gofiber/internal/domain/remapping/model"
)

func TestRemappingRepository(t *testing.T) {
	repository := New(remappingmodel.ResponseCode{
		Endpoint:              "/v1/test",
		OriginResponseCode:    "E001",
		OriginResponseMessage: "failed",
		ResponseCode:          "4000001",
		ResponseMessage:       "invalid request",
		StatusCode:            400,
	})

	responseCode, err := repository.Remapping(context.Background(), "/v1/test", "E001", "failed")
	if err != nil {
		t.Fatalf("remapping failed: %v", err)
	}
	if responseCode.ResponseCode != "4000001" {
		t.Fatalf("unexpected response code: %+v", responseCode)
	}

	if _, err := repository.Remapping(context.Background(), "/v1/test", "missing", "failed"); !errors.Is(err, ErrRemappingNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	if cacheKey("a", "b", "c") != "a|b|c" {
		t.Fatalf("unexpected cache key")
	}
}
