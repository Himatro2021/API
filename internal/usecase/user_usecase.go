package usecase

import (
	"context"
	"strconv"

	"github.com/Himatro2021/API/auth"
	"github.com/Himatro2021/API/internal/external/mailer"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo model.UserRepository
	mailer   *mailer.Mailer
}

// NewUserUsecase create new user usecase instance
func NewUserUsecase(userRepo model.UserRepository, mailer *mailer.Mailer) model.UserUsecase {
	return &userUsecase{
		userRepo,
		mailer,
	}
}

// CreateInvitation self explained
func (u *userUsecase) CreateInvitation(ctx context.Context, input *model.UserInvitationInput) (*model.UserInvitation, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(ctx),
	})

	user := auth.GetUserFromCtx(ctx)
	if !user.HasAccess(rbac.ResourceUser, rbac.ActionInvite) {
		return nil, ErrForbidden
	}

	if err := input.Validate(); err != nil {
		return nil, ErrValidation
	}

	invCode := strconv.FormatInt(utils.GenerateID(), 10)
	role, err := rbac.ParseStringToRole(input.Role)
	if err != nil {
		return nil, ErrValidation
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          input.Email,
		Name:           input.Name,
		InvitationCode: invCode,
		Role:           role,
		Status:         model.InvitationStatusPending,
	}

	// encrypt the invitation before saving to db
	if err := invitation.Encrypt(); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	isExists, err := u.userRepo.IsEmailAlreadyInvited(ctx, invitation.Email)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	// TODO implement reinvite member when feature sending invitation via email
	// is implemented
	if isExists {
		logger.Info("email exists")
		_, _ = handleReinviteMember(ctx, input.Email)
	}

	err = u.userRepo.CreateInvitation(ctx, invitation)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	// decrypt invitation before sending to user
	if err := invitation.Decrypt(); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	go u.mailer.SendInvitationEmail(invitation)

	return invitation, nil
}

// CheckIsInvitationExists check invitation from given invitation code
func (u *userUsecase) CheckIsInvitationExists(ctx context.Context, invitationCode string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"invitationCode": invitationCode,
	})

	code, err := helper.Cryptor().Encrypt(invitationCode)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	err = u.userRepo.CheckIsInvitationExists(ctx, code)
	switch err {
	default:
		logger.Error(err)
		return ErrInternal
	case repository.ErrNotFound:
		return ErrNotFound
	case nil:
		return nil
	}
}

// HandleRegistrationByInvitation not implemented
func (u *userUsecase) HandleRegistrationByInvitation(ctx context.Context, input *model.RegistrationInput) {
}

func handleReinviteMember(ctx context.Context, email string) (*model.UserInvitation, error) {
	// TODO implement resending email when corresponding feature are implemented

	return nil, nil
}
