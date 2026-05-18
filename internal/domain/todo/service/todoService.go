package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/internal/domain/todo/repository"
	"ms-gofiber/pkg/apperror"
)

type ITodoService interface {
	Create(ctx context.Context, title string, completed bool) (*todomodel.Todo, error)
	Get(ctx context.Context, id string) (*todomodel.Todo, error)
	List(ctx context.Context) ([]*todomodel.Todo, error)
	Update(ctx context.Context, id string, title string, completed bool) (*todomodel.Todo, error)
	Delete(ctx context.Context, id string) error
}

type TodoService struct {
	todoRepository repository.ITodoRepository
}

func New(todoRepository repository.ITodoRepository) ITodoService {
	return &TodoService{todoRepository: todoRepository}
}

func (s *TodoService) Create(ctx context.Context, title string, completed bool) (*todomodel.Todo, error) {
	now := time.Now().UTC()
	todo := &todomodel.Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.todoRepository.Create(ctx, todo)
}

func (s *TodoService) Get(ctx context.Context, id string) (*todomodel.Todo, error) {
	todo, err := s.todoRepository.Get(ctx, id)
	if errors.Is(err, todomodel.ErrTodoNotFound) {
		return nil, apperror.New(http.StatusNotFound, "todo not found")
	}
	return todo, err
}

func (s *TodoService) List(ctx context.Context) ([]*todomodel.Todo, error) {
	return s.todoRepository.List(ctx)
}

func (s *TodoService) Update(ctx context.Context, id string, title string, completed bool) (*todomodel.Todo, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	current.Title = title
	current.Completed = completed
	current.UpdatedAt = time.Now().UTC()
	return s.todoRepository.Update(ctx, current)
}

func (s *TodoService) Delete(ctx context.Context, id string) error {
	err := s.todoRepository.Delete(ctx, id)
	if errors.Is(err, todomodel.ErrTodoNotFound) {
		return apperror.New(http.StatusNotFound, "todo not found")
	}
	return err
}
