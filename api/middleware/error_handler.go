package middleware

import (
	"errors"

	"github.com/christiandoxa/welog/pkg/constant/generalkey"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"

	"ms-gofiber/api/respond"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/logging"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		ctx := c.UserContext()

		var aerr *apperror.Error
		if errors.As(err, &aerr) {
			status := respond.HTTPStatusFromCode(aerr.Code)
			if status >= 500 {
				logFiberError(c, aerr, "server error")
				apm.CaptureError(ctx, aerr).Send()
			}
			return c.Status(status).JSON(respond.ErrorEnvelope(aerr.Code, aerr.Message, aerr.Fields))
		}

		var fe *fiber.Error
		if errors.As(err, &fe) {
			if fe.Code >= 500 {
				logFiberError(c, err, "fiber error")
				apm.CaptureError(ctx, fe).Send()
			}
			return c.Status(fe.Code).JSON(respond.ErrorEnvelope(apperror.ErrBadRequest, fe.Message, nil))
		}

		logFiberError(c, err, "unknown error")
		apm.CaptureError(ctx, err).Send()
		return c.Status(fiber.StatusInternalServerError).
			JSON(respond.ErrorEnvelope(apperror.ErrInternal, "internal server error", nil))
	}
}

func logFiberError(c *fiber.Ctx, err error, message string) {
	fields := map[string]any{}
	if requestID, ok := c.Locals(generalkey.RequestID).(string); ok && requestID != "" {
		fields[generalkey.RequestID] = requestID
	}
	logging.Error(c.UserContext(), err, message, fields)
}
