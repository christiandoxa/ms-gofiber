package model

import "github.com/gofiber/fiber/v2"

type Router struct {
	ExternalRouter fiber.Router
	TodoRouter     fiber.Router
}
