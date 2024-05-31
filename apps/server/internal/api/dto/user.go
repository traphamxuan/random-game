package dto

type LoginPayload struct {
	User     string `json:"user"`
	Password string `json:"pass"`
}
