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
