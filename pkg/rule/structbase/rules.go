package structbase

import (
	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/domain/todo/model/dto"
	"ms-gofiber/pkg/constant/rulekey"
	"ms-gofiber/pkg/rule/structhelper"
)

func RegisterRule(validate *validator.Validate) {
	validate.RegisterStructValidation(ValidateTodoRequestRule, dto.TodoRequest{}) //nolint:errcheck
}

func ValidateTodoRequestRule(sl validator.StructLevel) {
	request, ok := sl.Current().Interface().(dto.TodoRequest)
	if !ok {
		return
	}
	if !structhelper.IsTrimmed(request.Title) {
		sl.ReportError(request.Title, "Title", "title", rulekey.TodoTitleTrimRule, "")
	}
}
