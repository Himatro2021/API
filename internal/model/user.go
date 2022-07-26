package model

import (
	"context"

	"github.com/Himatro2021/API/internal/rbac"
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
	ID       int64  `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     rbac.Role
}

// SetRole set role to a user. Is this needed?
func (u *User) SetRole(role rbac.Role) {
	u.Role = role
}

// GetRole return user role
func (u *User) GetRole() rbac.Role {
	return u.Role
}

// UserUsecase :nodoc:
type UserUsecase interface {
	CreateInvitation(ctx context.Context, input UserInvitationInput) (*UserInvitation, error)
}

// UserRepository :nodoc:
type UserRepository interface {
	CreateInvitation(ctx context.Context, name, email string) (*UserInvitation, error)
	IsEmailAlreadyInvited(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
}
