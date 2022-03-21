package models

type AnggotaBiasa struct {
	NPM  string `gorm:"primaryKey"`
	Nama string
}
