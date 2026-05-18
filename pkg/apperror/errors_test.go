package apperror

import (
	"errors"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	if got := New(400, "bad").Error(); got != "bad" {
		t.Fatalf("unexpected error: %s", got)
	}

	err := Wrap(500, "failed", errors.New("root"))
	if !strings.Contains(err.Error(), "root") {
		t.Fatalf("expected wrapped error, got %s", err.Error())
	}

	fields := map[string]string{"title": "required"}
	fieldErr := WithFields(400, "validation failed", fields)
	if fieldErr.Fields["title"] != "required" {
		t.Fatalf("unexpected fields: %+v", fieldErr.Fields)
	}
}
