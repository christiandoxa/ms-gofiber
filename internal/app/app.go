package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"

	"github.com/christiandoxa/welog"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/api/handler"
	"ms-gofiber/api/middleware"
	apirouter "ms-gofiber/api/router"
	apivalidation "ms-gofiber/api/validation"
	redisadapter "ms-gofiber/internal/app/adapter/repository/redis"
	sqliteadapter "ms-gofiber/internal/app/adapter/repository/sqlite"
	todousecase "ms-gofiber/internal/app/application/usecase"
	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/config"
	"ms-gofiber/pkg/cache"
	"ms-gofiber/pkg/db"
	"ms-gofiber/pkg/logging"
)

func Build(ctx context.Context, cfg *config.Config) (*fiber.App, func() error, error) {
	sqliteDB, err := db.NewSQLiteDB(ctx, cfg.SQLitePath)
	if err != nil {
		return nil, nil, err
	}

	redisClient, err := cache.NewRedis(ctx, cache.RedisOptions{
		Addr:        cfg.RedisAddr,
		Password:    cfg.RedisPassword,
		DB:          cfg.RedisDB,
		PingTimeout: cfg.RedisPingTimeout(),
	})
	if err != nil {
		return nil, nil, errors.Join(err, sqliteDB.Close())
	}

	closeApp := func() error {
		return errors.Join(sqliteDB.Close(), redisClient.Close())
	}

	validate, err := apivalidation.NewStructValidator()
	if err != nil {
		closeErr := closeApp()
		return nil, nil, errors.Join(err, closeErr)
	}

	fiberApp := newFiberApp(cfg)
	registerMiddleware(fiberApp, validate, redisClient, cfg.RedisDefaultTTLDuration())
	registerRoutes(fiberApp, validate, sqliteDB, redisClient, cfg.RedisDefaultTTLDuration())

	return fiberApp, closeApp, nil
}

func newFiberApp(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout(),
		WriteTimeout: cfg.WriteTimeout(),
		ErrorHandler: middleware.ErrorHandler(),
	})
}

func registerMiddleware(
	app *fiber.App,
	validate middleware.RequestValidator,
	redisClient *redis.Client,
	externalIDTTL time.Duration,
) {
	app.Use(cors.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.SecurityHeaders())
	app.Use(apmfiber.Middleware())
	app.Use(welog.NewFiber(fiber.Config{}))
	app.Use(middleware.HeaderGuard(validate))
	app.Use(middleware.ExternalIDGuard(redisClient, externalIDTTL))
}

func registerRoutes(
	app *fiber.App,
	validate handler.RequestValidator,
	sqliteDB *sql.DB,
	redisClient *redis.Client,
	cacheTTL time.Duration,
) {
	todoUC := todousecase.NewTodo(
		sqliteadapter.NewTodo(sqliteDB),
		redisadapter.NewTodo(redisClient),
		cacheTTL,
		reportTodoCacheError,
	)
	todoHandler := handler.NewTodo(todoUC, validate)
	internalHandler := &handler.Internal{}
	validationHandler := handler.NewValidation(validate)

	apirouter.Register(app, todoHandler, internalHandler, validationHandler)
}

func reportTodoCacheError(ctx context.Context, operation string, id domain.TodoID, err error) {
	logging.Warn(ctx, err, "todo cache operation failed", map[string]any{
		"operation": operation,
		"todo_id":   id,
	})
}
