package rcerror

import (
	"errors"
	"net/http"

	"ms-gofiber/pkg/responsecode/model"
)

var (
	ErrInvalidFieldFormat    = rcmodel.NewResponseCode(http.StatusBadRequest, "invalid field format", "4000001")
	ErrInvalidMandatoryField = rcmodel.NewResponseCode(http.StatusBadRequest, "invalid mandatory field", "4000002")

	ErrDuplicateExternalID = rcmodel.NewResponseCode(http.StatusConflict, "duplicate external id", "4090001")

	ErrGeneral = rcmodel.NewResponseCode(http.StatusInternalServerError, "general error", "5000000")
	ErrTimeout = rcmodel.NewResponseCode(http.StatusRequestTimeout, "timeout", "5040000")

	ErrDuplicateCacheKey = errors.New("cache key already exists")
	ErrFailedToStoreData = errors.New("failed to store data")
)
