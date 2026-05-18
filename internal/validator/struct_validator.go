package validator

import (
	"errors"
	"fmt"

	v10 "github.com/go-playground/validator/v10"

	"ms-gofiber/internal/validator/rule"
	"ms-gofiber/pkg/apperror"
)

type StructValidator struct {
	v *v10.Validate
}

func NewStructValidator() (*StructValidator, error) {
	v := v10.New()

	if err := rule.RegisterRule(v); err != nil {
		return nil, fmt.Errorf("register validation rules: %w", err)
	}
	return &StructValidator{v: v}, nil
}

func (sv *StructValidator) RegisterStructValidation(fn v10.StructLevelFunc, target any) {
	sv.v.RegisterStructValidation(fn, target)
}

func (sv *StructValidator) ValidateStruct(i any) error {
	if err := sv.v.Struct(i); err != nil {
		var errs v10.ValidationErrors
		if errors.As(err, &errs) {
			fields := map[string]string{}
			for _, fe := range errs {
				fields[fe.Field()] = fe.Tag()
			}
			return apperror.WithFields(apperror.ErrValidation, "validation error", fields)
		}
		return apperror.New(apperror.ErrValidation, "validation error")
	}
	return nil
}
