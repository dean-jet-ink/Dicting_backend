package dto

import "time"

type GetOutputsInput struct {
	EnglishItemId string    `form:"english_item_id" validate:"required"`
	CreatedAt     time.Time `form:"created_at" validate:"required"`
}

type GetOutputsOutput struct {
	Outputs []*Output `json:"outputs"`
}
