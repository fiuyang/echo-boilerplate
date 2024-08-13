package exception

import (
	"fmt"
	"net/http"
	"scylla/dto"
	"scylla/pkg/utils"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ExceptionHandlers(err error, ctx echo.Context) {
	if notFoundError(err, ctx) {
		return
	} else if validationError(err, ctx) {
		return
	} else if excelValidation(err, ctx) {
		return
	} else if badRequestError(err, ctx) {
		return
	} else if unauthorizedError(err, ctx) {
		return
	} else if forbiddenError(err, ctx) {
		return
	} else {
		internalServerError(err, ctx)
		return
	}
}

func validationError(err error, ctx echo.Context) bool {

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		report := make(map[string]string)
		var fieldName string

		for _, e := range castedObject {
			if len(e.Namespace()) > 0 && unicode.IsUpper(rune(e.Namespace()[0])) {
				dotIndex := strings.Index(e.Namespace(), ".")
				if dotIndex != -1 {
					fieldName = e.Namespace()[dotIndex+1:]
				}
			} else {
				fieldName = e.Field()
			}
			switch e.Tag() {
			case "required":
				report[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				report[fieldName] = fmt.Sprintf("%s is not valid email", fieldName)
			case "gte":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "lte":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "unique":
				report[fieldName] = fmt.Sprintf("%s has already been taken", fieldName)
			case "max":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "min":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "numeric":
				report[fieldName] = fmt.Sprintf("%s value must be numeric", fieldName)
			case "number":
				report[fieldName] = fmt.Sprintf("%s value must be number", fieldName)
			case "oneof":
				report[fieldName] = fmt.Sprintf("%s value must be %s", fieldName, e.Param())
			case "len":
				report[fieldName] = fmt.Sprintf("%s value must be exactly %s characters long", fieldName, e.Param())
			case "alphanum":
				report[fieldName] = fmt.Sprintf("%s value must be char and numeric %s", fieldName, e.Param())
			case "sliceString":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is string", fieldName)
			case "dive":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty", fieldName)
			case "datetime":
				report[fieldName] = fmt.Sprintf("%s value must be date (yyyy-mm-dd)", fieldName)
			case "required_if":
				report[fieldName] = fmt.Sprintf("%s must be filled in if %s", fieldName, e.Param())
			case "sliceInt":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is int", fieldName)
			case "equal":
				report[fieldName] = fmt.Sprintf("%s and %s do not match do not match", fieldName, e.Param())
			case "image":
				report[fieldName] = fmt.Sprintf("%s file must be of type jpg, jpeg, png", fieldName)
			case "base64Image":
				report[fieldName] = fmt.Sprintf("%s value must be base64 encoded image", fieldName)
			}
		}
		webResponse := dto.Error{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Errors: report,
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return true
	}
	return false
}

func excelValidation(err error, ctx echo.Context) bool {
	exception, ok := err.(*ExcelValidation)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Errors: exception.Errors,
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return true
	}
	return false
}

func notFoundError(err error, ctx echo.Context) bool {
	exception, ok := err.(*NotFoundStruct)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Errors: exception.Error(),
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusNotFound, webResponse)
		return true
	}
	return false
}

func badRequestError(err error, ctx echo.Context) bool {
	exception, ok := err.(*BadRequestStruct)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Errors: exception.Error(),
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusBadRequest, webResponse)
		return true
	}
	return false
}

func unauthorizedError(err error, ctx echo.Context) bool {
	exception, ok := err.(*UnauthorizedStruct)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Errors: exception.Error(),
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusUnauthorized, webResponse)
		return true
	}
	return false
}

func forbiddenError(err error, ctx echo.Context) bool {
	exception, ok := err.(*ForbiddenStruct)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Errors: exception.Error(),
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusForbidden, webResponse)
		return true
	}
	return false
}

func internalServerError(err error, ctx echo.Context) bool {
	exception, ok := err.(*InternalServerErrorStruct)
	if ok {
		webResponse := dto.Error{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Errors: exception.Error(),
		}
		utils.ErrorInterceptor(ctx, &webResponse)
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return true
	}
	return false

}
