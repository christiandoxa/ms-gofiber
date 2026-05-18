package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/api/dto"
	appvalidator "ms-gofiber/internal/validator"
)

func registerStructRules(validate *appvalidator.StructValidator) {
	validate.RegisterStructValidation(todoUpsertStructRule, dto.TodoUpsertRequest{})
	validate.RegisterStructValidation(prepareExampleStructRule, dto.PrepareExampleRequest{})
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
