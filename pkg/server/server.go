package server

import (
	"github.com/christiandoxa/welog"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/cmd/app/model"
	errorhandler "ms-gofiber/handler/error"
	middlewarehandler "ms-gofiber/handler/middleware"
	requestvalidatorservice "ms-gofiber/internal/domain/reqvalidator/service"
	todorepository "ms-gofiber/internal/domain/todo/repository"
	todoservice "ms-gofiber/internal/domain/todo/service"
	"ms-gofiber/pkg/config"
	"ms-gofiber/pkg/infrastructure/database"
	"ms-gofiber/router"
)

var (
	fiberConfig   fiber.Config
	recoverConfig recover.Config
	validate      requestvalidatorservice.IRequestValidator
	service       *model.Service
)

func init() {
	// load .env file
	config.Load()

	// init database
	db := database.Connect()

	// init rule
	validate = requestvalidatorservice.New()

	// init repository
	todoRepository := todorepository.New(db)

	// init service
	todoService := todoservice.New(todoRepository)

	// init service model
	service = &model.Service{
		RequestValidator: validate,
		TodoService:      todoService,
	}

	// set fiber config
	fiberConfig = fiber.Config{
		ErrorHandler: errorhandler.ErrorHandler(),
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	}

	// recover config
	recoverConfig = recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: errorhandler.StackTraceHandler,
	}
}

func NewServer() *fiber.App {
	// init fiber app
	app := fiber.New(fiberConfig)

	// init request id
	app.Use(requestid.New())

	// init logger handler
	app.Use(welog.NewFiber(fiberConfig))

	// recover from panic caused by handler
	app.Use(recover.New(recoverConfig))

	// init apm
	app.Use(apmfiber.Middleware(apmfiber.WithPanicPropagation()))

	// check header
	app.Use(middlewarehandler.CheckHeader(service))

	// register routes
	router.Register(app, service)

	// handle unregistered route
	app.Use(errorhandler.GeneralNotFound)

	return app
}
