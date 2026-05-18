package logging

import (
	"context"

	weloglogger "github.com/christiandoxa/welog/pkg/infrastructure/logger"
)

// Fatalf logs a fatal message through welog and terminates the process.
//
//go:noinline
func Fatalf(format string, v ...any) {
	weloglogger.Logger().Fatalf(format, v...)
}

func Info(ctx context.Context, message string, fields map[string]any) {
	entry := weloglogger.Logger().WithContext(ctx)
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	entry.Info(message)
}

func Error(ctx context.Context, err error, message string, fields map[string]any) {
	entry := weloglogger.Logger().WithContext(ctx).WithError(err)
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	entry.Error(message)
}

func Warn(ctx context.Context, err error, message string, fields map[string]any) {
	entry := weloglogger.Logger().WithContext(ctx).WithError(err)
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	entry.Warn(message)
}
