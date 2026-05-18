package model

import (
	"errors"
	"time"
)

var ErrTodoNotFound = errors.New("todo not found")

type Todo struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
