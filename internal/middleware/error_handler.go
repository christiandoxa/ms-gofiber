package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"

	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/respond"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		logger := loggerFromContext(c)
		ctx := c.UserContext()

		var aerr *apperror.Error
		if errors.As(err, &aerr) {
			status := respond.HTTPStatusFromCode(aerr.Code)
			if status >= 500 {
				logger.WithError(aerr).Error("server error")
				apm.CaptureError(ctx, aerr).Send()
			}
			return c.Status(status).JSON(respond.ErrorEnvelope(aerr.Code, aerr.Message, aerr.Fields))
		}

		var fe *fiber.Error
		if errors.As(err, &fe) {
			if fe.Code >= 500 {
				logger.WithError(err).Error("fiber error")
				apm.CaptureError(ctx, fe).Send()
			}
			return c.Status(fe.Code).JSON(respond.ErrorEnvelope(apperror.ErrBadRequest, fe.Message, nil))
		}

		logger.WithError(err).Error("unknown error")
		apm.CaptureError(ctx, err).Send()
		return c.Status(fiber.StatusInternalServerError).
			JSON(respond.ErrorEnvelope(apperror.ErrInternal, "internal server error", nil))
	}
}

func loggerFromContext(c *fiber.Ctx) *logrus.Entry {
	if logger, ok := c.Locals("logger").(*logrus.Entry); ok && logger != nil {
		return logger
	}
	return logrus.NewEntry(logrus.StandardLogger())
}
