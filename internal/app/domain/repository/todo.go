package repository

import (
	"context"

	"ms-gofiber/internal/app/domain"
)

// ITodo is interface of todo repository
type ITodo interface {
	Create(ctx context.Context, t *domain.Todo) (domain.TodoID, error)
	GetByID(ctx context.Context, id domain.TodoID) (*domain.Todo, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Todo, error)
	Update(ctx context.Context, t *domain.Todo) error
	Delete(ctx context.Context, id domain.TodoID) error
}
