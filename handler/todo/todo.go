package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	"ms-gofiber/internal/domain/todo/model/dto"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/response"
)

func Create(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &dto.TodoRequest{}
		if err := c.BodyParser(request); err != nil {
			return apperror.New(http.StatusBadRequest, "invalid request body")
		}
		if err := service.RequestValidator.ValidateStruct(request); err != nil {
			return err
		}
		todo, err := service.TodoService.Create(c.UserContext(), request.Title, request.Completed)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(response.Success(newTodoResponse(todo)))
	}
}

func Get(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		todo, err := service.TodoService.Get(c.UserContext(), c.Params("id"))
		if err != nil {
			return err
		}
		return c.JSON(response.Success(newTodoResponse(todo)))
	}
}

func List(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		todos, err := service.TodoService.List(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(response.Success(newTodoListResponse(todos)))
	}
}

func Update(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &dto.TodoRequest{}
		if err := c.BodyParser(request); err != nil {
			return apperror.New(http.StatusBadRequest, "invalid request body")
		}
		if err := service.RequestValidator.ValidateStruct(request); err != nil {
			return err
		}
		todo, err := service.TodoService.Update(c.UserContext(), c.Params("id"), request.Title, request.Completed)
		if err != nil {
			return err
		}
		return c.JSON(response.Success(newTodoResponse(todo)))
	}
}

func Delete(service *model.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := service.TodoService.Delete(c.UserContext(), c.Params("id")); err != nil {
			return err
		}
		return c.SendStatus(http.StatusNoContent)
	}
}
