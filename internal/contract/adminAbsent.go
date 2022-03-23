package contract

type CreateAbsentForm struct {
	Title                       string `json:"title" validate:"required"`
	StartAtDate                 string `json:"startAtDate" validate:"required"`
	StartAtTime                 string `json:"startAtTime" validate:"required"`
	FinishAtDate                string `json:"finishAtDate" validate:"required"`
	FinishAtTime                string `json:"finishAtTime" validate:"required"`
	RequireAttendanceImageProof bool   `json:"requireAttendanceImageProof" validate:"required,eq=true|eq=false"`
	RequireExecuseImageProof    bool   `json:"requireExecuseImageProof" validate:"required,eq=true|eq=false"`
	Participant                 string `json:"participant" validate:"required"`
}
