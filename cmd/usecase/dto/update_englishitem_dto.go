package dto

type UpdateEnglishItemRequest struct {
	Id            string     `json:"id" validate:"required"`
	Content       string     `json:"content" validate:"required,lt=255"`
	Translations  []string   `json:"translations" validate:"required"`
	EnExplanation string     `json:"en_explanation" validate:"required"`
	Imgs          []*Img     `json:"imgs" validate:"dive"`
	Examples      []*Example `json:"examples" validate:"dive"`
	Proficiency   string     `json:"proficiency" validate:"required"`
	Exp           uint       `json:"exp"`
	UserId        string
}
