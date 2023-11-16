package dto

import "time"

type DeleteOutputInput struct {
	EnglishItemId string    `form:"english_item_id" validate:"required"`
	CreatedAt     time.Time `form:"created_at" validate:"required"`
}
