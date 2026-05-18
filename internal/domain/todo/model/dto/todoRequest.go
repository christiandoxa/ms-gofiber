package dto

type TodoRequest struct {
	Title     string `json:"title" validate:"required,max=120,notBlankRule,alphanumWithSpaceRule"`
	Completed bool   `json:"completed"`
}
