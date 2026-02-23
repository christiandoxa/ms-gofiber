package plainbase

import (
	"testing"

	v10 "github.com/go-playground/validator/v10"
)

func mustRegister(t *testing.T) *v10.Validate {
	t.Helper()
	v := v10.New()
	_ = v.RegisterValidation("alphanum_with_space", ValidateAlphanumWithSpaceRule)
	_ = v.RegisterValidation("authorization_scope", ValidateAuthorizationScopeRule)
	_ = v.RegisterValidation("grant_type", ValidateGrantTypeRule)
	_ = v.RegisterValidation("payment_method_type", ValidatePaymentMethodTypeRule)
	_ = v.RegisterValidation("terminal_type", ValidateTerminalTypeRule)
	_ = v.RegisterValidation("time_rule", ValidateTimeRule)
	_ = v.RegisterValidation("xtrim", ValidateTrimRule)
	_ = v.RegisterValidation("xnotblank", ValidateNotBlankRule)
	return v
}

func TestValidateAlphanumWithSpaceRule(t *testing.T) {
	v := mustRegister(t)
	type S struct {
		V string `validate:"alphanum_with_space"`
	}
	if err := v.Struct(S{V: "Hello 123"}); err != nil {
		t.Fatalf("expected valid: %v", err)
	}
	if err := v.Struct(S{V: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := v.Struct(S{V: "Hello__"}); err == nil {
		t.Fatalf("expected invalid")
	}

	type Bad struct {
		V int `validate:"alphanum_with_space"`
	}
	if err := v.Struct(Bad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateAuthorizationScopeRule(t *testing.T) {
	v := mustRegister(t)
	type S1 struct {
		V string `validate:"authorization_scope"`
	}
	if err := v.Struct(S1{V: "AGREEMENT_PAY SEND_OTP"}); err != nil {
		t.Fatalf("expected valid string scope: %v", err)
	}
	if err := v.Struct(S1{V: "   "}); err != nil {
		t.Fatalf("blank-after-trim should be valid: %v", err)
	}
	if err := v.Struct(S1{V: "UNKNOWN"}); err == nil {
		t.Fatalf("expected invalid string scope")
	}

	type S2 struct {
		V []string `validate:"authorization_scope"`
	}
	if err := v.Struct(S2{V: []string{"AGREEMENT_PAY", "SEND_OTP"}}); err != nil {
		t.Fatalf("expected valid slice scope: %v", err)
	}
	if err := v.Struct(S2{V: []string{"UNKNOWN"}}); err == nil {
		t.Fatalf("expected invalid slice scope")
	}

	type Bad struct {
		V int `validate:"authorization_scope"`
	}
	if err := v.Struct(Bad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateGrantTypeRule(t *testing.T) {
	v := mustRegister(t)
	type S struct {
		V string `validate:"grant_type"`
	}
	if err := v.Struct(S{V: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := v.Struct(S{V: "AUTHORIZATION_CODE"}); err != nil {
		t.Fatalf("expected valid grant type: %v", err)
	}
	if err := v.Struct(S{V: "NOPE"}); err == nil {
		t.Fatalf("expected invalid grant type")
	}

	type Bad struct {
		V int `validate:"grant_type"`
	}
	if err := v.Struct(Bad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidatePaymentMethodTypeRule(t *testing.T) {
	v := mustRegister(t)
	type S struct {
		V string `validate:"payment_method_type"`
	}
	if err := v.Struct(S{V: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := v.Struct(S{V: "DANA"}); err != nil {
		t.Fatalf("expected valid payment type: %v", err)
	}
	if err := v.Struct(S{V: "CARD"}); err == nil {
		t.Fatalf("expected invalid payment type")
	}

	type Bad struct {
		V int `validate:"payment_method_type"`
	}
	if err := v.Struct(Bad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateTerminalTypeRule(t *testing.T) {
	v := mustRegister(t)
	type S struct {
		V string `validate:"terminal_type"`
	}
	if err := v.Struct(S{V: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := v.Struct(S{V: "APP"}); err != nil {
		t.Fatalf("expected valid terminal type: %v", err)
	}
	if err := v.Struct(S{V: "POS"}); err == nil {
		t.Fatalf("expected invalid terminal type")
	}

	type Bad struct {
		V int `validate:"terminal_type"`
	}
	if err := v.Struct(Bad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}

func TestValidateTimeTrimAndNotBlankRules(t *testing.T) {
	v := mustRegister(t)
	type Tm struct {
		V string `validate:"time_rule"`
	}
	if err := v.Struct(Tm{V: ""}); err != nil {
		t.Fatalf("empty should be valid: %v", err)
	}
	if err := v.Struct(Tm{V: "2026-02-23T10:30:00Z"}); err != nil {
		t.Fatalf("expected valid iso time: %v", err)
	}
	if err := v.Struct(Tm{V: "10:30:00"}); err == nil {
		t.Fatalf("expected invalid time")
	}
	type TmBad struct {
		V int `validate:"time_rule"`
	}
	if err := v.Struct(TmBad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}

	type Tr struct {
		V string `validate:"xtrim"`
	}
	if err := v.Struct(Tr{V: "abc"}); err != nil {
		t.Fatalf("expected trim valid: %v", err)
	}
	if err := v.Struct(Tr{V: " abc"}); err == nil {
		t.Fatalf("expected trim invalid")
	}
	type TrBad struct {
		V int `validate:"xtrim"`
	}
	if err := v.Struct(TrBad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}

	type Nb struct {
		V string `validate:"xnotblank"`
	}
	if err := v.Struct(Nb{V: "abc"}); err != nil {
		t.Fatalf("expected nonblank valid: %v", err)
	}
	if err := v.Struct(Nb{V: "   "}); err == nil {
		t.Fatalf("expected nonblank invalid")
	}
	type NbBad struct {
		V int `validate:"xnotblank"`
	}
	if err := v.Struct(NbBad{V: 1}); err == nil {
		t.Fatalf("expected invalid type")
	}
}
