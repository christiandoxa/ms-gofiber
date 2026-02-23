package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/christiandoxa/welog"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/internal/app/adapter/controller"
	redisadapter "ms-gofiber/internal/app/adapter/repository/redis"
	sqliteadapter "ms-gofiber/internal/app/adapter/repository/sqlite"
	todousecase "ms-gofiber/internal/app/application/usecase"
	"ms-gofiber/internal/config"
	"ms-gofiber/internal/middleware"
	"ms-gofiber/internal/validator"
	"ms-gofiber/pkg/cache"
	"ms-gofiber/pkg/db"
)

func Build(cfg *config.Config) (*fiber.App, func(), error) {
	// SQLite
	sqliteDB, err := db.NewSQLiteDB(context.Background(), cfg)
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

	// Middleware contoh dari project referensi:
	// 1) validasi header wajib
	// 2) guard duplikasi X-EXTERNAL-ID via redis
	skippedPaths := map[string]struct{}{
		"/v1/health":           {},
		"/v1/internal/echo":    {},
		"/v1/client/self-call": {},
	}
	app.Use(middleware.HeaderGuard(validate, skippedPaths))
	app.Use(middleware.ExternalIDGuard(redisClient, time.Duration(cfg.RedisDefaultTTL)*time.Second, skippedPaths))

	// Dependency wiring
	todoRepo := sqliteadapter.NewTodo(sqliteDB)
	todoCache := redisadapter.NewTodo(redisClient)
	todoUC := todousecase.NewTodo(todoRepo, todoCache, time.Duration(cfg.RedisDefaultTTL)*time.Second)
	todoController := controller.NewTodo(todoUC, validate)
	internalController := controller.NewInternal()
	validationController := controller.NewValidation(validate)

	// Router
	router := controller.NewRouter(app, todoController, internalController, validationController)
	router.RegisterRoutes()

	// closer
	closer := func() {
		_ = sqliteDB.Close()
		_ = redisClient.Close()
	}
	return app, closer, nil
}
