package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"scylla/adapter"
	"scylla/docs"
	"scylla/handler"
	"scylla/pkg/config"
	"scylla/pkg/connection"
	"scylla/pkg/exception"
	middlewares "scylla/pkg/middleware"
	"scylla/pkg/utils"
	"scylla/repo"
	"scylla/service"
)

//	@title			Boilerplate API
//	@version		1.0
//	@description	Boilerplate API in Go using Echo framework

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	conf := config.Get()

	//engine postgres
	db := connection.GetDatabase(conf.Database)

	//validate
	validate := utils.InitializeValidator()

	//adapter
	_, err := adapter.InitObsAdapter(conf.Obs)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	//environment swagger
	if conf.Swagger.Mode != "dev" {
		docs.SwaggerInfo.Host = conf.Swagger.Host
		docs.SwaggerInfo.BasePath = conf.Swagger.Url
	} else {
		docs.SwaggerInfo.Host = "localhost:3000"
		docs.SwaggerInfo.BasePath = "/api/v1"
	}
	//init repo
	customerRepo := repo.NewCustomerRepoImpl(db)
	//init service
	customerService := service.NewCustomerServiceImpl(customerRepo, validate)
	//init handler
	customerHandler := handler.NewCustomerHandler(customerService)

	//echo
	app := echo.New()
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 8, // 1 KB
		LogLevel:  log.ERROR,
	}))
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
	customerHandler.Route(app)

	//docs swagger
	app.GET("docs/*", echoSwagger.WrapHandler)

	//start
	app.Logger.Fatal(app.Start(":" + conf.Server.Port))
}
