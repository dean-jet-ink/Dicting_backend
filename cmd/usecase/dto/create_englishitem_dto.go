package dto

type CreateEnglishItemRequest struct {
	Content        string     `json:"content" validate:"required,lt=255"`
	JaTranslations []string   `json:"ja_translations,omitempty" validate:"required"`
	EnExplanation  string     `json:"en_explanation,omitempty" validate:"required"`
	ImgURLs        []string   `json:"img_urls" validate:"dive,http_url"`
	Examples       []*Example `json:"examples" validate:"dive,required"`
	UserId         string
}

type CreateEnglishItemResponse struct {
	Id             string     `json:"id"`
	Content        string     `json:"content"`
	JaTranslations []string   `json:"ja_translations,omitempty"`
	EnExplanation  string     `json:"en_explanation,omitempty"`
	Imgs           []*Img     `json:"imgs"`
	Examples       []*Example `json:"examples"`
}
