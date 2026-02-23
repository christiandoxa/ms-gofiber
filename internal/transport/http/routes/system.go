package routes

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/transport/http/handlers"
)

func RegisterSystemRoutes(api fiber.Router, internalHandler *handlers.InternalHandler) {
	api.Get("/internal/echo", internalHandler.Echo)
	api.Get("/client/self-call", internalHandler.SelfCall)
}
