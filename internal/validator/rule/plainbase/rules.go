package plainbase

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var reAlnumSpace = regexp.MustCompile(`^[A-Za-z0-9 ]+$`)

func ValidateAlphanumWithSpaceRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	return reAlnumSpace.MatchString(v)
}

func ValidateAuthorizationScopeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return true
	}
	parts := strings.Fields(v)
	allowed := map[string]struct{}{"read": {}, "write": {}}
	for _, p := range parts {
		if _, ok := allowed[p]; !ok {
			return false
		}
	}
	return true
}

func ValidateGrantTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	switch v {
	case "", "client_credentials", "authorization_code", "refresh_token":
		return true
	default:
		return false
	}
}

func ValidatePaymentMethodTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	switch v {
	case "", "CARD", "BANK_TRANSFER", "EWALLET":
		return true
	default:
		return false
	}
}

func ValidateTerminalTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	switch v {
	case "", "WEB", "MOBILE", "POS":
		return true
	default:
		return false
	}
}

func ValidateTimeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	_, err := time.Parse("15:04:05", v)
	return err == nil
}

func ValidateTrimRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return v == strings.TrimSpace(v)
}

func ValidateNotBlankRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(v) != ""
}
