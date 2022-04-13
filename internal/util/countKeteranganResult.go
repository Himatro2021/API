package util

import "himatro-api/internal/models"

func CountKeteranganResult(list []models.ReturnedAbsentList) (int, int, int) {
	var hadir, izin, alpha int = 0, 0, 0

	for _, absent := range list {
		switch absent.Keterangan {
		case "h":
			hadir++
		case "i":
			izin++
		case "?":
			alpha++
		default:
			LogErr("ERROR", "Invalid keterangan on absent list", "")
		}
	}

	return hadir, izin, alpha
}
