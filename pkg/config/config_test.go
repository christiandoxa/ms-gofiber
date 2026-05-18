package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"ms-gofiber/pkg/constant/envkey"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv(envkey.AppPort, "")
	t.Setenv(envkey.DatabasePath, "")
	t.Setenv(envkey.AppReadTimeout, "")
	t.Setenv(envkey.AppWriteTimeout, "")

	cfg := Load()
	if cfg.AppPort != "8080" || cfg.DatabasePath != ":memory:" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.ReadTimeout != 10*time.Second || cfg.WriteTimeout != 10*time.Second {
		t.Fatalf("unexpected timeout config: %+v", cfg)
	}
	if os.Getenv(envkey.AppPort) != "8080" {
		t.Fatalf("expected default app port env")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv(envkey.AppPort, "9090")
	t.Setenv(envkey.DatabasePath, "data.db")
	t.Setenv(envkey.AppReadTimeout, "1s")
	t.Setenv(envkey.AppWriteTimeout, "2s")

	cfg := Load()
	if cfg.AppPort != "9090" || cfg.DatabasePath != "data.db" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.ReadTimeout != time.Second || cfg.WriteTimeout != 2*time.Second {
		t.Fatalf("unexpected timeout config: %+v", cfg)
	}
}

func TestLoadDotenv(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_PORT=7070\n"), 0o600); err != nil {
		t.Fatalf("write dotenv: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	if err := os.Unsetenv(envkey.AppPort); err != nil {
		t.Fatalf("unset app port: %v", err)
	}

	if cfg := Load(); cfg.AppPort != "7070" {
		t.Fatalf("unexpected dotenv config: %+v", cfg)
	}
}
