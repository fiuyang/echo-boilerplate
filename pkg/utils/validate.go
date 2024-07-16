package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var db *gorm.DB

func InitializeValidator(db *gorm.DB) *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("notEmptyStringSlice", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]string)
		if len(slices) == 0 {
			return false
		}
		for _, s := range slices {
			if s == "" {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		dateRegex := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
		if !dateRegex.MatchString(fl.Field().String()) {
			return false
		}
		return true
	})

	_ = validate.RegisterValidation("notEmptyIntSlice", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]int)
		if len(slices) == 0 {
			return false
		}
		for _, val := range slices {
			if val == 0 {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("isString", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("isInt", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(int)
		if !ok {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		return ValidateUnique(db, fl)
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return name
		}
		return name
	})

	_ = validate.RegisterValidation("allowedMimeTypeExcel", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(*multipart.FileHeader)
		if !ok {
			return false
		}

		allowedExtensions := map[string]bool{
			".xls":  true,
			".xlsx": true,
		}

		fileExtension := getFileExtension(file.Filename)
		fmt.Println("fileExtension", fileExtension)
		return allowedExtensions[fileExtension]
	})

	_ = validate.RegisterValidation("allowedMimeTypeDoc", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(*multipart.FileHeader)
		if !ok {
			return false
		}

		allowedExtensions := map[string]bool{
			".doc": true,
		}

		fileExtension := getFileExtension(file.Filename)
		return allowedExtensions[fileExtension]
	})

	_ = validate.RegisterValidation("allowedMimeTypeImage", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		allowedExtensions := map[string]bool{
			".jpeg": true,
			".jpg":  true,
			".png":  true,
		}

		fileExtension := getFileExtension(file)
		return allowedExtensions[fileExtension]
	})

	return validate
}

func getFileExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext
}
