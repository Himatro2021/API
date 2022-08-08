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

// IsEmailAlreadyInvited checking is given email already invited or not based on db record
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

// GetUserByEmail get user information from given email address
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

// GetUserByID self explained
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

// GetUserInvitationByInvitationCode find invitation from given invitation code
func (r *userRepository) GetUserInvitationByInvitationCode(ctx context.Context, code string) (*model.UserInvitation, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"code": code,
	})

	userInvitation := &model.UserInvitation{}
	err := r.db.WithContext(ctx).Model(&model.UserInvitation{}).Where("invitation_code = ?", code).Take(&userInvitation).Error
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return userInvitation, nil
	}
}

// MarkInvitationStatus set invitation status to given model.InvitationStatus. This doesn't check if
// the invitation is exists or not in the first place
func (r *userRepository) MarkInvitationStatus(ctx context.Context, invitation *model.UserInvitation, status model.InvitationStatus) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"invitation": utils.Dump(invitation),
	})

	invitation.Status = status

	err := r.db.WithContext(ctx).Save(invitation).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CheckIsInvitationExists get invitation from invitation code
func (r *userRepository) CheckIsInvitationExists(ctx context.Context, invitationCode string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"invitationCode": invitationCode,
	})

	invitation := &model.UserInvitation{}
	err := r.db.Model(&model.UserInvitation{}).Where("invitation_code = ?", invitationCode).Take(invitation).Error
	switch err {
	default:
		logger.Error(err)
		return err
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	case nil:
		return nil
	}
}

// Register create a new user
func (r *userRepository) Register(ctx context.Context, user *model.User) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"user": utils.Dump(user),
	})

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
