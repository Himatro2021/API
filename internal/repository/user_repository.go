package repository

import (
	"context"

	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateInvitation(ctx context.Context, name, email string) (*model.UserInvitation, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"name":  name,
		"email": email,
	})

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          email,
		InvitationCode: "123",
		Name:           name,
	}

	err := r.db.WithContext(ctx).Create(invitation).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return invitation, nil
}
