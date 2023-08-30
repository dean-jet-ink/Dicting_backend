package repository

import "mime/multipart"

type FileStorageRepository interface {
	Upload(file *multipart.FileHeader, preURL string) (string, error)
	UploadFromURLs(urls []string, preURLs []string) ([]string, error)
}
