package tool

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func TestLoadDotenvIfExists(t *testing.T) {
	t.Run("missing env file", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if err := LoadDotenvIfExists(".env"); err != nil {
			t.Fatalf("missing .env should be optional: %v", err)
		}
	})

	t.Run("stat error", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
			return nil, errors.New("stat")
		})
		defer patches.Reset()

		if err := LoadDotenvIfExists(".env"); err == nil || !strings.Contains(err.Error(), "stat .env") {
			t.Fatalf("expected stat error, got %v", err)
		}
	})

	t.Run("load error", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if err := os.Mkdir(".env", 0o755); err != nil {
			t.Fatalf("create .env directory: %v", err)
		}
		if err := LoadDotenvIfExists(".env"); err == nil || !strings.Contains(err.Error(), "load .env") {
			t.Fatalf("expected load error, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if err := os.WriteFile(".env", []byte("TOOL_DOTENV_TEST=ok\n"), 0o644); err != nil {
			t.Fatalf("write .env: %v", err)
		}
		if err := LoadDotenvIfExists(".env"); err != nil {
			t.Fatalf("expected load success, got %v", err)
		}
	})
}
