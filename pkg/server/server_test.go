package server

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestNewServer(t *testing.T) {
	if app := NewServer(); app == nil {
		t.Fatalf("expected app")
	}
}

func TestInitValidator(t *testing.T) {
	if validate := initValidator(); validate == nil {
		t.Fatalf("expected validator")
	}
}

func TestInitValidatorError(t *testing.T) {
	originalRegisterRule := registerRule
	originalFatal := fatal
	t.Cleanup(func() {
		registerRule = originalRegisterRule
		fatal = originalFatal
	})

	expectedErr := errors.New("register")
	fatalCalled := false
	registerRule = func(*validator.Validate) error {
		return expectedErr
	}
	fatal = func(args ...any) {
		err, ok := args[0].(error)
		if !ok || !errors.Is(err, expectedErr) {
			t.Fatalf("unexpected fatal args: %+v", args)
		}
		fatalCalled = true
	}

	if validate := initValidator(); validate == nil {
		t.Fatalf("expected validator")
	}
	if !fatalCalled {
		t.Fatalf("expected fatal call")
	}
}
