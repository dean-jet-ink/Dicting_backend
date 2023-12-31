package entity

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	Id              string               `json:"id" gorm:"type:varchar(255);primaryKey"`
	Email           string               `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password        string               `json:"password" gorm:"type:varchar(255)"`
	Name            string               `json:"name" gorm:"type:varchar(255);not null"`
	ProfileImageURL string               `json:"profile_image_url" gorm:"type:varchar(2083)"`
	Iss             sql.NullString       `json:"iss" gorm:"type:varchar(255);uniqueIndex:idx_iss_sub"`
	Sub             sql.NullString       `json:"sub" gorm:"type:varchar(255);uniqueIndex:idx_iss_sub"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	EnglishItems    []*EnglishItemEntity `json:"english_items" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
}

func (u *UserEntity) TableName() string {
	return "users"
}
