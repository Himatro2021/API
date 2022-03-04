package models

import "github.com/jinzhu/gorm"

type AbsensiPengurus struct {
	gorm.Model
	FormAbsensiID uint
	NPM           string
	AnggotaBiasa  AnggotaBiasa `gorm:"foreignKey:NPM"`
	FormAbsensi   FormAbsensi  `gorm:"foreignKey:FormAbsensiID"`
}
