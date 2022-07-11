package usecase

import (
	"context"

	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo model.UserRepository
}

// NewUserUsecase create new user usecase instance
func NewUserUsecase(userRepo model.UserRepository) model.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// CreateInvitation self explained
func (u *userUsecase) CreateInvitation(ctx context.Context, input model.UserInvitationInput) (*model.UserInvitation, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(ctx),
	})

	if err := input.Validate(); err != nil {
		return nil, ErrValidation
	}

	isExists, err := u.userRepo.IsEmailAlreadyInvited(ctx, input.Email)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	// TODO implement reinvite member when feature sending invitation via email
	// is implemented
	if isExists {
		_, _ = handleReinviteMember(ctx, input.Email)
	}

	invitation, err := u.userRepo.CreateInvitation(ctx, input.Name, input.Email)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return invitation, nil
}

func handleReinviteMember(ctx context.Context, email string) (*model.UserInvitation, error) {
	// TODO implement resending email when corresponding feature are implemented

	return nil, nil
}
