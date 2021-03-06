package controller

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"himatro-api/internal/config"
	"himatro-api/internal/contract"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"himatro-api/internal/util"
)

type InitAbsentData struct {
	Title                       string    `json:"title" validate:"required"`
	Participant                 int       `json:"participant" validate:"required"`
	StartAt                     time.Time `json:"startAt" validate:"required"`
	FinishAt                    time.Time `json:"finishAt" validate:"required"`
	RequireAttendanceImageProof bool      `json:"requireAttendanceImageProof" validate:"required"`
	RequireExecuseImageProof    bool      `json:"requireExecuseImageProof" validate:"required"`
}

func ExtractInitAbsentPayload(payload contract.CreateAbsentForm) (InitAbsentData, error) {
	start, err := parseDate(payload.StartAtDate, payload.StartAtTime)

	if err != nil {
		util.LogErr("WARN", "field start time and date is invalid", err.Error())
		return InitAbsentData{}, errors.New("field start time and date is invalid")
	}

	end, err := parseDate(payload.FinishAtDate, payload.FinishAtTime)

	if err != nil {
		util.LogErr("WARN", "field finish time and date is invalid", err.Error())
		return InitAbsentData{}, errors.New("field finish time and date is invalid")
	}

	if start.After(end) {
		util.LogErr("WARN", "start date must happen before finish date", err.Error())
		return InitAbsentData{}, errors.New("start date must happen before finish date")
	}

	if time.Now().After(end) {
		util.LogErr("WARN", "absent form can't finish before current date", err.Error())
		return InitAbsentData{}, errors.New("absent form can't finish before current date")
	}

	if start.String() == end.String() {
		util.LogErr("WARN", "form absent cant't start and end in the same time", err.Error())
		return InitAbsentData{}, errors.New("form absent cant't start and end in the same time")
	}

	participantCode, err := validateParticipantCode(payload.Participant)

	if err != nil {
		util.LogErr("WARN", fmt.Sprintf("Invalid participant code used: %s", payload.Participant), err.Error())
		return InitAbsentData{}, err
	}

	initAbsentData := InitAbsentData{
		Title:                       payload.Title,
		Participant:                 participantCode,
		StartAt:                     start,
		FinishAt:                    end,
		RequireAttendanceImageProof: payload.RequireAttendanceImageProof,
		RequireExecuseImageProof:    payload.RequireExecuseImageProof,
	}

	return initAbsentData, nil
}

func RegisterNewAbsentForm(detail *InitAbsentData) (uint, error) {
	newAbsent := models.FormAbsensi{
		Title:                       detail.Title,
		Participant:                 detail.Participant,
		StartAt:                     detail.StartAt,
		FinishAt:                    detail.FinishAt,
		RequireAttendanceImageProof: detail.RequireAttendanceImageProof,
		RequireExecuseImageProof:    detail.RequireExecuseImageProof,
	}

	err := db.DB.Create(&newAbsent)

	if err.Error != nil {
		util.LogErr("ERROR", fmt.Sprintf("system failed to register new absent for %s", newAbsent.Title), err.Error.Error())
		return 0, errors.New("system failed to register new absent")
	}

	return newAbsent.ID, nil
}

func InitAbsentList(detail *InitAbsentData, absentID uint) error {
	listPengurus, err := getAllNPMFromDepartemenID(detail.Participant)

	if err != nil {
		util.LogErr("ERROR", fmt.Sprintf("failed to generate absent lists for: %d", absentID), err.Error())
		return errors.New("failed to generate absent lists")
	}

	absentLists := generateAbsentList(listPengurus, absentID)

	errs := db.DB.Create(&absentLists)

	if errs.Error != nil {
		util.LogErr("ERROR", fmt.Sprintf("failed to generate absent lists for: %d", absentID), err.Error())
		return errors.New("failed to generate absent lists")
	}

	return nil
}

func parseDate(startAtDate string, startAtTime string) (time.Time, error) {
	year, month, day, err := util.ExtractDateString(startAtDate)

	if err != nil {
		util.LogErr("WARN", fmt.Sprintf("invalid start date string received: %s", startAtDate), err.Error())
		return time.Now(), err
	}

	hour, minute, sec, err := util.ExtractTimeString(startAtTime)

	if err != nil {
		util.LogErr("WARN", fmt.Sprintf("invalid start time string received: %s", startAtTime), err.Error())
		return time.Now(), err
	}

	date := time.Date(year, time.Month(month), day, hour, minute, sec, 0, time.Local)

	tz, _ := time.LoadLocation(config.TimeZone())

	return date.In(tz), nil
}

func validateParticipantCode(participant string) (int, error) {
	switch strings.ToUpper(participant) {
	case "PH":
		return 1, nil
	case "PPD":
		return 2, nil
	case "KPO":
		return 3, nil
	case "KOMINFO":
		return 4, nil
	case "KWU":
		return 5, nil
	case "BANGTEK":
		return 6, nil
	case "ALL":
		return 0, nil // create absent for all
	default:
		util.LogErr("WARN", fmt.Sprintf("participant (departemenID) is invalid: %s", participant), "")
		return 0, errors.New("participant (departemenID) is invalid")
	}
}

func getAllNPMFromDepartemenID(departemenID int) ([]models.Pengurus, error) {
	pengurus := []models.Pengurus{}

	err := db.DB.Where(&models.Pengurus{DepartemenID: departemenID}).Find(&pengurus) // if departemenID 0, query all see: https://gorm.io/docs/query.html#Struct-amp-Map-Conditions

	if err.Error != nil {
		util.LogErr("ERROR", fmt.Sprintf("failed to instaniate new absent list departemenID: %d", departemenID), err.Error.Error())
		return pengurus, errors.New("failed to instaniate new absent list")
	}

	return pengurus, nil
}

func generateAbsentList(listPengurus []models.Pengurus, absentID uint) []models.AbsentList {
	absentLists := []models.AbsentList{}

	for _, pengurus := range listPengurus {
		absentList := models.AbsentList{
			FormAbsensiID: absentID,
			NPM:           pengurus.NPM,
		}

		absentLists = append(absentLists, absentList)
	}

	return absentLists
}
