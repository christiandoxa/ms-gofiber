package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/christiandoxa/welog"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/internal/app/adapter/controller"
	redisadapter "ms-gofiber/internal/app/adapter/repository/redis"
	sqliteadapter "ms-gofiber/internal/app/adapter/repository/sqlite"
	appvalidation "ms-gofiber/internal/app/adapter/validation"
	todousecase "ms-gofiber/internal/app/application/usecase"
	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/config"
	"ms-gofiber/internal/middleware"
	"ms-gofiber/internal/validator"
	"ms-gofiber/pkg/cache"
	"ms-gofiber/pkg/db"
)

type CloseFunc func() error

var newValidator = func() (controller.RequestValidator, error) {
	return validator.NewStructValidator(appvalidation.RegisterStructRules)
}

func Build(ctx context.Context, cfg *config.Config) (*fiber.App, CloseFunc, error) {
	if cfg == nil {
		return nil, nil, errors.New("config is nil")
	}

	infra, err := buildInfrastructure(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	validate, err := newValidator()
	if err != nil {
		closeErr := infra.Close()
		return nil, nil, errors.Join(err, closeErr)
	}

	fiberApp := newFiberApp(cfg)
	registerMiddleware(fiberApp, validate, infra.redisClient, cfg.RedisDefaultTTLDuration())
	registerRoutes(fiberApp, validate, infra, cfg.RedisDefaultTTLDuration())

	return fiberApp, infra.Close, nil
}

type infrastructure struct {
	sqliteDB    *sql.DB
	redisClient *redis.Client
}

func buildInfrastructure(ctx context.Context, cfg *config.Config) (*infrastructure, error) {
	sqliteDB, err := db.NewSQLiteDB(ctx, db.SQLiteOptions{Path: cfg.SQLitePath})
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

	return &infrastructure{
		sqliteDB:    sqliteDB,
		redisClient: redisClient,
	}, nil
}

func (i *infrastructure) Close() error {
	return errors.Join(i.sqliteDB.Close(), i.redisClient.Close())
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
	skippedPaths := middleware.DefaultSkippedPaths()

	app.Use(cors.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.SecurityHeaders())
	app.Use(apmfiber.Middleware())
	app.Use(welog.NewFiber(fiber.Config{}))
	app.Use(middleware.HeaderGuard(validate, skippedPaths))
	app.Use(middleware.ExternalIDGuard(redisClient, externalIDTTL, skippedPaths))
}

func registerRoutes(
	app *fiber.App,
	validate controller.RequestValidator,
	infra *infrastructure,
	cacheTTL time.Duration,
) {
	todoRepo := sqliteadapter.NewTodo(infra.sqliteDB)
	todoCache := redisadapter.NewTodo(infra.redisClient)
	todoUC := todousecase.NewTodo(
		todoRepo,
		todoCache,
		cacheTTL,
		todousecase.WithCacheErrorReporter(reportTodoCacheError),
	)
	todoController := controller.NewTodo(todoUC, validate)
	internalController := controller.NewInternal()
	validationController := controller.NewValidation(validate)

	router := controller.NewRouter(app, todoController, internalController, validationController)
	router.RegisterRoutes()
}

func reportTodoCacheError(ctx context.Context, operation string, id domain.TodoID, err error) {
	logrus.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
		"operation": operation,
		"todo_id":   id,
	}).Warn("todo cache operation failed")
}
