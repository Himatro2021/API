package rest

import (
	"net/http"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Service) handleCreateMemberInvitation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := model.UserInvitationInput{}

		if err := ctx.Bind(&request); err != nil {
			return ErrBadRequest
		}

		invitation, err := s.userUsecase.CreateInvitation(ctx.Request().Context(), request)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, invitation)
		case usecase.ErrValidation:
			return ErrValidation
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":     utils.DumpIncomingContext(ctx.Request().Context()),
				"payload": utils.Dump(request),
			})
			return ErrInternal
		}
	}
}
