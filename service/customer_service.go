package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"math"
	"scylla/dto"
	"scylla/entity"
	"scylla/pkg/exception"
	"scylla/pkg/helper"
	"scylla/repo"
	"strings"
	"time"
)

type CustomerService interface {
	Create(ctx context.Context, request dto.CreateCustomerRequest)
	CreateBatch(ctx context.Context, request dto.CreateCustomerBatchRequest)
	Update(ctx context.Context, request dto.UpdateCustomerRequest)
	DeleteBatch(ctx context.Context, request dto.DeleteBatchCustomerRequest)
	FindById(ctx context.Context, request dto.CustomerParams) (response dto.CustomerResponse)
	FindAll(ctx context.Context, dataFilter dto.CustomerQueryFilter) (response []dto.CustomerResponse, paging dto.Meta)
	Export(ctx context.Context, dataFilter dto.CustomerQueryFilter) (string, error)
	ExportStream(ctx context.Context, dataFilter dto.CustomerQueryFilter) (string, error)
	Import(ctx context.Context, request dto.UploadCustomerRequest) error
}

type CustomerServiceImpl struct {
	customerRepo repo.CustomerRepo
	validate     *validator.Validate
}

func NewCustomerServiceImpl(customerRepo repo.CustomerRepo, validate *validator.Validate) CustomerService {
	return &CustomerServiceImpl{
		customerRepo: customerRepo,
		validate:     validate,
	}
}

func (service *CustomerServiceImpl) Create(ctx context.Context, request dto.CreateCustomerRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := entity.Customer{
		Username: request.Username,
		Email:    request.Email,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	err = service.customerRepo.Insert(ctx, dataset)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}
}

func (service *CustomerServiceImpl) CreateBatch(ctx context.Context, request dto.CreateCustomerBatchRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	var customers []entity.Customer
	for _, req := range request.Customers {
		customer := entity.Customer{
			Username: req.Username,
			Email:    req.Email,
			Phone:    req.Phone,
			Address:  req.Address,
		}
		customers = append(customers, customer)
	}

	batchSize := len(request.Customers)

	err = service.customerRepo.InsertBatch(ctx, customers, batchSize)
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}
}

func (service *CustomerServiceImpl) Update(ctx context.Context, request dto.UpdateCustomerRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset, err := service.customerRepo.FindById(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	dataset.Username = request.Username
	dataset.Email = request.Email
	dataset.Phone = request.Phone
	dataset.Address = request.Address

	err = service.customerRepo.Update(ctx, dataset)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (service *CustomerServiceImpl) DeleteBatch(ctx context.Context, request dto.DeleteBatchCustomerRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	err = service.customerRepo.DeleteBatch(ctx, request.ID)
	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}
}

func (service *CustomerServiceImpl) FindById(ctx context.Context, request dto.CustomerParams) (response dto.CustomerResponse) {
	result, err := service.customerRepo.FindById(ctx, request.CustomerId)

	if err != nil {
		panic(exception.NewNotFoundHandler(err.Error()))
	}

	helper.Automapper(result, &response)
	return response
}

func (service *CustomerServiceImpl) FindAll(ctx context.Context, dataFilter dto.CustomerQueryFilter) (response []dto.CustomerResponse, paging dto.Meta) {
	result, total := service.customerRepo.FindAll(ctx, dataFilter)

	for _, value := range result {
		var res dto.CustomerResponse
		helper.Automapper(value, &res)

		response = append(response, res)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}

	paging.Page = dataFilter.Page
	paging.Limit = dataFilter.Limit
	paging.TotalData = int(total)
	paging.TotalPage = int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	return response, paging
}

func (service *CustomerServiceImpl) Export(ctx context.Context, dataFilter dto.CustomerQueryFilter) (string, error) {
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

	result, _ := service.customerRepo.FindAll(ctx, dataFilter)

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

func (service *CustomerServiceImpl) ExportStream(ctx context.Context, dataFilter dto.CustomerQueryFilter) (string, error) {
	excel := excelize.NewFile()
	defer func() {
		if err := excel.Close(); err != nil {
			panic(exception.NewInternalServerErrorHandler(err.Error()))
		}
	}()

	// Create a new sheet
	mstCustomer := "MST_CUSTOMER"
	index, err := excel.NewSheet(mstCustomer)
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	// Delete the default sheet
	err = excel.DeleteSheet("Sheet1")
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	// Get the data from the repository
	result, _ := service.customerRepo.FindAll(ctx, dataFilter)
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

	// Initialize StreamWriter
	streamWriter, err := excel.NewStreamWriter(mstCustomer)
	if err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	// Write headers
	headerRow := make([]interface{}, len(headers))
	for i, header := range headers {
		headerRow[i] = header
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := streamWriter.SetRow(cell, headerRow); err != nil {
			return "", exception.NewInternalServerErrorHandler(err.Error())
		}
		excel.SetCellStyle(mstCustomer, cell, cell, headerStyle)
	}

	// Write data rows
	for i, customer := range result {
		row := make([]interface{}, len(headers))
		row[0] = customer.ID
		row[1] = customer.Username
		row[2] = customer.Email
		row[3] = customer.Phone
		row[4] = customer.Address
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, row); err != nil {
			return "", exception.NewInternalServerErrorHandler(err.Error())
		}
	}

	// Flush and close the stream writer
	if err := streamWriter.Flush(); err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	// Set the active sheet
	excel.SetActiveSheet(index)

	// Save the file
	timestamp := time.Now().Format("2006-01-02_150405")
	filePath := fmt.Sprintf("customer_%s.xlsx", timestamp)
	if err := excel.SaveAs(filePath); err != nil {
		return "", exception.NewInternalServerErrorHandler(err.Error())
	}

	return filePath, nil
}

func (service *CustomerServiceImpl) Import(ctx context.Context, request dto.UploadCustomerRequest) error {
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
	chanUser := make(chan entity.Customer)
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

		// Check unique constraint in the engine
		for colIndex, rule := range helper.RulesExcelCustomer {
			fields := strings.Split(rule, ",")
			fieldName := fields[0]
			rules := fields[1:]

			for _, r := range rules {
				if r == "unique" {
					cell := row[colIndex]
					exists := service.customerRepo.CheckColumnExists(ctx, fieldName, cell)
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
			customer := entity.Customer{
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
	var customers []entity.Customer
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

	// Insert batch of customers into the engine
	if err := service.customerRepo.InsertBatch(ctx, customers, len(customers)); err != nil {
		return err
	}

	return nil
}
