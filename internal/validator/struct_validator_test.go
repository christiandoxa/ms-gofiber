package validator

import (
	"errors"
	"testing"

	"ms-gofiber/internal/dto"
	"ms-gofiber/pkg/apperror"
)

func TestStructValidator(t *testing.T) {
	sv := NewStructValidator()

	if err := sv.ValidateStruct(dto.TodoUpsertRequest{Title: "Todo123", Completed: false}); err != nil {
		t.Fatalf("expected valid struct: %v", err)
	}

	err := sv.ValidateStruct(dto.TodoUpsertRequest{Title: " ", Completed: false})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	var aerr *apperror.Error
	if !errors.As(err, &aerr) {
		t.Fatalf("expected apperror.Error, got %T", err)
	}
	if aerr.Code != apperror.ErrValidation {
		t.Fatalf("unexpected code: %s", aerr.Code)
	}
	if len(aerr.Fields) == 0 {
		t.Fatalf("expected fields in validation error")
	}

	err = sv.ValidateStruct(nil)
	if err == nil {
		t.Fatalf("expected generic validation error")
	}
	if !errors.As(err, &aerr) {
		t.Fatalf("expected apperror.Error for nil, got %T", err)
	}
}
