package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path/filepath"
	"scylla/entity"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/pkg/utils"
	"scylla/usecase"
	"time"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

// Note            godoc
//
// @Summary		Create customer
// @Description	Create customer.
// @Param		data	formData	entity.CreateCustomerRequest	true	"create customer"
// @Produce		application/json
// @Tags		customers
// @Success		201	{object}	entity.JsonCreated{data=nil}"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router		/customers [post]
func (handler *CustomerHandler) Create(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := new(entity.CreateCustomerRequest)
	err := ctx.Bind(request)
	helper.ErrorPanic(err)

	handler.customerUsecase.Create(c, *request)

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Successful",
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusCreated, webResponse)
}

// Note            godoc
//
// @Summary		Create customer batch
// @Description	Create customer batch.
// @Param		data	body	entity.CreateCustomerBatchRequest	true	"create customer batch"
// @Produce		application/json
// @Tags		customers
// @Success		201	{object}	entity.JsonCreated{data=nil}"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router		/customers/batch [post]
func (handler *CustomerHandler) CreateBatch(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := new(entity.CreateCustomerBatchRequest)
	err := ctx.Bind(request)
	helper.ErrorPanic(err)

	handler.customerUsecase.CreateBatch(c, *request)

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Batch Successful",
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusCreated, webResponse)
}

// Note            godoc
//
// @Summary		update customer
// @Description	update customer.
// @Param		data		body	entity.UpdateCustomerRequest	true	"update customer"
// @Param		customerId	path	string							true	"customer_id"
// @Produce		application/json
// @Tags		customers
// @Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router		/customers/{customerId} [patch]
func (handler *CustomerHandler) Update(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := new(entity.CustomerParams)
	if err := (&echo.DefaultBinder{}).BindPathParams(ctx, params); err != nil {
		ctx.Logger().Error("Handler : Param ID error : ", err.Error())
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	request := new(entity.UpdateCustomerRequest)
	if err := ctx.Bind(request); err != nil {
		ctx.Logger().Error("Handler : Request  error : ", err.Error())
		return exception.NewBadRequestHandler(err.Error())
	}

	request.ID = params.CustomerId

	handler.customerUsecase.Update(c, *request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Update Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusCreated, webResponse)
}

// Note             godoc
//
//	 @Summary		Delete batch customer
//		@Description	Delete batch customer.
//		@Param			data	body	entity.DeleteBatchCustomerRequest	true	"delete batch customer"
//		@Produce		application/json
//		@Tags			customers
//		@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//		@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//		@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//		@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//		@Router			/customers/batch [delete]
func (handler *CustomerHandler) DeleteBatch(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := new(entity.DeleteBatchCustomerRequest)
	err := ctx.Bind(request)
	helper.ErrorPanic(err)

	handler.customerUsecase.DeleteBatch(c, *request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete Batch Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusCreated, webResponse)
}

// Note 		    godoc
//
// @Summary		get customer by id.
// @Param		customerId	path	string	true	"customer_id"
// @Description	get customer by id.
// @Produce		application/json
// @Tags		customers
// @Success		200	{object}	entity.JsonSuccess{data=entity.CustomerResponse{}}	"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
// @Router		/customers/{customerId} [get]
func (handler *CustomerHandler) FindById(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := new(entity.CustomerParams)
	if err := ctx.Bind(params); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	data := handler.customerUsecase.FindById(c, *params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusOK, webResponse)
}

// Note             godoc
//
// @Summary		Get all customers.
// @Description	Get all customers.
// @Produce		application/json
// @Param		limit		query	string	false	"limit"
// @Param		page		query	string	false	"page"
// @Param		username	query	string	false	"username"
// @Param		email		query	string	false	"email"
// @Param		start_date	query	string	false	"start_date"
// @Param		end_date	query	string	false	"end_date"
// @Param		sort		query	string	false	"sort"
// @Tags		customers
// @Success		200	{object}	entity.Response{data=[]entity.CustomerResponse{}}	"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
// @Router		/customers [get]
func (handler *CustomerHandler) FindAllPaging(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter entity.CustomerQueryFilter

	if err := ctx.Bind(&dataFilter); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	response, paging := handler.customerUsecase.FindAllPaging(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusOK, webResponse)
}

//	    Note 		    godoc
//
//		@Summary		Export Excel customer.
//		@Description	Export Excel customer.
//		@Produce		application/json
//		@Tags			customers
//		@Param			start_date	query		string	false	"start_date"
//		@Param			end_date	query		string	false	"end_date"
//		@Param			username	query		string	false	"username"
//		@Param			email		query		string	false	"email"
//		@Success		200			{object}	entity.JsonSuccess{data=string}"Data"
//		@Failure		400			{object}	entity.JsonBadRequest{}				"Validation error"
//		@Failure		404			{object}	entity.JsonNotFound{}				"Data not found"
//		@Failure		500			{object}	entity.JsonInternalServerError{}	"Internal server error"
//		@Router			/customers/export [get]
func (handler *CustomerHandler) Export(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dataFilter entity.CustomerQueryFilter

	if err := ctx.Bind(&dataFilter); err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	filePath, err := handler.customerUsecase.Export(c, dataFilter)
	helper.ErrorPanic(err)
	defer os.Remove(filePath) // Remove the file after the function exits

	fileName := filepath.Base(filePath)
	// Set headers for the Excel file
	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// Read the Excel file and write to the response body
	data, err := os.ReadFile(filePath)
	helper.ErrorPanic(err)

	// Write data to the response body
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

//	    Note 		    godoc
//
//		@Summary		Import Excel customer.
//		@Description	Import Excel customer.
//		@Produce		application/json
//		@Accept			multipart/form-data
//		@Tags			customers
//		@Param			file	formData	file	true	"Import Excel customer"
//		@Success		200		{object}	entity.JsonSuccess{data=string}"Data"
//		@Failure		400		{object}	entity.JsonBadRequest{}				"Validation error"
//		@Failure		404		{object}	entity.JsonNotFound{}				"Data not found"
//		@Failure		500		{object}	entity.JsonInternalServerError{}	"Internal server error"
//		@Router			/customers/import [post]
func (handler *CustomerHandler) Import(ctx echo.Context) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := new(entity.UploadCustomerRequest)

	file, err := ctx.FormFile("file")
	if err != nil {
		panic(exception.NewBadRequestHandler(err.Error()))
	}

	fileExtension := filepath.Ext(file.Filename)
	if fileExtension != ".xlsx" && fileExtension != ".xls" {
		panic(exception.NewBadRequestHandler("Invalid file type. Only .xlsx and .xls are allowed"))
	}

	request.File = file

	error := handler.customerUsecase.Import(c, *request)
	helper.ErrorPanic(error)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Import Successful",
		Data:    nil,
	}
	utils.ResponseInterceptor(ctx, &webResponse)
	return ctx.JSON(http.StatusOK, webResponse)
}
