package utils

import (
	"github.com/labstack/echo/v4"
	"scylla/dto"
)

func ResponseInterceptor(ctx echo.Context, resp *dto.Response) {
	traceId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	resp.TraceID = traceId
}

func ErrorInterceptor(ctx echo.Context, resp *dto.Error) {
	traceId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	resp.TraceID = traceId
}
