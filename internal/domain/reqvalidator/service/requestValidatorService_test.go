package service

import (
	"errors"
	"testing"

	"ms-gofiber/pkg/apperror"
)

func TestRequestValidator(t *testing.T) {
	validate := New()

	type request struct {
		Title string `validate:"required"`
	}
	if err := validate.ValidateStruct(request{Title: "ok"}); err != nil {
		t.Fatalf("expected valid request: %v", err)
	}

	var appErr *apperror.Error
	if err := validate.ValidateStruct(request{}); !errors.As(err, &appErr) || appErr.Fields["Title"] != "required" {
		t.Fatalf("expected validation error, got %v", err)
	}
	if err := validate.ValidateStruct(nil); !errors.As(err, &appErr) {
		t.Fatalf("expected invalid validation error, got %v", err)
	}
}
