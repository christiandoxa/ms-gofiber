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

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return startuperror.Wrap(startuperror.ConfigLoad, err)
	}

	fbApp, closeApp, err := app.Build(ctx, cfg)
	if err != nil {
		return startuperror.Wrap(startuperror.AppBuild, err)
	}
	appClosed := false
	defer func() {
		if appClosed {
			return
		}
		if closeErr := closeApp(); closeErr != nil {
			logging.Error(ctx, startuperror.Wrap(startuperror.AppClose, closeErr), "app close failed", nil)
		}
	}()

	listenErr := make(chan error, 1)
	go func() {
		logging.Info(ctx, "server starting", map[string]any{"addr": cfg.ListenAddr()})
		listenErr <- fbApp.Listen(cfg.ListenAddr())
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

	if err := fbApp.ShutdownWithContext(shutdownCtx); err != nil {
		return startuperror.Wrap(startuperror.FiberShutdown, err)
	}

	closeErr := closeApp()
	appClosed = true
	if closeErr != nil {
		return startuperror.Wrap(startuperror.AppClose, closeErr)
	}

	logging.Info(ctx, "server gracefully stopped", nil)
	return nil
}
