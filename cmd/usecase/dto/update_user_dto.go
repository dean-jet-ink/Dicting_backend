package dto

type UpdateUserRequest struct {
	Id            string
	Email         string `json:"email" validate:"email"`
	Name          string `json:"name" validate:"gte=1,lt=30"`
	ProfileImgURL string
}

type UpdateUserResponse struct {
	Email         string `json:"email,omitempty"`
	Name          string `json:"name,omitempty"`
	ProfileImgURL string `json:"profile_img_url,omitempty"`
}
