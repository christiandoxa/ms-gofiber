package rule

import (
	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/validator/rule/plainbase"
	"ms-gofiber/internal/validator/rulekey"
)

var customRules = map[string]validator.Func{
	rulekey.AlphanumWithSpaceRule:  plainbase.ValidateAlphanumWithSpaceRule,
	rulekey.AuthorizationScopeRule: plainbase.ValidateAuthorizationScopeRule,
	rulekey.GrantTypeRule:          plainbase.ValidateGrantTypeRule,
	rulekey.PaymentMethodTypeRule:  plainbase.ValidatePaymentMethodTypeRule,
	rulekey.TerminalTypeRule:       plainbase.ValidateTerminalTypeRule,
	rulekey.TimeRule:               plainbase.ValidateTimeRule,
	rulekey.TrimRule:               plainbase.ValidateTrimRule,
	rulekey.NotBlankRule:           plainbase.ValidateNotBlankRule,
}

func RegisterRule(validate *validator.Validate) error {
	for name, fn := range customRules {
		if err := validate.RegisterValidation(name, fn); err != nil {
			return err
		}
	}
	return nil
}
