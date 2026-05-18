package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ms-gofiber/internal/app/domain"
	"ms-gofiber/pkg/apperror"
)

type mockRepo struct {
	create func(context.Context, *domain.Todo) (domain.TodoID, error)
	get    func(context.Context, domain.TodoID) (*domain.Todo, error)
	list   func(context.Context, int, int) ([]*domain.Todo, error)
	update func(context.Context, *domain.Todo) error
	delete func(context.Context, domain.TodoID) error
}

func (m mockRepo) Create(ctx context.Context, t *domain.Todo) (domain.TodoID, error) {
	return m.create(ctx, t)
}
func (m mockRepo) GetByID(ctx context.Context, id domain.TodoID) (*domain.Todo, error) {
	return m.get(ctx, id)
}
func (m mockRepo) List(ctx context.Context, limit, offset int) ([]*domain.Todo, error) {
	return m.list(ctx, limit, offset)
}
func (m mockRepo) Update(ctx context.Context, t *domain.Todo) error   { return m.update(ctx, t) }
func (m mockRepo) Delete(ctx context.Context, id domain.TodoID) error { return m.delete(ctx, id) }

type mockCache struct {
	get    func(context.Context, domain.TodoID) (*domain.Todo, bool, error)
	set    func(context.Context, *domain.Todo, time.Duration) error
	delete func(context.Context, domain.TodoID) error
}

func (m mockCache) GetTodo(ctx context.Context, id domain.TodoID) (*domain.Todo, bool, error) {
	if m.get == nil {
		return nil, false, errors.New("cache miss")
	}
	return m.get(ctx, id)
}
func (m mockCache) SetTodo(ctx context.Context, todo *domain.Todo, ttl time.Duration) error {
	if m.set == nil {
		return nil
	}
	return m.set(ctx, todo, ttl)
}
func (m mockCache) DeleteTodo(ctx context.Context, id domain.TodoID) error {
	if m.delete == nil {
		return nil
	}
	return m.delete(ctx, id)
}

func assertAppErrorCode(t *testing.T, err error, code apperror.Code) {
	t.Helper()
	var aerr *apperror.Error
	if !errors.As(err, &aerr) {
		t.Fatalf("expected apperror, got %T (%v)", err, err)
	}
	if aerr.Code != code {
		t.Fatalf("expected code %s got %s", code, aerr.Code)
	}
}

func TestTodoUsecaseCreate(t *testing.T) {
	u := NewTodo(mockRepo{create: func(_ context.Context, td *domain.Todo) (domain.TodoID, error) {
		if td.ID == "" || td.CreatedAt.IsZero() || td.UpdatedAt.IsZero() {
			t.Fatalf("expected generated id and timestamps")
		}
		return td.ID, nil
	}}, nil, time.Second, nil)

	out, err := u.Create(context.Background(), &domain.Todo{Title: "x"})
	if err != nil || out.ID == "" {
		t.Fatalf("create unexpected result: out=%+v err=%v", out, err)
	}

	uErr := NewTodo(mockRepo{create: func(context.Context, *domain.Todo) (domain.TodoID, error) {
		return "", errors.New("db")
	}}, nil, time.Second, nil)
	_, err = uErr.Create(context.Background(), &domain.Todo{Title: "x"})
	assertAppErrorCode(t, err, apperror.ErrDB)
}

func TestNewTodoNilCacheErrorReporter(t *testing.T) {
	noopCacheErrorReporter(context.Background(), "set", "1", errors.New("cache"))
	NewTodo(mockRepo{}, nil, time.Second, nil)
}

func TestTodoUsecaseGet(t *testing.T) {
	now := time.Now().UTC()

	// cache hit valid JSON
	uCacheHit := NewTodo(
		mockRepo{get: func(context.Context, domain.TodoID) (*domain.Todo, error) {
			t.Fatalf("repo should not be called")
			return nil, nil
		}},
		mockCache{get: func(context.Context, domain.TodoID) (*domain.Todo, bool, error) {
			return &domain.Todo{ID: "1", Title: "a", Completed: true, CreatedAt: now, UpdatedAt: now}, true, nil
		}},
		time.Second,
		nil,
	)
	out, err := uCacheHit.Get(context.Background(), "1")
	if err != nil || out.ID != "1" {
		t.Fatalf("cache hit get failed: out=%+v err=%v", out, err)
	}

	// cache error -> fallback repo + cache set
	setCalled := false
	uFallback := NewTodo(
		mockRepo{get: func(context.Context, domain.TodoID) (*domain.Todo, error) {
			return &domain.Todo{ID: "2", Title: "b", CreatedAt: now, UpdatedAt: now}, nil
		}},
		mockCache{
			get: func(context.Context, domain.TodoID) (*domain.Todo, bool, error) {
				return nil, false, errors.New("cache error")
			},
			set: func(context.Context, *domain.Todo, time.Duration) error {
				setCalled = true
				return nil
			},
		},
		time.Second,
		nil,
	)
	out, err = uFallback.Get(context.Background(), "2")
	if err != nil || out.ID != "2" || !setCalled {
		t.Fatalf("fallback get failed: out=%+v err=%v set=%v", out, err, setCalled)
	}

	reportedSet := false
	uSetErr := NewTodo(
		mockRepo{get: func(context.Context, domain.TodoID) (*domain.Todo, error) {
			return &domain.Todo{ID: "3", Title: "c", CreatedAt: now, UpdatedAt: now}, nil
		}},
		mockCache{
			get: func(context.Context, domain.TodoID) (*domain.Todo, bool, error) {
				return nil, false, errors.New("cache error")
			},
			set: func(context.Context, *domain.Todo, time.Duration) error {
				return errors.New("cache set")
			},
		},
		time.Second,
		func(_ context.Context, operation string, id domain.TodoID, err error) {
			reportedSet = operation == "set" && id == "3" && err != nil
		},
	)
	out, err = uSetErr.Get(context.Background(), "3")
	if err != nil || out.ID != "3" {
		t.Fatalf("expected cache set error to be ignored, out=%+v err=%v", out, err)
	}
	if !reportedSet {
		t.Fatalf("expected cache set error to be reported")
	}

	uNotFound := NewTodo(mockRepo{get: func(context.Context, domain.TodoID) (*domain.Todo, error) {
		return nil, domain.ErrTodoNotFound
	}}, nil, time.Second, nil)
	_, err = uNotFound.Get(context.Background(), "x")
	assertAppErrorCode(t, err, apperror.ErrNotFound)

	uDBErr := NewTodo(mockRepo{get: func(context.Context, domain.TodoID) (*domain.Todo, error) {
		return nil, errors.New("db")
	}}, nil, time.Second, nil)
	_, err = uDBErr.Get(context.Background(), "x")
	assertAppErrorCode(t, err, apperror.ErrDB)
}

func TestTodoUsecaseListUpdateDelete(t *testing.T) {
	now := time.Now().UTC()
	u := NewTodo(
		mockRepo{
			list: func(context.Context, int, int) ([]*domain.Todo, error) {
				return []*domain.Todo{{ID: "1", Title: "a", CreatedAt: now, UpdatedAt: now}}, nil
			},
			update: func(context.Context, *domain.Todo) error { return nil },
			delete: func(context.Context, domain.TodoID) error { return nil },
		},
		mockCache{delete: func(context.Context, domain.TodoID) error { return nil }},
		time.Second,
		nil,
	)

	list, err := u.List(context.Background(), 10, 0)
	if err != nil || len(list) != 1 {
		t.Fatalf("list failed: %v %+v", err, list)
	}

	updated, err := u.Update(context.Background(), &domain.Todo{ID: "1", Title: "u"})
	if err != nil || updated.UpdatedAt.IsZero() {
		t.Fatalf("update failed: %v %+v", err, updated)
	}

	if err := u.Delete(context.Background(), "1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	reportedDelete := 0
	reportDelete := func(_ context.Context, operation string, id domain.TodoID, err error) {
		if operation == "delete" && id == "cache" && err != nil {
			reportedDelete++
		}
	}

	uUpdateCacheErr := NewTodo(
		mockRepo{update: func(context.Context, *domain.Todo) error { return nil }},
		mockCache{delete: func(context.Context, domain.TodoID) error { return errors.New("cache delete") }},
		time.Second,
		reportDelete,
	)
	if _, err := uUpdateCacheErr.Update(context.Background(), &domain.Todo{ID: "cache", Title: "u"}); err != nil {
		t.Fatalf("expected update cache delete error to be ignored, got %v", err)
	}

	uDeleteCacheErr := NewTodo(
		mockRepo{delete: func(context.Context, domain.TodoID) error { return nil }},
		mockCache{delete: func(context.Context, domain.TodoID) error { return errors.New("cache delete") }},
		time.Second,
		reportDelete,
	)
	if err := uDeleteCacheErr.Delete(context.Background(), "cache"); err != nil {
		t.Fatalf("expected delete cache error to be ignored, got %v", err)
	}
	if reportedDelete != 2 {
		t.Fatalf("expected two cache delete reports, got %d", reportedDelete)
	}

	uListErr := NewTodo(mockRepo{list: func(context.Context, int, int) ([]*domain.Todo, error) {
		return nil, errors.New("db")
	}}, nil, time.Second, nil)
	_, err = uListErr.List(context.Background(), 10, 0)
	assertAppErrorCode(t, err, apperror.ErrDB)

	uUpdateNotFound := NewTodo(mockRepo{update: func(context.Context, *domain.Todo) error {
		return domain.ErrTodoNotFound
	}}, nil, time.Second, nil)
	_, err = uUpdateNotFound.Update(context.Background(), &domain.Todo{ID: "x"})
	assertAppErrorCode(t, err, apperror.ErrNotFound)

	uUpdateErr := NewTodo(mockRepo{update: func(context.Context, *domain.Todo) error {
		return errors.New("db")
	}}, nil, time.Second, nil)
	_, err = uUpdateErr.Update(context.Background(), &domain.Todo{ID: "x"})
	assertAppErrorCode(t, err, apperror.ErrDB)

	uDeleteNotFound := NewTodo(mockRepo{delete: func(context.Context, domain.TodoID) error {
		return domain.ErrTodoNotFound
	}}, nil, time.Second, nil)
	err = uDeleteNotFound.Delete(context.Background(), "x")
	assertAppErrorCode(t, err, apperror.ErrNotFound)

	uDeleteErr := NewTodo(mockRepo{delete: func(context.Context, domain.TodoID) error {
		return errors.New("db")
	}}, nil, time.Second, nil)
	err = uDeleteErr.Delete(context.Background(), "x")
	assertAppErrorCode(t, err, apperror.ErrDB)
}
