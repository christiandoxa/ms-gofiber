package handler

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/response"
	"ms-gofiber/pkg/responsecode/model"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if responseCode, ok := errors.AsType[*rcmodel.ResponseCode](err); ok {
			return c.Status(responseCode.StatusCode).JSON(responseCode)
		}

		if appErr, ok := errors.AsType[*apperror.Error](err); ok {
			return c.Status(appErr.Status).JSON(response.Error(appErr.Message, appErr.Fields))
		}

		if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
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
