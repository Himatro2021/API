package controller

import (
	"errors"
	"fmt"
	"himatro-api/auth"
	"himatro-api/db"
	"himatro-api/models"
	"net/http"
	"time"
)

func FillAbsentForm(absentID int, NPM string, keterangan string) (string, error) {
	pengurus, err := getPengurusData(NPM)

	if err != nil {
		return "", err
	}

	formDetail, err := getFormAbsentDetail(absentID)

	if err != nil {
		return "", err
	}

	if formDetail.Participant != pengurus.DepartemenID && formDetail.Participant != 0 {
		return "", fmt.Errorf("you are not the expected attendance of this absent form")
	}

	if err := isAlreadyAttend(absentID, NPM); err != nil {
		return "", err
	}

	saveAttendanceRecord(absentID, NPM, keterangan)

	return nil
}

func IsFormWriteable(absentID int) error {
	formAbsensi, err := getFormAbsentDetail(absentID)

	if err != nil {
		return err
	}

	if formAbsensi.StartAt.After(time.Now()) {
		return fmt.Errorf("absent form with ID: %d is not open yet", absentID)
	}

	if formAbsensi.FinishAt.Before(time.Now()) {
		return fmt.Errorf("absent form with ID: %d is already closed", absentID)
	}

	return nil
}

func getFormAbsentDetail(absentID int) (models.FormAbsensi, error) {
	formAbsent := models.FormAbsensi{}

	res := db.DB.Model(&models.FormAbsensi{}).
		Where("id = ?", absentID).
		First(&formAbsent)

	if res.Error != nil {
		return models.FormAbsensi{}, fmt.Errorf("absent form with ID: %d is not exists", absentID)
	}

	return formAbsent, nil
}

func getPengurusData(NPM string) (models.Pengurus, error) {
	pengurus := models.Pengurus{
		NPM: NPM,
	}

	res := db.DB.Model(&models.Pengurus{}).
		Where(&models.Pengurus{
			NPM: NPM,
		}).
		First(&pengurus)

	if res.Error != nil {
		return pengurus, fmt.Errorf("pengurus with NPM: %s is not found", NPM)
	}

	return pengurus, nil
}

func isAlreadyAttend(absentID int, NPM string) error {
	absentList := models.AbsentList{}

	res := db.DB.Model(&models.AbsentList{}).
		Where(&models.AbsentList{
			FormAbsensiID: uint(absentID),
			NPM:           NPM,
		}).First(&absentList)

	if res.Error != nil {
		return errors.New("failed to fill attendance record")
	}

	if absentList.Keterangan != "?" {
		return fmt.Errorf("attendant with NPM: %s is alredy filled this form", NPM)
	}

	return nil
}

func saveAttendanceRecord(absentID int, NPM string, keterangan string) error {
	res := db.DB.Model(&models.AbsentList{}).
		Where(&models.AbsentList{
			FormAbsensiID: uint(absentID),
			NPM:           NPM,
		}).
		Update("keterangan", keterangan)

	if res.Error != nil {
		return errors.New("system failure to save absent record")
	}

	return nil
}
