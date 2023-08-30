package entity

import "time"

type EnglishItemEntity struct {
	Id             string           `json:"id" gorm:"type:varchar(255);uniqueIndex:idx_id_content;primaryKey"`
	Content        string           `json:"content" gorm:"type:varchar(255);uniqueIndex:idx_id_content;not null"`
	JaTranslations string           `json:"ja_translations" gorm:"type:varchar(1000);not null"`
	EnExplanation  string           `json:"en_explanation" gorm:"type:varchar(1000);not null"`
	UserId         string           `json:"user_id" gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Examples       []*ExampleEntity `gorm:"foreignKey:EnglishItemId"`
	Imgs           []*ImgEntity     `gorm:"foreignKey:EnglishItemId"`
}

func (e *EnglishItemEntity) TableName() string {
	return "english_items"
}
