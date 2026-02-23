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

var (
	loadConfig = config.Load
	buildApp   = func(cfg *config.Config) (server, func(), error) {
		return app.Build(cfg)
	}
	notifySignal = signal.Notify
	withTimeout  = context.WithTimeout
	runMain      = run
	fatalf       = log.Fatalf
)

func main() {
	if err := runMain(); err != nil {
		fatalf("%v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config load error: %w", err)
	}

	fbApp, closer, err := buildApp(cfg)
	if err != nil {
		return fmt.Errorf("app build error: %w", err)
	}
	defer func() {
		if closer != nil {
			closer()
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
	}

	log.Println("shutdown signal received")
	ctx, cancel := withTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fbApp.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("fiber shutdown error: %w", err)
	}
	log.Println("server gracefully stopped")
	return nil
}
