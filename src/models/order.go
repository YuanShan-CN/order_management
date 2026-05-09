package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID   uint    `json:"user_id" gorm:"index;not null"`
	Date     string  `json:"date" gorm:"type:date;not null" validate:"required,datetime=2006-01-02"`
	Location string  `json:"location" gorm:"type:varchar(200)" validate:"required"`
	Content  string  `json:"content" gorm:"type:text" validate:"required"`
	Income   float64 `json:"income" gorm:"type:decimal(10,2)" validate:"numeric"`
	Deposit  float64 `json:"deposit" gorm:"type:decimal(10,2);default:0.00" validate:"numeric"`
	Balance  float64 `json:"balance" gorm:"type:decimal(10,2);default:0.00" validate:"numeric"`
	Settled  bool    `json:"settled" gorm:"default:false"`
	Shop     string  `json:"shop" gorm:"type:varchar(100)" validate:"required"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) CalculateIncome() {
	if o.Settled {
		o.Income = o.Deposit + o.Balance
	} else {
		o.Income = o.Deposit
	}
}
