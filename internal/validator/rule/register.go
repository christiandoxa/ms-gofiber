package rule

import (
	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/app/adapter/dto"
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

type structRule struct {
	fn     validator.StructLevelFunc
	target any
}

var customStructRules = []structRule{
	{
		fn:     structbase.TodoUpsertStructRule,
		target: dto.TodoUpsertRequest{},
	},
	{
		fn:     structbase.PrepareExampleStructRule,
		target: dto.PrepareExampleRequest{},
	},
}

func RegisterRule(validate *validator.Validate) error {
	for name, fn := range customRules {
		if err := validate.RegisterValidation(name, fn); err != nil {
			return err
		}
	}
	for _, r := range customStructRules {
		validate.RegisterStructValidation(r.fn, r.target)
	}
	return nil
}
