package repository

import (
	"context"

	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

// NewSessionRepository return new instance of sessionRepo
func NewSessionRepository(db *gorm.DB) model.SessionRepository {
	return &sessionRepo{
		db: db,
	}
}

// Create create a session
func (r *sessionRepo) Create(ctx context.Context, session *model.Session) error {
	err := r.db.WithContext(ctx).Create(session).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.DumpIncomingContext(ctx),
			"sess": utils.Dump(session),
		}).Error(err)

		return err
	}

	return nil
}

// FindByAccessToken find a session based of it's access token
func (r *sessionRepo) FindByAccessToken(ctx context.Context, token string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"token": token,
	})

	session := &model.Session{}

	err := r.db.WithContext(ctx).Model(&model.Session{}).
		Where("access_token = ?", token).
		Take(session).Error
	switch err {
	case nil:
		return session, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}
}
