package structhelper

import "strings"

func IsTrimmed(value string) bool {
	return value == strings.TrimSpace(value)
}

func IsNotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}
