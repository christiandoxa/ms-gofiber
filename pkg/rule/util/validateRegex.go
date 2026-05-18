package util

import "regexp"

func ValidateRegex(value string, expression string) bool {
	matched, err := regexp.MatchString(expression, value)
	return err == nil && matched
}
