package validator

import (
	"strconv"

	"ms-gofiber/pkg/apperror"
)

func ValidatePagination(limitStr, offsetStr string, max int) (int, int, error) {
	limit := 10
	offset := 0

	if limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v <= 0 {
			return 0, 0, apperror.New(apperror.ErrBadRequest, "limit must be positive integer")
		}
		limit = v
	}
	if offsetStr != "" {
		v, err := strconv.Atoi(offsetStr)
		if err != nil || v < 0 {
			return 0, 0, apperror.New(apperror.ErrBadRequest, "offset must be non-negative integer")
		}
		offset = v
	}
	if limit > max {
		limit = max
	}
	return limit, offset, nil
}
