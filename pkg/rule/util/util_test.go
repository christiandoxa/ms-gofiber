package util

import "testing"

func TestValidateRegex(t *testing.T) {
	if !ValidateRegex("abc123", `^[a-z0-9]+$`) {
		t.Fatalf("expected regex match")
	}
	if ValidateRegex("abc-123", `^[a-z0-9]+$`) {
		t.Fatalf("expected regex mismatch")
	}
	if ValidateRegex("abc", `[`) {
		t.Fatalf("invalid regex should fail")
	}
}

func TestValidateArrayOfString(t *testing.T) {
	allowed := map[string]struct{}{
		"read":  {},
		"write": {},
	}
	if !ValidateArrayOfString([]string{"read", "write"}, allowed) {
		t.Fatalf("expected valid values")
	}
	if ValidateArrayOfString([]string{"read", "delete"}, allowed) {
		t.Fatalf("expected invalid values")
	}
}
