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
func (r *userRepository) CreateInvitation(ctx context.Context, invitation *model.UserInvitation) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"invitation": utils.Dump(invitation),
	})

	err := r.db.WithContext(ctx).Create(invitation).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
	}

	err = r.db.WithContext(ctx).Create(invitation).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// return unencrypted form of invitation code
	// because this will only visible to admin who done the
	// invitation
	invitation.InvitationCode = invCode

	return invitation, nil
}

func (r *userRepository) IsEmailAlreadyInvited(ctx context.Context, email string) (bool, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	invitation := &model.UserInvitation{}

	err := r.db.WithContext(ctx).Model(&model.UserInvitation{}).
		Where("email = ?", email).Take(invitation).Error
	switch err {
	case nil:
		return true, nil
	case gorm.ErrRecordNotFound:
		return false, nil
	default:
		logger.Error(err)
		return false, err
	}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	user := &model.User{}
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Take(user).Error
	switch err {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"userID": id,
	})

	user := &model.User{}

	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Take(user).Error
	switch err {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}
}
