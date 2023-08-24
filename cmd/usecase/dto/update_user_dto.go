package dto

type UpdateUserRequest struct {
	Id    string
	Email string `json:"email" validate:"email"`
	Name  string `json:"name" validate:"gte=1,lt=30"`
}

type UpdateUserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
