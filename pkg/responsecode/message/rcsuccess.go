package rcsuccess

import (
	"net/http"

	"ms-gofiber/pkg/responsecode/model"
)

var GeneralSuccess = rcmodel.NewResponseCode(http.StatusOK, "success", "2000000")
