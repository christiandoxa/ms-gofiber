package handler

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	"ms-gofiber/pkg/responsecode/message"
)

func Flush(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := service.CacheService.Flush(c.UserContext()); err != nil {
			return err
		}
		return c.Status(rcsuccess.GeneralSuccess.StatusCode).JSON(rcsuccess.GeneralSuccess)
	}
}
