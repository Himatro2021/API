package model

import (
	"context"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// CreateAbsentInput define input for creating absent
type CreateAbsentInput struct {
	Title          string `json:"title" validate:"required,min=5,max=255"`
	StartAtDate    string `json:"start_at_date" validate:"date,required"`
	StartAtTime    string `json:"start_at_time" validate:"required,time"`
	FinishedAtDate string `json:"finished_at_date" validate:"required,date"`
	FinishedAtTime string `json:"finished_at_time" validate:"required,time"`
	GroupMemberID  int64  `json:"group_member_id" validate:"required"`
}

// Validate return whether the CreateAbsentInput struct values are valid based on the
// defined struct tag
func (i *CreateAbsentInput) Validate() error {
	if err := Validator.Struct(i); err != nil {
		return err
	}

	return nil
}

// AbsentUsecase :nodoc:
type AbsentUsecase interface {
	GetAbsentFormByID(ctx context.Context, id int64) (*AbsentForm, error)
	GetAbsentResultByFormID(ctx context.Context, id int64) (*AbsentResult, error)
	GetAllAbsentForm(ctx context.Context, limit, offset int) ([]AbsentForm, error)
	CreateAbsentForm(ctx context.Context, input *CreateAbsentInput) (*AbsentForm, error)
	CreateConfirmationOnAbsentForm(ctx context.Context, absentFormID int64, status, reason string) (*AbsentList, error)
	FillAbsentFormByAttendee(ctx context.Context, absentFormID int64, status, reason string) (*AbsentList, error)
	UpdateAbsentListByAttendee(ctx context.Context, absentListID int64, input *UpdateAbsentListInput) (*AbsentList, error)
	UpdateAbsentForm(ctx context.Context, absentFormID int64, input *CreateAbsentInput) (*AbsentForm, error)
}

// AbsentRepository :nodoc:
type AbsentRepository interface {
	GetAbsentFormByID(ctx context.Context, id int64) (*AbsentForm, error)
	GetAllAbsentForm(ctx context.Context, limit, offset int) ([]AbsentForm, error)
	GetParticipantsByFormID(ctx context.Context, id int64) ([]Participant, error)
	GetAbsentListByID(ctx context.Context, formID, absentListID int64) (*AbsentList, error)
	GetAbsentListByCreatorID(ctx context.Context, formID, creatorID int64) (*AbsentList, error)
	CreateAbsentForm(ctx context.Context, title string, start, finish time.Time, groupID int64) (*AbsentForm, error)
	FillAbsentFormByAttendee(ctx context.Context, userID, formID int64, status Status, reason string) (*AbsentList, error)
	UpdateAbsentForm(ctx context.Context, absentFormID int64, title string, start, finish time.Time, groupID int64) (*AbsentForm, error)
	UpdateAbsentListByAttendee(ctx context.Context, absentList *AbsentList) (*AbsentList, error)
}

// AbsentForm :nodoc:
type AbsentForm struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	ParticipantGroupID int64     `json:"participant_group_id" gorm:"not null"`
	StartAt            time.Time `json:"start_at" gorm:"not null"`
	FinishedAt         time.Time `json:"finished_at" gorm:"not null"`
	Title              string    `json:"title" gorm:"not null"`

	// TODO: create feature for allow update by attende
	// https://github.com/Himatro2021/API/issues/18
	AllowUpdateByAttendee bool `json:"allow_update_by_attendee" gorm:"default:false"`

	// TODO: create feature for create confirmation
	// https://github.com/Himatro2021/API/issues/17
	AllowCreateConfirmation bool `json:"allow_create_confirmation" gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"default:null"`

	CreatedBy int64    `json:"created_by" gorm:"not null"`
	UpdatedBy int64    `json:"updated_by" gorm:"not null"`
	DeletedBy null.Int `json:"deleted_by" gorm:"default:null"`
}

// TableName return what name should be used by gorm
func (a AbsentForm) TableName() string {
	return "absent_forms"
}

// IsStillOpen checks the form is it open to write by attendee to make an absent list
// based on whether the form is deleted, and by timing
func (a *AbsentForm) IsStillOpen() bool {
	if a == nil {
		return false
	}

	// check if deleted, form must not able to write absent list again
	if a.DeletedBy.Valid {
		return false
	}

	now := time.Now()

	if a.StartAt.Before(now) && a.FinishedAt.Before(now) {
		return false
	}

	return true
}

// Status define absent status
type Status string

const (
	// Present used when participant declaring present on absent list
	Present Status = "PRESENT"

	// Absent used when participant declaring absent on absent list
	Absent Status = "ABSENT"

	// Execuse used when participant declaring execuse on absent list
	Execuse Status = "EXECUSE"

	// PendingPresent used when participant want to / will present absent on absent list
	PendingPresent Status = "PENDING_PRESENT"

	// PendingExecuse used when participant want to / will make an execuse on absent list
	PendingExecuse Status = "PENDING_EXECUSE"
)

// FillAbsentFormInput used to define filling absent form payload
type FillAbsentFormInput struct {
	Status        string `json:"status" validate:"required"`
	ExecuseReason string `json:"execuse_reason"`
}

// Validate return whether the FillAbsentFormInput struct values are valid based on the
// defined struct tag
func (f *FillAbsentFormInput) Validate() error {
	return Validator.Struct(f)
}

// UpdateAbsentListInput represent input payload for update absent list
type UpdateAbsentListInput struct {
	AbsentFormID int64  `json:"absent_form_id" validate:"required"`
	Status       string `json:"status" validate:"required"`
	Reason       string `json:"reason"`
}

// Validate return whether the UpdateAbsentListInput struct values are valid based on the
// defined struct tag
func (u *UpdateAbsentListInput) Validate() error {
	return Validator.Struct(u)
}

// AbsentList represent model for database table for absent list
type AbsentList struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	AbsentFormID  int64     `json:"absent_form_id"`
	CreatedBy     int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:null"`
	Status        Status    `json:"status" gorm:"not null"`
	ExecuseReason string    `json:"reason" gorm:"default:null"`
}

// Participant define what structure should be on absent result
type Participant struct {
	Name     string `json:"name" gorm:"column:name"`
	FilledAt string `json:"filled_at" gorm:"column:updated_at"`
	Status   Status `json:"status" gorm:"column:status"`
	Reason   string `json:"reason" gorm:"column:execuse_reason"`
}

// AbsentResult don't confuse with AbsentList
// This is used to render the absent list to user
// This also act like translation from AbsentList to be rendered to user
type AbsentResult struct {
	Title        string        `json:"title" gorm:"column:title"`
	StartAt      time.Time     `json:"start_at" gorm:"column:start_at"`
	FinishedAt   time.Time     `json:"finished_at" gorm:"column:finished_at"`
	Participants []Participant `json:"participants"`
}
