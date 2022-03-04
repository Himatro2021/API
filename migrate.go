package main

import (
	"himatro-api/db"
	"himatro-api/models"
)

func Migrate() {
	db.Connect()

	db.DB.AutoMigrate(&models.AnggotaBiasa{})
	db.DB.AutoMigrate(&models.Jabatan{})
	db.DB.AutoMigrate(&models.Pengurus{})
	db.DB.AutoMigrate(&models.FormAbsensi{})
	db.DB.AutoMigrate(&models.AbsensiPengurus{})
	db.DB.AutoMigrate(&models.Departemen{})
	db.DB.AutoMigrate(&models.User{})
}
