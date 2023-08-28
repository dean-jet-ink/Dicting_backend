package main

import (
	"english/algo"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func main() {
	resp, err := http.Get("https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpD8COglFbxKT_nMlWCdoC7Mj52S4jMAoesTEZPLJjqyNfYFU&s")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	fmt.Println(contentType)

	ext, err := getExtensionByContentType(contentType)
	if err != nil {
		fmt.Println(err)
		return
	}

	ulid, _ := algo.GenerateULID()
	fileName := fmt.Sprintf("%v.%v", ulid, ext)
	out, err := os.Create(fmt.Sprintf("./static/img/english/%v", fileName))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}
