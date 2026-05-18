package handler

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

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/pkg/apperror"
)

type mockTodoService struct {
	create func(context.Context, string, bool) (*todomodel.Todo, error)
	get    func(context.Context, string) (*todomodel.Todo, error)
	list   func(context.Context) ([]*todomodel.Todo, error)
	update func(context.Context, string, string, bool) (*todomodel.Todo, error)
	delete func(context.Context, string) error
}

type mockValidator struct {
	err error
}

func (m mockValidator) ValidateStruct(any) error {
	return m.err
}

func (m mockTodoService) Create(ctx context.Context, title string, completed bool) (*todomodel.Todo, error) {
	return m.create(ctx, title, completed)
}
func (m mockTodoService) Get(ctx context.Context, id string) (*todomodel.Todo, error) {
	return m.get(ctx, id)
}
func (m mockTodoService) List(ctx context.Context) ([]*todomodel.Todo, error) {
	return m.list(ctx)
}
func (m mockTodoService) Update(ctx context.Context, id string, title string, completed bool) (*todomodel.Todo, error) {
	return m.update(ctx, id, title, completed)
}
func (m mockTodoService) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func setupApp(todoService mockTodoService, validate mockValidator) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: errorhandler.ErrorHandler()})
	service := &model.Service{
		RequestValidator: validate,
		TodoService:      todoService,
	}
	app.Post("/todos", Create(service))
	app.Get("/todos", List(service))
	app.Get("/todos/:id", Get(service))
	app.Put("/todos/:id", Update(service))
	app.Delete("/todos/:id", Delete(service))
	return app
}

func TestTodoHandler(t *testing.T) {
	now := time.Now().UTC()
	expectedErr := apperror.New(http.StatusInternalServerError, "failed")
	service := mockTodoService{
		create: func(_ context.Context, title string, completed bool) (*todomodel.Todo, error) {
			return &todomodel.Todo{ID: "1", Title: title, Completed: completed, CreatedAt: now, UpdatedAt: now}, nil
		},
		get: func(_ context.Context, id string) (*todomodel.Todo, error) {
			if id == "error" {
				return nil, expectedErr
			}
			return &todomodel.Todo{ID: id, Title: "title", CreatedAt: now, UpdatedAt: now}, nil
		},
		list: func(context.Context) ([]*todomodel.Todo, error) {
			return []*todomodel.Todo{{ID: "1", Title: "title", CreatedAt: now, UpdatedAt: now}}, nil
		},
		update: func(_ context.Context, id string, title string, completed bool) (*todomodel.Todo, error) {
			if id == "error" {
				return nil, expectedErr
			}
			return &todomodel.Todo{ID: id, Title: title, Completed: completed, CreatedAt: now, UpdatedAt: now}, nil
		},
		delete: func(_ context.Context, id string) error {
			if id == "error" {
				return expectedErr
			}
			return nil
		},
	}

	app := setupApp(service, mockValidator{})
	assertStatus(t, app, jsonRequest(http.MethodPost, "/todos", "{"), http.StatusBadRequest)
	assertStatus(t, setupApp(service, mockValidator{err: apperror.New(http.StatusBadRequest, "invalid")}), jsonRequest(http.MethodPost, "/todos", `{"title":"a"}`), http.StatusBadRequest)
	assertStatus(t, app, jsonRequest(http.MethodPost, "/todos", `{"title":"a","completed":true}`), http.StatusCreated)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos/1", nil), http.StatusOK)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos/error", nil), http.StatusInternalServerError)
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos", nil), http.StatusOK)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/1", "{"), http.StatusBadRequest)
	assertStatus(t, setupApp(service, mockValidator{err: apperror.New(http.StatusBadRequest, "invalid")}), jsonRequest(http.MethodPut, "/todos/1", `{"title":"a"}`), http.StatusBadRequest)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/1", `{"title":"b","completed":true}`), http.StatusOK)
	assertStatus(t, app, jsonRequest(http.MethodPut, "/todos/error", `{"title":"b","completed":true}`), http.StatusInternalServerError)
	assertStatus(t, app, httptest.NewRequest(http.MethodDelete, "/todos/1", nil), http.StatusNoContent)
	assertStatus(t, app, httptest.NewRequest(http.MethodDelete, "/todos/error", nil), http.StatusInternalServerError)
}

func TestTodoHandlerListErrorAndBody(t *testing.T) {
	now := time.Now().UTC()
	app := setupApp(mockTodoService{
		create: func(context.Context, string, bool) (*todomodel.Todo, error) { return nil, nil },
		get:    func(context.Context, string) (*todomodel.Todo, error) { return nil, nil },
		list: func(context.Context) ([]*todomodel.Todo, error) {
			return nil, apperror.New(http.StatusInternalServerError, "failed")
		},
		update: func(context.Context, string, string, bool) (*todomodel.Todo, error) { return nil, nil },
		delete: func(context.Context, string) error { return nil },
	}, mockValidator{})
	assertStatus(t, app, httptest.NewRequest(http.MethodGet, "/todos", nil), http.StatusInternalServerError)

	app = setupApp(mockTodoService{
		create: func(_ context.Context, title string, completed bool) (*todomodel.Todo, error) {
			return &todomodel.Todo{ID: "1", Title: title, Completed: completed, CreatedAt: now, UpdatedAt: now}, nil
		},
		get:    func(context.Context, string) (*todomodel.Todo, error) { return nil, nil },
		list:   func(context.Context) ([]*todomodel.Todo, error) { return nil, nil },
		update: func(context.Context, string, string, bool) (*todomodel.Todo, error) { return nil, nil },
		delete: func(context.Context, string) error { return nil },
	}, mockValidator{})
	res := performRequest(t, app, jsonRequest(http.MethodPost, "/todos", `{"title":"a"}`))
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	body := map[string]any{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "success" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestTodoHandlerCreateError(t *testing.T) {
	app := setupApp(mockTodoService{
		create: func(context.Context, string, bool) (*todomodel.Todo, error) {
			return nil, errors.New("create")
		},
		get:    func(context.Context, string) (*todomodel.Todo, error) { return nil, nil },
		list:   func(context.Context) ([]*todomodel.Todo, error) { return nil, nil },
		update: func(context.Context, string, string, bool) (*todomodel.Todo, error) { return nil, nil },
		delete: func(context.Context, string) error { return nil },
	}, mockValidator{})
	assertStatus(t, app, jsonRequest(http.MethodPost, "/todos", `{"title":"a"}`), http.StatusInternalServerError)
}

func assertStatus(t *testing.T, app *fiber.App, req *http.Request, expected int) {
	t.Helper()
	res := performRequest(t, app, req)
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("close body: %v", err)
		}
	}()
	if res.StatusCode != expected {
		t.Fatalf("expected %d got %d", expected, res.StatusCode)
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

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
