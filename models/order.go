package models

import "time"

type Order struct {
	ID        uint      `json:"ID" gorm:"primary_key"`
	Date      string    `json:"date" gorm:"type:date;not null"`
	Location  string    `json:"location" gorm:"type:varchar(200)"`
	Content   string    `json:"content" gorm:"type:text"`
	Fee       float64   `json:"fee" gorm:"type:decimal(10,2)"`
	Settled   bool      `json:"settled" gorm:"default:false"`
	Edited    bool      `json:"edited" gorm:"default:false"`
	Shop      string    `json:"shop" gorm:"type:varchar(100)"`
	DeletedAt time.Time `json:"deleted_at" gorm:"default:null"`
}

func (Order) TableName() string {
	return "orders"
}