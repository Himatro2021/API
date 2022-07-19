package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func TestAbsentUsecase_GetAbsentFormByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	formID := utils.GenerateID()

	t.Run("ok - found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(&model.AbsentForm{
			ID: formID,
		}, nil)

		absentForm, err := uc.GetAbsentFormByID(ctx, formID)

		assert.NoError(t, err)
		assert.Equal(t, absentForm.ID, formID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.GetAbsentFormByID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err from db", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, errors.New("err ffrom db"))

		_, err := uc.GetAbsentFormByID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentUsecase_GetAllAbsentForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	absentForms := []model.AbsentForm{
		{
			ID:                      utils.GenerateID(),
			ParticipantGroupID:      0,
			StartAt:                 time.Time{},
			FinishedAt:              time.Time{},
			Title:                   "",
			AllowUpdateByAttendee:   false,
			AllowCreateConfirmation: false,
			CreatedAt:               time.Time{},
			UpdatedAt:               time.Time{},
			DeletedAt:               gorm.DeletedAt{},
			CreatedBy:               0,
			UpdatedBy:               0,
			DeletedBy:               null.Int{},
		},
	}

	t.Run("ok - found", func(t *testing.T) {
		limit := 1
		offset := 0

		repo.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(absentForms, nil)

		result, err := uc.GetAllAbsentForm(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, result, absentForms)
	})

	t.Run("ok - return zero absent forms", func(t *testing.T) {
		limit := 1
		offset := 0

		repo.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return([]model.AbsentForm{}, nil)

		result, err := uc.GetAllAbsentForm(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(result))

	})

	t.Run("err from repo", func(t *testing.T) {
		limit := 1
		offset := 0

		repo.EXPECT().GetAllAbsentForm(ctx, limit, offset).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.GetAllAbsentForm(ctx, limit, offset)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentUsecase_CreateAbsentForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	input := &model.CreateAbsentInput{
		Title:          "ini title",
		StartAtDate:    "2001-10-20",
		StartAtTime:    "10:11",
		FinishedAtDate: "3002-01-23",
		FinishedAtTime: "12:09",
		GroupMemberID:  utils.GenerateID(),
	}
	start, _ := helper.ParseDateAndTimeStringToTime(input.StartAtDate, input.StartAtTime)
	finish, _ := helper.ParseDateAndTimeStringToTime(input.FinishedAtDate, input.FinishedAtTime)
	absentForm := &model.AbsentForm{
		ID:                 utils.GenerateID(),
		ParticipantGroupID: input.GroupMemberID,
		Title:              input.Title,
		StartAt:            start,
		FinishedAt:         finish,
	}

	t.Run("ok - created", func(t *testing.T) {
		repo.EXPECT().CreateAbsentForm(ctx, input.Title, start, finish, input.GroupMemberID).Times(1).Return(absentForm, nil)

		result, err := uc.CreateAbsentForm(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, result.ID, absentForm.ID)
	})

	t.Run("err when creating form", func(t *testing.T) {
		repo.EXPECT().CreateAbsentForm(ctx, input.Title, start, finish, input.GroupMemberID).Times(1).Return(nil, errors.New("err from db"))

		_, err := uc.CreateAbsentForm(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - start date invalid format", func(t *testing.T) {
		inputCopy := input
		inputCopy.StartAtDate = "2001.12.1"

		_, err := uc.CreateAbsentForm(ctx, inputCopy)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - finish date invalid format", func(t *testing.T) {
		inputCopy := input
		inputCopy.FinishedAtDate = "2001.12.1"

		_, err := uc.CreateAbsentForm(ctx, inputCopy)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - finish time invalid format", func(t *testing.T) {
		inputCopy := input
		inputCopy.FinishedAtTime = "24:09"

		_, err := uc.CreateAbsentForm(ctx, inputCopy)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - start time invalid format", func(t *testing.T) {
		inputCopy := input
		inputCopy.StartAtTime = "24:09"

		_, err := uc.CreateAbsentForm(ctx, inputCopy)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - start and finish time backward", func(t *testing.T) {
		backwardTime := &model.CreateAbsentInput{
			StartAtDate:    "3004-10-01",
			StartAtTime:    "12:34",
			FinishedAtDate: "3003-11-02",
			FinishedAtTime: "11:23",
		}

		_, err := uc.CreateAbsentForm(ctx, backwardTime)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})
}

func TestAbsentUsecase_FillAbsentFormByAttendee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	formID := utils.GenerateID()
	statusPresent := model.Present
	status := "present"
	reason := "nothing"
	absentForm := &model.AbsentForm{
		ID: formID,
		// default to always open, unless needed to
		// mock when form is closed
		StartAt:    time.Now().Add(-time.Hour * 1),
		FinishedAt: time.Now().Add(time.Hour * 1),
	}
	absentList := &model.AbsentList{
		ID:           utils.GenerateID(),
		AbsentFormID: formID,
	}

	t.Run("ok - filled", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		// TODO use ctx based to get user id
		repo.EXPECT().GetAbsentListByCreatorID(ctx, formID, int64(1)).Times(1).Return(nil, repository.ErrNotFound)
		repo.EXPECT().FillAbsentFormByAttendee(ctx, int64(1), formID, statusPresent, reason).Return(absentList, nil)

		result, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.NoError(t, err)
		assert.Equal(t, result.AbsentFormID, formID)
	})

	t.Run("ok - form is not found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("error when getting form", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("user already fill form", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByCreatorID(ctx, formID, int64(1)).Times(1).Return(absentList, nil)

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrAlreadyExists)
	})

	t.Run("error failed to check is user already fill", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByCreatorID(ctx, formID, int64(1)).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - form is closed", func(t *testing.T) {
		closedForm := &model.AbsentForm{
			ID:         absentForm.ID,
			StartAt:    time.Now().Add(-time.Hour * 7),
			FinishedAt: time.Now().Add(-time.Hour * 6),
		}

		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(closedForm, nil)

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("ok - received error status string", func(t *testing.T) {
		errorStatus := "error status?"

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, errorStatus, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("error failed to fill absent from db", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		// TODO use ctx based to get user id
		repo.EXPECT().GetAbsentListByCreatorID(ctx, formID, int64(1)).Times(1).Return(nil, repository.ErrNotFound)
		repo.EXPECT().FillAbsentFormByAttendee(ctx, int64(1), formID, statusPresent, reason).Return(nil, errors.New("err db"))

		_, err := uc.FillAbsentFormByAttendee(ctx, formID, status, reason)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentUsecase_UpdateAbsentListByAttendee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	formID := utils.GenerateID()
	//statusPresent := model.PRESENT
	status := "present"
	reason := "nothing"
	absentForm := &model.AbsentForm{
		ID: formID,
		// default to always open, unless needed to
		// mock when form is closed
		StartAt:    time.Now().Add(-time.Hour * 1),
		FinishedAt: time.Now().Add(time.Hour * 1),

		// default to make the form always open to update
		AllowUpdateByAttendee: true,
	}
	absentList := &model.AbsentList{
		ID:           utils.GenerateID(),
		AbsentFormID: formID,
	}
	updateAbsentListInput := &model.UpdateAbsentListInput{
		AbsentFormID: formID,
		Status:       status,
		Reason:       reason,
	}

	t.Run("ok - updated", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByID(ctx, formID, absentList.ID).Times(1).Return(absentList, nil)
		repo.EXPECT().UpdateAbsentListByAttendee(ctx, absentList).Times(1).Return(absentList, nil)

		result, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.NoError(t, err)
		assert.Equal(t, result.ID, absentList.ID)
	})

	t.Run("ok - input is invalid", func(t *testing.T) {
		invalidInput := &model.UpdateAbsentListInput{
			AbsentFormID: utils.GenerateID(),
			Status:       "",
			Reason:       "inimah ga dicek",
		}

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, invalidInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - input status is invalid", func(t *testing.T) {
		invalidInput := &model.UpdateAbsentListInput{
			AbsentFormID: utils.GenerateID(),
			Status:       "invalid status xixixi",
			Reason:       "inimah ga dicek",
		}

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, invalidInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - absent form not found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("error - when getting absent form", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("ok - form is not allowed to update", func(t *testing.T) {
		updateDisabledForm := &model.AbsentForm{
			AllowUpdateByAttendee: false,
		}

		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(updateDisabledForm, nil)

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("ok - form is not started yet", func(t *testing.T) {
		notStartedForm := &model.AbsentForm{
			StartAt: time.Now().Add(time.Hour * 1),
		}

		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(notStartedForm, nil)

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrForbidden)
	})

	t.Run("ok - absent list not found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByID(ctx, formID, absentList.ID).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("error - failed to get absent list", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByID(ctx, formID, absentList.ID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("error - failed to insert absent list", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetAbsentListByID(ctx, formID, absentList.ID).Times(1).Return(absentList, nil)
		repo.EXPECT().UpdateAbsentListByAttendee(ctx, absentList).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.UpdateAbsentListByAttendee(ctx, absentList.ID, updateAbsentListInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentUsecase_UpdateAbsentForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	formID := utils.GenerateID()
	groupID := utils.GenerateID()
	title := "ini judul"

	updateInput := &model.CreateAbsentInput{
		Title:          title,
		GroupMemberID:  groupID,
		StartAtDate:    "2001-10-02",
		StartAtTime:    "23:00",
		FinishedAtDate: "6001-12-01",
		FinishedAtTime: "13:45",
	}

	absentForm := &model.AbsentForm{
		ID:                 formID,
		ParticipantGroupID: groupID,
		Title:              title,
	}

	start, _ := helper.ParseDateAndTimeStringToTime(updateInput.StartAtDate, updateInput.StartAtTime)
	finish, _ := helper.ParseDateAndTimeStringToTime(updateInput.FinishedAtDate, updateInput.FinishedAtTime)

	t.Run("ok - updated", func(t *testing.T) {
		repo.EXPECT().UpdateAbsentForm(ctx, formID, title, start, finish, groupID).Times(1).Return(absentForm, nil)

		result, err := uc.UpdateAbsentForm(ctx, formID, updateInput)

		assert.NoError(t, err)
		assert.Equal(t, result.ID, formID)
	})

	t.Run("ok - invalid input start date received", func(t *testing.T) {
		invalidInput := &model.CreateAbsentInput{
			StartAtDate: "2001-13-09",
		}

		_, err := uc.UpdateAbsentForm(ctx, formID, invalidInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - invalid input finish date received", func(t *testing.T) {
		invalidInput := &model.CreateAbsentInput{
			FinishedAtDate: "2001-13-09",
		}

		_, err := uc.UpdateAbsentForm(ctx, formID, invalidInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - invalid input start and finish date received", func(t *testing.T) {
		invalidInput := &model.CreateAbsentInput{
			StartAtDate:    "2021-12-09",
			StartAtTime:    "23:11",
			FinishedAtDate: "2021-12-12",
			FinishedAtTime: "22:10",
		}

		_, err := uc.UpdateAbsentForm(ctx, formID, invalidInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("ok - backward start and finish time received", func(t *testing.T) {
		backwardTime := &model.CreateAbsentInput{
			StartAtDate:    "5023-12-09",
			StartAtTime:    "23:11",
			FinishedAtDate: "4023-12-12",
			FinishedAtTime: "22:10",
		}

		_, err := uc.UpdateAbsentForm(ctx, formID, backwardTime)

		assert.Error(t, err)
		assert.Equal(t, err, ErrValidation)
	})

	t.Run("err from db", func(t *testing.T) {
		repo.EXPECT().UpdateAbsentForm(ctx, formID, title, start, finish, groupID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.UpdateAbsentForm(ctx, formID, updateInput)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})
}

func TestAbsentUsecase_GetAbsentResultByFormID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	repo := mock.NewMockAbsentRepository(ctrl)
	uc := absentUsecase{
		absentRepo: repo,
	}

	formID := utils.GenerateID()
	participants := []model.Participant{
		{
			Name:     "lucky",
			FilledAt: time.Now().Format(time.RFC3339),
			Status:   model.Present,
			Reason:   "",
		},
	}
	absentForm := &model.AbsentForm{
		ID: formID,
	}

	t.Run("ok", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetParticipantsByFormID(ctx, formID).Times(1).Return(participants, nil)

		_, err := uc.GetAbsentResultByFormID(ctx, formID)

		assert.NoError(t, err)
	})

	t.Run("ok - form not found", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, repository.ErrNotFound)

		_, err := uc.GetAbsentResultByFormID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - db err when get absent form", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.GetAbsentResultByFormID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrInternal)
	})

	t.Run("err - db err when getting participants", func(t *testing.T) {
		repo.EXPECT().GetAbsentFormByID(ctx, formID).Times(1).Return(absentForm, nil)
		repo.EXPECT().GetParticipantsByFormID(ctx, formID).Times(1).Return(nil, errors.New("err db"))

		_, err := uc.GetAbsentResultByFormID(ctx, formID)

		assert.Error(t, err)
		assert.Error(t, err, ErrInternal)
	})
}