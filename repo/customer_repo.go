package repo

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"scylla/dto"
	"scylla/entity"
	"scylla/pkg/helper"
	"strings"
)

type CustomerRepo interface {
	Insert(ctx context.Context, data entity.Customer) error
	InsertBatch(ctx context.Context, data []entity.Customer, batchSize int) error
	Update(ctx context.Context, data entity.Customer) error
	DeleteBatch(ctx context.Context, Id []int) error
	FindById(ctx context.Context, Id int) (data entity.Customer, err error)
	FindAll(ctx context.Context, dataFilter dto.CustomerQueryFilter) (domain []dto.CustomerResponse, count int64)
	CheckColumnExists(ctx context.Context, column string, value interface{}) bool
}

type CustomerRepoImpl struct {
	db *gorm.DB
}

func NewCustomerRepoImpl(db *gorm.DB) CustomerRepo {
	return &CustomerRepoImpl{db: db}
}

func (repo *CustomerRepoImpl) Insert(ctx context.Context, data entity.Customer) error {
	result := repo.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CustomerRepoImpl) InsertBatch(ctx context.Context, data []entity.Customer, batchSize int) error {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.CreateInBatches(&data, batchSize).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (repo *CustomerRepoImpl) Update(ctx context.Context, data entity.Customer) error {
	result := repo.db.WithContext(ctx).Updates(&data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CustomerRepoImpl) DeleteBatch(ctx context.Context, Id []int) error {
	var data entity.Customer
	result := repo.db.WithContext(ctx).Where("id IN (?)", Id).Delete(&data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CustomerRepoImpl) FindById(ctx context.Context, Id int) (data entity.Customer, err error) {
	result := repo.db.WithContext(ctx).First(&data, Id)
	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (repo *CustomerRepoImpl) FindAlls(ctx context.Context, dataFilter dto.CustomerQueryFilter) (domain []dto.CustomerResponse) {
	rawQuery := `
		SELECT 
			id, username, email, phone, address, created_at
		FROM 
			customers
	`

	var filters []string
	var args []interface{}

	if dataFilter.Username != "" {
		filters = append(filters, "username LIKE ?")
		args = append(args, "%"+dataFilter.Username+"%")
	}
	if dataFilter.Email != "" {
		filters = append(filters, "email LIKE ?")
		args = append(args, "%"+dataFilter.Email+"%")
	}
	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		filters = append(filters, "created_at BETWEEN ? AND ?")
		args = append(args, dataFilter.StartDate, dataFilter.EndDate)
	}

	if len(filters) > 0 {
		rawQuery += " WHERE " + strings.Join(filters, " AND ")
	}

	sortBy := "id DESC"
	if dataFilter.Sort != "" {
		var sortClauses []string
		for _, row := range strings.Split(dataFilter.Sort, ",") {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", colSort[0], colSort[1]))
			}
		}
		if len(sortClauses) > 0 {
			sortBy = strings.Join(sortClauses, ", ")
		}
	}
	rawQuery += " ORDER BY " + sortBy

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	if dataFilter.Limit > 0 || dataFilter.Page > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		rawQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	result := repo.db.Raw(rawQuery, args...).WithContext(ctx).Scan(&domain)
	helper.ErrorPanic(result.Error)

	return domain
}

func (repo *CustomerRepoImpl) FindAll(ctx context.Context, dataFilter dto.CustomerQueryFilter) (domain []dto.CustomerResponse, count int64) {
	rawQuery := `
        SELECT 
            id, username, email, phone, address, created_at
        FROM 
            customers
    `

	var filters []string
	var args []interface{}

	if dataFilter.Username != "" {
		filters = append(filters, "username LIKE ?")
		args = append(args, "%"+dataFilter.Username+"%")
	}
	if dataFilter.Email != "" {
		filters = append(filters, "email LIKE ?")
		args = append(args, "%"+dataFilter.Email+"%")
	}
	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		filters = append(filters, "created_at BETWEEN ? AND ?")
		args = append(args, dataFilter.StartDate, dataFilter.EndDate)
	}

	if len(filters) > 0 {
		rawQuery += " WHERE " + strings.Join(filters, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM (" + rawQuery + ") AS subquery"
	resultCount := repo.db.Raw(countQuery, args...).WithContext(ctx).Scan(&count)
	helper.ErrorPanic(resultCount.Error)

	sortBy := "id DESC"
	if dataFilter.Sort != "" {
		var sortClauses []string
		for _, row := range strings.Split(dataFilter.Sort, ",") {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", colSort[0], colSort[1]))
			}
		}
		if len(sortClauses) > 0 {
			sortBy = strings.Join(sortClauses, ", ")
		}
	}
	rawQuery += " ORDER BY " + sortBy

	if dataFilter.All == true {
		result := repo.db.Raw(rawQuery, args...).WithContext(ctx).Scan(&domain)
		helper.ErrorPanic(result.Error)
	} else {
		if dataFilter.Page == 0 {
			dataFilter.Page = 1
		}
		if dataFilter.Limit == 0 {
			dataFilter.Limit = 10
		}

		offset := (dataFilter.Page - 1) * dataFilter.Limit
		rawQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)

		result := repo.db.Raw(rawQuery, args...).WithContext(ctx).Scan(&domain)
		helper.ErrorPanic(result.Error)
	}

	return domain, count
}

func (repo *CustomerRepoImpl) CheckColumnExists(ctx context.Context, column string, value interface{}) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM customers WHERE %s = ?)", column)
	err := repo.db.WithContext(ctx).Raw(query, value).Scan(&exists).Error
	if err != nil {
		return false
	}
	return exists
}
