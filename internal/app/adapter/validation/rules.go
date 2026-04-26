package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/app/adapter/dto"
)

type structRule struct {
	fn     validator.StructLevelFunc
	target any
}

var customStructRules = []structRule{
	{
		fn:     todoUpsertStructRule,
		target: dto.TodoUpsertRequest{},
	},
	{
		fn:     prepareExampleStructRule,
		target: dto.PrepareExampleRequest{},
	},
}

func RegisterStructRules(validate *validator.Validate) error {
	for _, r := range customStructRules {
		validate.RegisterStructValidation(r.fn, r.target)
	}
	return nil
}

func todoUpsertStructRule(sl validator.StructLevel) {
	req, ok := sl.Current().Interface().(dto.TodoUpsertRequest)
	if !ok {
		return
	}

	if strings.TrimSpace(req.Title) != req.Title {
		sl.ReportError(req.Title, "Title", "Title", "xtrim", "")
	}
	if strings.TrimSpace(req.Title) == "" {
		sl.ReportError(req.Title, "Title", "Title", "xnotblank", "")
	}
}

func prepareExampleStructRule(sl validator.StructLevel) {
	req, ok := sl.Current().Interface().(dto.PrepareExampleRequest)
	if !ok {
		return
	}

	checkOsType(sl, req.TerminalType, req.OsType, req.OsVersion)
}
