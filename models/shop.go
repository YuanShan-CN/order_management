package models

type Shop struct {
	ID   uint   `json:"ID" gorm:"primary_key"`
	Name string `json:"name" gorm:"type:varchar(100);not null;unique"`
}

func (Shop) TableName() string {
	return "shops"
}