package entity

import "time"

type OutputEntity struct {
	Id            string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	Question      string    `json:"question" gorm:"type:varchar(255);not null"`
	Answer        string    `json:"answer" gorm:"type:varchar(255);not null"`
	Advice        string    `json:"advice" gorm:"type:varchar(1000); not null"`
	Index         uint      `json:"index" gorm:"not null"`
	EnglishItemId string    `json:"english_item_id" gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time `json:"created_at"`
}

func (e *OutputEntity) TableName() string {
	return "outputs"
}
