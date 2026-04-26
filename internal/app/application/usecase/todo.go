package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"ms-gofiber/internal/app/domain"
	"ms-gofiber/internal/app/domain/repository"
	"ms-gofiber/pkg/apperror"
)

type TodoCache interface {
	GetTodo(ctx context.Context, id domain.TodoID) (*domain.Todo, bool, error)
	SetTodo(ctx context.Context, todo *domain.Todo, ttl time.Duration) error
	DeleteTodo(ctx context.Context, id domain.TodoID) error
}

type CacheErrorReporter func(ctx context.Context, operation string, id domain.TodoID, err error)

type TodoOption func(*todo)

func WithCacheErrorReporter(reporter CacheErrorReporter) TodoOption {
	return func(u *todo) {
		if reporter != nil {
			u.reportCacheError = reporter
		}
	}
}

type TodoUseCase interface {
	Create(ctx context.Context, in *domain.Todo) (*domain.Todo, error)
	Get(ctx context.Context, id domain.TodoID) (*domain.Todo, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Todo, error)
	Update(ctx context.Context, in *domain.Todo) (*domain.Todo, error)
	Delete(ctx context.Context, id domain.TodoID) error
}

type todo struct {
	repo             repository.TodoRepository
	cache            TodoCache
	ttl              time.Duration
	reportCacheError CacheErrorReporter
}

func NewTodo(repo repository.TodoRepository, cache TodoCache, defaultTTL time.Duration, opts ...TodoOption) TodoUseCase {
	u := &todo{
		repo:             repo,
		cache:            cache,
		ttl:              defaultTTL,
		reportCacheError: noopCacheErrorReporter,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func noopCacheErrorReporter(context.Context, string, domain.TodoID, error) {
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
	if u.cache != nil {
		t, found, err := u.cache.GetTodo(ctx, id)
		if err == nil && found {
			return t, nil
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
		if err := u.cache.SetTodo(ctx, t, u.ttl); err != nil {
			u.reportCacheError(ctx, "set", t.ID, err)
			return t, nil
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
		if err := u.cache.DeleteTodo(ctx, in.ID); err != nil {
			u.reportCacheError(ctx, "delete", in.ID, err)
			return in, nil
		}
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
		if err := u.cache.DeleteTodo(ctx, id); err != nil {
			u.reportCacheError(ctx, "delete", id, err)
			return nil
		}
	}
	return nil
}
