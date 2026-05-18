package app

import (
	"context"
	"database/sql"
	"errors"
	"io"
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

type Runtime struct {
	*fiber.App
	resources []io.Closer
}

func NewRuntime(fiberApp *fiber.App, resources ...io.Closer) *Runtime {
	return &Runtime{
		App:       fiberApp,
		resources: resources,
	}
}

func (r *Runtime) Close() error {
	return closeAll(r.resources...)
}

func closeAll(resources ...io.Closer) error {
	var err error
	for _, resource := range resources {
		err = errors.Join(err, resource.Close())
	}
	return err
}

func Build(ctx context.Context, cfg *config.Config) (*Runtime, error) {
	sqliteDB, err := db.NewSQLiteDB(ctx, cfg.SQLitePath)
	if err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedis(ctx, cache.RedisOptions{
		Addr:        cfg.RedisAddr,
		Password:    cfg.RedisPassword,
		DB:          cfg.RedisDB,
		PingTimeout: cfg.RedisPingTimeout(),
	})
	if err != nil {
		return nil, errors.Join(err, sqliteDB.Close())
	}

	validate, err := apivalidation.NewStructValidator()
	if err != nil {
		return nil, errors.Join(err, closeAll(sqliteDB, redisClient))
	}

	fiberApp := newFiberApp(cfg)
	registerMiddleware(fiberApp, validate, redisClient, cfg.RedisDefaultTTLDuration())
	registerRoutes(fiberApp, validate, sqliteDB, redisClient, cfg.RedisDefaultTTLDuration())

	return NewRuntime(fiberApp, sqliteDB, redisClient), nil
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
