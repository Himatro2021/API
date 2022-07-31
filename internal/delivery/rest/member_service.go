package rest

import (
	"net/http"
	"strconv"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Service) handleCreateMemberInvitation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := &model.UserInvitationInput{}

		if err := ctx.Bind(request); err != nil {
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

func (s *Service) handleCheckInvitation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("invitation_code")
		if param == "" {
			return ErrBadRequest
		}

		_, err := strconv.Atoi(param)
		if err != nil {
			return ErrBadRequest
		}

		err = s.userUsecase.CheckIsInvitationExists(ctx.Request().Context(), param)
		switch err {
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":             utils.DumpIncomingContext(ctx.Request().Context()),
				"invitation_code": param,
			})
			return ErrInternal
		case usecase.ErrNotFound:
			return ErrNotFound
		case nil:
			return ctx.NoContent(http.StatusOK)
		}
	}
}
