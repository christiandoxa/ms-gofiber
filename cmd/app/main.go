package main

import (
	"github.com/christiandoxa/welog/pkg/infrastructure/logger"
	"ms-gofiber/pkg/constant/envkey"
	"ms-gofiber/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// new server
	app := server.NewServer()

	// handles os signals for shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown
	go func() {
		<-c
		app.Shutdown() //nolint:errcheck
	}()

	// start the server
	if err := app.Listen(":" + os.Getenv(envkey.AppPort)); err != nil {
		logger.Logger().Fatal(err)
	}
}
