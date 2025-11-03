package structbase

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"ms-gofiber/internal/dto"
)

func TodoUpsertStructRule(sl validator.StructLevel) {
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
