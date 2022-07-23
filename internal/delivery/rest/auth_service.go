package rest

import (
	"net/http"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Service) handleLoginByEmailAndPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {

		req := &model.LoginByEmailAndPasswordInput{}
		if err := ctx.Bind(req); err != nil {
			return ErrBadRequest
		}

		if err := req.Validate(); err != nil {
			return ErrValidation
		}

		session, err := s.authUsecase.LoginByEmailAndPassword(ctx.Request().Context(), req.Email, req.Password)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, session)
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrUnauthorized:
			return ErrUnauthorized
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":   utils.DumpIncomingContext(ctx.Request().Context()),
				"email": req.Email,
			}).Error(err)

			return ErrInternal
		}
	}
}
