package structbase

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/domain/todo/model/dto"
	"ms-gofiber/pkg/constant/rulekey"
	"ms-gofiber/pkg/rule/plainbase"
)

func mustRegister(t *testing.T) *validator.Validate {
	t.Helper()
	validate := validator.New()
	if err := validate.RegisterValidation(rulekey.AlphanumWithSpaceRule, plainbase.ValidateAlphanumWithSpaceRule); err != nil {
		t.Fatalf("register alphanum rule: %v", err)
	}
	if err := validate.RegisterValidation(rulekey.NotBlankRule, plainbase.ValidateNotBlankRule); err != nil {
		t.Fatalf("register not blank rule: %v", err)
	}
	RegisterRule(validate)
	return validate
}

func TestValidateTodoRequestRule(t *testing.T) {
	validate := mustRegister(t)

	if err := validate.Struct(dto.TodoRequest{Title: "Task"}); err != nil {
		t.Fatalf("expected valid request: %v", err)
	}
	if err := validate.Struct(dto.TodoRequest{Title: " Task "}); err == nil {
		t.Fatalf("expected trim error")
	}
}

func TestValidateTodoRequestRuleWrongType(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(ValidateTodoRequestRule, struct{ Title string }{}) //nolint:errcheck

	if err := validate.Struct(struct{ Title string }{Title: " Task "}); err != nil {
		t.Fatalf("wrong type should be ignored: %v", err)
	}
}

func TestRuleKeys(t *testing.T) {
	if rulekey.TodoTitleTrimRule == "" {
		t.Fatalf("expected todo title trim rule key")
	}
}
