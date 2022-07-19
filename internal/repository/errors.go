package repository

import "errors"

var (
	// ErrNotFound used to indicate that repo layer can't found the requested data.
	ErrNotFound = errors.New("data not found")
)
