package gateway

import (
	"context"
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/client"
	"english/config"
	"english/lib"
	"english/myerror"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type FileStorageGCSRepository struct {
	gcsClient *client.GCSClient
}

func NewFileStorageGCSRepository(gcsClient *client.GCSClient) repository.FileStorageRepository {
	return &FileStorageGCSRepository{
		gcsClient: gcsClient,
	}
}

func (r *FileStorageGCSRepository) Upload(file *model.ImgFile, preImg string) error {
	ctx := context.Background()

	r.uploadFile(ctx, file, false)

	r.deleteImg(ctx, preImg)

	return nil
}

func (r *FileStorageGCSRepository) UploadImgs(imgs []*model.Img, preImgs []*model.Img) error {
	wg := &sync.WaitGroup{}

	errChan := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// GCSのドメインを持つURL、またはローカルドメインを持つURLの場合、既にアップロード済みの画像であるため除外
	filteredImgs := []*model.Img{}
	gcsPath := fmt.Sprintf("https://storage.cloud.google.com/%v", client.BUCKET_NAME)
	for _, img := range imgs {
		if !strings.Contains(img.URL(), gcsPath) &&
			!strings.Contains(img.URL(), config.StaticFilePath()) {
			filteredImgs = append(filteredImgs, img)
		}
	}

	// imgsのURLと一致するpreImgについては削除しないため、除外
	imgSet := map[string]bool{}
	for _, img := range imgs {
		imgSet[img.URL()] = true
	}

	preURLs := []string{}
	for _, preImg := range preImgs {
		url := preImg.URL()
		if !imgSet[url] {
			preURLs = append(preURLs, url)
		}
	}

	for _, img := range filteredImgs {
		wg.Add(1)
		go r.uploadFileFromURL(ctx, wg, img, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		cancel()
		return err
	}

	if err := r.deleteImgs(ctx, preURLs); err != nil {
		cancel()
		return err
	}

	return nil
}

func (r *FileStorageGCSRepository) DeleteImgs(imgs []*model.Img) error {
	urls := make([]string, 0)

	for _, img := range imgs {
		urls = append(urls, img.URL())
	}

	if err := r.deleteImgs(context.Background(), urls); err != nil {
		return err
	}

	return nil
}

func (r *FileStorageGCSRepository) uploadFileFromURL(ctx context.Context, wg *sync.WaitGroup, img *model.Img, errChan chan error) {
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
	ulid, err := lib.GenerateULID()
	if err != nil {
		errChan <- err
		return
	}

	fileName := fmt.Sprintf("%v%v", ulid, ext)

	imgFile := &model.ImgFile{
		Body:     resp.Body,
		FileName: fileName,
	}

	if err := r.uploadFile(ctx, imgFile, true); err != nil {
		errChan <- err
	}

	img.SetURL(imgFile.URL)
}

func (r *FileStorageGCSRepository) uploadFile(ctx context.Context, file *model.ImgFile, isEnglish bool) error {
	defer file.Body.Close()

	var path string

	if isEnglish {
		path = "english"
	} else {
		path = "user"
	}

	if config.GoEnv() == "dev" {
		out, err := os.Create(fmt.Sprintf("./static/img/%v/%v", path, file.FileName))
		if err != nil {
			return err
		}

		if _, err = io.Copy(out, file.Body); err != nil {
			return err
		}

		file.URL = fmt.Sprintf("%v/img/%v/%v", config.StaticFilePath(), path, file.FileName)
	} else {
		if err := r.gcsClient.Save(ctx, file); err != nil {
			return err
		}
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

func (r *FileStorageGCSRepository) deleteImgs(ctx context.Context, urls []string) error {
	for _, url := range urls {
		r.deleteImg(ctx, url)
	}

	return nil
}

func (r *FileStorageGCSRepository) deleteImg(ctx context.Context, url string) error {
	if config.GoEnv() == "dev" {
		params := strings.Split(url, "/")
		preFileName := params[len(params)-1]
		preFilePath := fmt.Sprintf("./static/img/english/%v", preFileName)

		if err := os.Remove(preFilePath); err != nil {
			return err
		}
	} else {
		if err := r.gcsClient.Delete(ctx, url); err != nil {
			return err
		}
	}

	return nil
}
