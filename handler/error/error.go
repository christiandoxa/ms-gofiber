package handler

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/response"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		appErr := &apperror.Error{}
		if errors.As(err, &appErr) {
			return c.Status(appErr.Status).JSON(response.Error(appErr.Message, appErr.Fields))
		}

		fiberErr := &fiber.Error{}
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(response.Error(fiberErr.Message, nil))
		}

		return c.Status(http.StatusInternalServerError).JSON(response.Error("internal server error", nil))
	}
}

func StackTraceHandler(_ *fiber.Ctx, e any) {
	_ = fmt.Sprintf("panic: %v\n%s", e, debug.Stack())
}

func GeneralNotFound(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(response.Error("route not found", nil))
}
