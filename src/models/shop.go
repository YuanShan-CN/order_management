package models

import "time"

type Shop struct {
	ID        uint       `json:"ID" gorm:"primary_key"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null;unique" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-" gorm:"default:null"`
}

func (Shop) TableName() string {
	return "shops"
}
