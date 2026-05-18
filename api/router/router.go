package router

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/api/handler"
)

func Register(app *fiber.App, todo *handler.Todo, internal *handler.Internal, validation *handler.Validation) {
	api := app.Group("/v1")

	todos := api.Group("/todos")
	todos.Post("/", todo.Create)
	todos.Get("/", todo.List)
	todos.Get("/:id", todo.Get)
	todos.Put("/:id", todo.Update)
	todos.Delete("/:id", todo.Delete)

	api.Get("/internal/echo", internal.Echo)
	api.Get("/client/self-call", internal.SelfCall)
	api.Post("/internal/validator/prepare-example", validation.PrepareExample)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
