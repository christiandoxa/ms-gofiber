package plainbase

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/pkg/constant/rulekey"
)

func mustRegister(t *testing.T) *validator.Validate {
	t.Helper()
	validate := validator.New()
	if err := validate.RegisterValidation(rulekey.AlphanumWithSpaceRule, ValidateAlphanumWithSpaceRule); err != nil {
		t.Fatalf("register alphanum rule: %v", err)
	}
	if err := validate.RegisterValidation(rulekey.NotBlankRule, ValidateNotBlankRule); err != nil {
		t.Fatalf("register not blank rule: %v", err)
	}
	if err := validate.RegisterValidation(rulekey.TrimRule, ValidateTrimRule); err != nil {
		t.Fatalf("register trim rule: %v", err)
	}
	if err := validate.RegisterValidation(rulekey.TimeRule, ValidateTimeRule); err != nil {
		t.Fatalf("register time rule: %v", err)
	}
	return validate
}

func TestValidateAlphanumWithSpaceRule(t *testing.T) {
	validate := mustRegister(t)
	type request struct {
		Value string `validate:"alphanumWithSpaceRule"`
	}
	if err := validate.Struct(request{Value: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := validate.Struct(request{Value: "Task 123"}); err != nil {
		t.Fatalf("expected valid value: %v", err)
	}
	if err := validate.Struct(request{Value: "Task-123"}); err == nil {
		t.Fatalf("expected invalid value")
	}

	type badRequest struct {
		Value int `validate:"alphanumWithSpaceRule"`
	}
	if err := validate.Struct(badRequest{Value: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateNotBlankRule(t *testing.T) {
	validate := mustRegister(t)
	type request struct {
		Value string `validate:"notBlankRule"`
	}
	if err := validate.Struct(request{Value: "content"}); err != nil {
		t.Fatalf("expected valid value: %v", err)
	}
	if err := validate.Struct(request{Value: "  "}); err == nil {
		t.Fatalf("expected blank value invalid")
	}

	type badRequest struct {
		Value int `validate:"notBlankRule"`
	}
	if err := validate.Struct(badRequest{Value: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateTrimRule(t *testing.T) {
	validate := mustRegister(t)
	type request struct {
		Value string `validate:"trimRule"`
	}
	if err := validate.Struct(request{Value: "content"}); err != nil {
		t.Fatalf("expected valid value: %v", err)
	}
	if err := validate.Struct(request{Value: " content "}); err == nil {
		t.Fatalf("expected untrimmed value invalid")
	}

	type badRequest struct {
		Value int `validate:"trimRule"`
	}
	if err := validate.Struct(badRequest{Value: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateTimeRule(t *testing.T) {
	validate := mustRegister(t)
	type request struct {
		Value string `validate:"timeRule"`
	}
	if err := validate.Struct(request{Value: "2026-05-19T10:30:00Z"}); err != nil {
		t.Fatalf("expected valid time: %v", err)
	}
	if err := validate.Struct(request{Value: "2026-99-99T10:30:00Z"}); err == nil {
		t.Fatalf("expected invalid time")
	}

	type badRequest struct {
		Value int `validate:"timeRule"`
	}
	if err := validate.Struct(badRequest{Value: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}
