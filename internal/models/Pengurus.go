package models

type Pengurus struct {
	NPM          string `gorm:"primaryKey"`
	DepartemenID int
	JabatanID    int
	AnggotaBiasa AnggotaBiasa `gorm:"foreignKey:NPM"`
	Departemen   Departemen   `gorm:"foreignKey:DepartemenID"`
	Jabatan      Jabatan      `gorm:"foreignKey:JabatanID"`
}
