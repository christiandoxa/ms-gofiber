package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.elastic.co/apm/v2"

	"ms-gofiber/internal/app/domain"
)

var rowsAffected = func(res sql.Result) (int64, error) {
	return res.RowsAffected()
}

// Todo is sqlite implementation of todo repository
type Todo struct {
	db *sql.DB
}

// NewTodo is constructor of sqlite todo repository
func NewTodo(db *sql.DB) *Todo {
	return &Todo{db: db}
}

func (r *Todo) Create(ctx context.Context, t *domain.Todo) (domain.TodoID, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Create", "dbrepo")
	defer span.End()

	const q = `INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(
		ctx,
		q,
		t.ID,
		t.Title,
		boolToInt(t.Completed),
		t.CreatedAt.UTC().Format(time.RFC3339Nano),
		t.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	return t.ID, err
}

func (r *Todo) GetByID(ctx context.Context, id domain.TodoID) (*domain.Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.GetByID", "dbrepo")
	defer span.End()

	const q = `SELECT id, title, completed, created_at, updated_at FROM todos WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, id)

	t, err := scanTodo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *Todo) List(ctx context.Context, limit, offset int) ([]*domain.Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.List", "dbrepo")
	defer span.End()

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	const q = `SELECT id, title, completed, created_at, updated_at FROM todos ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*domain.Todo
	for rows.Next() {
		t, scanErr := scanTodo(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *Todo) Update(ctx context.Context, t *domain.Todo) error {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Update", "dbrepo")
	defer span.End()

	const q = `UPDATE todos SET title = ?, completed = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(
		ctx,
		q,
		t.Title,
		boolToInt(t.Completed),
		t.UpdatedAt.UTC().Format(time.RFC3339Nano),
		t.ID,
	)
	if err != nil {
		return err
	}

	affected, err := rowsAffected(res)
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrTodoNotFound
	}
	return nil
}

func (r *Todo) Delete(ctx context.Context, id domain.TodoID) error {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Delete", "dbrepo")
	defer span.End()

	const q = `DELETE FROM todos WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	affected, err := rowsAffected(res)
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrTodoNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTodo(s scanner) (*domain.Todo, error) {
	var (
		id        string
		title     string
		completed int
		createdAt string
		updatedAt string
	)
	if err := s.Scan(&id, &title, &completed, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	createdAtTime, err := parseRFC3339(createdAt)
	if err != nil {
		return nil, err
	}
	updatedAtTime, err := parseRFC3339(updatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:        domain.TodoID(id),
		Title:     title,
		Completed: completed != 0,
		CreatedAt: createdAtTime,
		UpdatedAt: updatedAtTime,
	}, nil
}

func parseRFC3339(v string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, v)
	if err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, v)
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
