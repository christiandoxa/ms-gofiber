package logging

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	weloglogger "github.com/christiandoxa/welog/pkg/infrastructure/logger"
	"github.com/sirupsen/logrus"
)

func patchLogger(t *testing.T) (*bytes.Buffer, *logrus.Logger) {
	t.Helper()

	buf := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = buf
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	patches := gomonkey.ApplyFunc(weloglogger.Logger, func() *logrus.Logger {
		return logger
	})
	t.Cleanup(patches.Reset)

	return buf, logger
}

func TestInfoErrorWarn(t *testing.T) {
	buf, _ := patchLogger(t)

	Info(context.Background(), "server starting", map[string]any{"addr": "127.0.0.1:8080"})
	if out := buf.String(); !strings.Contains(out, "server starting") || !strings.Contains(out, "addr") {
		t.Fatalf("unexpected info log: %s", out)
	}

	buf.Reset()
	Error(context.Background(), errors.New("boom"), "server error", map[string]any{"request_id": "rid"})
	if out := buf.String(); !strings.Contains(out, "server error") || !strings.Contains(out, "boom") || !strings.Contains(out, "request_id") {
		t.Fatalf("unexpected error log: %s", out)
	}

	buf.Reset()
	Warn(context.Background(), errors.New("cache"), "cache warning", map[string]any{"operation": "get"})
	if out := buf.String(); !strings.Contains(out, "cache warning") || !strings.Contains(out, "cache") || !strings.Contains(out, "operation") {
		t.Fatalf("unexpected warn log: %s", out)
	}
}

func TestFatalf(t *testing.T) {
	_, logger := patchLogger(t)

	exitCode := -1
	logger.ExitFunc = func(code int) {
		exitCode = code
	}

	Fatalf("fatal %s", "error")
	if exitCode != 1 {
		t.Fatalf("expected fatal exit code 1, got %d", exitCode)
	}
}
