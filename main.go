package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"scylla/docs"
	"scylla/handler"
	"scylla/pkg/config"
	"scylla/pkg/exception"
	"scylla/pkg/middlewares"
	"scylla/pkg/utils"
	"scylla/repo"
	"scylla/routes"
	"scylla/usecase"
)

//	@title			Boilerplate API
//	@version		1.0
//	@description	Boilerplate API in Go using Echo framework

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	//database gorm postgres
	db := config.ConnectionGormPostgres(&loadConfig)

	//validate
	validate := utils.InitializeValidator(db)

	//environment swagger
	if loadConfig.Environment != "dev" {
		docs.SwaggerInfo.Host = loadConfig.SwaggerHost
		docs.SwaggerInfo.BasePath = loadConfig.SwaggerUrl
	} else {
		docs.SwaggerInfo.Host = "localhost:3000"
		docs.SwaggerInfo.BasePath = "/api/v1"
	}
	//init repo
	customerRepo := repo.NewCustomerRepoImpl(db)
	//init usecase
	customerUsecase := usecase.NewCustomerUsecaseImpl(customerRepo, validate)
	//init handler
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	//echo
	app := echo.New()
	app.Binder = utils.NewBindFile(app.Binder)
	app.Use(middleware.Recover())
	app.HTTPErrorHandler = exception.ExceptionHandlers
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${id}] ${remote_ip} - ${method} ${status} ${error} ${uri} - ${latency} - ${user_agent} - ${time_rfc3339_nano}\n",
	}))
	app.Use(middlewares.NotFoundMiddleware)
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
	}))

	//routes v1
	routes.NewRoutesV1(
		app,
		customerHandler,
	)

	//docs swagger
	app.GET("docs/*", echoSwagger.WrapHandler)

	//start
	app.Logger.Fatal(app.Start(":" + loadConfig.ServerPort))
}
