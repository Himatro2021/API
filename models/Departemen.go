package models

type Departemen struct {
	ID   int `gorm:"primaryKey"`
	Nama string
}
