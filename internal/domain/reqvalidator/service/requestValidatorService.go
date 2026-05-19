package service

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/pkg/apperror"
)

type IRequestValidator interface {
	ValidateStruct(request any) error
}

type RequestValidator struct {
	validate *validator.Validate
}

func New(validate *validator.Validate) IRequestValidator {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &RequestValidator{
		validate: validate,
	}
}

func (r *RequestValidator) ValidateStruct(request any) error {
	if err := r.validate.Struct(request); err != nil {
		if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
			fields := map[string]string{}
			for _, fieldError := range validationErrors {
				fields[fieldError.Field()] = fieldError.Tag()
			}
			return apperror.WithFields(http.StatusBadRequest, "validation failed", fields)
		}
		return apperror.Wrap(http.StatusBadRequest, "validation failed", err)
	}
	return nil
}
