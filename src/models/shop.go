package models

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(100);not null;unique" validate:"required"`
}

func (Shop) TableName() string {
	return "shops"
}
