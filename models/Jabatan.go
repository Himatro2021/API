package models

type Jabatan struct {
	ID             uint   `gorm:"primaryKey"`
	PrivilegeLevel int    `gorm:"not null"`
	Name           string `gorm:"not null"`
}
