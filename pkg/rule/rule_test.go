package rule

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/domain/todo/model/dto"
)

func TestRegisterRule(t *testing.T) {
	validate := validator.New()
	if err := RegisterRule(validate); err != nil {
		t.Fatalf("register rule: %v", err)
	}

	if err := validate.Struct(dto.TodoRequest{Title: "Task 1"}); err != nil {
		t.Fatalf("expected valid request: %v", err)
	}
	if err := validate.Struct(dto.TodoRequest{Title: "Task-1"}); err == nil {
		t.Fatalf("expected plainbase error")
	}
	if err := validate.Struct(dto.TodoRequest{Title: " Task 1 "}); err == nil {
		t.Fatalf("expected structbase error")
	}
}

func TestRegisterRuleError(t *testing.T) {
	expectedErr := errors.New("register")
	patches := gomonkey.ApplyGlobalVar(&registerValidation, func(*validator.Validate, string, validator.Func) error {
		return expectedErr
	})
	defer patches.Reset()

	validate := validator.New()
	if err := RegisterRule(validate); !errors.Is(err, expectedErr) {
		t.Fatalf("expected register error, got %v", err)
	}
}
