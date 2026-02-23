package controller

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app        *fiber.App
	todo       *Todo
	internal   *Internal
	validation *Validation
}

func NewRouter(app *fiber.App, todo *Todo, internal *Internal, validation *Validation) *Router {
	return &Router{
		app:        app,
		todo:       todo,
		internal:   internal,
		validation: validation,
	}
}

func (r *Router) RegisterRoutes() {
	api := r.app.Group("/v1")

	todos := api.Group("/todos")
	todos.Post("/", r.todo.Create)
	todos.Get("/", r.todo.List)
	todos.Get("/:id", r.todo.Get)
	todos.Put("/:id", r.todo.Update)
	todos.Delete("/:id", r.todo.Delete)

	api.Get("/internal/echo", r.internal.Echo)
	api.Get("/client/self-call", r.internal.SelfCall)
	api.Post("/internal/validator/prepare-example", r.validation.PrepareExample)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
