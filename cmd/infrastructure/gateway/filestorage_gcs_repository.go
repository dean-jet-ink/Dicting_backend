package gateway

import (
	"english/algo"
	"english/cmd/domain/repository"
	"english/config"
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

func (r *FileStorageGCSRepository) Upload(file *multipart.FileHeader, preURL string) (string, error) {
	return "", nil
}

func (r *FileStorageGCSRepository) UploadFromURLs(urls []string, preURLs []string) ([]string, error) {
	wg := &sync.WaitGroup{}

	imgFileChan := make(chan *ImgFile, len(urls))
	errChan := make(chan error)

	for _, url := range urls {
		wg.Add(1)
		go r.fetchFileFromURL(wg, url, imgFileChan, errChan)
	}

	go func() {
		wg.Wait()
		close(imgFileChan)
	}()

	imgFiles := []*ImgFile{}

outer:
	for {
		select {
		case imgFile, ok := <-imgFileChan:
			if ok {
				imgFiles = append(imgFiles, imgFile)
			} else {
				break outer
			}
		case err := <-errChan:
			return nil, err
		}
	}

	newURLs := []string{}
	for _, imgFile := range imgFiles {
		if err := r.uploadFile(imgFile); err != nil {
			return nil, err
		}

		newURLs = append(newURLs, imgFile.URL)
	}

	if err := r.deletePreFiles(preURLs); err != nil {
		return nil, err
	}

	return newURLs, nil
}

func (r *FileStorageGCSRepository) fetchFileFromURL(wg *sync.WaitGroup, url string, imgFileChan chan *ImgFile, errChan chan error) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("failed to get image file: %v", resp.StatusCode)
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
	imgFileChan <- imgFile
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

func (r *FileStorageGCSRepository) deletePreFiles(preURLs []string) error {
	if config.GoEnv() == "dev" {
		for _, url := range preURLs {
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
