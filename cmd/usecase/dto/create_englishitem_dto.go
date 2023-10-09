package dto

type CreateEnglishItemRequest struct {
	Content       string     `json:"content" validate:"required,lt=255"`
	Translations  []string   `json:"translations,omitempty" validate:"required"`
	EnExplanation string     `json:"en_explanation,omitempty" validate:"required"`
	Imgs          []*Img     `json:"imgs" validate:"dive"`
	Examples      []*Example `json:"examples" validate:"dive,required"`
	UserId        string
}
