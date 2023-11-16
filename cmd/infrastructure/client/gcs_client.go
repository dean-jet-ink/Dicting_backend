package client

import (
	"context"
	"english/cmd/domain/model"
	"english/config"
	"fmt"
	"io"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	CREDENTIALS_FILE = "/usr/src/app/keyFile.json"
	BUCKET_NAME      = "dicting_bucket"
)

type GCSClient struct {
	bucket *storage.BucketHandle
	path   string
}

func NewGCSClient() *GCSClient {
	credentialJSON := config.GCSServiceKey()

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(credentialJSON)))
	if err != nil {
		log.Fatal(err)
	}

	bucket := client.Bucket(BUCKET_NAME)

	return &GCSClient{
		bucket: bucket,
		path:   "https://storage.cloud.google.com/dicting_bucket",
	}
}

func (c *GCSClient) Save(ctx context.Context, file *model.ImgFile) error {
	obj := c.bucket.Object(file.FileName)

	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, file.Body); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	file.URL = fmt.Sprintf("%v/%v", c.path, file.FileName)

	return nil
}

func (c *GCSClient) Delete(ctx context.Context, path string) error {
	params := strings.Split(path, "/")

	objName := params[len(params)-1]

	obj := c.bucket.Object(objName)

	if err := obj.Delete(ctx); err != nil {
		return err
	}

	return nil
}
