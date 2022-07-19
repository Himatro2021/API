package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
)

type absentUsecase struct {
	absentRepo model.AbsentRepository
}

// NewAbsentUsecase return absent usecase instance
func NewAbsentUsecase(absentRepo model.AbsentRepository) model.AbsentUsecase {
	return &absentUsecase{
		absentRepo: absentRepo,
	}
}

// GetAbsentFormByID self explained
func (au *absentUsecase) GetAbsentFormByID(ctx context.Context, id int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": id,
	})

	absentForm, err := au.absentRepo.GetAbsentFormByID(ctx, id)
	switch err {
	case nil:
		return absentForm, nil
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

}

// GetAbsentResultByFormID self explained
func (au *absentUsecase) GetAbsentResultByFormID(ctx context.Context, id int64) (*model.AbsentResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": id,
	})

	absentResult := &model.AbsentResult{}

	form, err := au.absentRepo.GetAbsentFormByID(ctx, id)
	switch err {
	case nil:
		break
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

	absentResult.Title = form.Title
	absentResult.StartAt = form.StartAt
	absentResult.FinishedAt = form.FinishedAt

	participants, err := au.absentRepo.GetParticipantsByFormID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	absentResult.Participants = participants

	return absentResult, nil
}

// GetAllAbsentForm get all absent form based on value given in limit and offset. if no value sent, default
// to get all form
func (au *absentUsecase) GetAllAbsentForm(ctx context.Context, limit, offset int) ([]model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"limit":  limit,
		"offset": offset,
	})

	absentForms, err := au.absentRepo.GetAllAbsentForm(ctx, limit, offset)
	if err != nil {
		logger.Error(err)
		return absentForms, ErrInternal
	}

	return absentForms, nil
}

// CreateAbsentForm self explained
func (au *absentUsecase) CreateAbsentForm(ctx context.Context, input *model.CreateAbsentInput) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		logger.Error(err)
		return nil, ErrValidation
	}

	start, err := helper.ParseDateAndTimeStringToTime(input.StartAtDate, input.StartAtTime)
	if err != nil {
		logger.Warn(err)
		return nil, ErrValidation
	}

	finish, err := helper.ParseDateAndTimeStringToTime(input.FinishedAtDate, input.FinishedAtTime)
	if err != nil {
		logger.Warn(err)
		return nil, ErrValidation
	}

	if !helper.IsStartAndFinishTimeValid(start, finish) {
		return nil, ErrValidation
	}

	absentForm, err := au.absentRepo.CreateAbsentForm(ctx, input.Title, start, finish, input.GroupMemberID)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return absentForm, nil
}

// CreateConfirmationOnAbsentForm not implemented
// https://github.com/Himatro2021/API/issues/17
func (au *absentUsecase) CreateConfirmationOnAbsentForm(ctx context.Context, absentFormID int64, status, reason string) (*model.AbsentList, error) {
	return nil, nil
}

// FillAbsentFormByAttendee self explained
func (au *absentUsecase) FillAbsentFormByAttendee(ctx context.Context, absentFormID int64, status, reason string) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"status": status,
		"reason": reason,
		"formID": absentFormID,
	})

	presenceStatus, err := au.transformStringToStatus(status)
	if err != nil {
		return nil, ErrValidation
	}

	absentForm, err := au.absentRepo.GetAbsentFormByID(ctx, absentFormID)
	switch err {
	case nil:
		break
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

	if !absentForm.IsStillOpen() {
		return nil, ErrForbidden
	}

	// TODO use ctx based user id
	// https://github.com/Himatro2021/API/issues/21
	absentList, err := au.absentRepo.GetAbsentListByCreatorID(ctx, absentFormID, int64(1))
	if err != nil && err != repository.ErrNotFound {
		logger.Error(err)
		return nil, ErrInternal
	}

	if absentList != nil {
		return nil, ErrAlreadyExists
	}

	// TODO add group checking between absent form group and user id
	// https://github.com/Himatro2021/API/issues/20

	// TODO use ctx based user id
	// https://github.com/Himatro2021/API/issues/21
	absentList, err = au.absentRepo.FillAbsentFormByAttendee(ctx, int64(1), absentFormID, presenceStatus, reason)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return absentList, nil
}

// UpdateAbsentListByAttendee self explained
func (au *absentUsecase) UpdateAbsentListByAttendee(ctx context.Context, absentListID int64, input *model.UpdateAbsentListInput) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":          utils.DumpIncomingContext(ctx),
		"absentListID": absentListID,
		"input":        utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		return nil, ErrValidation
	}

	presenceStatus, err := au.transformStringToStatus(input.Status)
	if err != nil {
		return nil, ErrValidation
	}

	absentForm, err := au.absentRepo.GetAbsentFormByID(ctx, input.AbsentFormID)
	switch err {
	case nil:
		break
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

	if !absentForm.AllowUpdateByAttendee || !absentForm.IsStillOpen() {
		return nil, ErrForbidden
	}

	absentList, err := au.absentRepo.GetAbsentListByID(ctx, input.AbsentFormID, absentListID)
	switch err {
	case nil:
		break
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}

	// TODO check whether the user is the owner of that absent list

	absentList.Status = presenceStatus
	absentList.ExecuseReason = input.Reason

	result, err := au.absentRepo.UpdateAbsentListByAttendee(ctx, absentList)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return result, nil
}

// UpdateAbsentForm used by admin to update absent form
func (au *absentUsecase) UpdateAbsentForm(ctx context.Context, absentFormID int64, input *model.CreateAbsentInput) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"input":  utils.Dump(input),
		"formID": absentFormID,
	})

	if err := input.Validate(); err != nil {
		return nil, ErrValidation
	}

	start, err := helper.ParseDateAndTimeStringToTime(input.StartAtDate, input.StartAtTime)
	if err != nil {
		logger.Warn(err)
		return nil, ErrValidation
	}

	finish, err := helper.ParseDateAndTimeStringToTime(input.FinishedAtDate, input.FinishedAtTime)
	if err != nil {
		logger.Warn(err)
		return nil, ErrValidation
	}

	if !helper.IsStartAndFinishTimeValid(start, finish) {
		return nil, ErrValidation
	}

	newAbsentForm, err := au.absentRepo.UpdateAbsentForm(ctx, absentFormID, input.Title, start, finish, input.GroupMemberID)
	switch err {
	case nil:
		return newAbsentForm, nil
	case repository.ErrNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, ErrInternal
	}
}

// transformStringToStatus convert string to equal model.Status. call panic if convertion is other than model.Status defined
func (au *absentUsecase) transformStringToStatus(s string) (model.Status, error) {
	status := strings.ToLower(s)
	switch status {
	case "present":
		return model.Present, nil
	case "absent":
		return model.Absent, nil
	case "execuse":
		return model.Execuse, nil
	case "pending_present":
		return model.PendingPresent, nil
	case "pending_execuse":
		return model.PendingExecuse, nil
	default:
		return model.PendingPresent, errors.New("invalid status")
	}
}

// func (au *absentUsecase) transformAbsentListToAbsentResult(absentForm *model.AbsentForm, absentList []model.AbsentList) (*model.AbsentResult, error) {
// 	participants := []model.Participant{}

// 	for _, list := range absentList {
// 		participants = append(participants, model.Participant{
// 			Name: list.,
// 		})
// 	}
// }