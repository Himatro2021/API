package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/Himatro2021/API/internal/config"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/kumparan/go-utils"
	"github.com/mattheath/base62"
	"github.com/sirupsen/logrus"
)

type authUsecase struct {
	sessionRepo model.SessionRepository
	userRepo    model.UserRepository
}

// NewAuthUsecase returns an instance of authUsecase
func NewAuthUsecase(sessionRepo model.SessionRepository, userRepo model.UserRepository) model.AuthUsecase {
	return &authUsecase{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// LoginByEmailAndPassword self explained
func (au *authUsecase) LoginByEmailAndPassword(ctx context.Context, email, plainPassword string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	user, err := au.userRepo.GetUserByEmail(ctx, email)
	switch err {
	case nil:
		break
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

	if !helper.IsHashedStringMatch([]byte(plainPassword), []byte(user.Password)) {
		return nil, ErrUnauthorized
	}

	session := &model.Session{
		ID:                    utils.GenerateID(),
		UserID:                user.ID,
		AccessToken:           generateToken(user.ID),
		RefreshToken:          generateToken(utils.GenerateID()),
		AccessTokenExpiredAt:  time.Now().Add(config.DefaultAccessTokenExpiry),
		RefreshTokenExpiredAt: time.Now().Add(config.DefaultRefreshTokenExpiry),
	}

	if err := au.sessionRepo.Create(ctx, session); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return session, nil
}

// TODO move this func to helper package
func generateToken(uniqueID int64) string {
	sb := strings.Builder{}

	encodedID := base62.EncodeInt64(uniqueID)
	sb.WriteString(encodedID)
	sb.WriteString("-__-")

	randString := utils.GenerateRandomAlphanumeric(config.DefaultTokenLength)
	sb.WriteString(randString)

	return sb.String()
}
