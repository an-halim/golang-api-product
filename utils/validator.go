package utils

import "github.com/go-playground/validator/v10"

type Product struct {
	Name      string `validate:"required"`
	Price     float64    `validate:"required"`
	Stock     int    `validate:"required"`
	Image_url string `validate:"required"`
}

func ValidateProduct(product *Product) error {
	validate := validator.New()
	return validate.Struct(product)
}