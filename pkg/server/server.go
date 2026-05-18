package server

import (
	"github.com/christiandoxa/welog"
	"github.com/christiandoxa/welog/pkg/infrastructure/logger"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber/v2"

	"ms-gofiber/cmd/app/model"
	echorepository "ms-gofiber/external/domain/echo/repository"
	echoservice "ms-gofiber/external/domain/echo/service"
	errorhandler "ms-gofiber/handler/error"
	middlewarehandler "ms-gofiber/handler/middleware"
	requestvalidatorservice "ms-gofiber/internal/domain/reqvalidator/service"
	todorepository "ms-gofiber/internal/domain/todo/repository"
	todoservice "ms-gofiber/internal/domain/todo/service"
	"ms-gofiber/pkg/config"
	"ms-gofiber/pkg/constant/envkey"
	"ms-gofiber/pkg/infrastructure/database"
	"ms-gofiber/pkg/rule"
	"ms-gofiber/router"
	"os"
)

var (
	fiberConfig   fiber.Config
	recoverConfig recover.Config
	validate      *validator.Validate
	service       *model.Service
)

var (
	registerRule = rule.RegisterRule
	fatal        = func(args ...any) {
		logger.Logger().Fatal(args...)
	}
)

func init() {
	// load .env file
	cfg := config.Load()

	// init welog config
	welogConfig := welog.Config{
		ElasticIndex:    os.Getenv(envkey.ElasticIndex),
		ElasticURL:      os.Getenv(envkey.ElasticURL),
		ElasticUsername: os.Getenv(envkey.ElasticUsername),
		ElasticPassword: os.Getenv(envkey.ElasticPassword),
	}

	// set welog config
	welog.SetConfig(welogConfig)

	// init database
	db := database.Connect(cfg.DatabasePath)

	// init rule
	validate = initValidator()

	// init repository
	echoRepository := echorepository.New()
	todoRepository := todorepository.New(db)

	// init service
	echoService := echoservice.New(echoRepository)
	requestValidatorService := requestvalidatorservice.New(validate)
	todoService := todoservice.New(todoRepository)

	// init service model
	service = &model.Service{
		EchoService:      echoService,
		RequestValidator: requestValidatorService,
		TodoService:      todoService,
	}

	// set fiber config
	fiberConfig = fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
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
	app.Use(middlewarehandler.RequestID())

	// init logger handler
	app.Use(welog.NewFiber(fiberConfig))

	// init security headers
	app.Use(middlewarehandler.SecurityHeaders())

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

func initValidator() *validator.Validate {
	validate := validator.New()
	if err := registerRule(validate); err != nil {
		fatal(err)
	}
	return validate
}
