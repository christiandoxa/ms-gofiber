package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christiandoxa/welog"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/api/handler"
	mw "ms-gofiber/api/middleware"
	"ms-gofiber/internal/app/domain"
)

type mockTodoUC struct {
	create func(*domain.Todo) (*domain.Todo, error)
	get    func(domain.TodoID) (*domain.Todo, error)
	list   func(int, int) ([]*domain.Todo, error)
	update func(*domain.Todo) (*domain.Todo, error)
	delete func(domain.TodoID) error
}

func (m mockTodoUC) Create(_ context.Context, in *domain.Todo) (*domain.Todo, error) {
	return m.create(in)
}
func (m mockTodoUC) Get(_ context.Context, id domain.TodoID) (*domain.Todo, error) {
	return m.get(id)
}
func (m mockTodoUC) List(_ context.Context, limit, offset int) ([]*domain.Todo, error) {
	return m.list(limit, offset)
}
func (m mockTodoUC) Update(_ context.Context, in *domain.Todo) (*domain.Todo, error) {
	return m.update(in)
}
func (m mockTodoUC) Delete(_ context.Context, id domain.TodoID) error {
	return m.delete(id)
}

type mockValidator struct{}

func (mockValidator) ValidateStruct(any) error { return nil }

func TestRouterRegisterRoutes(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(welog.NewFiber(fiber.Config{ErrorHandler: mw.ErrorHandler()}))

	uc := mockTodoUC{
		create: func(td *domain.Todo) (*domain.Todo, error) { return td, nil },
		get:    func(id domain.TodoID) (*domain.Todo, error) { return &domain.Todo{ID: id}, nil },
		list:   func(int, int) ([]*domain.Todo, error) { return []*domain.Todo{}, nil },
		update: func(td *domain.Todo) (*domain.Todo, error) { return td, nil },
		delete: func(domain.TodoID) error { return nil },
	}

	Register(app, handler.NewTodo(uc, mockValidator{}), &handler.Internal{}, handler.NewValidation(mockValidator{}))

	assertStatus(t, app, httptest.NewRequest("GET", "/v1/health", nil), 200)
}

func assertStatus(t *testing.T, app *fiber.App, req *http.Request, expected int) {
	t.Helper()

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != expected {
		t.Fatalf("expected %d got %d", expected, res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
}
