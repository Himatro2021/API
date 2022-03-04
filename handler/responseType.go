package handler

import "himatro-api/models"

type AbsentListSuccessMessage struct {
	OK     bool                  `json:"ok"`
	Status int                   `json:"status"`
	Result []models.AnggotaBiasa `json:"result"`
}

type ErrorMessage struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type LoginTokenResp struct {
	OK    bool   `json:"ok"`
	Token string `json:"token"`
}
