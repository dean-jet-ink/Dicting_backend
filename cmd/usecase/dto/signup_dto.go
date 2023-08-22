package dto

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=4,lt=30"`
	Name     string `json:"name" validate:"required"`
}
