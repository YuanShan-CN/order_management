package models

import "time"

type Order struct {
	ID        uint       `json:"ID" gorm:"primary_key"`
	Date      string     `json:"date" gorm:"type:date;not null" validate:"required,datetime=2006-01-02"`
	Location  string     `json:"location" gorm:"type:varchar(200)" validate:"required"`
	Content   string     `json:"content" gorm:"type:text" validate:"required"`
	Fee       float64    `json:"fee" gorm:"type:decimal(10,2)" validate:"required,numeric"`
	Settled   bool       `json:"settled" gorm:"default:false"`
	Shop      string     `json:"shop" gorm:"type:varchar(100)" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-" gorm:"default:null"`
}

func (Order) TableName() string {
	return "orders"
}
