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

	invitation, err := u.userRepo.CreateInvitation(ctx, input.Name, input.Email)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return invitation, nil
}
