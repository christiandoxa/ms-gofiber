package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"

	"ms-gofiber/api/dto"
	"ms-gofiber/api/respond"
	"ms-gofiber/pkg/apperror"
)

type Validation struct {
	validate RequestValidator
}

func NewValidation(v RequestValidator) *Validation {
	return &Validation{validate: v}
}

func (h *Validation) PrepareExample(c *fiber.Ctx) error {
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
