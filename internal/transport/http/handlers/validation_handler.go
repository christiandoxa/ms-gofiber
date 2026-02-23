package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"

	"ms-gofiber/internal/dto"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/respond"
)

type ValidationHandler struct {
	validate RequestValidator
}

func NewValidationHandler(v RequestValidator) *ValidationHandler {
	return &ValidationHandler{validate: v}
}

func (h *ValidationHandler) PrepareExample(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, _ := apm.StartSpan(ctx, "ValidationHandler.PrepareExample", "handler")
	defer span.End()

	var req dto.PrepareExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.ErrBadRequest, "invalid JSON body")
	}
	if err := h.validate.ValidateStruct(req); err != nil {
		return err
	}

	return c.JSON(respond.SuccessEnvelope(req, nil))
}
