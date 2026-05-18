package router

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
	middlewarehandler "ms-gofiber/handler/middleware"
	todomodel "ms-gofiber/internal/domain/todo/model"
)

type mockTodoService struct{}

func (mockTodoService) Create(context.Context, string, bool) (*todomodel.Todo, error) {
	now := time.Now().UTC()
	return &todomodel.Todo{ID: "1", Title: "title", CreatedAt: now, UpdatedAt: now}, nil
}
func (mockTodoService) Get(_ context.Context, id string) (*todomodel.Todo, error) {
	now := time.Now().UTC()
	return &todomodel.Todo{ID: id, Title: "title", CreatedAt: now, UpdatedAt: now}, nil
}
func (mockTodoService) List(context.Context) ([]*todomodel.Todo, error) {
	return []*todomodel.Todo{}, nil
}
func (mockTodoService) Update(_ context.Context, id string, title string, completed bool) (*todomodel.Todo, error) {
	now := time.Now().UTC()
	return &todomodel.Todo{ID: id, Title: title, Completed: completed, CreatedAt: now, UpdatedAt: now}, nil
}
func (mockTodoService) Delete(context.Context, string) error { return nil }

func TestRegister(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	Register(app, &model.Service{
		RequestValidator: mockValidator{},
		TodoService:      mockTodoService{},
	})

	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}
}

type mockValidator struct{}

func (mockValidator) ValidateStruct(any) error { return nil }

func TestRegisterTodoRoutes(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	app.Use(middlewarehandler.CheckHeader(&model.Service{RequestValidator: mockValidator{}}))
	Register(app, &model.Service{
		RequestValidator: mockValidator{},
		TodoService:      mockTodoService{},
	})

	req := httptest.NewRequest(fiber.MethodGet, "/v1/todos", nil)
	req.Header.Set("X-CLIENT-ID", "client")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close body: %v", err)
	}
}
