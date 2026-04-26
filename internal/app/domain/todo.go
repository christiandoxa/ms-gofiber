package domain

import (
	"errors"
	"time"
)

type TodoID string

type Todo struct {
	ID        TodoID
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrTodoNotFound = errors.New("todo not found")
