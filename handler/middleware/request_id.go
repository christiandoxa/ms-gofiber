package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ms-gofiber/pkg/constant/generalkey"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(generalkey.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(generalkey.RequestIDHeader, requestID)
		c.Locals(generalkey.RequestID, requestID)
		return c.Next()
	}
}
