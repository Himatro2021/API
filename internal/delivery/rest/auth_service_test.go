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

func TestAuthService_handleLoginByEmailAndPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAuthUsecase(ctrl)
	service := &Service{
		authUsecase: mock,
	}

	email := "luc@email.com"
	password := "inipassword"

	ctx := context.TODO()

	sess := &model.Session{
		ID: utils.GenerateID(),
	}

	t.Run("ok", func(t *testing.T) {
		mock.EXPECT().LoginByEmailAndPassword(ctx, email, password).Times(1).Return(sess, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"password": "%s"
			}
		`, email, password)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - email not found", func(t *testing.T) {
		mock.EXPECT().LoginByEmailAndPassword(ctx, email, password).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"password": "%s"
			}
		`, email, password)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("ok - password mismatch", func(t *testing.T) {
		mock.EXPECT().LoginByEmailAndPassword(ctx, email, password).Times(1).Return(sess, usecase.ErrUnauthorized)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"password": "%s"
			}
		`, email, password)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrUnauthorized)
	})

	t.Run("err - internal error", func(t *testing.T) {
		mock.EXPECT().LoginByEmailAndPassword(ctx, email, password).Times(1).Return(sess, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"password": "%s"
			}
		`, email, password)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - email invalid", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "emailnya invalid",
				"password": "%s"
			}
		`, password)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - password too short", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"email": "%s",
				"password": "pwd"
			}
		`, email)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		err := service.handleLoginByEmailAndPassword()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})
}
