package rule

import (
	"testing"

	v10 "github.com/go-playground/validator/v10"
)

func TestRegisterRule(t *testing.T) {
	v := v10.New()
	if err := RegisterRule(v); err != nil {
		t.Fatalf("expected first register success: %v", err)
	}

	orig, exists := customRules["__bad__"]
	customRules["__bad__"] = nil
	t.Cleanup(func() {
		if exists {
			customRules["__bad__"] = orig
			return
		}
		delete(customRules, "__bad__")
	})

	if err := RegisterRule(v10.New()); err == nil {
		t.Fatalf("expected register error for nil function")
	}
}
