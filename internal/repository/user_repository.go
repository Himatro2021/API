package repository

import (
	"context"
	"strconv"

	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository initialize user repository
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateInvitation create invitation for member
func (r *userRepository) CreateInvitation(ctx context.Context, name, email string) (*model.UserInvitation, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"name":  name,
		"email": email,
	})

	invCode, err := helper.HashString(strconv.FormatInt(utils.GenerateID(), 10))
	if err != nil {
		return nil, err
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          email,
		Name:           name,
		InvitationCode: invCode,
	}

	err = r.db.WithContext(ctx).Create(invitation).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return invitation, nil
}
