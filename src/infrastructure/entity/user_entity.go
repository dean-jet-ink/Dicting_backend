package entity

import "time"

type UserEntity struct {
	Id              int       `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"unique"`
	Password        string    `json:"password"`
	Name            string    `json:"name" gorm:"not null"`
	ProfileImageURL string    `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (ue *UserEntity) TableName() string {
	return "users"
}
