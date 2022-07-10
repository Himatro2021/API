package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_CreateInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockUserRepository(ctrl)
	uc := userUsecase{
		userRepo: repo,
	}

	input := model.UserInvitationInput{
		Name:  "lucky",
		Email: "lucky@test.ting",
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          input.Email,
		Name:           input.Name,
		InvitationCode: "utils.GenerateID()",
	}

	t.Run("ok - created", func(t *testing.T) {
		repo.EXPECT().CreateInvitation(ctx, input.Name, input.Email).Times(1).Return(invitation, nil)

		result, err := uc.CreateInvitation(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, result.Email, input.Email)
	})

	t.Run("err from db", func(t *testing.T) {
		repo.EXPECT().CreateInvitation(ctx, input.Name, input.Email).Times(1).Return(nil, errors.New("err from db"))

		_, err := uc.CreateInvitation(ctx, input)

		assert.Error(t, err)
	})
}
