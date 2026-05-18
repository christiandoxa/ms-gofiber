package service

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/domain/todo/model/dto"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/rule"
)

func TestRequestValidator(t *testing.T) {
	validate := newTestValidator(t)

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

	type hiddenRequest struct {
		Hidden string `json:"-" validate:"required"`
	}
	if err := validate.ValidateStruct(hiddenRequest{}); !errors.As(err, &appErr) || len(appErr.Fields) == 0 {
		t.Fatalf("expected hidden validation error, got %v", err)
	}
}

func TestRequestValidatorCustomRules(t *testing.T) {
	validate := newTestValidator(t)

	if err := validate.ValidateStruct(dto.TodoRequest{Title: "Task 1"}); err != nil {
		t.Fatalf("expected valid todo request: %v", err)
	}

	var appErr *apperror.Error
	if err := validate.ValidateStruct(dto.TodoRequest{Title: "   "}); !errors.As(err, &appErr) || appErr.Fields["title"] == "" {
		t.Fatalf("expected blank title validation error, got %v", err)
	}
	if err := validate.ValidateStruct(dto.TodoRequest{Title: "Task-1"}); !errors.As(err, &appErr) || appErr.Fields["title"] != "alphanumWithSpaceRule" {
		t.Fatalf("expected alphanum title validation error, got %v", err)
	}
	if err := validate.ValidateStruct(dto.TodoRequest{Title: " Task 1 "}); !errors.As(err, &appErr) || appErr.Fields["title"] == "" {
		t.Fatalf("expected struct-level title validation error, got %v", err)
	}
}

func newTestValidator(t *testing.T) IRequestValidator {
	t.Helper()
	validate := validator.New()
	if err := rule.RegisterRule(validate); err != nil {
		t.Fatalf("register rule: %v", err)
	}
	return New(validate)
}
