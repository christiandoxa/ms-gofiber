package domain

import (
	"errors"
	"time"
)

type TodoID string

type Todo struct {
	ID        TodoID    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrTodoNotFound = errors.New("todo not found")
