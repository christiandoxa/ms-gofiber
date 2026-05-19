package model

import (
	echoservice "ms-gofiber/external/domain/echo/service"
	cacheservice "ms-gofiber/internal/domain/cache/service"
	externalidservice "ms-gofiber/internal/domain/externalid/service"
	remappingservice "ms-gofiber/internal/domain/remapping/service"
	requestvalidatorservice "ms-gofiber/internal/domain/reqvalidator/service"
	todoservice "ms-gofiber/internal/domain/todo/service"
)

type Service struct {
	CacheService      cacheservice.ICacheService
	EchoService       echoservice.IEchoService
	ExternalIDService externalidservice.IExternalIDService
	RemappingService  remappingservice.IRemappingService
	RequestValidator  requestvalidatorservice.IRequestValidator
	TodoService       todoservice.ITodoService
}
