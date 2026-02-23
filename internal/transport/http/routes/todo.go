package routes

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/transport/http/handlers"
)

func RegisterTodoRoutes(api fiber.Router, todoHandler *handlers.TodoHandler) {
	todos := api.Group("/todos")
	todos.Post("/", todoHandler.Create)
	todos.Get("/", todoHandler.List)
	todos.Get("/:id", todoHandler.Get)
	todos.Put("/:id", todoHandler.Update)
	todos.Delete("/:id", todoHandler.Delete)
}
