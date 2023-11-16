package repository

import (
	"english/cmd/domain/model"
)

type FileStorageRepository interface {
	Upload(file *model.ImgFile, preImg string) error
	UploadImgs(imgs []*model.Img, preImgs []*model.Img) error
	DeleteImgs(imgs []*model.Img) error
}
