package dto

import "mime/multipart"

type CustomerResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
}

type CreateCustomerBatchRequest struct {
	Customers []CreateCustomerRequest `json:"customers" validate:"required,dive"`
}

type CreateCustomerRequest struct {
	Username string `form:"username" json:"username" validate:"required" example:"admin"`
	Email    string `form:"email" json:"email" validate:"required,email,unique=customers;email" example:"admin@gmail.com"`
	Phone    string `form:"phone" json:"phone" validate:"required" example:"1234567890"`
	Address  string `form:"address" json:"address" validate:"required" example:"123 Main St, Anytown, USA"`
}

type UpdateCustomerRequest struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,unique=customers;email;id"`
	Phone    string `json:"phone" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type DeleteBatchCustomerRequest struct {
	ID []int `json:"id" validate:"required,notEmptyIntSlice"`
}

type UploadCustomerRequest struct {
	File *multipart.FileHeader `form:"file" json:"file" validate:"required"`
}

type CustomerParams struct {
	CustomerId int `param:"customerId" validate:"required"`
}

type CustomerQueryFilter struct {
	All       bool   `query:"all"`
	Limit     int    `query:"limit"`
	Page      int    `query:"page"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Username  string `query:"username"`
	Email     string `query:"email"`
	Sort      string `query:"sort"`
}
