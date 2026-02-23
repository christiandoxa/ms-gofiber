package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"ms-gofiber/internal/app/domain"
	mw "ms-gofiber/internal/middleware"
	"ms-gofiber/pkg/apperror"
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
func (m mockTodoUC) Get(_ context.Context, id domain.TodoID) (*domain.Todo, error) { return m.get(id) }
func (m mockTodoUC) List(_ context.Context, limit, offset int) ([]*domain.Todo, error) {
	return m.list(limit, offset)
}
func (m mockTodoUC) Update(_ context.Context, in *domain.Todo) (*domain.Todo, error) {
	return m.update(in)
}
func (m mockTodoUC) Delete(_ context.Context, id domain.TodoID) error { return m.delete(id) }

type mockValidator struct{ err error }

func (m mockValidator) ValidateStruct(any) error { return m.err }

func setupControllerApp() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: mw.ErrorHandler()})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logrus.NewEntry(logrus.New()))
		return c.Next()
	})
	return app
}

func TestTodoControllerCreate(t *testing.T) {
	// invalid body
	app := setupControllerApp()
	h := NewTodo(mockTodoUC{create: func(*domain.Todo) (*domain.Todo, error) {
		return nil, nil
	}}, mockValidator{})
	app.Post("/todos", h.Create)
	res, err := app.Test(httptest.NewRequest("POST", "/todos", strings.NewReader("{")))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// validator error
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(*domain.Todo) (*domain.Todo, error) {
		return nil, nil
	}}, mockValidator{err: apperror.New(apperror.ErrValidation, "invalid")})
	app.Post("/todos", h.Create)
	req := httptest.NewRequest("POST", "/todos", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// usecase error
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(*domain.Todo) (*domain.Todo, error) {
		return nil, apperror.New(apperror.ErrInternal, "boom")
	}}, mockValidator{})
	app.Post("/todos", h.Create)
	req = httptest.NewRequest("POST", "/todos", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}

	// success
	now := time.Now().UTC()
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(td *domain.Todo) (*domain.Todo, error) {
		return &domain.Todo{ID: "1", Title: td.Title, Completed: td.Completed, CreatedAt: now, UpdatedAt: now}, nil
	}}, mockValidator{})
	app.Post("/todos", h.Create)
	req = httptest.NewRequest("POST", "/todos", strings.NewReader(`{"title":"abc","completed":true}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 201 {
		t.Fatalf("expected 201 got %d", res.StatusCode)
	}
}

func TestTodoControllerGetListUpdateDelete(t *testing.T) {
	now := time.Now().UTC()
	uc := mockTodoUC{
		get: func(id domain.TodoID) (*domain.Todo, error) {
			if id == "e" {
				return nil, apperror.New(apperror.ErrNotFound, "not found")
			}
			return &domain.Todo{ID: id, Title: "x", CreatedAt: now, UpdatedAt: now}, nil
		},
		list: func(limit, offset int) ([]*domain.Todo, error) {
			if limit == 99 {
				return nil, apperror.New(apperror.ErrInternal, "db")
			}
			return []*domain.Todo{{ID: "1", Title: "x", CreatedAt: now, UpdatedAt: now}}, nil
		},
		update: func(td *domain.Todo) (*domain.Todo, error) {
			if td.ID == "e" {
				return nil, apperror.New(apperror.ErrInternal, "db")
			}
			td.UpdatedAt = now
			return td, nil
		},
		delete: func(id domain.TodoID) error {
			if id == "e" {
				return apperror.New(apperror.ErrInternal, "db")
			}
			return nil
		},
	}

	app := setupControllerApp()
	h := NewTodo(uc, mockValidator{})
	hBadValidator := NewTodo(uc, mockValidator{err: apperror.New(apperror.ErrValidation, "invalid")})
	app.Get("/todos/:id", h.Get)
	app.Get("/todos", h.List)
	app.Put("/todos/:id", h.Update)
	app.Put("/todos-update", h.Update) // cover missing id branch
	app.Put("/todos-validation/:id", hBadValidator.Update)
	app.Delete("/todos/:id", h.Delete)
	app.Delete("/todos-delete", h.Delete) // cover missing id branch
	app.Get("/todos-get", h.Get)          // cover missing id branch

	// get not found branch
	res, err := app.Test(httptest.NewRequest("GET", "/todos/e", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 404 { // mapped from apperror.ErrNotFound
		t.Fatalf("expected 404 got %d", res.StatusCode)
	}
	// get missing id branch
	res, err = app.Test(httptest.NewRequest("GET", "/todos-get", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	res, err = app.Test(httptest.NewRequest("GET", "/todos/1", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	// list pagination error
	res, err = app.Test(httptest.NewRequest("GET", "/todos?limit=abc", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}
	// list usecase error
	res, err = app.Test(httptest.NewRequest("GET", "/todos?limit=99", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}
	// list success
	res, err = app.Test(httptest.NewRequest("GET", "/todos?limit=10&offset=0", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	// update invalid body
	req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}
	// update error
	req = httptest.NewRequest("PUT", "/todos/e", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}
	// update validation error
	req = httptest.NewRequest("PUT", "/todos-validation/1", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}
	// update success
	req = httptest.NewRequest("PUT", "/todos/1", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	// update missing id
	req = httptest.NewRequest("PUT", "/todos-update", strings.NewReader(`{"title":"abc","completed":false}`))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// delete error
	res, err = app.Test(httptest.NewRequest("DELETE", "/todos/e", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 500 {
		t.Fatalf("expected 500 got %d", res.StatusCode)
	}
	// delete success
	res, err = app.Test(httptest.NewRequest("DELETE", "/todos/1", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 204 {
		t.Fatalf("expected 204 got %d", res.StatusCode)
	}
	// delete missing id
	res, err = app.Test(httptest.NewRequest("DELETE", "/todos-delete", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 got %d", res.StatusCode)
	}

	// quick sanity decode one response to touch JSON body path
	res, err = app.Test(httptest.NewRequest("GET", "/todos/1", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if body["code"] != "OK" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

var _ = errors.New
