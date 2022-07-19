package model

import (
	"context"
)

// UserInvitationInput input for create user invitation
type UserInvitationInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// Validate validate struct
func (usi *UserInvitationInput) Validate() error {
	if err := Validator.Struct(usi); err != nil {
		return err
	}

	return nil
}

// UserInvitation :nodoc:
type UserInvitation struct {
	ID             int64  `json:"id" gorm:"primaryKey"`
	Email          string `json:"email" gorm:"unique,not null"`
	Name           string `json:"name" gorm:"not null"`
	InvitationCode string `json:"invitation_code"`
}

// User :nodoc:
type User struct {
	ID    int64  `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

// UserUsecase :nodoc:
type UserUsecase interface {
	CreateInvitation(ctx context.Context, input UserInvitationInput) (*UserInvitation, error)
}

// UserRepository :nodoc:
type UserRepository interface {
	CreateInvitation(ctx context.Context, name, email string) (*UserInvitation, error)
	IsEmailAlreadyInvited(ctx context.Context, email string) (bool, error)
}
