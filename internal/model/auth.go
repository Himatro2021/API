package model

import "context"

// AuthUsecase :nodoc:
type AuthUsecase interface {
	LoginByEmailAndPassword(ctx context.Context, email, password string) (*Session, error)
}

// LoginByEmailAndPasswordInput define input for request in login by email and password
type LoginByEmailAndPasswordInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate validate struct
func (i *LoginByEmailAndPasswordInput) Validate() error {
	return Validator.Struct(i)
}
