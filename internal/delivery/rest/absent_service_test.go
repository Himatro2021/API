package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Himatro2021/API/auth"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAbsentService_handleGetFormByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	ctx := context.TODO()
	formID := utils.GenerateID()

	absentForm := &model.AbsentForm{
		ID: formID,
	}

	t.Run("ok - form found", func(t *testing.T) {
		mock.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetFormByID()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - params id is invalid int", func(t *testing.T) {
		ec := echo.New()
		invalidID := "akusukabaububukmesiudipagihari"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/")
		ectx.SetParamNames("id")
		ectx.SetParamValues(invalidID)

		err := service.handleGetFormByID()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		mock.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetFormByID()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err internal from usecase", func(t *testing.T) {
		mock.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetFormByID()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentService_handleCreateAbsentForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	input := &model.CreateAbsentInput{
		Title:          "ini title",
		StartAtDate:    "2009-10-23",
		StartAtTime:    "12:00",
		FinishedAtDate: "9001-12-21",
		FinishedAtTime: "12:00",
		GroupMemberID:  utils.GenerateID(),
	}

	absentForm := &model.AbsentForm{
		ID: utils.GenerateID(),
	}

	admin := auth.User{
		ID:    utils.GenerateID(),
		Email: "",
		Role:  rbac.RoleAdmin,
	}
	ctx := auth.SetUserToCtx(context.TODO(), admin)

	t.Run("ok - created", func(t *testing.T) {
		mock.EXPECT().CreateAbsentForm(ctx, input).Times(1).Return(absentForm, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, input.Title, input.StartAtDate, input.StartAtTime, input.FinishedAtDate, input.FinishedAtTime, input.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		err := service.handleCreateAbsentForm()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - payload error", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": invalid payload here,
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, input.Title, input.StartAtTime, input.FinishedAtDate, input.FinishedAtTime, input.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		err := service.handleCreateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - validation error in usecase", func(t *testing.T) {
		mock.EXPECT().CreateAbsentForm(ctx, input).Times(1).Return(nil, usecase.ErrValidation)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, input.Title, input.StartAtDate, input.StartAtTime, input.FinishedAtDate, input.FinishedAtTime, input.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		err := service.handleCreateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("err - internal error from usecase", func(t *testing.T) {
		mock.EXPECT().CreateAbsentForm(ctx, input).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, input.Title, input.StartAtDate, input.StartAtTime, input.FinishedAtDate, input.FinishedAtTime, input.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		err := service.handleCreateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - member don't have permission", func(t *testing.T) {
		member := auth.User{
			ID:   utils.GenerateID(),
			Role: rbac.RoleMember,
		}
		memberCtx := auth.SetUserToCtx(context.TODO(), member)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, input.Title, input.StartAtDate, input.StartAtTime, input.FinishedAtDate, input.FinishedAtTime, input.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(memberCtx))

		mock.EXPECT().CreateAbsentForm(memberCtx, input).Times(1).Return(nil, usecase.ErrForbidden)

		err := service.handleCreateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})
}

func TestAbsentUsecase_handleGetAllAbsentForms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	admin := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleAdmin,
	}
	ctx := auth.SetUserToCtx(context.TODO(), admin)
	limit := 0
	offset := 0
	absentForms := []model.AbsentForm{
		{
			ID: utils.GenerateID(),
		},
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(absentForms, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - limit params is invalid", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(absentForms, nil)

		ec := echo.New()

		q := make(url.Values)
		q.Set("limit", "ayeaye")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - offset params is invalid", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(absentForms, nil)

		ec := echo.New()

		q := make(url.Values)
		q.Set("offset", "ayeaye")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - query limit is set", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, 10, offset).Times(1).Return(absentForms, nil)

		ec := echo.New()

		q := make(url.Values)
		q.Set("limit", "10")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - query offset is set", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, limit, 10).Times(1).Return(absentForms, nil)

		ec := echo.New()

		q := make(url.Values)
		q.Set("offset", "10")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - query limit & offset is set", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, 10, 10).Times(1).Return(absentForms, nil)

		ec := echo.New()

		q := make(url.Values)
		q.Set("limit", "10")
		q.Set("offset", "10")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.NoError(t, err)
	})

	t.Run("err - internal err from usecase", func(t *testing.T) {
		mock.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - non admin can't get all absent", func(t *testing.T) {
		member := auth.User{
			Role: rbac.RoleMember,
		}
		memberCtx := auth.SetUserToCtx(context.TODO(), member)
		mock.EXPECT().GetAllAbsentForm(memberCtx, limit, offset).Times(1).Return(nil, usecase.ErrForbidden)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(memberCtx))

		ectx.SetPath("/absent/form/")

		err := service.handleGetAllAbsentForms()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})
}

func TestAbsentService_handleUpdateAbsentForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	admin := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleAdmin,
	}
	ctx := auth.SetUserToCtx(context.TODO(), admin)

	formID := utils.GenerateID()
	updateInput := &model.CreateAbsentInput{
		Title:          "ini title",
		StartAtDate:    "2009-10-23",
		StartAtTime:    "12:00",
		FinishedAtDate: "9001-12-21",
		FinishedAtTime: "12:00",
		GroupMemberID:  utils.GenerateID(),
	}

	newForm := &model.AbsentForm{
		ID: formID,
	}

	t.Run("ok - updated", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentForm(ctx, formID, updateInput).Times(1).Return(newForm, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleUpdateAbsentForm()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - form id invalid", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues("ini invalid id")

		err := service.handleUpdateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - payload invalid json", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s", loh kok ada ini?
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleUpdateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentForm(ctx, formID, updateInput).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleUpdateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentForm(ctx, formID, updateInput).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleUpdateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("ok - non admin can't update absent form", func(t *testing.T) {
		member := auth.User{
			ID:   utils.GenerateID(),
			Role: rbac.RoleAdmin,
		}
		memberCtx := auth.SetUserToCtx(context.TODO(), member)

		mock.EXPECT().UpdateAbsentForm(memberCtx, formID, updateInput).Times(1).Return(nil, usecase.ErrForbidden)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(fmt.Sprintf(`
			{
				"title": "%s",
				"start_at_date": "%s",
				"start_at_time": "%s",
				"finished_at_date": "%s",
				"finished_at_time": "%s",
				"group_member_id": %d
			}
		`, updateInput.Title, updateInput.StartAtDate, updateInput.StartAtTime, updateInput.FinishedAtDate, updateInput.FinishedAtTime, updateInput.GroupMemberID)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(memberCtx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleUpdateAbsentForm()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})
}

func TestAbsentService_handleFillAbsentFormByAttendee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	member := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleMember,
	}
	ctx := auth.SetUserToCtx(context.TODO(), member)

	formID := utils.GenerateID()
	status := "PRESENT"
	reason := ""

	absentList := &model.AbsentList{
		ID:           utils.GenerateID(),
		AbsentFormID: formID,
		CreatedBy:    member.ID,
	}

	t.Run("ok - filled", func(t *testing.T) {
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, status, reason).Times(1).Return(absentList, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - id invalid", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues("fmt.Sprintf(")
		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - payload invalid json", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s", loh heh
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - input invalid", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "",
				"execuse_reason": "%s"
			}
		`, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - form not exits", func(t *testing.T) {
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, status, reason).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("ok - status invalid", func(t *testing.T) {
		invalidStatus := "invalid status"
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, invalidStatus, reason).Times(1).Return(nil, usecase.ErrValidation)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, invalidStatus, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - form already closed", func(t *testing.T) {
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, status, reason).Times(1).Return(nil, usecase.ErrForbidden)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("ok - user already filled", func(t *testing.T) {
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, status, reason).Times(1).Return(nil, usecase.ErrAlreadyExists)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrAlreadyExists)
	})

	t.Run("err - internal err from usecase", func(t *testing.T) {
		mock.EXPECT().FillAbsentFormByAttendee(ctx, formID, status, reason).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`
			{
				"status": "%s",
				"execuse_reason": "%s"
			}
		`, status, reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleFillAbsentFormByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentService_handleUpdateAbsentListByAttendee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	member := auth.User{
		ID:   utils.GenerateID(),
		Role: rbac.RoleMember,
	}
	ctx := auth.SetUserToCtx(context.TODO(), member)

	absentID := utils.GenerateID()
	updateInput := &model.UpdateAbsentListInput{
		AbsentFormID: absentID,
		Status:       "EXECUSE",
		Reason:       "",
	}
	absentList := &model.AbsentList{
		ID: absentID,
	}

	t.Run("ok - updated", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(absentList, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - id invalid", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues("invalid id")

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrBadRequest)
	})

	t.Run("ok - json payload err", func(t *testing.T) {
		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,, <-
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Error(t, err, ErrBadRequest)
	})

	t.Run("ok - validation error", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(nil, usecase.ErrValidation)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("ok - updating is forbidden", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(nil, usecase.ErrForbidden)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("err - internal", func(t *testing.T) {
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - updating another user absent list", func(t *testing.T) {
		// improve: how do we really mock that the usecase detect the user trying to update
		// another user list? the err returned is the same type with case when user try to fill already closed form
		// any ideas?
		mock.EXPECT().UpdateAbsentListByAttendee(ctx, absentID, updateInput).Times(1).Return(nil, usecase.ErrForbidden)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(fmt.Sprintf(`
			{
				"absent_form_id": %d,
				"status": "%s",
				"reason": "%s"
			}
		`, updateInput.AbsentFormID, updateInput.Status, updateInput.Reason)))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)
		ectx.SetRequest(ectx.Request().WithContext(ctx))

		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", absentID))

		err := service.handleUpdateAbsentListByAttendee()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})
}

func TestAbsentUsecase_handleGetParticipantsByFormID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockAbsentUsecase(ctrl)
	service := &Service{
		absentUsecase: mock,
	}

	ctx := context.TODO()
	formID := utils.GenerateID()

	result := &model.AbsentResult{
		Title: "ok",
	}

	t.Run("ok", func(t *testing.T) {
		mock.EXPECT().GetAbsentResultByFormID(ctx, formID).Times(1).Return(result, nil)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/:id/result")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetParticipantsByFormID()(ectx)

		assert.NoError(t, err)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		mock.EXPECT().GetAbsentResultByFormID(ctx, formID).Times(1).Return(nil, usecase.ErrNotFound)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/:id/result")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetParticipantsByFormID()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - internal err", func(t *testing.T) {
		mock.EXPECT().GetAbsentResultByFormID(ctx, formID).Times(1).Return(nil, usecase.ErrInternal)

		ec := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		ectx.SetPath("/absent/form/:id/result")
		ectx.SetParamNames("id")
		ectx.SetParamValues(fmt.Sprintf("%d", formID))

		err := service.handleGetParticipantsByFormID()(ectx)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}
