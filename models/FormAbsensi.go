package models

import (
	"time"

	"gorm.io/gorm"
)

type FormAbsensi struct {
	gorm.Model

	Title                       string    `gorm:"not null"`
	Participant                 int       `gorm:"not null"`
	StartAt                     time.Time `gorm:"not null"`
	FinishAt                    time.Time `gorm:"not null"`
	RequireAttendanceImageProof bool      `gorm:"not null"`
	RequireExecuseImageProof    bool      `gorm:"not null"`
}

type ReturnedFormAbsentDetails struct {
	FormID                      uint      `json:"form_id"`
	Title                       string    `json:"title"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	ParticipantCode             int       `json:"participant_code"`
	StartAt                     time.Time `json:"start_at"`
	FinishAt                    time.Time `json:"finish_at"`
	RequireAttendanceImageProof bool      `json:"require_attendance_image_proof"`
	RequireExecuseImageProof    bool      `json:"require_execuse_image_proof"`
	TotalParticipant            int       `json:"total_participant"`
	Hadir                       int       `json:"hadir"`
	Izin                        int       `json:"izin"`
	TanpaKeterangan             int       `json:"tanpa_keterangan"`
}
