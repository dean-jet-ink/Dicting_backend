package entity

import "time"

type ExampleEntity struct {
	Id            string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	Example       string    `json:"example" gorm:"type:varchar(1000);not null"`
	Translation   string    `json:"translation" gorm:"type:varchar(1000);not null"`
	EnglishItemId string    `json:"english_item_id" gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (e *ExampleEntity) TableName() string {
	return "examples"
}
