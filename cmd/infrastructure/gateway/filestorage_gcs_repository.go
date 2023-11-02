package gateway

import (
	"english/algo"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/config"
	"english/myerror"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
)

type FileStorageGCSRepository struct {
	// GCSクライアント
}

type ImgFile struct {
	Body     io.ReadCloser
	FileName string
	URL      string
}

func NewFileStorageGCSRepository() repository.FileStorageRepository {
	return &FileStorageGCSRepository{}
}

func (r *FileStorageGCSRepository) Upload(file *multipart.FileHeader, preImg *model.Img) error {
	return nil
}

func (r *FileStorageGCSRepository) UploadImgs(imgs []*model.Img, preImgs []*model.Img) error {
	wg := &sync.WaitGroup{}

	errChan := make(chan error)

	for _, img := range imgs {
		wg.Add(1)
		go r.fetchFileFromURL(wg, img, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return err
	}

	var urls []string

	for _, img := range preImgs {
		urls = append(urls, img.URL())
	}

	if err := r.deleteImgs(urls); err != nil {
		return err
	}

	return nil
}

func (r *FileStorageGCSRepository) DeleteImgs(imgs []*model.Img) error {
	urls := make([]string, 0)

	for _, img := range imgs {
		urls = append(urls, img.URL())
	}

	if err := r.deleteImgs(urls); err != nil {
		return err
	}

	return nil
}

func (r *FileStorageGCSRepository) fetchFileFromURL(wg *sync.WaitGroup, img *model.Img, errChan chan error) {
	defer wg.Done()

	resp, err := http.Get(img.URL())
	if err != nil {
		errChan <- err
		return
	}

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("%v: %w", myerror.ErrImgNotFound, errors.New("failed to get image file from URL"))
		return
	}

	contentType := resp.Header.Get("Content-Type")
	ext, err := getExtensionByContentType(contentType)
	if err != nil {
		errChan <- err
		return
	}
	ulid, err := algo.GenerateULID()
	if err != nil {
		errChan <- err
		return
	}

	fileName := fmt.Sprintf("%v%v", ulid, ext)

	imgFile := &ImgFile{
		Body:     resp.Body,
		FileName: fileName,
	}

	if err := r.uploadFile(imgFile); err != nil {
		errChan <- err
	}

	img.SetURL(imgFile.URL)
}

func (r *FileStorageGCSRepository) uploadFile(file *ImgFile) error {
	defer file.Body.Close()

	if config.GoEnv() == "dev" {
		out, err := os.Create(fmt.Sprintf("./static/img/english/%v", file.FileName))
		if err != nil {
			return err
		}

		_, err = io.Copy(out, file.Body)
		if err != nil {
			return err
		}

		file.URL = fmt.Sprintf("%v/img/english/%v", config.FilePath(), file.FileName)
	} else {
		// GCS
	}

	return nil
}

func getExtensionByContentType(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func (r *FileStorageGCSRepository) deleteImgs(urls []string) error {
	if config.GoEnv() == "dev" {
		for _, url := range urls {
			params := strings.Split(url, "/")
			preFileName := params[len(params)-1]
			preFilePath := fmt.Sprintf("./static/img/english/%v", preFileName)

			if err := os.Remove(preFilePath); err != nil {
				return err
			}
		}
	} else {
		// GCS
	}

	return nil
}
