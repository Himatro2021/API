package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_LoginByEmailAndPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	sessRepo := mock.NewMockSessionRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	uc := authUsecase{
		sessionRepo: sessRepo,
		userRepo:    userRepo,
	}

	email := "lucky@amil.cos"
	pw := "ini password"
	password, _ := helper.HashString(pw)

	user := &model.User{
		ID:       utils.GenerateID(),
		Email:    email,
		Password: password,
	}

	t.Run("ok", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(ctx, email).Times(1).Return(user, nil)
		sessRepo.EXPECT().Create(ctx, gomock.Any()).Times(1).Return(nil)

		_, err := uc.LoginByEmailAndPassword(ctx, email, pw)

		assert.NoError(t, err)
	})

	t.Run("ok - email not found", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(ctx, email).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.LoginByEmailAndPassword(ctx, email, pw)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - err get user by email", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(ctx, email).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.LoginByEmailAndPassword(ctx, email, pw)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - password mismatch", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(ctx, email).Times(1).Return(user, nil)

		_, err := uc.LoginByEmailAndPassword(ctx, email, "passwordnya salah")

		assert.Error(t, err)
		assert.Equal(t, err, ErrUnauthorized)
	})

	t.Run("err - err when creating session", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(ctx, email).Times(1).Return(user, nil)
		sessRepo.EXPECT().Create(ctx, gomock.Any()).Times(1).Return(errors.New("err db"))

		_, err := uc.LoginByEmailAndPassword(ctx, email, pw)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}
