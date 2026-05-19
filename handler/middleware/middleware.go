package handler

import (
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/constant/generalkey"
)

var (
	allowedSkippedPath = map[string]bool{
		"/v1/flush/cache": true,
		"/v1/health":      true,
	}
)

func CheckHeader(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if allowedSkippedPath[c.Path()] {
			return c.Next()
		}
		header := model.Header{}
		if err := c.ReqHeaderParser(&header); err != nil {
			return apperror.New(fiber.StatusBadRequest, "invalid request headers")
		}
		if err := service.RequestValidator.ValidateStruct(header); err != nil {
			return err
		}
		c.Locals(generalkey.RequestHeader, header)
		return c.Next()
	}
}

func ExternalID(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if allowedSkippedPath[c.Path()] {
			return c.Next()
		}
		if err := service.ExternalIDService.StoreExternalID(c.UserContext(), c.Get(generalkey.ExternalID)); err != nil {
			return err
		}
		return c.Next()
	}
}
