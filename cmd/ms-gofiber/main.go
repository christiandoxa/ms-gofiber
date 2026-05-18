package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms-gofiber/internal/app"
	"ms-gofiber/internal/config"
	"ms-gofiber/internal/startuperror"
	"ms-gofiber/pkg/logging"
)

func main() {
	if err := run(context.Background()); err != nil {
		logging.Fatalf("%v", err)
	}
}

func run(ctx context.Context) (err error) {
	cfg, err := config.Load()
	if err != nil {
		return startuperror.Wrap(startuperror.ConfigLoad, err)
	}

	runtime, err := app.Build(ctx, cfg)
	if err != nil {
		return startuperror.Wrap(startuperror.AppBuild, err)
	}
	defer func() {
		err = closeRuntime(ctx, runtime, err)
	}()

	listenErr := make(chan error, 1)
	go func() {
		logging.Info(ctx, "server starting", map[string]any{"addr": cfg.ListenAddr()})
		listenErr <- runtime.Listen(cfg.ListenAddr())
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		if err != nil {
			return startuperror.Wrap(startuperror.FiberListen, err)
		}
	case <-quit:
	case <-ctx.Done():
	}

	logging.Info(ctx, "shutdown signal received", nil)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := runtime.ShutdownWithContext(shutdownCtx); err != nil {
		return startuperror.Wrap(startuperror.FiberShutdown, err)
	}

	logging.Info(ctx, "server gracefully stopped", nil)
	return nil
}

func closeRuntime(ctx context.Context, runtime *app.Runtime, result error) error {
	err := runtime.Close()
	if err == nil {
		return result
	}

	wrapped := startuperror.Wrap(startuperror.AppClose, err)
	if result != nil {
		logging.Error(ctx, wrapped, "app close failed", nil)
		return result
	}
	return wrapped
}
