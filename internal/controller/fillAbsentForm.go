package controller

import (
	"errors"
	"fmt"
	"himatro-api/internal/auth"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"himatro-api/internal/util"
	"net/http"
	"time"
)

func FillAbsentForm(absentID int, NPM string, keterangan string) (string, error) {
	pengurus, err := getPengurusData(NPM)

	if err != nil {
		util.LogErr("WARN", "Absent filling failed", err.Error())
		return "", err
	}

	formDetail, err := getFormAbsentDetail(absentID)

	if err != nil {
		util.LogErr("WARN", "Absent filling failed", err.Error())
		return "", err
	}

	if formDetail.Participant != pengurus.DepartemenID && formDetail.Participant != 0 {
		util.LogErr("WARN", fmt.Sprintf("Unexpected attendance on absentID: %d", absentID), NPM)
		return "", fmt.Errorf("you are not the expected attendance of this absent form")
	}

	if err := isAlreadyAttend(absentID, NPM); err != nil {
		util.LogErr("WARN", fmt.Sprintf("%s already attend absentID: %d", NPM, absentID), err.Error())
		return "", err
	}

	saveAttendanceRecord(absentID, NPM, keterangan)
	updateToken, err := auth.CreateUpdateAbsentListToken(absentID, NPM)

	if err != nil {
		util.LogErr("ERROR", "Failed to create update absent list token", err.Error())
		return "", nil
	}

	return updateToken, nil
}

func IsFormWriteable(absentID int) error {
	formAbsensi, err := getFormAbsentDetail(absentID)

	if err != nil {
		util.LogErr("WARN", "Failed to lookup the absent form", err.Error())
		return err
	}

	if formAbsensi.StartAt.After(time.Now()) {
		util.LogErr("WARN", fmt.Sprintf("Accessing to early absent form: %d", absentID), "")
		return fmt.Errorf("absent form with ID: %d is not open yet", absentID)
	}

	if formAbsensi.FinishAt.Before(time.Now()) {
		util.LogErr("WARN", fmt.Sprintf("Accessing closed absent form: %d", absentID), "")
		return fmt.Errorf("absent form with ID: %d is already closed", absentID)
	}

	return nil
}

func UpdateAbsentListByAttendant(absentID int, keterangan string, cookie *http.Cookie) error {
	tokenPayload := auth.UpdateAbsentListClaims{}

	if err := auth.ExtractJWTPayload(cookie.Value, &tokenPayload); err != nil {
		util.LogErr("WARN", "failed to update absent list by attendant", err.Error())
		return fmt.Errorf("update absent failed because: %s", err.Error())
	}

	if absentID != int(tokenPayload.AbsentID) {
		util.LogErr("WARN", "token mismatch with absentID requested", fmt.Sprintf("absentID: %d", absentID))
		return fmt.Errorf("token mismatch with absentID requested")
	}

	updateAttendanceRecord(absentID, tokenPayload.NPM, keterangan)

	return nil
}

func getFormAbsentDetail(absentID int) (models.FormAbsensi, error) {
	formAbsent := models.FormAbsensi{}

	res := db.DB.Model(&models.FormAbsensi{}).
		Where("id = ?", absentID).
		First(&formAbsent)

	if res.Error != nil {
		util.LogErr("WARN", fmt.Sprintf("absent form not found at absentID: %d", absentID), "")
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
		util.LogErr("WARN", fmt.Sprintf("pengurus with NPM: %s is not found", NPM), res.Error.Error())
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
		util.LogErr("ERROR", fmt.Sprintf("failed to fill attendance record for %d - %s", absentID, NPM), res.Error.Error())
		return errors.New("failed to fill attendance record")
	}

	if absentList.Keterangan != "?" {
		util.LogErr("WARN", fmt.Sprintf("attendant with NPM: %s is alredy filled this form", NPM), "")
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
		util.LogErr("ERROR", fmt.Sprintf("failed to fill attendance record for %d - %s", absentID, NPM), res.Error.Error())
		return errors.New("system failure to save absent record")
	}

	return nil
}

func updateAttendanceRecord(absentID int, NPM string, keterangan string) error {
	absentList := models.AbsentList{
		FormAbsensiID: uint(absentID),
		NPM:           NPM,
	}

	res := db.DB.Model(&models.AbsentList{}).
		Where(&absentList).
		First(&absentList).
		Update("keterangan", keterangan)

	if res.Error != nil {
		util.LogErr("ERROR", fmt.Sprintf("failed to fill attendance record for %d - %s", absentID, NPM), res.Error.Error())
		return errors.New("server failed to update absent list")
	}

	return nil
}
