package rule

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/internal/dto"
	"ms-gofiber/internal/validator/rule/plainbase"
	"ms-gofiber/internal/validator/rule/structbase"
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

var customStructRules = []fiber.Map{
	{
		"func": structbase.TodoUpsertStructRule,
		"type": dto.TodoUpsertRequest{},
	},
	{
		"func": structbase.PrepareExampleStructRule,
		"type": dto.PrepareExampleRequest{},
	},
}

func RegisterRule(validate *validator.Validate) error {
	for name, fn := range customRules {
		if err := validate.RegisterValidation(name, fn); err != nil {
			return err
		}
	}
	for _, r := range customStructRules {
		validate.RegisterStructValidation(r["func"].(func(validator.StructLevel)), r["type"])
	}
	return nil
}
