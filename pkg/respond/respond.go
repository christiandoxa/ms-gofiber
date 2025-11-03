package respond

import (
	"net/http"

	"ms-gofiber/pkg/apperror"
)

type Envelope struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Meta    map[string]any    `json:"meta,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func SuccessEnvelope(data any, meta map[string]any) Envelope {
	return Envelope{Code: "OK", Message: "success", Data: data, Meta: meta}
}

func ErrorEnvelope(code apperror.Code, message string, fields map[string]string) Envelope {
	return Envelope{Code: string(code), Message: message, Fields: fields}
}

var httpStatusMap = map[apperror.Code]int{
	apperror.ErrBadRequest:   http.StatusBadRequest,
	apperror.ErrValidation:   http.StatusBadRequest,
	apperror.ErrUnauthorized: http.StatusUnauthorized,
	apperror.ErrForbidden:    http.StatusForbidden,
	apperror.ErrNotFound:     http.StatusNotFound,
	apperror.ErrConflict:     http.StatusConflict,
	apperror.ErrDB:           http.StatusInternalServerError,
	apperror.ErrInternal:     http.StatusInternalServerError,
}

func HTTPStatusFromCode(code apperror.Code) int {
	if s, ok := httpStatusMap[code]; ok {
		return s
	}
	return http.StatusInternalServerError
}
