package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {

		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("The %s field is required.", err.Field())
		case "min":
			message = fmt.Sprintf("The %s field must be at least %s characters long.", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("The %s field must not exceed %s characters.", err.Field(), err.Param())
		case "hostname_rfc1123":
			message = fmt.Sprintf("The %s field must be a valid hostname.", err.Field())
		default:
			message = fmt.Sprintf("The %s field is invalid.", err.Field())
		}
		errors[err.Field()] = message
	}
	return errors
}
