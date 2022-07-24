package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type absentRepository struct {
	db *gorm.DB
}

// NewAbsentRepository create instance for absent repository
func NewAbsentRepository(db *gorm.DB) model.AbsentRepository {
	return &absentRepository{
		db: db,
	}
}

// GetAbsentFormByID self explained
func (r *absentRepository) GetAbsentFormByID(ctx context.Context, id int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": id,
	})

	form := &model.AbsentForm{}
	err := r.db.WithContext(ctx).Table(form.TableName()).Where("id = ?", id).Take(form).Error
	switch err {
	case nil:
		return form, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}
}

// GetAllAbsentForm get all absent form based on value given in limit and offset. if no value sent, default
// to get all form
func (r *absentRepository) GetAllAbsentForm(ctx context.Context, limit, offset int) ([]model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"limit":  limit,
		"offset": offset,
	})

	absentForms := []model.AbsentForm{}

	err := r.db.Model(&model.AbsentForm{}).Limit(limit).Offset(offset).Scan(&absentForms).Error
	if err != nil {
		logger.Error(err)
		return absentForms, err
	}

	return absentForms, nil
}

// GetParticipantsByFormID self explained
func (r *absentRepository) GetParticipantsByFormID(ctx context.Context, id int64) ([]model.Participant, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	participants := []model.Participant{}

	query := fmt.Sprintf("select * from absent_lists al inner join users u ON u.id = al.created_by where al.absent_form_id = %d", id)
	err := r.db.WithContext(ctx).Raw(query).Scan(&participants).Error
	if err != nil {
		logger.Error(err)
		return participants, err
	}

	return participants, nil
}

// GetAbsentListByID self explained
func (r *absentRepository) GetAbsentListByID(ctx context.Context, formID, absentListID int64) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":          utils.DumpIncomingContext(ctx),
		"formID":       formID,
		"absentListID": absentListID,
	})

	absentList := &model.AbsentList{}

	err := r.db.WithContext(ctx).Model(&model.AbsentList{}).
		Where("id = ? AND absent_form_id = ?", absentListID, formID).Take(absentList).Error
	switch err {
	case nil:
		return absentList, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}

}

// GetAbsentListByCreatorID self explained
func (r *absentRepository) GetAbsentListByCreatorID(ctx context.Context, formID, creatorID int64) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"formID":    formID,
		"creatorID": creatorID,
	})

	absentList := &model.AbsentList{}
	err := r.db.WithContext(ctx).Model(&model.AbsentList{}).
		Where("absent_form_id = ? AND created_by = ?", formID, creatorID).
		Take(absentList).Error

	switch err {
	case nil:
		return absentList, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}
}

// CreateAbsentForm self explained
func (r *absentRepository) CreateAbsentForm(ctx context.Context, title string, start, finish time.Time, groupID, userID int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"title":   title,
		"groupID": groupID,
		"start":   utils.Dump(start),
		"finish":  utils.Dump(finish),
	})

	now := time.Now()

	form := &model.AbsentForm{
		ID:                 utils.GenerateID(),
		ParticipantGroupID: groupID,
		StartAt:            start,
		FinishedAt:         finish,
		Title:              title,

		// TODO implement allow update and allow confirmation
		AllowUpdateByAttendee:   false,
		AllowCreateConfirmation: false,

		CreatedAt: now,
		UpdatedAt: now,

		// TODO implement user authentication
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err := r.db.WithContext(ctx).Table(form.TableName()).Create(form).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return form, nil
}

// FillAbsentFormByAttendee self explained
func (r *absentRepository) FillAbsentFormByAttendee(ctx context.Context, userID, formID int64, status model.Status, reason string) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": formID,
		"status": string(status),
		"reason": reason,
	})

	now := time.Now()
	absentList := &model.AbsentList{
		ID:            utils.GenerateID(),
		AbsentFormID:  formID,
		CreatedAt:     now,
		UpdatedAt:     now,
		Status:        status,
		ExecuseReason: reason,
		CreatedBy:     userID,
	}

	err := r.db.WithContext(ctx).Model(&model.AbsentList{}).Create(absentList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return absentList, nil
}

// UpdateAbsentForm used by admin to update data about absent form
func (r *absentRepository) UpdateAbsentForm(ctx context.Context, absentFormID int64, title string, start, finish time.Time, groupID, userID int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"id":      absentFormID,
		"title":   title,
		"start":   start,
		"finish":  finish,
		"groupID": groupID,
		"formID":  absentFormID,
	})

	absentForm := &model.AbsentForm{}
	now := time.Now()

	err := r.db.WithContext(ctx).Model(&model.AbsentForm{}).Where("id = ?", absentFormID).Take(absentForm).Error
	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		logger.Error(err)
		return nil, err
	}

	// TODO add implementation for allow update and create confirmation
	absentForm.UpdatedAt = now

	absentForm.UpdatedBy = userID
	absentForm.StartAt = start
	absentForm.FinishedAt = finish
	absentForm.ParticipantGroupID = groupID
	absentForm.Title = title

	err = r.db.WithContext(ctx).Save(&absentForm).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return absentForm, nil
}

// UpdateAbsentListByAttendee self explained
func (r *absentRepository) UpdateAbsentListByAttendee(ctx context.Context, absentList *model.AbsentList) (*model.AbsentList, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"absentList": utils.Dump(absentList),
	})

	err := r.db.WithContext(ctx).Save(absentList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return absentList, nil
}
