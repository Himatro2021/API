package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AbsentList struct {
	gorm.Model
	FormAbsensiID uint
	NPM           string
	Keterangan    string       `gorm:"default:'?'"`
	AnggotaBiasa  AnggotaBiasa `gorm:"foreignKey:NPM"`
	FormAbsensi   FormAbsensi  `gorm:"foreignKey:FormAbsensiID"`
}

type ReturnedAbsentList struct {
	NPM            string    `json:"npm"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Keterangan     string    `json:"keterangan"`
	Nama           string    `json:"nama"`
	NamaDepartemen string    `json:"departemen"`
}
