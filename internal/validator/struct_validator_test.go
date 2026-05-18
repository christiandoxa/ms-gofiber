package validator

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	v10 "github.com/go-playground/validator/v10"

	"ms-gofiber/internal/validator/rule"
	"ms-gofiber/pkg/apperror"
)

func TestStructValidator(t *testing.T) {
	sv, err := NewStructValidator()
	if err != nil {
		t.Fatalf("new struct validator: %v", err)
	}

	sv.RegisterStructValidation(func(sl v10.StructLevel) {}, struct{}{})

	type sample struct {
		Name string `validate:"required,alphanum_with_space"`
	}

	if err := sv.ValidateStruct(sample{Name: "Todo 123"}); err != nil {
		t.Fatalf("expected valid struct: %v", err)
	}

	err = sv.ValidateStruct(sample{Name: "!"})
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

func TestNewStructValidatorDefaultRegisterError(t *testing.T) {
	patches := gomonkey.ApplyFunc(rule.RegisterRule, func(*v10.Validate) error {
		return errors.New("default register error")
	})
	defer patches.Reset()

	if _, err := NewStructValidator(); err == nil {
		t.Fatalf("expected default register error")
	}
}
