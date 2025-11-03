package validator

import (
	"errors"

	v10 "github.com/go-playground/validator/v10"

	"ms-gofiber/internal/validator/rule"
	"ms-gofiber/pkg/apperror"
)

type StructValidator struct {
	v *v10.Validate
}

func NewStructValidator() *StructValidator {
	v := v10.New()
	_ = rule.RegisterRule(v) // daftar semua custom rules
	return &StructValidator{v: v}
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
