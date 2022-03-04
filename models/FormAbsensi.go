package models

import (
	"time"

	"gorm.io/gorm"
)

type FormAbsensi struct {
	gorm.Model

	ClosedAt time.Time
}
