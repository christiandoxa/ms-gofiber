package handler

type RequestValidator interface {
	ValidateStruct(any) error
}
