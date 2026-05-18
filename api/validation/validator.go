package validation

import appvalidator "ms-gofiber/internal/validator"

func NewStructValidator() (*appvalidator.StructValidator, error) {
	validate, err := appvalidator.NewStructValidator()
	if err != nil {
		return nil, err
	}
	registerStructRules(validate)
	return validate, nil
}
