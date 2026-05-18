package model

import (
	requestvalidatorservice "ms-gofiber/internal/domain/reqvalidator/service"
	todoservice "ms-gofiber/internal/domain/todo/service"
)

type Service struct {
	RequestValidator requestvalidatorservice.IRequestValidator
	TodoService      todoservice.ITodoService
}
