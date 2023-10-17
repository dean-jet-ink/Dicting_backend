package entity

import "time"

type ImgEntity struct {
	Id            string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	URL           string    `json:"url" gorm:"type:varchar(2083);not null"`
	IsThumbnail   bool      `json:"is_thumbnail" gorm:"default:false"`
	EnglishItemId string    `json:"english_item_id" gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (e *ImgEntity) TableName() string {
	return "imgs"
}
