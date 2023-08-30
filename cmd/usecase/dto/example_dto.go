package dto

type Example struct {
	Id          string `json:"id"`
	Example     string `json:"example" validate:"required"`
	Translation string `json:"translation" validate:"required"`
}
