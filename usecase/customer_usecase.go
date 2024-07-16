package usecase

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"math"
	"scylla/entity"
	"scylla/model"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/repo"
	"strings"
	"time"
)

type CustomerUsecase interface {
	Create(ctx context.Context, request entity.CreateCustomerRequest)
	CreateBatch(ctx context.Context, request entity.CreateCustomerBatchRequest)
	Update(ctx context.Context, request entity.UpdateCustomerRequest)
	DeleteBatch(ctx context.Context, request entity.DeleteBatchCustomerRequest)
	FindById(ctx context.Context, request entity.CustomerParams) (response entity.CustomerResponse)
	FindAll(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse)
	FindAllPaging(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse, paging entity.Meta)
	Export(ctx context.Context, dataFilter entity.CustomerQueryFilter) (string, error)
	Import(ctx context.Context, request entity.UploadCustomerRequest) error
}

type CustomerUsecaseImpl struct {
	customerRepo repo.CustomerRepo
	validate     *validator.Validate
}

func NewCustomerUsecaseImpl(customerRepo repo.CustomerRepo, validate *validator.Validate) CustomerUsecase {
	return &CustomerUsecaseImpl{
		customerRepo: customerRepo,
		validate:     validate,
	}
}

func (usecase *CustomerUsecaseImpl) Create(ctx context.Context, request entity.CreateCustomerRequest) {
	err := usecase.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.Customer{
		Username: request.Username,
		Email:    request.Email,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	err = usecase.customerRepo.Insert(ctx, dataset)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}
}

func (usecase *CustomerUsecaseImpl) CreateBatch(ctx context.Context, request entity.CreateCustomerBatchRequest) {
	err := usecase.validate.Struct(request)
	helper.ErrorPanic(err)

	var customers []model.Customer
	for _, req := range request.Customers {
		customer := model.Customer{
			Username: req.Username,
			Email:    req.Email,
			Phone:    req.Phone,
			Address:  req.Address,
		}
		customers = append(customers, customer)
	}

	batchSize := len(request.Customers)

	err = usecase.customerRepo.InsertBatch(ctx, customers, batchSize)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}
}

func (usecase *CustomerUsecaseImpl) Update(ctx context.Context, request entity.UpdateCustomerRequest) {
	err := usecase.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset, err := usecase.customerRepo.FindById(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	dataset.Username = request.Username
	dataset.Email = request.Email
	dataset.Phone = request.Phone
	dataset.Address = request.Address

	err = usecase.customerRepo.Update(ctx, dataset)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (usecase *CustomerUsecaseImpl) DeleteBatch(ctx context.Context, request entity.DeleteBatchCustomerRequest) {
	err := usecase.validate.Struct(request)
	helper.ErrorPanic(err)

	err = usecase.customerRepo.DeleteBatch(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (usecase *CustomerUsecaseImpl) FindById(ctx context.Context, request entity.CustomerParams) (response entity.CustomerResponse) {
	result, err := usecase.customerRepo.FindById(ctx, request.CustomerId)

	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	helper.Automapper(result, &response)
	return response
}

func (usecase *CustomerUsecaseImpl) FindAll(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse) {
	result, err := usecase.customerRepo.FindAll(ctx, dataFilter)

	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	for _, row := range result {
		var res entity.CustomerResponse
		helper.Automapper(row, &res)
		response = append(response, res)
	}
	return response
}

func (usecase *CustomerUsecaseImpl) FindAllPaging(ctx context.Context, dataFilter entity.CustomerQueryFilter) (response []entity.CustomerResponse, paging entity.Meta) {

	result := usecase.customerRepo.FindAllPaging(ctx, dataFilter)

	for _, value := range result {
		var res entity.CustomerResponse
		helper.Automapper(value, &res)

		response = append(response, res)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}

	var total int
	if len(result) > 0 {
		total = len(result)
	} else {
		total = 0
	}

	paging.Page = dataFilter.Page
	paging.Limit = dataFilter.Limit
	paging.TotalData = total
	paging.TotalPage = int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	return response, paging
}

func (usecase *CustomerUsecaseImpl) Export(ctx context.Context, dataFilter entity.CustomerQueryFilter) (string, error) {
	excel := excelize.NewFile()
	defer func() {
		if err := excel.Close(); err != nil {
			panic(exception.NewInternalServerErrorHandler(err.Error()))
		}
	}()

	//sheet
	mstCustomer := "MST_CUSTOMER"
	index, err := excel.NewSheet(mstCustomer)
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	err = excel.DeleteSheet("Sheet1")
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	result, err := usecase.customerRepo.FindAll(ctx, dataFilter)
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	// Define headers
	headers := []string{"ID", "Name", "Email", "Phone", "Address"}

	// Set headers and apply styles
	headerStyle, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})

	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		excel.SetCellValue(mstCustomer, cell, header)
		excel.SetCellStyle(mstCustomer, cell, cell, headerStyle)
	}

	// customer the sheet with data
	for i, customer := range result {
		excel.SetCellValue(mstCustomer, fmt.Sprintf("A%d", i+2), customer.ID)
		excel.SetCellValue(mstCustomer, fmt.Sprintf("B%d", i+2), customer.Username)
		excel.SetCellValue(mstCustomer, fmt.Sprintf("C%d", i+2), customer.Email)
		excel.SetCellValue(mstCustomer, fmt.Sprintf("D%d", i+2), customer.Phone)
		excel.SetCellValue(mstCustomer, fmt.Sprintf("E%d", i+2), customer.Address)
	}

	excel.SetActiveSheet(index)

	timestamp := time.Now().Format("2006-01-02_150405")
	filePath := fmt.Sprintf("customer_%s.xlsx", timestamp)

	if err := excel.SaveAs(filePath); err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	return filePath, nil
}

func (usecase *CustomerUsecaseImpl) Import(ctx context.Context, request entity.UploadCustomerRequest) error {
	// Open the Excel file from the request
	src, err := request.File.Open()
	if err != nil {
		return exception.NewInternalServerErrorHandler(err.Error())
	}
	defer src.Close()

	// Initialize Excel reader
	xlFile, err := excelize.OpenReader(src)
	if err != nil {
		return exception.NewInternalServerErrorHandler(err.Error())
	}

	// Define the sheet name to read from
	sheetName := "MST_CUSTOMER"

	// Read all rows from the sheet
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		return exception.NewInternalServerErrorHandler(err.Error())
	}

	// Use error group to manage concurrent operations
	g, _ := errgroup.WithContext(ctx)
	chanUser := make(chan model.Customer)
	excelValidation := exception.ExcelValidation{}
	uniqueTracker := make(map[string]map[string]bool)
	rowErrors := map[string][]string{}

	// Initialize uniqueTracker for each field based on validation rules
	for _, rule := range helper.RulesExcelCustomer {
		fieldName := strings.Split(rule, ",")[0]
		uniqueTracker[fieldName] = make(map[string]bool)
	}

	// Validate each row and cell dynamically
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue // Skip header row
		}

		// Validate each cell in the row based on rules
		for colIndex, rule := range helper.RulesExcelCustomer {
			fields := strings.Split(rule, ",")
			fieldName := fields[0]
			cell := row[colIndex]
			rules := fields[1:]

			for _, r := range rules {
				switch r {
				case "required":
					if cell == "" {
						rowErrors[fieldName] = append(rowErrors[fieldName], fmt.Sprintf("%s row %d is required", fieldName, rowIndex+1))
					}
				case "unique":
					if uniqueTracker[fieldName][cell] {
						rowErrors[fieldName] = append(rowErrors[fieldName], fmt.Sprintf("%s '%s' is not unique row %d", fieldName, cell, rowIndex+1))
					}
					uniqueTracker[fieldName][cell] = true
				}
			}
		}

		// Check unique constraint in the database
		for colIndex, rule := range helper.RulesExcelCustomer {
			fields := strings.Split(rule, ",")
			fieldName := fields[0]
			rules := fields[1:]

			for _, r := range rules {
				if r == "unique" {
					cell := row[colIndex]
					exists := usecase.customerRepo.CheckColumnExists(ctx, fieldName, cell)
					if exists != false {
						excelValidation.AddHandler(fieldName, rowIndex+1, fmt.Sprintf("%s '%s' already taken", fieldName, cell))
					}
				}
			}
		}

		// If there are validation errors for this row, skip further processing
		if len(rowErrors) > 0 {
			for field, errs := range rowErrors {
				for _, err := range errs {
					excelValidation.AddHandler(field, rowIndex+1, err)
				}
			}
			continue
		}

		// Process valid rows concurrently
		g.Go(func() error {
			if len(row) < 4 {
				return fmt.Errorf("invalid row length: %v", row)
			}
			customer := model.Customer{
				Username: row[0],
				Email:    row[1],
				Phone:    row[2],
				Address:  row[3],
			}
			chanUser <- customer
			return nil
		})
	}

	// Close the channel once all goroutines are done
	go func() {
		err := g.Wait()
		if err != nil {
			return
		}
		close(chanUser)
	}()

	// Collect customers from the channel
	var customers []model.Customer
	for customer := range chanUser {
		customers = append(customers, customer)
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return err
	}

	// If there are any validation errors, return them
	if len(excelValidation.Errors) > 0 {
		return &excelValidation
	}

	// Insert batch of customers into the database
	if err := usecase.customerRepo.InsertBatch(ctx, customers, len(customers)); err != nil {
		return err
	}

	return nil
}
