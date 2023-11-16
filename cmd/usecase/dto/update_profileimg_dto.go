package dto

import "mime/multipart"

type UpdateProfileImgInput struct {
	Id         string
	FileHeader *multipart.FileHeader `validate:"required"`
}

type UpdateProfileImgOutput struct {
	ProfileImgURL string `json:"profile_img_url"`
}
