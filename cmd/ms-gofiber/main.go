package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms-gofiber/internal/app"
	"ms-gofiber/internal/config"
)

type server interface {
	Listen(addr string) error
	ShutdownWithContext(ctx context.Context) error
}

type closeFunc = app.CloseFunc

var (
	loadConfig = config.Load
	buildApp   = func(ctx context.Context, cfg *config.Config) (server, closeFunc, error) {
		return app.Build(ctx, cfg)
	}
	notifySignal = signal.Notify
	withTimeout  = context.WithTimeout
	runMain      = runBackground
	fatalf       = log.Fatalf
)

func runBackground() error {
	return run(context.Background())
}

func main() {
	if err := runMain(); err != nil {
		fatalf("%v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config load error: %w", err)
	}

	fbApp, closer, err := buildApp(ctx, cfg)
	if err != nil {
		return fmt.Errorf("app build error: %w", err)
	}
	defer func() {
		if closer != nil {
			if closeErr := closer(); closeErr != nil {
				log.Printf("app close error: %v", closeErr)
			}
		}
	}()

	listenErr := make(chan error, 1)
	go func() {
		log.Printf("server starting at %s", cfg.ListenAddr())
		listenErr <- fbApp.Listen(cfg.ListenAddr())
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	notifySignal(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		if err != nil {
			return fmt.Errorf("fiber listen error: %w", err)
		}
	case <-quit:
	case <-ctx.Done():
	}

	log.Println("shutdown signal received")
	shutdownCtx, cancel := withTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fbApp.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("fiber shutdown error: %w", err)
	}

	closeErr := closer()
	closer = nil
	if closeErr != nil {
		return fmt.Errorf("app close error: %w", closeErr)
	}

	log.Println("server gracefully stopped")
	return nil
}
