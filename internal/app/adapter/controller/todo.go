package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"

	"ms-gofiber/internal/app/adapter/presenter"
	"ms-gofiber/internal/app/application/usecase"
	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/dto"
	"ms-gofiber/internal/validator"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/respond"
)

type Todo struct {
	usecase  usecase.ITodo
	validate RequestValidator
}

func NewTodo(usecase usecase.ITodo, v RequestValidator) *Todo {
	return &Todo{
		usecase:  usecase,
		validate: v,
	}
}

func (h *Todo) Create(c *fiber.Ctx) error {
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

	t := &domain.Todo{Title: req.Title, Completed: req.Completed}
	out, err := h.usecase.Create(ctx, t)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *Todo) Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Get", "handler")
	defer span.End()

	id := c.Params("id")
	if id == "" {
		return apperror.New(apperror.ErrBadRequest, "missing id")
	}

	out, err := h.usecase.Get(ctx, domain.TodoID(id))
	if err != nil {
		return err
	}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *Todo) List(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, _ := apm.StartSpan(ctx, "TodoHandler.List", "handler")
	defer span.End()

	limit, offset, err := validator.ValidatePagination(c.Query("limit"), c.Query("offset"), 100)
	if err != nil {
		return err
	}
	out, err := h.usecase.List(ctx, limit, offset)
	if err != nil {
		return err
	}
	meta := map[string]any{"limit": limit, "offset": offset, "ts": time.Now().UTC()}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoListData(out), meta))
}

func (h *Todo) Update(c *fiber.Ctx) error {
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

	t := &domain.Todo{ID: domain.TodoID(id), Title: req.Title, Completed: req.Completed}
	out, err := h.usecase.Update(ctx, t)
	if err != nil {
		return err
	}
	return c.JSON(respond.SuccessEnvelope(presenter.TodoData(out), nil))
}

func (h *Todo) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "TodoHandler.Delete", "handler")
	defer span.End()

	id := c.Params("id")
	if id == "" {
		return apperror.New(apperror.ErrBadRequest, "missing id")
	}
	if err := h.usecase.Delete(ctx, domain.TodoID(id)); err != nil {
		return err
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
