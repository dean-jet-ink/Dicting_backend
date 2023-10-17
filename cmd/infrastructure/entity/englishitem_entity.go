package entity

import "time"

type EnglishItemEntity struct {
	Id            string           `json:"id" gorm:"type:varchar(255);primaryKey"`
	Content       string           `json:"content" gorm:"type:varchar(255);uniqueIndex:idx_id_content;not null"`
	Translations  string           `json:"translations" gorm:"type:varchar(1000);not null"`
	EnExplanation string           `json:"en_explanation" gorm:"type:varchar(1000);not null"`
	UserId        string           `json:"user_id" gorm:"type:varchar(255);uniqueIndex:idx_id_content;not null"`
	Proficiency   string           `json:"proficiency" gorm:"type:enum('Learning', 'Understand', 'Mastered');default:'Learning';not null"`
	Exp           uint             `json:"exp" gorm:"default:0"`
	Examples      []*ExampleEntity `gorm:"foreignKey:EnglishItemId; constraint:OnDelete:CASCADE"`
	Imgs          []*ImgEntity     `gorm:"foreignKey:EnglishItemId; constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (e *EnglishItemEntity) TableName() string {
	return "english_items"
}
