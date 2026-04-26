package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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
	assertStatus(t, app, jsonRequest(http.MethodPost, "/todos", "{"), http.StatusBadRequest)

	// validator error
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(*domain.Todo) (*domain.Todo, error) {
		return nil, nil
	}}, mockValidator{err: apperror.New(apperror.ErrValidation, "invalid")})
	app.Post("/todos", h.Create)
	req := jsonRequest(http.MethodPost, "/todos", `{"title":"abc","completed":false}`)
	assertStatus(t, app, req, http.StatusBadRequest)

	// usecase error
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(*domain.Todo) (*domain.Todo, error) {
		return nil, apperror.New(apperror.ErrInternal, "boom")
	}}, mockValidator{})
	app.Post("/todos", h.Create)
	req = jsonRequest(http.MethodPost, "/todos", `{"title":"abc","completed":false}`)
	assertStatus(t, app, req, http.StatusInternalServerError)

	// success
	now := time.Now().UTC()
	app = setupControllerApp()
	h = NewTodo(mockTodoUC{create: func(td *domain.Todo) (*domain.Todo, error) {
		return &domain.Todo{ID: "1", Title: td.Title, Completed: td.Completed, CreatedAt: now, UpdatedAt: now}, nil
	}}, mockValidator{})
	app.Post("/todos", h.Create)
	req = jsonRequest(http.MethodPost, "/todos", `{"title":"abc","completed":true}`)
	assertStatus(t, app, req, http.StatusCreated)
}

func TestTodoControllerGet(t *testing.T) {
	app := todoFixture(t)

	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos/e", nil), http.StatusNotFound)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos-get", nil), http.StatusBadRequest)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos/1", nil), http.StatusOK)
}

func TestTodoControllerList(t *testing.T) {
	app := todoFixture(t)

	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos?limit=abc", nil), http.StatusBadRequest)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos?limit=99", nil), http.StatusInternalServerError)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos?limit=10&offset=0", nil), http.StatusOK)
}

func TestTodoControllerUpdate(t *testing.T) {
	app := todoFixture(t)

	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/1", "{"), http.StatusBadRequest)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/e", `{"title":"abc","completed":false}`), http.StatusInternalServerError)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos-validation/1", `{"title":"abc","completed":false}`), http.StatusBadRequest)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/1", `{"title":"abc","completed":false}`), http.StatusOK)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos-update", `{"title":"abc","completed":false}`), http.StatusBadRequest)
}

func TestTodoControllerDelete(t *testing.T) {
	app := todoFixture(t)

	assertStatus(t, app, httptest.NewRequest(http.MethodDelete, "/todos/e", nil), http.StatusInternalServerError)
	assertStatus(t, app, httptest.NewRequest(http.MethodDelete, "/todos/1", nil), http.StatusNoContent)
	assertStatus(t, app, httptest.NewRequest(http.MethodDelete, "/todos-delete", nil), http.StatusBadRequest)
}

func TestTodoControllerResponseBody(t *testing.T) {
	app := todoFixture(t)
	res := performRequest(t, app, httptest.NewRequest(http.MethodGet, "/todos/1", nil))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, res.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
	if body["code"] != "OK" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func todoFixture(t *testing.T) *fiber.App {
	t.Helper()

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

	return app
}

func assertStatus(t *testing.T, app *fiber.App, req *http.Request, expected int) {
	t.Helper()

	res := performRequest(t, app, req)
	if res.StatusCode != expected {
		t.Fatalf("expected %d got %d", expected, res.StatusCode)
	}
	if err := res.Body.Close(); err != nil {
		t.Fatalf("close response body: %v", err)
	}
}

func performRequest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	t.Helper()

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return res
}

func jsonRequest(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

var _ = errors.New
