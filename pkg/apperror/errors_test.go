package apperror

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorString(t *testing.T) {
	err := &Error{Code: ErrBadRequest, Message: "bad input"}
	if got := err.Error(); got != "BAD_REQUEST: bad input" {
		t.Fatalf("unexpected error string: %s", got)
	}

	wrapped := &Error{Code: ErrDB, Message: "db failed", Err: errors.New("boom")}
	if got := wrapped.Error(); !strings.Contains(got, "DB_ERROR: db failed: boom") {
		t.Fatalf("unexpected wrapped string: %s", got)
	}
}

func TestConstructors(t *testing.T) {
	e1 := New(ErrConflict, "conflict")
	if e1.Code != ErrConflict || e1.Message != "conflict" || e1.Err != nil {
		t.Fatalf("unexpected New result: %+v", e1)
	}

	fields := map[string]string{"a": "b"}
	e2 := WithFields(ErrValidation, "invalid", fields)
	if e2.Code != ErrValidation || e2.Fields["a"] != "b" {
		t.Fatalf("unexpected WithFields result: %+v", e2)
	}

	base := errors.New("base")
	e3 := Wrap(ErrInternal, "wrapped", base)
	if !errors.Is(e3.Err, base) || e3.Code != ErrInternal {
		t.Fatalf("unexpected Wrap result: %+v", e3)
	}
}
