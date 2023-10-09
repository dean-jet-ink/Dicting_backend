package entity

import "time"

type ProficiencyEntity struct {
	Id          string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	Proficiency string    `json:"proficiency" gorm:"type:enum('Learning', 'Understand', 'Mastered');default:'Learning';not null"`
	Exp         uint      `json:"exp" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e *ProficiencyEntity) TableName() string {
	return "proficiencies"
}
