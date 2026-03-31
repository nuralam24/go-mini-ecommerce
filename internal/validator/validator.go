package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func Validate(s interface{}) error {
	if validate == nil {
		Init()
	}
	return validate.Struct(s)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var message string
			field := strings.ToLower(string(e.Field()[0])) + e.Field()[1:]

			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s must be a valid email", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			case "gt":
				message = fmt.Sprintf("%s must be greater than %s", field, e.Param())
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
			case "lt":
				message = fmt.Sprintf("%s must be less than %s", field, e.Param())
			case "lte":
				message = fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
			default:
				message = fmt.Sprintf("%s is invalid", field)
			}

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return errors
}
