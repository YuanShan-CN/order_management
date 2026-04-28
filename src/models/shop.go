package models

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
	UserID uint   `json:"user_id" gorm:"index;not null"`
	Name   string `json:"name" gorm:"type:varchar(100);not null" validate:"required"`
}

func (Shop) TableName() string {
	return "shops"
}
