package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TODO implement presenter layer unit testing
func TestREST_handleCreateMemberInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := mock.NewMockUserUsecase(ctrl)
	service := &Service{
		userUsecase: mockUserUsecase,
	}

	validEmail := "lucky@testing.com"
	invalidEmail := "lucky@,ees.com"
	validName := "lucky"

	invitationResult := &model.UserInvitation{
		ID:             utils.GenerateID(),
		Email:          validEmail,
		Name:           validName,
		InvitationCode: "anything",
	}

	t.Run("ok - invited", func(t *testing.T) {
		ec := echo.New()
		invitationInput := &model.UserInvitationInput{
			Email: validEmail,
			Name:  validName,
		}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s"
			}
		`, validEmail, validName)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockUserUsecase.EXPECT().CreateInvitation(ctx, invitationInput).Times(1).Return(invitationResult, nil)

		err := service.handleCreateMemberInvitation()(ectx)

		assert.NoError(t, err)
	})

	t.Run("err - failed binding because payload invalid", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				1"email": "%s",
				"name": "%s"
			}
		`, invalidEmail, validName)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleCreateMemberInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - email invalid sent by user", func(t *testing.T) {
		ec := echo.New()
		invitationInput := &model.UserInvitationInput{
			Email: invalidEmail,
			Name:  validName,
		}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s"
			}
		`, invalidEmail, validName)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockUserUsecase.EXPECT().CreateInvitation(ctx, invitationInput).Times(1).Return(nil, usecase.ErrValidation)

		err := service.handleCreateMemberInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("error - internal err from usecase", func(t *testing.T) {
		ec := echo.New()
		invitationInput := &model.UserInvitationInput{
			Email: validEmail,
			Name:  validName,
		}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s"
			}
		`, validEmail, validName)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockUserUsecase.EXPECT().CreateInvitation(ctx, invitationInput).Times(1).Return(nil, usecase.ErrInternal)

		err := service.handleCreateMemberInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - action is forbidden", func(t *testing.T) {
		ec := echo.New()
		invitationInput := &model.UserInvitationInput{
			Email: validEmail,
			Name:  validName,
		}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s"
			}
		`, validEmail, validName)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockUserUsecase.EXPECT().CreateInvitation(ctx, invitationInput).Times(1).Return(nil, usecase.ErrForbidden)

		err := service.handleCreateMemberInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})
}

func TestMemberService_handleCheckInvitation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockUserUsecase := mock.NewMockUserUsecase(ctrl)
	service := &Service{
		userUsecase: mockUserUsecase,
	}

	t.Run("ok - found", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		ectx := ec.NewContext(req, rec)
		id := fmt.Sprintf("%d", utils.GenerateID())

		ectx.SetPath("/members/invitation/")
		ectx.SetParamNames("invitation_code")
		ectx.SetParamValues(id)

		mockUserUsecase.EXPECT().CheckIsInvitationExists(ctx, id).Times(1).Return(nil)

		err := service.handleCheckInvitation()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - not found", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		ectx := ec.NewContext(req, rec)
		id := fmt.Sprintf("%d", utils.GenerateID())

		ectx.SetPath("/members/invitation/")
		ectx.SetParamNames("invitation_code")
		ectx.SetParamValues(id)

		mockUserUsecase.EXPECT().CheckIsInvitationExists(ctx, id).Times(1).Return(usecase.ErrNotFound)

		err := service.handleCheckInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - internal err", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		ectx := ec.NewContext(req, rec)
		id := fmt.Sprintf("%d", utils.GenerateID())

		ectx.SetPath("/members/invitation/")
		ectx.SetParamNames("invitation_code")
		ectx.SetParamValues(id)

		mockUserUsecase.EXPECT().CheckIsInvitationExists(ctx, id).Times(1).Return(usecase.ErrInternal)

		err := service.handleCheckInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - invitation code is empty", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		ectx := ec.NewContext(req, rec)
		ectx.SetPath("/members/invitation/")
		ectx.SetParamNames("invitation_code")
		ectx.SetParamValues("")

		err := service.handleCheckInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - id is not a number", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		ectx := ec.NewContext(req, rec)
		id := "ini tuh bukan int, harusnya error"

		ectx.SetPath("/members/invitation/")
		ectx.SetParamNames("invitation_code")
		ectx.SetParamValues(id)

		err := service.handleCheckInvitation()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})
}

func TestMemberService_handleRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := mock.NewMockUserUsecase(ctrl)
	service := &Service{
		userUsecase: mockUserUsecase,
	}
	ctx := context.TODO()

	t.Run("ok - registered", func(t *testing.T) {
		registrationInput := &model.RegistrationInput{
			Email:                "lucky@akbar.mail.id",
			Name:                 "lucky",
			Password:             "testpassword",
			PasswordConfirmation: "testpassword",
			InvitationCode:       "testonly",
		}

		user := &model.User{
			ID:    utils.GenerateID(),
			Name:  registrationInput.Name,
			Email: registrationInput.Email,
			Role:  rbac.RoleMember,
		}

		mockUserUsecase.EXPECT().Register(ctx, registrationInput).Times(1).Return(user, nil)
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s",
				"password": "%s",
				"password_confirmation": "%s",
				"invitation_code": "%s"
			}
		`, registrationInput.Email, registrationInput.Name, registrationInput.Password, registrationInput.PasswordConfirmation, registrationInput.InvitationCode)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleRegistration()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - err validation", func(t *testing.T) {
		invalidInput := &model.RegistrationInput{
			Email:                "luckyemailinvalidakbar.mail.id",
			Name:                 "lucky",
			Password:             "testpassword",
			PasswordConfirmation: "testpassword",
			InvitationCode:       "testonly",
		}

		mockUserUsecase.EXPECT().Register(ctx, invalidInput).Times(1).Return(nil, usecase.ErrValidation)
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s",
				"password": "%s",
				"password_confirmation": "%s",
				"invitation_code": "%s"
			}
		`, invalidInput.Email, invalidInput.Name, invalidInput.Password, invalidInput.PasswordConfirmation, invalidInput.InvitationCode)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleRegistration()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - err forbidden", func(t *testing.T) {
		invalidInput := &model.RegistrationInput{
			Email:                "luckyemail@invalidakbar.mail.id",
			Name:                 "lucky",
			Password:             "testpassword",
			PasswordConfirmation: "testpassword",
			InvitationCode:       "testonly",
		}

		mockUserUsecase.EXPECT().Register(ctx, invalidInput).Times(1).Return(nil, usecase.ErrForbidden)
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s",
				"password": "%s",
				"password_confirmation": "%s",
				"invitation_code": "%s"
			}
		`, invalidInput.Email, invalidInput.Name, invalidInput.Password, invalidInput.PasswordConfirmation, invalidInput.InvitationCode)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleRegistration()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("err - internal error", func(t *testing.T) {
		invalidInput := &model.RegistrationInput{
			Email:                "luckyemail@invalidakbar.mail.id",
			Name:                 "lucky",
			Password:             "testpassword",
			PasswordConfirmation: "testpassword",
			InvitationCode:       "testonly",
		}

		mockUserUsecase.EXPECT().Register(ctx, invalidInput).Times(1).Return(nil, usecase.ErrInternal)
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"name": "%s",
				"password": "%s",
				"password_confirmation": "%s",
				"invitation_code": "%s"
			}
		`, invalidInput.Email, invalidInput.Name, invalidInput.Password, invalidInput.PasswordConfirmation, invalidInput.InvitationCode)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleRegistration()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}
