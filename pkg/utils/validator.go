package utils

import (
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("projtype", validateProjectType)
}

func validateProjectType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "mobile" || value == "desktop"
}
