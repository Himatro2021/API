package model

import "context"

// AuthUsecase :nodoc:
type AuthUsecase interface {
	LoginByEmailAndPassword(ctx context.Context, email string, password string) error
}
