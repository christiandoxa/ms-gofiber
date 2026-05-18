package router

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	"ms-gofiber/pkg/constant/pathkey"
	"ms-gofiber/pkg/response"
)

func Register(app *fiber.App, service *model.Service) {
	apiVersion := app.Group("v1")

	router := &model.Router{
		ExternalRouter: apiVersion.Group(pathkey.ExternalBasePath),
		TodoRouter:     apiVersion.Group(pathkey.TodoBasePath),
	}

	apiVersion.Get(pathkey.HealthPath, func(c *fiber.Ctx) error {
		return c.JSON(response.Success(fiber.Map{"status": "ok"}))
	})

	externalRouter(router, service)
	todoRouter(router, service)
}
