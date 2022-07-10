package model

import (
	validate "github.com/go-playground/validator/v10"
)

// Validator singleton for validating struct
var Validator = validate.New()
