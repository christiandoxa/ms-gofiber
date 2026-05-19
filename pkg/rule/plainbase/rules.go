package plainbase

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/pkg/rule/util"
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

func ValidateTimeRule(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	expression := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	return util.ValidateRegex(value, expression)
}
