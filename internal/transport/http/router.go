package http

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/transport/http/handlers"
	"ms-gofiber/internal/transport/http/routes"
)

type Router struct {
	app               *fiber.App
	todoHandler       *handlers.TodoHandler
	internalHandler   *handlers.InternalHandler
	validationHandler *handlers.ValidationHandler
}

func NewRouter(
	app *fiber.App,
	todoHandler *handlers.TodoHandler,
	internalHandler *handlers.InternalHandler,
	validationHandler *handlers.ValidationHandler,
) *Router {
	return &Router{
		app:               app,
		todoHandler:       todoHandler,
		internalHandler:   internalHandler,
		validationHandler: validationHandler,
	}
}

func (r *Router) RegisterRoutes() {
	api := r.app.Group("/v1")

	routes.RegisterTodoRoutes(api, r.todoHandler)
	routes.RegisterSystemRoutes(api, r.internalHandler)
	routes.RegisterValidationRoutes(api, r.validationHandler)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
