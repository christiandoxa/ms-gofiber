package todo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ms-gofiber/internal/domain/todo"
	"ms-gofiber/pkg/apperror"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type Service interface {
	Create(ctx context.Context, in *todo.Todo) (*todo.Todo, error)
	Get(ctx context.Context, id todo.ID) (*todo.Todo, error)
	List(ctx context.Context, limit, offset int) ([]*todo.Todo, error)
	Update(ctx context.Context, in *todo.Todo) (*todo.Todo, error)
	Delete(ctx context.Context, id todo.ID) error
}

type service struct {
	repo  todo.Repository
	cache Cache
	ttl   time.Duration
}

func NewService(repo todo.Repository, cache Cache, defaultTTL time.Duration) Service {
	return &service{repo: repo, cache: cache, ttl: defaultTTL}
}

func (s *service) Create(ctx context.Context, in *todo.Todo) (*todo.Todo, error) {
	in.ID = todo.ID(uuid.New().String())
	in.CreatedAt = time.Now().UTC()
	in.UpdatedAt = in.CreatedAt

	id, err := s.repo.Create(ctx, in)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDB, "failed to create todo", err)
	}
	in.ID = id
	return in, nil
}

func (s *service) Get(ctx context.Context, id todo.ID) (*todo.Todo, error) {
	cacheKey := fmt.Sprintf("todo:%s", id)
	if s.cache != nil {
		if b, err := s.cache.Get(ctx, cacheKey); err == nil && len(b) > 0 {
			var t todo.Todo
			if e := json.Unmarshal(b, &t); e == nil {
				return &t, nil
			}
		}
	}

	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return nil, apperror.Wrap(apperror.ErrDB, "failed to get todo", err)
	}

	if s.cache != nil {
		if b, e := json.Marshal(t); e == nil {
			_ = s.cache.Set(ctx, cacheKey, b, s.ttl)
		}
	}
	return t, nil
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*todo.Todo, error) {
	ts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrDB, "failed to list todos", err)
	}
	return ts, nil
}

func (s *service) Update(ctx context.Context, in *todo.Todo) (*todo.Todo, error) {
	in.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, in); err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			return nil, apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return nil, apperror.Wrap(apperror.ErrDB, "failed to update todo", err)
	}

	if s.cache != nil {
		_ = s.cache.Delete(ctx, fmt.Sprintf("todo:%s", in.ID))
	}
	return in, nil
}

func (s *service) Delete(ctx context.Context, id todo.ID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			return apperror.New(apperror.ErrNotFound, "todo not found")
		}
		return apperror.Wrap(apperror.ErrDB, "failed to delete todo", err)
	}

	if s.cache != nil {
		_ = s.cache.Delete(ctx, fmt.Sprintf("todo:%s", id))
	}
	return nil
}
