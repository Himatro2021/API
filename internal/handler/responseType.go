package handler

import (
	"himatro-api/internal/models"
	"time"
)

type AbsentListSuccessMessage struct {
	OK     bool                  `json:"ok"`
	Status int                   `json:"status"`
	Result []models.AnggotaBiasa `json:"result"`
}

type ErrorMessage struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type JSONPayloadValidationError struct {
	OK      bool     `json:"ok"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type LoginTokenResp struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}

type SuccessCreateAbsent struct {
	OK                          bool      `json:"ok"`
	AbsentID                    uint      `json:"absentID"`
	Title                       string    `json:"title"`
	Participant                 int       `json:"participant"`
	StartAt                     time.Time `json:"startAt"`
	FinishAt                    time.Time `json:"finishAt"`
	RequireAttendanceImageProof bool      `json:"requireAttendanceImageProof"`
	RequireExecuseImageProof    bool      `json:"requireExecuseImageProof"`
}

type SuccessListAbsent struct {
	OK         bool                        `json:"ok"`
	FormID     int                         `json:"formID"`
	Total      int                         `json:"total"`
	TotalHadir int                         `json:"total_hadir"`
	TotalIzin  int                         `json:"total_izin"`
	TotalAlpha int                         `json:"total_tanpa_keterangan"`
	List       []models.ReturnedAbsentList `json:"list"`
}

type SuccessUpdateForm struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

type SuccessFormAbsentDetails struct {
	OK      bool                               `json:"ok"`
	Message string                             `json:"message"`
	Total   int                                `json:"total"`
	List    []models.ReturnedFormAbsentDetails `json:"list"`
}
