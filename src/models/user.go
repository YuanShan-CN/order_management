package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(50);not null;unique" validate:"required,min=3,max=50"`
	Password string `json:"-" gorm:"type:varchar(255);not null" validate:"required,min=6"`
	Timezone string `json:"timezone" gorm:"type:varchar(50);not null;default:'Asia/Shanghai'" validate:"timezone"`
}

func (User) TableName() string {
	return "users"
}