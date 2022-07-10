package rest

import (
	"net/http"

	"github.com/Himatro2021/API/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Service) handleCreateMemberInvitation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := model.UserInvitationInput{}

		if err := ctx.Bind(&request); err != nil {
			return ErrBadRequest
		}

		if err := request.Validate(); err != nil {
			return customValidationErrMessage(err.Error())
		}

		invitation, err := s.userUsecase.CreateInvitation(ctx.Request().Context(), request)
		if err != nil {
			logrus.Error(err)
			return ErrInternal
		}

		return ctx.JSON(http.StatusOK, invitation)
	}
}
