package structhelper

import "testing"

func TestStringHelpers(t *testing.T) {
	if !IsTrimmed("value") {
		t.Fatalf("expected trimmed value")
	}
	if IsTrimmed(" value ") {
		t.Fatalf("expected untrimmed value")
	}
	if !IsNotBlank("value") {
		t.Fatalf("expected not blank value")
	}
	if IsNotBlank("  ") {
		t.Fatalf("expected blank value")
	}
}
