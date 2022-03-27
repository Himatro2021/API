package util

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

var Validator *validator.Validate = validator.New()

func ExtractValidationErrorMsg(err error) []string {
	var errorStack []string

	for _, err := range err.(validator.ValidationErrors) {
		LogErr("WARN", fmt.Sprintf("%s is invalid, status: %s, got: '%s'.", err.StructField(), err.Tag(), err.Value()), "")
		errorStack = append(errorStack, fmt.Sprintf("%s is invalid, status: %s, got: '%s'.", err.StructField(), err.Tag(), err.Value()))
	}

	return errorStack
}
