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
		invitationInput := model.UserInvitationInput{
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
		invitationInput := model.UserInvitationInput{
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
		invitationInput := model.UserInvitationInput{
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
}
