package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"

	tododomain "ms-gofiber/internal/domain/todo"
	"ms-gofiber/internal/dto"
	"ms-gofiber/internal/transport/http/presenter"
	todousecase "ms-gofiber/internal/usecase/todo"
	"ms-gofiber/internal/validator"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/respond"
)

type TodoHandler struct {
	svc      todousecase.Service
	validate RequestValidator
}

type RequestValidator interface {
	ValidateStruct(any) error
}

func NewTodoHandler(svc todousecase.Service, v RequestValidator) *TodoHandler {
	return &TodoHandler{svc: svc, validate: v}
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Create", "handler")
	defer span.End()

	var req dto.TodoUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		c.Locals("logger").(*logrus.Entry).Error("Something went wrong")
		apm.CaptureError(ctx, err).Send()
		return apperror.New(apperror.ErrBadRequest, "invalid JSON body")
	}
	if err := h.validate.ValidateStruct(req); err != nil {
		return err
	}

	t := &tododomain.Todo{Title: req.Title, Completed: req.Completed}
	out, err := h.svc.Create(ctx, t)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *TodoHandler) Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Get", "handler")
	defer span.End()

	id := c.Params("id")
	if id == "" {
		return apperror.New(apperror.ErrBadRequest, "missing id")
	}

	out, err := h.svc.Get(ctx, tododomain.ID(id))
	if err != nil {
		return err
	}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *TodoHandler) List(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, _ := apm.StartSpan(ctx, "TodoHandler.List", "handler")
	defer span.End()

	limit, offset, err := validator.ValidatePagination(c.Query("limit"), c.Query("offset"), 100)
	if err != nil {
		return err
	}
	out, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		return err
	}
	meta := map[string]any{"limit": limit, "offset": offset, "ts": time.Now().UTC()}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoListData(out), meta))
}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Update", "handler")
	defer span.End()

	id := c.Params("id")
	if id == "" {
		return apperror.New(apperror.ErrBadRequest, "missing id")
	}

	var req dto.TodoUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		apm.CaptureError(ctx, err).Send()
		return apperror.New(apperror.ErrBadRequest, "invalid JSON body")
	}
	if err := h.validate.ValidateStruct(req); err != nil {
		return err
	}

	t := &tododomain.Todo{ID: tododomain.ID(id), Title: req.Title, Completed: req.Completed}
	out, err := h.svc.Update(ctx, t)
	if err != nil {
		return err
	}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Delete", "handler")
	defer span.End()

	id := c.Params("id")
	if id == "" {
		return apperror.New(apperror.ErrBadRequest, "missing id")
	}
	if err := h.svc.Delete(ctx, tododomain.ID(id)); err != nil {
		return err
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
