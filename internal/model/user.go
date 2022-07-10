package model

import (
	"context"
)

// UserInvitationInput input for create user invitation
type UserInvitationInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

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

type User struct {
}

// UserUsecase :nodoc:
type UserUsecase interface {
	CreateInvitation(ctx context.Context, input UserInvitationInput) (*UserInvitation, error)
}

// UserRepository :nodoc:
type UserRepository interface {
	CreateInvitation(ctx context.Context, name string, email string) (*UserInvitation, error)
}
