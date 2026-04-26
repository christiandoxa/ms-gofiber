package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"

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
	infra, err := buildInfrastructure(context.Background(), cfg)
	if err != nil {
		return nil, nil, err
	}

	validate := validator.NewStructValidator()
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

	redisClient := cache.NewRedis(cache.RedisOptions{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &infrastructure{
		sqliteDB:    sqliteDB,
		redisClient: redisClient,
	}, nil
}

func (i *infrastructure) Close() {
	_ = i.sqliteDB.Close()
	_ = i.redisClient.Close()
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
	app.Use(middleware.HeaderGuard(validate, middleware.DefaultSkippedPaths()))
	app.Use(middleware.ExternalIDGuard(redisClient, externalIDTTL, middleware.DefaultSkippedPaths()))
}

func registerRoutes(
	app *fiber.App,
	validate controller.RequestValidator,
	infra *infrastructure,
	cacheTTL time.Duration,
) {
	todoRepo := sqliteadapter.NewTodo(infra.sqliteDB)
	todoCache := redisadapter.NewTodo(infra.redisClient)
	todoUC := todousecase.NewTodo(todoRepo, todoCache, cacheTTL)
	todoController := controller.NewTodo(todoUC, validate)
	internalController := controller.NewInternal()
	validationController := controller.NewValidation(validate)

	router := controller.NewRouter(app, todoController, internalController, validationController)
	router.RegisterRoutes()
}
