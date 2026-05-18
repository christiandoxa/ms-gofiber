package respond

import (
	"testing"

	"ms-gofiber/pkg/apperror"
)

func TestEnvelopes(t *testing.T) {
	s := SuccessEnvelope(map[string]any{"ok": true}, map[string]any{"meta": 1})
	if s.Code != "OK" || s.Message != "success" {
		t.Fatalf("unexpected success envelope: %+v", s)
	}

	e := ErrorEnvelope(apperror.ErrValidation, "invalid", map[string]string{"f": "required"})
	if e.Code != string(apperror.ErrValidation) || e.Fields["f"] != "required" {
		t.Fatalf("unexpected error envelope: %+v", e)
	}
}

func TestHTTPStatusFromCode(t *testing.T) {
	cases := map[apperror.Code]int{
		apperror.ErrBadRequest:   400,
		apperror.ErrValidation:   400,
		apperror.ErrUnauthorized: 401,
		apperror.ErrForbidden:    403,
		apperror.ErrNotFound:     404,
		apperror.ErrConflict:     409,
		apperror.ErrDB:           500,
		apperror.ErrInternal:     500,
	}

	for code, expected := range cases {
		if got := HTTPStatusFromCode(code); got != expected {
			t.Fatalf("code %s expected %d got %d", code, expected, got)
		}
	}

	if got := HTTPStatusFromCode("UNKNOWN_CODE"); got != 500 {
		t.Fatalf("expected fallback 500 got %d", got)
	}
}
