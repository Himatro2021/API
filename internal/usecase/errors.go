package usecase

import "errors"

var (
	// ErrValidation returned when failed on data validation
	ErrValidation = errors.New("validation error")

	// ErrInternal self explained
	ErrInternal = errors.New("internal server error")

	// ErrNotFound used when usecase received not found error from repo layer
	ErrNotFound = errors.New("error record not found")

	// ErrForbidden used when tried to do something that ruled to forbidden
	ErrForbidden = errors.New("action is forbidden")

	// ErrAlreadyExists used when user tried to create same / duplicate entry
	ErrAlreadyExists = errors.New("record already exists")

	// ErrUnauthorized used when authorization process return error
	ErrUnauthorized = errors.New("unauthorized")
)
