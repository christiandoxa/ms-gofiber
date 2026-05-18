package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/app"
	"ms-gofiber/internal/config"
	"ms-gofiber/internal/startuperror"
	"ms-gofiber/pkg/logging"
)

func TestRunBranches(t *testing.T) {
	cfg := &config.Config{AppHost: "127.0.0.1", AppPort: 18080}

	t.Run("config error", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(config.Load, func() (*config.Config, error) {
			return nil, errors.New("cfg")
		})
		defer patches.Reset()

		assertStartupCode(t, run(context.Background()), startuperror.ConfigLoad)
	})

	t.Run("build error", func(t *testing.T) {
		patches := patchRunDependencies(cfg, nil, nil, errors.New("build"))
		defer patches.Reset()

		assertStartupCode(t, run(context.Background()), startuperror.AppBuild)
	})

	t.Run("listen error", func(t *testing.T) {
		patches := patchRunDependencies(cfg, fiber.New(), func() error {
			return errors.New("close")
		}, nil)
		var fiberApp *fiber.App
		patches.ApplyMethod(fiberApp, "Listen", func(*fiber.App, string) error {
			return errors.New("listen")
		})
		defer patches.Reset()

		assertStartupCode(t, run(context.Background()), startuperror.FiberListen)
	})

	t.Run("shutdown error", func(t *testing.T) {
		patches := patchRunDependencies(cfg, fiber.New(), func() error { return nil }, nil)
		patchFiberLifecycle(patches, nil, errors.New("shutdown"))
		patches.ApplyFunc(signal.Notify, func(c chan<- os.Signal, sig ...os.Signal) {
			go func() { c <- os.Interrupt }()
		})
		defer patches.Reset()

		assertStartupCode(t, run(context.Background()), startuperror.FiberShutdown)
	})

	t.Run("close error", func(t *testing.T) {
		patches := patchRunDependencies(cfg, fiber.New(), func() error {
			return errors.New("close")
		}, nil)
		patchFiberLifecycle(patches, nil, nil)
		defer patches.Reset()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assertStartupCode(t, run(ctx), startuperror.AppClose)
	})

	t.Run("success", func(t *testing.T) {
		patches := patchRunDependencies(cfg, fiber.New(), func() error {
			return nil
		}, nil)
		patchFiberLifecycle(patches, nil, nil)
		defer patches.Reset()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := run(ctx); err != nil {
			t.Fatalf("expected success run, got %v", err)
		}
	})
}

func assertStartupCode(t *testing.T, err error, want startuperror.Code) {
	t.Helper()

	got, ok := startuperror.CodeOf(err)
	if !ok || got != want {
		t.Fatalf("expected startup code %s, got code=%s ok=%v err=%v", want, got, ok, err)
	}
}

func patchRunDependencies(
	cfg *config.Config,
	fiberApp *fiber.App,
	closeResource func() error,
	err error,
) *gomonkey.Patches {
	patches := gomonkey.NewPatches()
	patches.ApplyFunc(config.Load, func() (*config.Config, error) {
		return cfg, nil
	})
	patches.ApplyFunc(app.Build, func(context.Context, *config.Config) (*app.Runtime, error) {
		if err != nil {
			return nil, err
		}
		return app.NewRuntime(fiberApp, closeFunc(closeResource)), nil
	})
	return patches
}

type closeFunc func() error

func (fn closeFunc) Close() error {
	return fn()
}

func patchFiberLifecycle(patches *gomonkey.Patches, listenErr, shutdownErr error) {
	listenBlock := make(chan struct{})
	var fiberApp *fiber.App
	patches.ApplyMethod(fiberApp, "Listen", func(*fiber.App, string) error {
		<-listenBlock
		return listenErr
	})
	patches.ApplyMethod(fiberApp, "ShutdownWithContext", func(*fiber.App, context.Context) error {
		close(listenBlock)
		return shutdownErr
	})
}

func TestDefaultBuildApp(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	defer mr.Close()

	cfg := &config.Config{
		AppHost:            "127.0.0.1",
		AppPort:            18080,
		AppReadTimeout:     1,
		AppWriteTimeout:    1,
		SQLitePath:         filepath.Join(t.TempDir(), "db", "app.db"),
		RedisAddr:          mr.Addr(),
		RedisDefaultTTL:    1,
		RedisPingTimeoutMs: 10,
	}

	runtime, err := app.Build(context.Background(), cfg)
	if err != nil {
		t.Fatalf("build app: %v", err)
	}
	if runtime == nil {
		t.Fatalf("expected app runtime")
	}
	if err := runtime.Close(); err != nil {
		t.Fatalf("close app: %v", err)
	}
}

func TestMainFunction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fatalCalled := false
		patches := gomonkey.NewPatches()
		patches.ApplyFunc(run, func(context.Context) error { return nil })
		patches.ApplyFunc(logging.Fatalf, func(format string, v ...any) {
			fatalCalled = true
		})
		defer patches.Reset()

		main()
		if fatalCalled {
			t.Fatalf("fatal should not be called on nil error")
		}
	})

	t.Run("fatal", func(t *testing.T) {
		fatalCalled := false
		fatalMsg := ""
		patches := gomonkey.NewPatches()
		patches.ApplyFunc(run, func(context.Context) error { return errors.New("boom") })
		patches.ApplyFunc(logging.Fatalf, func(format string, v ...any) {
			fatalCalled = true
			fatalMsg = format
		})
		defer patches.Reset()

		main()
		if !fatalCalled || !strings.Contains(fatalMsg, "%v") {
			t.Fatalf("expected fatal called with format, called=%v msg=%s", fatalCalled, fatalMsg)
		}
	})
}
