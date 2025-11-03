package dto

type TodoUpsertRequest struct {
	Title     string `json:"title" validate:"required,min=3,max=100,alphanum_with_space"`
	Completed bool   `json:"completed"`
}
