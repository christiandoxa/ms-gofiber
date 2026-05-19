package main

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/christiandoxa/welog/pkg/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"ms-gofiber/pkg/server"
)

func TestMainFunction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		patches.ApplyFunc(server.NewServer, func() *fiber.App { return fiber.New() })
		var fiberApp *fiber.App
		patches.ApplyMethod(fiberApp, "Listen", func(*fiber.App, string) error {
			process, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatalf("find process: %v", err)
			}
			if err := process.Signal(os.Interrupt); err != nil {
				t.Fatalf("send interrupt: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		defer patches.Reset()

		main()
	})

	t.Run("listen error", func(t *testing.T) {
		fatalCalled := false
		log := logrus.New()
		log.SetOutput(io.Discard)
		log.ExitFunc = func(int) { fatalCalled = true }

		patches := gomonkey.NewPatches()
		patches.ApplyFunc(server.NewServer, func() *fiber.App { return fiber.New() })
		patches.ApplyFunc(logger.Logger, func() *logrus.Logger { return log })
		var fiberApp *fiber.App
		patches.ApplyMethod(fiberApp, "Listen", func(*fiber.App, string) error { return errors.New("listen") })
		defer patches.Reset()

		main()
		if !fatalCalled {
			t.Fatalf("expected fatal called")
		}
	})

	t.Run("shutdown error", func(t *testing.T) {
		log := logrus.New()
		log.SetOutput(io.Discard)

		patches := gomonkey.NewPatches()
		patches.ApplyFunc(server.NewServer, func() *fiber.App { return fiber.New() })
		patches.ApplyFunc(logger.Logger, func() *logrus.Logger { return log })
		var fiberApp *fiber.App
		patches.ApplyMethod(fiberApp, "Shutdown", func(*fiber.App) error { return errors.New("shutdown") })
		patches.ApplyMethod(fiberApp, "Listen", func(*fiber.App, string) error {
			process, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatalf("find process: %v", err)
			}
			if err := process.Signal(os.Interrupt); err != nil {
				t.Fatalf("send interrupt: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		defer patches.Reset()

		main()
	})
}
