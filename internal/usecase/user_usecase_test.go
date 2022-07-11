package usecase

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/Himatro2021/API/internal/helper"
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
	_ = os.Setenv("PRIVATE_KEY", "7F8r47rRaaVJUUpcbpn5xLSKj7zL1lr4jemeI4ARqPTRBlcEPNl3mDp7m1lY")
	_ = os.Setenv("IV_KEY", "4e6064d3814c2cd22c5501aswepoicn8")
	repo := mock.NewMockUserRepository(ctrl)
	uc := userUsecase{
		userRepo: repo,
	}

	invCode := strconv.FormatInt(utils.GenerateID(), 10)
	encrypted, _ := helper.HashString(invCode)

	input := model.UserInvitationInput{
		Name:  "lucky",
		Email: "lucky@test.ting",
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          input.Email,
		Name:           input.Name,
		InvitationCode: encrypted,
	}

	t.Run("ok - created", func(t *testing.T) {
		repo.EXPECT().IsEmailAlreadyInvited(ctx, input.Email).Times(1).Return(true, nil)
		repo.EXPECT().CreateInvitation(ctx, input.Name, input.Email).Times(1).Return(invitation, nil)

		_, err := uc.CreateInvitation(ctx, input)

		assert.NoError(t, err)
		// assert.Equal(t, result.Email, input.Email)
	})

	t.Run("ok - doing reinvite", func(t *testing.T) {
		// TODO implement unit test when reinvite member feature is implemented
		// dont forget to also add unit test on the related edge cases
		// e.g err from db when checking is exists
	})

	t.Run("err from db when inviting", func(t *testing.T) {
		repo.EXPECT().IsEmailAlreadyInvited(ctx, input.Email).Times(1).Return(true, nil)
		repo.EXPECT().CreateInvitation(ctx, input.Name, input.Email).Times(1).Return(nil, errors.New("err from db"))

		_, err := uc.CreateInvitation(ctx, input)

		assert.Error(t, err)
	})
}
