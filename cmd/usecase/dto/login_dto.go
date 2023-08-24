package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=4,lt=30"`
	Iss      string
	Sub      string
}
