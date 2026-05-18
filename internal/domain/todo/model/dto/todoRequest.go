package dto

type TodoRequest struct {
	Title     string `json:"title" validate:"required,max=120"`
	Completed bool   `json:"completed"`
}
