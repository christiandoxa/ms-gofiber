package router

import (
	"ms-gofiber/cmd/app/model"
	"ms-gofiber/handler/echo"
	"ms-gofiber/pkg/constant/pathkey"
)

func externalRouter(router *model.Router, service *model.Service) {
	router.ExternalRouter.Get(pathkey.ExternalEchoPath, handler.Echo(service))
}
