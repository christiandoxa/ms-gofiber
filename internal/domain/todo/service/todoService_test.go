package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/pkg/apperror"
)

var errExpected = errors.New("repo")

type fakeTodoRepository struct {
	now time.Time
}

func (f fakeTodoRepository) Create(_ context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	if todo.ID == "" || todo.CreatedAt.IsZero() || todo.UpdatedAt.IsZero() {
		return nil, errExpected
	}
	return todo, nil
}

func (f fakeTodoRepository) Get(_ context.Context, id string) (*todomodel.Todo, error) {
	switch id {
	case "missing":
		return nil, todomodel.ErrTodoNotFound
	case "error":
		return nil, errExpected
	default:
		return &todomodel.Todo{ID: id, Title: "old", CreatedAt: f.now, UpdatedAt: f.now}, nil
	}
}

func (fakeTodoRepository) List(context.Context) ([]*todomodel.Todo, error) {
	return []*todomodel.Todo{{ID: "1"}}, nil
}

func (fakeTodoRepository) Update(_ context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	return todo, nil
}

func (fakeTodoRepository) Delete(_ context.Context, id string) error {
	switch id {
	case "missing":
		return todomodel.ErrTodoNotFound
	case "error":
		return errExpected
	default:
		return nil
	}
}

func newTestService() ITodoService {
	return New(fakeTodoRepository{now: time.Now().UTC()})
}

func TestTodoServiceCreateAndList(t *testing.T) {
	service := newTestService()
	if todo, err := service.Create(context.Background(), "title", false); err != nil || todo.Title != "title" {
		t.Fatalf("unexpected create result: %+v %v", todo, err)
	}
	if list, err := service.List(context.Background()); err != nil || len(list) != 1 {
		t.Fatalf("unexpected list result: %+v %v", list, err)
	}
}

func TestTodoServiceGet(t *testing.T) {
	service := newTestService()
	if todo, err := service.Get(context.Background(), "1"); err != nil || todo.ID != "1" {
		t.Fatalf("unexpected get result: %+v %v", todo, err)
	}
	var appErr *apperror.Error
	if _, err := service.Get(context.Background(), "missing"); !errors.As(err, &appErr) || appErr.Status != http.StatusNotFound {
		t.Fatalf("expected not found app error, got %v", err)
	}
	if _, err := service.Get(context.Background(), "error"); !errors.Is(err, errExpected) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestTodoServiceUpdate(t *testing.T) {
	service := newTestService()
	todo, err := service.Update(context.Background(), "1", "new", true)
	if err != nil || todo.Title != "new" || !todo.Completed {
		t.Fatalf("unexpected update result: %+v %v", todo, err)
	}
	if _, err := service.Update(context.Background(), "missing", "new", true); err == nil {
		t.Fatalf("expected update error")
	}
}

func TestTodoServiceDelete(t *testing.T) {
	service := newTestService()
	if err := service.Delete(context.Background(), "1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	var appErr *apperror.Error
	if err := service.Delete(context.Background(), "missing"); !errors.As(err, &appErr) || appErr.Status != http.StatusNotFound {
		t.Fatalf("expected delete not found app error, got %v", err)
	}
	if err := service.Delete(context.Background(), "error"); !errors.Is(err, errExpected) {
		t.Fatalf("expected delete repo error, got %v", err)
	}
}
