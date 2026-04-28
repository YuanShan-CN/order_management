package models

import "time"

type Shop struct {
	ID        uint       `json:"ID" gorm:"primary_key"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null;unique"`
	DeletedAt *time.Time `json:"-" gorm:"default:null"`
}

func (Shop) TableName() string {
	return "shops"
}