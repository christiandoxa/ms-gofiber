package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/response"
)

func Echo(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		target := c.Query("target")
		if target == "" {
			return apperror.New(http.StatusBadRequest, "target is required")
		}
		body, err := service.EchoService.Echo(c.UserContext(), target)
		if err != nil {
			return err
		}
		return c.JSON(response.Success(fiber.Map{
			"body": body,
		}))
	}
}
