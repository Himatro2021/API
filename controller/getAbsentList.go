package controller

import (
	"errors"
	"fmt"
	"himatro-api/db"
	"himatro-api/models"
)

func GetAbsentListResult(absentID int) ([]models.ReturnedAbsentList, error) {
	if err := isFormAbsentExists(absentID); err != nil {
		return []models.ReturnedAbsentList{}, err
	}

	absentList, err := getAbsentListFromFormID(absentID)

	if err != nil {
		return []models.ReturnedAbsentList{}, err
	}

	return absentList, nil
}

func isFormAbsentExists(absentID int) error {
	formAbsent := models.FormAbsensi{}

	res := db.DB.Model(&models.FormAbsensi{}).
		Where("id = ?", absentID).
		First(&formAbsent)

	if res.Error != nil {
		return fmt.Errorf("absent form with ID: %d is not exists", absentID)
	}

	return nil
}

func getAbsentListFromFormID(absentID int) ([]models.ReturnedAbsentList, error) {
	absentLists := []models.ReturnedAbsentList{}

	res := db.DB.Model(&models.AbsentList{}).Select("anggota_biasas.nama, absent_lists.npm, absent_lists.updated_at, absent_lists.keterangan, departemens.nama as nama_departemen").Where(&models.AbsentList{FormAbsensiID: uint(absentID)}).Joins("inner join anggota_biasas on anggota_biasas.npm = absent_lists.npm").Joins("inner join pengurus on pengurus.npm = anggota_biasas.npm").Joins("inner join departemens on departemens.id = pengurus.departemen_id").Find(&absentLists)

	if res.Error != nil {
		return absentLists, errors.New("server failed to fetch requested absent list")
	}

	return absentLists, nil
}
