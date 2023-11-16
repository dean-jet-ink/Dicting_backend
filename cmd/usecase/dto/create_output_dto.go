package dto

import "time"

type Output struct {
	Id        string
	Index     uint   `json:"index"`
	Question  string `json:"question" validate:"required"`
	Answer    string `json:"answer" validate:"required"`
	Advice    string `json:"advice" validate:"required"`
	CreatedAt time.Time
}

type CreateOutputInput struct {
	EnglishItemId string    `json:"english_item_id" validate:"required"`
	Outputs       []*Output `json:"outputs" validate:"dive"`
}
