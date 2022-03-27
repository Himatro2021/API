package controller

import (
	"errors"
	"fmt"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"himatro-api/internal/util"
)

func UpdateFormTitle(absentID int, title string) (string, error) {
	newAbsentForm := models.FormAbsensi{
		Title: title,
	}

	if err := updateAbsentFormDetail(newAbsentForm, absentID); err != nil {
		util.LogErr("ERROR", "failed to update absent form title", err.Error())
		return "", err
	}

	return title, nil
}

func UpdateParticipant(formID int, newParticipant string) error {
	formDetail, err := getFormDetail(formID)

	if err != nil {
		util.LogErr("ERROR", "failed to get form detail", err.Error())
		return err
	}

	participantCode, err := validateParticipantCode(newParticipant)

	if err != nil {
		util.LogErr("WARN", fmt.Sprintf("participant with code: %s is invalid", newParticipant), err.Error())
		return fmt.Errorf("participant with code: %s is invalid", newParticipant)
	}

	if formDetail.Participant == participantCode {
		return nil
	}

	if err := isParticipantChangeable(formID); err != nil {
		util.LogErr("WARN", fmt.Sprintf("participant on absentID: %d is not changeable", formID), err.Error())
		return err
	}

	updateFormParticipant(formID, participantCode)
	deleteOldAbsentList(formID)

	if err := createNewAbsentList(formID, participantCode); err != nil {
		util.LogErr("ERROR", "Failed to generate new absent list", err.Error())
		return err
	}

	return nil
}

func UpdateAbsentFormStartAt(formID int, startAtDate, startAtTime string) (string, error) {
	newStartAt, err := parseDate(startAtDate, startAtTime)

	if err != nil {
		util.LogErr("WARN", "invalid date time string received", err.Error())
		return "", fmt.Errorf("invalid date time string received")
	}

	formDetail, err := getFormDetail(formID)

	if err != nil {
		util.LogErr("WARN", "Failed to get form detail", err.Error())
		return "", err
	}

	if newStartAt.After(formDetail.FinishAt) {
		util.LogErr("WARN", "Failed to create new form absent can't start after it's end date", "")
		return "", errors.New("form absent can't start after it's end date")
	}

	if newStartAt == formDetail.StartAt {
		return "", nil
	}

	if newStartAt.String() == formDetail.FinishAt.String() {
		util.LogErr("WARN", "form absent cant't start and end in the same time", "")
		return "", errors.New("form absent cant't start and end in the same time")
	}

	newAbsentForm := models.FormAbsensi{
		StartAt: newStartAt,
	}

	if err := updateAbsentFormDetail(newAbsentForm, formID); err != nil {
		util.LogErr("ERROR", "Failed to update form start at", err.Error())
		return "", errors.New("server failure to update form details")
	}

	return newStartAt.String(), nil
}

func UpdateAbsentFormFinishAt(formID int, finishAtDate, finishAtTime string) (string, error) {
	newFinishAt, err := parseDate(finishAtDate, finishAtTime)

	if err != nil {
		util.LogErr("WARN", "invalid date time string received", err.Error())
		return "", fmt.Errorf("invalid date time string received")
	}

	formDetail, err := getFormDetail(formID)

	if err != nil {
		util.LogErr("WARN", "Failed to get form detail", err.Error())
		return "", err
	}

	if newFinishAt.Before(formDetail.StartAt) {
		util.LogErr("WARN", "form absent can't end before it's start date", "")
		return "", errors.New("form absent can't end before it's start date")
	}

	if newFinishAt == formDetail.FinishAt {
		return "", nil
	}

	if newFinishAt.String() == formDetail.StartAt.String() {
		util.LogErr("WARN", "form absent cant't start and end in the same time", "")
		return "", errors.New("form absent cant't start and end in the same time")
	}

	newAbsentForm := models.FormAbsensi{
		FinishAt: newFinishAt,
	}

	if err := updateAbsentFormDetail(newAbsentForm, formID); err != nil {
		util.LogErr("ERROR", "Failed to update form start at", err.Error())
		return "", errors.New("server failure to update form details")
	}

	return newFinishAt.String(), nil
}

func UpdateAbsentFormExecuseImageProof(formID int, proof bool) error {
	absentForm := models.FormAbsensi{}

	err := db.DB.Model(&absentForm).Where("id = ?", formID).First(&absentForm)

	if err.Error != nil {
		util.LogErr("WARN", fmt.Sprintf("form with ID: %d is not exists", formID), err.Error.Error())
		return fmt.Errorf("form with ID: %d is not exists", formID)
	}

	absentForm.RequireExecuseImageProof = proof

	db.DB.Save(&absentForm)

	return nil
}

func UpdateAbsentFormAttendanceImageProof(formID int, proof bool) error {
	absentForm := models.FormAbsensi{}

	err := db.DB.Model(&absentForm).Where("id = ?", formID).First(&absentForm)

	if err.Error != nil {
		util.LogErr("WARN", fmt.Sprintf("form with ID: %d is not exists", formID), err.Error.Error())
		return fmt.Errorf("form with ID: %d is not exists", formID)
	}

	absentForm.RequireAttendanceImageProof = proof

	db.DB.Save(&absentForm)

	return nil
}

func updateAbsentFormDetail(absentForm models.FormAbsensi, absentID int) error {
	res := db.DB.Model(&models.FormAbsensi{}).Where("id = ?", absentID).Updates(&absentForm)

	if res.RowsAffected == 0 {
		util.LogErr("WARN", fmt.Sprintf("absent form for ID: %d is not found", absentID), res.Error.Error())
		return fmt.Errorf("absent form for ID: %d is not found", absentID)
	}

	return nil
}

func deleteOldAbsentList(formID int) {
	absentList := &models.AbsentList{
		FormAbsensiID: uint(formID),
	}

	db.DB.Model(&models.AbsentList{}).
		Where(&models.AbsentList{FormAbsensiID: uint(formID)}).
		Delete(absentList)
}

func createNewAbsentList(formID int, participantCode int) error {
	NPMs, err := getAllNPMFromDepartemenID(participantCode)

	if err != nil {
		util.LogErr("ERROR", "absent list creation failed", err.Error())
		return errors.New("absent list creation failed")
	}

	absentList := generateAbsentList(NPMs, uint(formID))

	db.DB.Create(&absentList)

	return nil
}

func getFormDetail(formID int) (models.FormAbsensi, error) {
	formAbsent := models.FormAbsensi{}

	res := db.DB.Model(formAbsent).Where("id = ?", formID).Find(&formAbsent)

	if res.RowsAffected == 0 {
		util.LogErr("WARN", fmt.Sprintf("absent form for ID: %d is not found", formID), res.Error.Error())
		return formAbsent, fmt.Errorf("form with ID: %d is not found", formID)
	}

	return formAbsent, nil
}

func updateFormParticipant(formID int, newParticipant int) {
	absentForm := models.FormAbsensi{}

	db.DB.Model(&absentForm).Where("id = ?", formID).First(&absentForm)

	absentForm.Participant = newParticipant

	db.DB.Save(&absentForm)
}

func isParticipantChangeable(absentID int) error {
	absentLists := []models.AbsentList{}

	res := db.DB.Model(&models.AbsentList{}).
		Where(&models.AbsentList{
			FormAbsensiID: uint(absentID),
		}).Find(&absentLists)

	if res.Error != nil {
		util.LogErr("ERROR", "failed to change participant of an absent form", res.Error.Error())
		return errors.New("failed to change participant of an absent form")
	}

	for _, absentList := range absentLists {
		if absentList.Keterangan != "?" {
			util.LogErr("ERROR", fmt.Sprintf("participant of absent form with ID: %d can't be changed because some participants are already fill it", absentID), "")
			return fmt.Errorf("participant of absent form with ID: %d can't be changed because some participants are already fill it", absentID)
		}
	}

	return nil
}
