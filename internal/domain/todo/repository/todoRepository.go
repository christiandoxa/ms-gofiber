package repository

import (
	"context"

	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/pkg/infrastructure/database"
)

type ITodoRepository interface {
	Create(ctx context.Context, todo *todomodel.Todo) (*todomodel.Todo, error)
	Get(ctx context.Context, id string) (*todomodel.Todo, error)
	List(ctx context.Context) ([]*todomodel.Todo, error)
	Update(ctx context.Context, todo *todomodel.Todo) (*todomodel.Todo, error)
	Delete(ctx context.Context, id string) error
}

type TodoRepository struct {
	db *database.DB
}

func New(db *database.DB) ITodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	record, err := r.db.CreateTodo(ctx, toRecord(todo))
	if err != nil {
		return nil, err
	}
	return toModel(record), nil
}

func (r *TodoRepository) Get(ctx context.Context, id string) (*todomodel.Todo, error) {
	record, ok, err := r.db.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, todomodel.ErrTodoNotFound
	}
	return toModel(record), nil
}

func (r *TodoRepository) List(ctx context.Context) ([]*todomodel.Todo, error) {
	records, err := r.db.ListTodos(ctx)
	if err != nil {
		return nil, err
	}
	todos := make([]*todomodel.Todo, 0, len(records))
	for _, record := range records {
		todos = append(todos, toModel(record))
	}
	return todos, nil
}

func (r *TodoRepository) Update(ctx context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	record, ok, err := r.db.UpdateTodo(ctx, toRecord(todo))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, todomodel.ErrTodoNotFound
	}
	return toModel(record), nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	ok, err := r.db.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return todomodel.ErrTodoNotFound
	}
	return nil
}

func toRecord(todo *todomodel.Todo) database.TodoRecord {
	return database.TodoRecord{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func toModel(record database.TodoRecord) *todomodel.Todo {
	return &todomodel.Todo{
		ID:        record.ID,
		Title:     record.Title,
		Completed: record.Completed,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}
