package dto

type GetEnglishItemResponse struct {
	EnglishItems []*EnglishItem `json:"english_items"`
}

type EnglishItem struct {
	Id            string     `json:"id"`
	Content       string     `json:"content"`
	Translations  []string   `json:"translations"`
	EnExplanation string     `json:"en_explanation"`
	Examples      []*Example `json:"examples"`
	Imgs          []*Img     `json:"imgs"`
}
