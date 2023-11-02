package dto

type UpdateUserRequest struct {
	Id    string
	Email string `json:"email" validate:"email"`
	Name  string `json:"name" validate:"gte=1,lt=30"`
	Image string
}

type UpdateUserResponse struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}
