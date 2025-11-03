package app

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/christiandoxa/welog"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/internal/config"
	"ms-gofiber/internal/middleware"
	"ms-gofiber/internal/transport/http"
	"ms-gofiber/internal/validator"
	"ms-gofiber/pkg/cache"
	"ms-gofiber/pkg/db"
)

func Build(cfg *config.Config) (*fiber.App, func(), error) {
	// Postgres (APM instrument di pkg/db)
	pool, err := db.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		return nil, nil, err
	}

	// Redis (APM hook di pkg/cache)
	redisClient := cache.NewRedis(cfg)

	// Validator dengan register rule custom
	validate := validator.NewStructValidator()

	// Fiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout(),
		WriteTimeout: cfg.WriteTimeout(),
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Middlewares
	app.Use(cors.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.SecurityHeaders())

	// APM inbound tracing + auto recover
	app.Use(apmfiber.Middleware())

	// Welog access log + per-request logger di c.Locals("logger")
	app.Use(welog.NewFiber(fiber.Config{}))

	// Router
	router := http.NewRouter(app, pool, redisClient, validate)
	router.RegisterRoutes()

	// closer
	closer := func() {
		pool.Close()
		_ = redisClient.Close()
	}
	return app, closer, nil
}
