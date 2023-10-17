package repository

import (
	"english/cmd/domain/model"
	"mime/multipart"
)

type FileStorageRepository interface {
	Upload(file *multipart.FileHeader, preImg *model.Img) error
	UploadImgs(imgs []*model.Img, preImgs []*model.Img) error
}
