package controller

import (
	"errors"
	"fmt"
	"himatro-api/db"
	"himatro-api/models"
	"time"
)

func UpdateFormTitle(absentID int, title string) (string, error) {
	newAbsentForm := models.FormAbsensi{
		Title: title,
	}

	if err := updateAbsentFormDetail(newAbsentForm, absentID); err != nil {
		return "", err
	}

	return title, nil
}

func UpdateParticipant(formID int, newParticipant string) error {
	formDetail, err := getFormDetail(formID)

	if err != nil {
		return err
	}

	participantCode, err := validateParticipantCode(newParticipant)

	if err != nil {
		return fmt.Errorf("participant with code: %s is invalid", newParticipant)
	}

	if formDetail.Participant == participantCode {
		return nil
	}

	updateFormParticipant(formID, participantCode)
	deleteOldAbsentList(formID)

	if err := createNewAbsentList(formID, participantCode); err != nil {
		return err
	}

	return nil
}

func UpdateAbsentFormStartAt(formID int, startAtDate, startAtTime string) (string, error) {
	newStartAt, err := parseDate(startAtDate, startAtTime)

	if err != nil {
		return "", fmt.Errorf("invalid date time string received")
	}

	formDetail, err := getFormDetail(formID)

	if err != nil {
		return "", err
	}

	if newStartAt.After(formDetail.FinishAt) {
		return "", errors.New("form absent can't start after it's end date")
	}

	if newStartAt == formDetail.StartAt {
		return "", nil
	}

	if newStartAt.String() == formDetail.FinishAt.String() {
		return "", errors.New("form absent cant't start and end in the same time")
	}

	newAbsentForm := models.FormAbsensi{
		StartAt: newStartAt,
	}

	if err := updateAbsentFormDetail(newAbsentForm, formID); err != nil {
		return "", errors.New("server failure to update form details")
	}

	return newStartAt.String(), nil
}

func UpdateAbsentFormFinishAt(formID int, finishAtDate, finishAtTime string) (string, error) {
	newFinishAt, err := parseDate(finishAtDate, finishAtTime)

	if err != nil {
		return "", fmt.Errorf("invalid date time string received")
	}

	formDetail, err := getFormDetail(formID)

	if err != nil {
		return "", err
	}

	if newFinishAt.Before(formDetail.StartAt) {
		return "", errors.New("form absent can't end before it's start date")
	}

	if newFinishAt.Before(time.Now()) {
		return "", errors.New("form absent can't end before now")
	}

	if newFinishAt == formDetail.FinishAt {
		return "", nil
	}

	if newFinishAt.String() == formDetail.StartAt.String() {
		return "", errors.New("form absent cant't start and end in the same time")
	}

	newAbsentForm := models.FormAbsensi{
		FinishAt: newFinishAt,
	}

	if err := updateAbsentFormDetail(newAbsentForm, formID); err != nil {
		return "", errors.New("server failure to update form details")
	}

	return newFinishAt.String(), nil
}

func UpdateAbsentFormExecuseImageProof(formID int, proof bool) error {
	absentForm := models.FormAbsensi{}

	err := db.DB.Model(&absentForm).Where("id = ?", formID).First(&absentForm)

	if err.Error != nil {
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
		return fmt.Errorf("form with ID: %d is not exists", formID)
	}

	absentForm.RequireAttendanceImageProof = proof

	db.DB.Save(&absentForm)

	return nil
}

func updateAbsentFormDetail(absentForm models.FormAbsensi, absentID int) error {
	res := db.DB.Model(&models.FormAbsensi{}).Where("id = ?", absentID).Updates(&absentForm)

	if res.RowsAffected == 0 {
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
