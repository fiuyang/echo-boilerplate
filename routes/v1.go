package routes

import (
	"github.com/labstack/echo/v4"
	"scylla/handler"
)

func NewRoutesV1(
	app *echo.Echo,
	customerHandler *handler.CustomerHandler,
) {
	routes := app.Group("/api/v1")
	//customer
	customerRouter := routes.Group("/customers")
	customerRouter.GET("", customerHandler.FindAllPaging)
	customerRouter.GET("/:customerId", customerHandler.FindById)
	customerRouter.GET("/export", customerHandler.Export)
	customerRouter.POST("/import", customerHandler.Import)
	customerRouter.POST("", customerHandler.Create)
	customerRouter.POST("/batch", customerHandler.CreateBatch)
	customerRouter.PATCH("/:customerId", customerHandler.Update)
	customerRouter.DELETE("/batch", customerHandler.DeleteBatch)

}
