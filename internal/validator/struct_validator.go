package validator

import (
	"errors"
	"fmt"

	v10 "github.com/go-playground/validator/v10"

	"ms-gofiber/internal/validator/rule"
	"ms-gofiber/pkg/apperror"
)

type RuleRegistrar func(*v10.Validate) error

type StructValidator struct {
	v *v10.Validate
}

func NewStructValidator(registrars ...RuleRegistrar) (*StructValidator, error) {
	v := v10.New()
	allRegistrars := make([]RuleRegistrar, 0, len(registrars)+1)
	allRegistrars = append(allRegistrars, rule.RegisterRule)
	allRegistrars = append(allRegistrars, registrars...)

	for _, register := range allRegistrars {
		if err := register(v); err != nil {
			return nil, fmt.Errorf("register validation rules: %w", err)
		}
	}
	return &StructValidator{v: v}, nil
}

func (sv *StructValidator) ValidateStruct(i any) error {
	if err := sv.v.Struct(i); err != nil {
		var verrs v10.ValidationErrors
		if errors.As(err, &verrs) {
			fields := map[string]string{}
			for _, fe := range verrs {
				fields[fe.Field()] = fe.Tag()
			}
			return apperror.WithFields(apperror.ErrValidation, "validation error", fields)
		}
		return apperror.New(apperror.ErrValidation, "validation error")
	}
	return nil
}
