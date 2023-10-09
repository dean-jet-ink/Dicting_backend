package dto

type Img struct {
	Id          string `json:"id"`
	URL         string `json:"url" validate:"http_url"`
	IsThumbnail bool   `json:"is_thumbnail"`
}
