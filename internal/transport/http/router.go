package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"ms-gofiber/internal/adapter/repository/postgres"
	"ms-gofiber/internal/domain/todo"
	"ms-gofiber/internal/transport/http/handlers"
	"ms-gofiber/internal/validator"
)

type Router struct {
	app      *fiber.App
	pool     *pgxpool.Pool
	redis    *redis.Client
	validate *validator.StructValidator
}

func NewRouter(app *fiber.App, pool *pgxpool.Pool, redis *redis.Client, validate *validator.StructValidator) *Router {
	return &Router{app: app, pool: pool, redis: redis, validate: validate}
}

func (r *Router) RegisterRoutes() {
	api := r.app.Group("/v1")

	repo := postgres.NewTodoRepo(r.pool)
	svc := todo.NewService(repo, r.redis, 60*time.Second)
	todoHandler := handlers.NewTodoHandler(svc, r.validate)

	todos := api.Group("/todos")
	todos.Post("/", todoHandler.Create)
	todos.Get("/", todoHandler.List)
	todos.Get("/:id", todoHandler.Get)
	todos.Put("/:id", todoHandler.Update)
	todos.Delete("/:id", todoHandler.Delete)

	// Internal + Self-call wajib
	internal := handlers.NewInternalHandler()
	api.Get("/internal/echo", internal.Echo)
	api.Get("/client/self-call", internal.SelfCall)

	// health
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
