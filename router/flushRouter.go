package router

import (
	"ms-gofiber/cmd/app/model"
	"ms-gofiber/handler/flush"
	"ms-gofiber/pkg/constant/pathkey"
)

func flushRouter(router *model.Router, service *model.Service) {
	router.FlushRouter.Get(pathkey.FlushCachePath, handler.Flush(service))
}
