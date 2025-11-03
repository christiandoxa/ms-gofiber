package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.elastic.co/apm/v2"

	"ms-gofiber/pkg/apperror"
)

type Service interface {
	Create(ctx context.Context, in *Todo) (*Todo, error)
	Get(ctx context.Context, id ID) (*Todo, error)
	List(ctx context.Context, limit, offset int) ([]*Todo, error)
	Update(ctx context.Context, in *Todo) (*Todo, error)
	Delete(ctx context.Context, id ID) error
}

type service struct {
	repo  Repository
	cache *redis.Client
	ttl   time.Duration
}

func NewService(repo Repository, cache *redis.Client, defaultTTL time.Duration) Service {
	return &service{repo: repo, cache: cache, ttl: defaultTTL}
}

func (s *service) Create(ctx context.Context, in *Todo) (*Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoService.Create", "service")
	defer span.End()

	in.ID = ID(uuid.New().String())
	in.CreatedAt = time.Now().UTC()
	in.UpdatedAt = in.CreatedAt

	id, err := s.repo.Create(ctx, in)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, apperror.Wrap(apperror.ErrDB, "failed to create todo", err)
	}
	in.ID = id
	return in, nil
}

func (s *service) Get(ctx context.Context, id ID) (*Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoService.Get", "service")
	defer span.End()

	cacheKey := fmt.Sprintf("todo:%s", id)
	if s.cache != nil {
		if b, err := s.cache.Get(ctx, cacheKey).Bytes(); err == nil && len(b) > 0 {
			var t Todo
			if e := json.Unmarshal(b, &t); e == nil {
				return &t, nil
			}
		}
	}
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, apperror.Wrap(apperror.ErrNotFound, "todo not found", err)
	}
	if s.cache != nil {
		if b, e := json.Marshal(t); e == nil {
			_ = s.cache.Set(ctx, cacheKey, b, s.ttl).Err()
		}
	}
	return t, nil
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoService.List", "service")
	defer span.End()

	ts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, apperror.Wrap(apperror.ErrDB, "failed to list todos", err)
	}
	return ts, nil
}

func (s *service) Update(ctx context.Context, in *Todo) (*Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoService.Update", "service")
	defer span.End()

	in.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, in); err != nil {
		apm.CaptureError(ctx, err).Send()
		return nil, apperror.Wrap(apperror.ErrDB, "failed to update todo", err)
	}
	if s.cache != nil {
		_ = s.cache.Del(ctx, fmt.Sprintf("todo:%s", in.ID)).Err()
	}
	return in, nil
}

func (s *service) Delete(ctx context.Context, id ID) error {
	span, ctx := apm.StartSpan(ctx, "TodoService.Delete", "service")
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		apm.CaptureError(ctx, err).Send()
		return apperror.Wrap(apperror.ErrDB, "failed to delete todo", err)
	}
	if s.cache != nil {
		_ = s.cache.Del(ctx, fmt.Sprintf("todo:%s", id)).Err()
	}
	return nil
}
