package handler

import (
	todomodel "ms-gofiber/internal/domain/todo/model"
	"ms-gofiber/internal/domain/todo/model/dto"
)

func newTodoResponse(todo *todomodel.Todo) dto.TodoResponse {
	return dto.TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func newTodoListResponse(todos []*todomodel.Todo) []dto.TodoResponse {
	response := make([]dto.TodoResponse, 0, len(todos))
	for _, todo := range todos {
		response = append(response, newTodoResponse(todo))
	}
	return response
}
