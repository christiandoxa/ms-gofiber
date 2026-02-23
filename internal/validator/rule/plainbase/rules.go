package plainbase

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	reAlnumSpace = regexp.MustCompile(`^[a-zA-Z0-9]+(?: [a-zA-Z0-9]+)*$`)
	reISO8601    = regexp.MustCompile(`^(?:[1-9]\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d(?:\.\d{1,9})?(?:Z|[+-][01]\d:[0-5]\d)$`)
)

var authorizationScope = map[string]struct{}{
	"AGREEMENT_PAY":           {},
	"USER_LOGIN_ID":           {},
	"BASE_USER_INFO":          {},
	"HASH_LOGIN_ID":           {},
	"SEND_OTP":                {},
	"PLAINTEXT_USER_LOGIN_ID": {},
}

var grantType = map[string]struct{}{
	"AUTHORIZATION_CODE": {},
	"REFRESH_TOKEN":      {},
}

var paymentMethodType = map[string]struct{}{
	"TRUEMONEY":       {},
	"ALIPAY_HK":       {},
	"TNG":             {},
	"ALIPAY_CN":       {},
	"GCASH":           {},
	"DANA":            {},
	"KAKAOPAY":        {},
	"RABBIT_LINE_PAY": {},
	"BPI":             {},
	"CONNECT_WALLET":  {},
}

var terminalType = map[string]struct{}{
	"WEB":      {},
	"WAP":      {},
	"APP":      {},
	"MINI_APP": {},
}

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
	switch v := fl.Field().Interface().(type) {
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return true
		}
		for _, p := range strings.Fields(v) {
			if _, ok := authorizationScope[p]; !ok {
				return false
			}
		}
		return true
	case []string:
		for _, p := range v {
			if _, ok := authorizationScope[p]; !ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func ValidateGrantTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	_, exists := grantType[v]
	return exists
}

func ValidatePaymentMethodTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	_, exists := paymentMethodType[v]
	return exists
}

func ValidateTerminalTypeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	_, exists := terminalType[v]
	return exists
}

func ValidateTimeRule(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if v == "" {
		return true
	}
	return reISO8601.MatchString(v)
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
