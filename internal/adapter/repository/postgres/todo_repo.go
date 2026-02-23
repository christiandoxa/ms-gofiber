package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.elastic.co/apm/v2"

	"ms-gofiber/internal/domain/todo"
)

type TodoRepo struct {
	pool *pgxpool.Pool
}

func NewTodoRepo(pool *pgxpool.Pool) *TodoRepo {
	return &TodoRepo{pool: pool}
}

func (r *TodoRepo) Create(ctx context.Context, t *todo.Todo) (todo.ID, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Create", "dbrepo")
	defer span.End()

	const q = `INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.Title, t.Completed, t.CreatedAt, t.UpdatedAt)
	return t.ID, err
}

func (r *TodoRepo) GetByID(ctx context.Context, id todo.ID) (*todo.Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.GetByID", "dbrepo")
	defer span.End()

	const q = `SELECT id, title, completed, created_at, updated_at FROM todos WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	var t todo.Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, todo.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepo) List(ctx context.Context, limit, offset int) ([]*todo.Todo, error) {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.List", "dbrepo")
	defer span.End()

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	const q = `SELECT id, title, completed, created_at, updated_at
			   FROM todos ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, &t)
	}
	return res, rows.Err()
}

func (r *TodoRepo) Update(ctx context.Context, t *todo.Todo) error {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Update", "dbrepo")
	defer span.End()

	const q = `UPDATE todos SET title=$2, completed=$3, updated_at=$4 WHERE id=$1`
	ct, err := r.pool.Exec(ctx, q, t.ID, t.Title, t.Completed, t.UpdatedAt)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return todo.ErrNotFound
	}
	return nil
}

func (r *TodoRepo) Delete(ctx context.Context, id todo.ID) error {
	span, ctx := apm.StartSpan(ctx, "TodoRepo.Delete", "dbrepo")
	defer span.End()

	const q = `DELETE FROM todos WHERE id=$1`
	ct, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return todo.ErrNotFound
	}
	return nil
}
