package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	NPM      string
	Password string

	AnggotaBiasa AnggotaBiasa `gorm:"foreignKey:NPM"`
}
