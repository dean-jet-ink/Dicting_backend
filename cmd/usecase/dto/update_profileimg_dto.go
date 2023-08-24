package dto

import "mime/multipart"

type UpdateProfileImgRequest struct {
	Id   string
	File *multipart.FileHeader `validate:"required"`
}

type UpdateProfileImgResponse struct {
	ProfileImgURL string `json:"profile_img_url"`
}
