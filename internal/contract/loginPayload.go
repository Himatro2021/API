package contract

type LoginPayload struct {
	NPM      string `json:"NPM"`
	Password string `json:"password"`
}
