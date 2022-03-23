package contract

type FillAbsentList struct {
	NPM        string `json:"NPM" validate:"required"`
	Keterangan string `json:"keterangan" validate:"eq=h|eq=i"`
}

type UpdateKeteranganAbsent struct {
	Keterangan string `json:"keterangan" validate:"eq=h|eq=i"`
}
