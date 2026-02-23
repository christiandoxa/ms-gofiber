package routes

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/transport/http/handlers"
)

func RegisterValidationRoutes(api fiber.Router, validationHandler *handlers.ValidationHandler) {
	api.Post("/internal/validator/prepare-example", validationHandler.PrepareExample)
}
