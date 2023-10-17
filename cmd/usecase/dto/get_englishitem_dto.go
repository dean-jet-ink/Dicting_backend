package dto

type GetEnglishItemsResponse struct {
	EnglishItems []*EnglishItem `json:"english_items"`
}

type EnglishItem struct {
	Id            string   `json:"id"`
	Content       string   `json:"content"`
	Translations  []string `json:"translations"`
	EnExplanation string   `json:"en_explanation"`
	Img           string   `json:"img"`
	Proficiency   string   `json:"proficiency"`
	Exp           uint     `json:"exp"`
}

type GetEnglishItemResponse struct {
	Id            string     `json:"id"`
	Content       string     `json:"content"`
	Translations  []string   `json:"translations"`
	EnExplanation string     `json:"en_explanation"`
	Examples      []*Example `json:"examples"`
	Imgs          []*Img     `json:"imgs"`
	Proficiency   string     `json:"proficiency"`
	Exp           uint       `json:"exp"`
}
