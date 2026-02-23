package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"ms-gofiber/internal/app/domain"
	mw "ms-gofiber/internal/middleware"
)

func TestRouterRegisterRoutes(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})

	uc := mockTodoUC{
		create: func(td *domain.Todo) (*domain.Todo, error) { return td, nil },
		get:    func(id domain.TodoID) (*domain.Todo, error) { return &domain.Todo{ID: id}, nil },
		list:   func(int, int) ([]*domain.Todo, error) { return []*domain.Todo{}, nil },
		update: func(td *domain.Todo) (*domain.Todo, error) { return td, nil },
		delete: func(domain.TodoID) error { return nil },
	}

	r := NewRouter(app, NewTodo(uc, mockValidator{}), NewInternal(), NewValidation(mockValidator{}))
	r.RegisterRoutes()

	res, err := app.Test(httptest.NewRequest("GET", "/v1/health", nil))
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected health 200 got %d", res.StatusCode)
	}
}
