package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"scylla/entity"
)

func NotFoundMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err == echo.ErrNotFound {
			return c.JSON(http.StatusNotFound, entity.Error{
				Code:    http.StatusNotFound,
				Status:  "NOT FOUND",
				Errors:  "Endpoint not found",
				TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
			})
		}
		return err
	}
}
