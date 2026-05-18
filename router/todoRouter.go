package router

import (
	"ms-gofiber/cmd/app/model"
	"ms-gofiber/handler/todo"
	"ms-gofiber/pkg/constant/pathkey"
)

func todoRouter(router *model.Router, service *model.Service) {
	router.TodoRouter.Post(pathkey.TodoRootPath, handler.Create(service))
	router.TodoRouter.Get(pathkey.TodoRootPath, handler.List(service))
	router.TodoRouter.Get(pathkey.TodoIDPath, handler.Get(service))
	router.TodoRouter.Put(pathkey.TodoIDPath, handler.Update(service))
	router.TodoRouter.Delete(pathkey.TodoIDPath, handler.Delete(service))
}
