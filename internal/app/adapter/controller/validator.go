package controller

type RequestValidator interface {
	ValidateStruct(any) error
}
