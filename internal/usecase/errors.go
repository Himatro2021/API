package usecase

import "errors"

var (
	// ErrValidation returned when failed on data validation
	ErrValidation = errors.New("validation error")

	// ErrInternal self explained
	ErrInternal = errors.New("internal server error")
)
