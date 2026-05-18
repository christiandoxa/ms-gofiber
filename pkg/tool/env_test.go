package tool

import (
	"os"
	"testing"
)

func TestStringFromEnv(t *testing.T) {
	if err := os.Unsetenv("X_NOPE"); err != nil {
		t.Fatalf("unset env: %v", err)
	}
	if got := StringFromEnv("X_NOPE", "default"); got != "default" {
		t.Fatalf("expected default, got %s", got)
	}

	t.Setenv("X_VALUE", "custom")
	if got := StringFromEnv("X_VALUE", "default"); got != "custom" {
		t.Fatalf("expected env value, got %s", got)
	}
}

func TestIntFromEnv(t *testing.T) {
	if err := os.Unsetenv("X_NOPE_INT"); err != nil {
		t.Fatalf("unset env: %v", err)
	}
	got, err := IntFromEnv("X_NOPE_INT", 123)
	if err != nil || got != 123 {
		t.Fatalf("expected fallback int, got %d err=%v", got, err)
	}

	t.Setenv("X_INT", "456")
	got, err = IntFromEnv("X_INT", 123)
	if err != nil || got != 456 {
		t.Fatalf("expected env int, got %d err=%v", got, err)
	}

	t.Setenv("X_BAD_INT", "nan")
	if _, err := IntFromEnv("X_BAD_INT", 123); err == nil {
		t.Fatalf("expected invalid integer error")
	}
}
