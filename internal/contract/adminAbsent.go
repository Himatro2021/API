package contract

type CreateAbsentForm struct {
	Title                       string `json:"title" validate:"required"`
	StartAtDate                 string `json:"startAtDate" validate:"required"`
	StartAtTime                 string `json:"startAtTime" validate:"required"`
	FinishAtDate                string `json:"finishAtDate" validate:"required"`
	FinishAtTime                string `json:"finishAtTime" validate:"required"`
	RequireAttendanceImageProof bool   `json:"requireAttendanceImageProof,omitempty"`
	RequireExecuseImageProof    bool   `json:"requireExecuseImageProof,omitempty"`
	Participant                 string `json:"participant" validate:"required"`
}

type UpdateFormTitle struct {
	Title string `json:"title" validate:"required"`
}

type UpdateFormParticipant struct {
	Participant string `json:"participant" validate:"required"`
}

type UpdateFormImageProof struct {
	Status bool `json:"status"`
}

type UpdateFormTime struct {
	Date string `json:"date" validate:"required"`
	Time string `json:"time" validate:"required"`
}
