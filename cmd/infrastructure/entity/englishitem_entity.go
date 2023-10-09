package entity

import "time"

type EnglishItemEntity struct {
	Id            string             `json:"id" gorm:"type:varchar(255);primaryKey"`
	Content       string             `json:"content" gorm:"type:varchar(255);uniqueIndex:idx_id_content;not null"`
	Translations  string             `json:"translations" gorm:"type:varchar(1000);not null"`
	EnExplanation string             `json:"en_explanation" gorm:"type:varchar(1000);not null"`
	UserId        string             `json:"user_id" gorm:"type:varchar(255);uniqueIndex:idx_id_content;not null"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Examples      []*ExampleEntity   `gorm:"foreignKey:EnglishItemId; constraint:OnDelete:CASCADE"`
	Imgs          []*ImgEntity       `gorm:"foreignKey:EnglishItemId; constraint:OnDelete:CASCADE"`
	Proficiency   *ProficiencyEntity `gorm:"foreignKey:Id; constraint:OnDelete:CASCADE"`
}

func (e *EnglishItemEntity) TableName() string {
	return "english_items"
}
