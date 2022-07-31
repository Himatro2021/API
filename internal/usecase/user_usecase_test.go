package usecase

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/Himatro2021/API/auth"
	"github.com/Himatro2021/API/internal/external/mailer"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_CreateInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	admin := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleAdmin,
	}
	ctx := auth.SetUserToCtx(context.TODO(), admin)
	_ = os.Setenv("PRIVATE_KEY", "supersecret")
	_ = os.Setenv("IV_KEY", "4e6064d3814c2cd22c550155655fefc6")
	_ = os.Setenv("MAIL_SERVICE_URL", "https://staging.mail.service.luckyakbar.tech/rest/mail/free/")
	repo := mock.NewMockUserRepository(ctrl)
	uc := userUsecase{
		userRepo: repo,
		mailer:   mailer.NewMailer(repo),
	}

	invCode := strconv.FormatInt(utils.GenerateID(), 10)
	encrypted, _ := helper.HashString(invCode)

	input := &model.UserInvitationInput{
		Name:  "lucky",
		Email: "lucky@test.ting",
		Role:  "ADMIN",
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          input.Email,
		Name:           input.Name,
		InvitationCode: encrypted,
		MailServiceID:  utils.GenerateID(),
		Role:           rbac.RoleAdmin,
		Status:         model.InvitationStatusPending,
	}
	email, _ := helper.Cryptor().Encrypt(invitation.Email)
	helper.PanicIfErr(invitation.Encrypt())

	t.Run("ok - created", func(t *testing.T) {
		repo.EXPECT().IsEmailAlreadyInvited(ctx, email).Times(1).Return(false, nil)
		repo.EXPECT().CreateInvitation(ctx, gomock.Any()).Times(1).Return(nil)

		_, err := uc.CreateInvitation(ctx, input)

		assert.NoError(t, err)
	})

	t.Run("ok - doing reinvite", func(t *testing.T) {
		// TODO implement unit test when reinvite member feature is implemented
		// dont forget to also add unit test on the related edge cases
		// e.g err from db when checking is exists
	})

	t.Run("err from db when inviting", func(t *testing.T) {
		repo.EXPECT().IsEmailAlreadyInvited(ctx, invitation.Email).Times(1).Return(false, nil)
		repo.EXPECT().CreateInvitation(ctx, gomock.Any()).Times(1).Return(errors.New("err from db"))

		_, err := uc.CreateInvitation(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, ErrInternal, err)
	})

	t.Run("ok - non admin can't invite", func(t *testing.T) {
		member := auth.User{
			Role: rbac.RoleMember,
		}
		memberCtx := auth.SetUserToCtx(context.TODO(), member)

		_, err := uc.CreateInvitation(memberCtx, input)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("ok - input invalid", func(t *testing.T) {
		invalidInputEmail := &model.UserInvitationInput{
			Name:  "ok di nama",
			Email: "emailTidak@okeyagesya",
			Role:  "ADMIN",
		}

		_, err := uc.CreateInvitation(ctx, invalidInputEmail)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - invalid input on role", func(t *testing.T) {
		invalidInputRole := &model.UserInvitationInput{
			Name:  "ok di nama",
			Email: "ok@di.email",
			Role:  "nah ini role apa?",
		}

		_, err := uc.CreateInvitation(ctx, invalidInputRole)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})
}

func TestUserUsecase_CheckIsInvitationExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_ = os.Setenv("PRIVATE_KEY", "supersecret")
	_ = os.Setenv("IV_KEY", "4e6064d3814c2cd22c550155655fefc6")

	repo := mock.NewMockUserRepository(ctrl)
	uc := userUsecase{
		userRepo: repo,
	}

	admin := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleAdmin,
	}
	invCode := "iniinvcode"
	code, _ := helper.Cryptor().Encrypt(invCode)
	ctx := auth.SetUserToCtx(context.TODO(), admin)

	t.Run("ok - exists", func(t *testing.T) {
		repo.EXPECT().CheckIsInvitationExists(ctx, code).Times(1).Return(nil)

		err := uc.CheckIsInvitationExists(ctx, invCode)

		assert.NoError(t, err)
	})

	t.Run("ok - invitation not found", func(t *testing.T) {
		repo.EXPECT().CheckIsInvitationExists(ctx, code).Times(1).Return(repository.ErrNotFound)

		err := uc.CheckIsInvitationExists(ctx, invCode)

		assert.Error(t, err)
		assert.Error(t, err, ErrNotFound)
	})

	t.Run("err - err db", func(t *testing.T) {
		repo.EXPECT().CheckIsInvitationExists(ctx, code).Times(1).Return(errors.New("err db"))

		err := uc.CheckIsInvitationExists(ctx, invCode)

		assert.Error(t, err)
		assert.Error(t, err, ErrInternal)
	})
}
