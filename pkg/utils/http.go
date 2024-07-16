package utils

import (
	"github.com/labstack/echo/v4"
	"scylla/entity"
)

func ResponseInterceptor(ctx echo.Context, resp *entity.Response) {
	traceId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	resp.TraceID = traceId
}

func ErrorInterceptor(ctx echo.Context, resp *entity.Error) {
	traceId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	resp.TraceID = traceId
}
