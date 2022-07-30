package model

import (
	"context"
	"fmt"

	"github.com/Himatro2021/API/internal/config"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/sirupsen/logrus"
)

// UserInvitationInput input for create user invitation
type UserInvitationInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}

// Validate validate struct
func (usi *UserInvitationInput) Validate() error {
	if err := Validator.Struct(usi); err != nil {
		return err
	}

	return nil
}

// InvitationStatus represent current invitation status
type InvitationStatus string

// Enum for InvitationStatus type
var (
	InvitationStatusSent    InvitationStatus = "SENT"
	InvitationStatusPending InvitationStatus = "PENDING"
	InvitationStatusFailed  InvitationStatus = "FAILED"
)

// UserInvitation :nodoc:
type UserInvitation struct {
	ID             int64            `json:"id" gorm:"primaryKey"`
	MailServiceID  int64            `json:"-"`
	Email          string           `json:"email" gorm:"unique,not null"`
	Name           string           `json:"name" gorm:"not null"`
	InvitationCode string           `json:"invitation_code"`
	Role           rbac.Role        `json:"role"`
	Status         InvitationStatus `json:"invitation_status"`
}

// GenerateInvitationLink generate invitation link based on user invitation base url
// and it's invitation code
func (u *UserInvitation) GenerateInvitationLink() string {
	return fmt.Sprintf("%s/%s/", config.UserInvitationBaseURL(), u.InvitationCode)
}

// Encrypt encrypt invitation content that must be encrypted
func (u *UserInvitation) Encrypt() error {
	cryptor := helper.Cryptor()

	mail, err := cryptor.Encrypt(u.Email)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.Email = mail

	name, err := cryptor.Encrypt(u.Name)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.Name = name

	invitationCode, err := cryptor.Encrypt(u.InvitationCode)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.InvitationCode = invitationCode

	return nil
}

// Decrypt decrypt invitation content that previously was encrypted
func (u *UserInvitation) Decrypt() error {
	cryptor := helper.Cryptor()

	mail, err := cryptor.Decrypt(u.Email)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.Email = mail

	name, err := cryptor.Decrypt(u.Name)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.Name = name

	invitationCode, err := cryptor.Decrypt(u.InvitationCode)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.InvitationCode = invitationCode

	return nil
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

// RegistrationInput represent registration input
type RegistrationInput struct {
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,eqfield=Password"`
}

// Validate validate registration input
func (r *RegistrationInput) Validate() error {
	return Validator.Struct(r)
}

// UserUsecase :nodoc:
type UserUsecase interface {
	CreateInvitation(ctx context.Context, input *UserInvitationInput) (*UserInvitation, error)
	CheckIsInvitationExists(ctx context.Context, invitationCode string) error
}

// UserRepository :nodoc:
type UserRepository interface {
	CreateInvitation(ctx context.Context, invitation *UserInvitation) error
	CheckIsInvitationExists(ctx context.Context, invitationCode string) error
	IsEmailAlreadyInvited(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	MarkInvitationStatus(ctx context.Context, invitation *UserInvitation, status InvitationStatus) error
}
