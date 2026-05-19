package rule

import (
	"github.com/go-playground/validator/v10"

	"ms-gofiber/pkg/constant/rulekey"
	"ms-gofiber/pkg/rule/plainbase"
	"ms-gofiber/pkg/rule/structbase"
)

var customRules = map[string]validator.Func{
	rulekey.AlphanumWithSpaceRule: plainbase.ValidateAlphanumWithSpaceRule,
	rulekey.NotBlankRule:          plainbase.ValidateNotBlankRule,
	rulekey.TimeRule:              plainbase.ValidateTimeRule,
	rulekey.TrimRule:              plainbase.ValidateTrimRule,
}

var registerValidation = func(validate *validator.Validate, rule string, function validator.Func) error {
	return validate.RegisterValidation(rule, function)
}

func RegisterRule(validate *validator.Validate) error {
	for rule, function := range customRules {
		if err := registerValidation(validate, rule, function); err != nil {
			return err
		}
	}
	structbase.RegisterRule(validate)
	return nil
}
