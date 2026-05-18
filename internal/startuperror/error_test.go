package startuperror

import (
	"errors"
	"testing"
)

func TestWrapAndCodeOf(t *testing.T) {
	cause := errors.New("boom")
	err := Wrap(ConfigLoad, cause)

	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause")
	}
	if got, ok := CodeOf(err); !ok || got != ConfigLoad {
		t.Fatalf("expected code %s, got code=%s ok=%v", ConfigLoad, got, ok)
	}
	if err.Error() != "startup.config_load: boom" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
}

func TestWrapNil(t *testing.T) {
	if err := Wrap(AppBuild, nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCodeOfUnknown(t *testing.T) {
	if got, ok := CodeOf(errors.New("plain")); ok || got != "" {
		t.Fatalf("expected no startup code, got code=%s ok=%v", got, ok)
	}
}
