package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/pkg/infrastructure/database"
)

func TestTodoRepository(t *testing.T) {
	ctx := context.Background()
	db := database.Connect(":memory:")
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	}()
	repo := New(db)
	now := time.Now().UTC()

	todo, err := repo.Create(ctx, &todomodel.Todo{ID: "1", Title: "one", CreatedAt: now, UpdatedAt: now})
	if err != nil || todo.ID != "1" {
		t.Fatalf("unexpected create result: %+v %v", todo, err)
	}

	got, err := repo.Get(ctx, "1")
	if err != nil || got.Title != "one" {
		t.Fatalf("unexpected get result: %+v %v", got, err)
	}

	if _, err := repo.Get(ctx, "missing"); !errors.Is(err, todomodel.ErrTodoNotFound) {
		t.Fatalf("expected missing get error, got %v", err)
	}

	list, err := repo.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("unexpected list result: %+v %v", list, err)
	}

	got.Title = "updated"
	got.UpdatedAt = now.Add(time.Second)
	updated, err := repo.Update(ctx, got)
	if err != nil || updated.Title != "updated" {
		t.Fatalf("unexpected update result: %+v %v", updated, err)
	}

	if _, err := repo.Update(ctx, &todomodel.Todo{ID: "missing"}); !errors.Is(err, todomodel.ErrTodoNotFound) {
		t.Fatalf("expected missing update error, got %v", err)
	}
	if err := repo.Delete(ctx, "1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if err := repo.Delete(ctx, "missing"); !errors.Is(err, todomodel.ErrTodoNotFound) {
		t.Fatalf("expected missing delete error, got %v", err)
	}
}

func TestTodoRepositoryDatabaseErrors(t *testing.T) {
	ctx := context.Background()
	db := database.Connect(":memory:")
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	repo := New(db)
	todo := &todomodel.Todo{ID: "1", Title: "one", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}

	if _, err := repo.Create(ctx, todo); err == nil {
		t.Fatalf("expected create error")
	}
	if _, err := repo.Get(ctx, "1"); err == nil {
		t.Fatalf("expected get error")
	}
	if _, err := repo.List(ctx); err == nil {
		t.Fatalf("expected list error")
	}
	if _, err := repo.Update(ctx, todo); err == nil {
		t.Fatalf("expected update error")
	}
	if err := repo.Delete(ctx, "1"); err == nil {
		t.Fatalf("expected delete error")
	}
}
