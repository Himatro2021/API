package controller

import (
	"errors"
	"fmt"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"himatro-api/internal/util"
)

func GetAbsentFormsDetails(limit int) ([]models.ReturnedFormAbsentDetails, error) {
	absentFormsDetails := []models.ReturnedFormAbsentDetails{}

	res := db.DB.Model(&models.FormAbsensi{}).
		Select(`
			form_absensis.id as form_id,
			form_absensis.title,
			form_absensis.created_at,
			form_absensis.updated_at,
			form_absensis.participant as participant_code,
			form_absensis.start_at,
			form_absensis.finish_at,
			form_absensis.require_attendance_image_proof,
			form_absensis.require_execuse_image_proof,
			count(absent_lists.id) as total_participant,
			count(absent_lists.keterangan) filter(where absent_lists.keterangan = 'h') as hadir,
			count(absent_lists.keterangan) filter(where absent_lists.keterangan = 'i') as izin,
			count(absent_lists.keterangan) filter (where absent_lists.keterangan = '?') as tanpa_keterangan
		`).
		Limit(limit).
		Joins("inner join absent_lists on absent_lists.form_absensi_id = form_absensis.id").
		Group("form_id").
		Scan(&absentFormsDetails)

	if res.Error != nil {
		util.LogErr("ERROR", "failed to query absent forms details", res.Error.Error())
		return []models.ReturnedFormAbsentDetails{}, errors.New("failed to query absent forms details")
	}

	return absentFormsDetails, nil

}

func GetAbsentListResult(absentID int) ([]models.ReturnedAbsentList, error) {
	if err := isFormAbsentExists(absentID); err != nil {
		util.LogErr("WARN", fmt.Sprintf("Form absent not found ID: %d", absentID), err.Error())
		return []models.ReturnedAbsentList{}, err
	}

	absentList, err := getAbsentListFromFormID(absentID)

	if err != nil {
		util.LogErr("ERROR", fmt.Sprintf("failed to get absent result list ID: %d", absentID), err.Error())
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
		util.LogErr("WARN", fmt.Sprintf("absent form with ID: %d is not exists", absentID), res.Error.Error())
		return fmt.Errorf("absent form with ID: %d is not exists", absentID)
	}

	return nil
}

func getAbsentListFromFormID(absentID int) ([]models.ReturnedAbsentList, error) {
	absentLists := []models.ReturnedAbsentList{}

	res := db.DB.Model(&models.AbsentList{}).Select(`
			anggota_biasas.nama,
			absent_lists.npm,
			absent_lists.updated_at,
			absent_lists.keterangan,
			departemens.nama as nama_departemen
		`).
		Where(&models.AbsentList{FormAbsensiID: uint(absentID)}).
		Joins(`
			inner join anggota_biasas on anggota_biasas.npm = absent_lists.npm
		`).
		Joins(`
			inner join pengurus on pengurus.npm = anggota_biasas.npm
		`).
		Joins(`
			inner join departemens on departemens.id = pengurus.departemen_id
		`).
		Find(&absentLists)

	if res.Error != nil {
		util.LogErr("ERROR", fmt.Sprintf("server failed to fetch requested absent list ID: %d", absentID), res.Error.Error())
		return absentLists, errors.New("server failed to fetch requested absent list")
	}

	return absentLists, nil
}
