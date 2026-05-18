package plainbase

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var reAlnumSpace = regexp.MustCompile(`^[a-zA-Z0-9]+(?: [a-zA-Z0-9]+)*$`)

func ValidateAlphanumWithSpaceRule(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if value == "" {
		return true
	}
	return reAlnumSpace.MatchString(value)
}

func ValidateNotBlankRule(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(value) != ""
}

func ValidateTrimRule(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return value == strings.TrimSpace(value)
}
