package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/app/domain/repository"
	"ms-gofiber/pkg/apperror"
)

// ITodoCache is interface of todo cache
type ITodoCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// ITodo is interface of todo usecase
type ITodo interface {
	Create(ctx context.Context, in *domain.Todo) (*domain.Todo, error)
	Get(ctx context.Context, id domain.TodoID) (*domain.Todo, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Todo, error)
	Update(ctx context.Context, in *domain.Todo) (*domain.Todo, error)
	Delete(ctx context.Context, id domain.TodoID) error
}

type todo struct {
	repo  repository.ITodo
	cache ITodoCache
	ttl   time.Duration
}

// NewTodo is constructor of todo usecase
func NewTodo(repo repository.ITodo, cache ITodoCache, defaultTTL time.Duration) ITodo {
	return &todo{
		repo:  repo,
		cache: cache,
		ttl:   defaultTTL,
	}
}

func (u *todo) Create(ctx context.Context, in *domain.Todo) (*domain.Todo, error) {
	in.ID = domain.TodoID(uuid.New().String())
	in.CreatedAt = time.Now().UTC()
	in.UpdatedAt = in.CreatedAt

	id, err := u.repo.Create(ctx, in)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDB, "failed to create todo", err)
	}
	in.ID = id
	return in, nil
}

func (u *todo) Get(ctx context.Context, id domain.TodoID) (*domain.Todo, error) {
	cacheKey := fmt.Sprintf("todo:%s", id)
	if u.cache != nil {
		if b, err := u.cache.Get(ctx, cacheKey); err == nil && len(b) > 0 {
			var t domain.Todo
			if e := json.Unmarshal(b, &t); e == nil {
				return &t, nil
			}
		}
	}

	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return nil, apperror.Wrap(apperror.ErrDB, "failed to get todo", err)
	}

	if u.cache != nil {
		if b, e := json.Marshal(t); e == nil {
			_ = u.cache.Set(ctx, cacheKey, b, u.ttl)
		}
	}
	return t, nil
}

func (u *todo) List(ctx context.Context, limit, offset int) ([]*domain.Todo, error) {
	todos, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDB, "failed to list todos", err)
	}
	return todos, nil
}

func (u *todo) Update(ctx context.Context, in *domain.Todo) (*domain.Todo, error) {
	in.UpdatedAt = time.Now().UTC()
	if err := u.repo.Update(ctx, in); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return nil, apperror.Wrap(apperror.ErrDB, "failed to update todo", err)
	}

	if u.cache != nil {
		_ = u.cache.Delete(ctx, fmt.Sprintf("todo:%s", in.ID))
	}
	return in, nil
}

func (u *todo) Delete(ctx context.Context, id domain.TodoID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return apperror.Wrap(apperror.ErrDB, "failed to delete todo", err)
	}

	if u.cache != nil {
		_ = u.cache.Delete(ctx, fmt.Sprintf("todo:%s", id))
	}
	return nil
}
