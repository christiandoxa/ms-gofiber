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

func (r *TodoRepository) Create(_ context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	record := r.db.CreateTodo(toRecord(todo))
	return toModel(record), nil
}

func (r *TodoRepository) Get(_ context.Context, id string) (*todomodel.Todo, error) {
	record, ok := r.db.GetTodo(id)
	if !ok {
		return nil, todomodel.ErrTodoNotFound
	}
	return toModel(record), nil
}

func (r *TodoRepository) List(_ context.Context) ([]*todomodel.Todo, error) {
	records := r.db.ListTodos()
	todos := make([]*todomodel.Todo, 0, len(records))
	for _, record := range records {
		todos = append(todos, toModel(record))
	}
	return todos, nil
}

func (r *TodoRepository) Update(_ context.Context, todo *todomodel.Todo) (*todomodel.Todo, error) {
	record, ok := r.db.UpdateTodo(toRecord(todo))
	if !ok {
		return nil, todomodel.ErrTodoNotFound
	}
	return toModel(record), nil
}

func (r *TodoRepository) Delete(_ context.Context, id string) error {
	if !r.db.DeleteTodo(id) {
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
