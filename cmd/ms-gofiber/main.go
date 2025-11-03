package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms-gofiber/internal/app"
	"ms-gofiber/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	fbApp, closer, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("app build error: %v", err)
	}
	defer func() {
		if closer != nil {
			closer()
		}
	}()

	go func() {
		log.Printf("server starting at %s", cfg.ListenAddr())
		if err := fbApp.Listen(cfg.ListenAddr()); err != nil {
			log.Fatalf("fiber listen error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fbApp.ShutdownWithContext(ctx); err != nil {
		log.Printf("fiber shutdown error: %v", err)
	}
	log.Println("server gracefully stopped")
}
