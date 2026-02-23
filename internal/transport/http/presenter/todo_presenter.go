package presenter

import (
	"time"

	"ms-gofiber/internal/domain/todo"
)

type Todo struct {
	ID        todo.ID   `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TodoData(in *todo.Todo) Todo {
	if in == nil {
		return Todo{}
	}
	return Todo{
		ID:        in.ID,
		Title:     in.Title,
		Completed: in.Completed,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func TodoListData(in []*todo.Todo) []Todo {
	out := make([]Todo, 0, len(in))
	for _, item := range in {
		out = append(out, TodoData(item))
	}
	return out
}
