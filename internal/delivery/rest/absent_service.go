package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Service) handleGetFormByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("id")

		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrBadRequest
		}

		form, err := s.absentUsecase.GetAbsentFormByID(ctx.Request().Context(), id)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, form)
		case usecase.ErrNotFound:
			return ErrNotFound
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":    utils.DumpIncomingContext(ctx.Request().Context()),
				"formID": id,
			}).Error(err)

			return ErrInternal
		}

	}
}

func (s *Service) handleCreateAbsentForm() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := &model.CreateAbsentInput{}

		if err := ctx.Bind(request); err != nil {
			return ErrBadRequest
		}

		form, err := s.absentUsecase.CreateAbsentForm(ctx.Request().Context(), request)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, form)
		case usecase.ErrValidation:
			return ErrBadRequest
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":     utils.DumpIncomingContext(ctx.Request().Context()),
				"request": utils.Dump(request),
			}).Error(err)

			return ErrInternal
		}
	}
}

func (s *Service) handleGetAllAbsentForms() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		limit := ctx.QueryParam("limit")
		offset := ctx.QueryParam("offset")

		if limit == "" {
			limit = "0"
		}

		if offset == "" {
			offset = "0"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			limitInt = 0
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			offsetInt = 0
		}

		absentForms, err := s.absentUsecase.GetAllAbsentForm(ctx.Request().Context(), limitInt, offsetInt)
		if errors.Is(err, usecase.ErrInternal) {
			logrus.WithFields(logrus.Fields{
				"ctx":    utils.DumpIncomingContext(ctx.Request().Context()),
				"limit":  limit,
				"offset": offset,
			}).Error(err)

			return ErrInternal
		}

		return ctx.JSON(http.StatusOK, absentForms)
	}
}

func (s *Service) handleUpdateAbsentForm() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrValidation
		}

		req := &model.CreateAbsentInput{}
		if err := ctx.Bind(req); err != nil {
			return ErrBadRequest
		}

		updatedForm, err := s.absentUsecase.UpdateAbsentForm(ctx.Request().Context(), id, req)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, updatedForm)
		case usecase.ErrNotFound:
			return ErrNotFound
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":   utils.DumpIncomingContext(ctx.Request().Context()),
				"id":    id,
				"input": utils.Dump(req),
			}).Error(err)
			return ErrInternal
		}
	}
}

func (s *Service) handleFillAbsentFormByAttendee() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrValidation
		}

		req := &model.FillAbsentFormInput{}
		if err := ctx.Bind(req); err != nil {
			return ErrBadRequest
		}

		if err := req.Validate(); err != nil {
			return ErrValidation
		}

		result, err := s.absentUsecase.FillAbsentFormByAttendee(ctx.Request().Context(), id, req.Status, req.ExecuseReason)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, result)
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrValidation:
			return ErrValidation
		case usecase.ErrForbidden:
			return ErrForbidden
		case usecase.ErrAlreadyExists:
			return ErrAlreadyExists
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":   utils.DumpIncomingContext(ctx.Request().Context()),
				"id":    id,
				"input": utils.Dump(req),
			}).Error(err)
			return ErrInternal
		}
	}
}

func (s *Service) handleUpdateAbsentListByAttendee() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("id")

		absentListID, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrBadRequest
		}

		req := &model.UpdateAbsentListInput{}
		if err := ctx.Bind(req); err != nil {
			logrus.Info(err)
			return ErrBadRequest
		}

		absentList, err := s.absentUsecase.UpdateAbsentListByAttendee(ctx.Request().Context(), absentListID, req)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, absentList)
		case usecase.ErrValidation:
			return ErrValidation
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrForbidden:
			return ErrForbidden
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":          utils.DumpIncomingContext(ctx.Request().Context()),
				"absentListID": absentListID,
				"input":        utils.Dump(req),
			}).Error(err)

			return ErrInternal
		}
	}
}

func (s *Service) handleGetParticipantsByFormID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		param := ctx.Param("id")

		formID, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return ErrBadRequest
		}

		result, err := s.absentUsecase.GetAbsentResultByFormID(ctx.Request().Context(), formID)
		switch err {
		case nil:
			return ctx.JSON(http.StatusOK, result)
		case usecase.ErrNotFound:
			return ErrNotFound
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":    utils.DumpIncomingContext(ctx.Request().Context()),
				"formID": formID,
			}).Error(err)

			return ErrInternal
		}
	}
}
