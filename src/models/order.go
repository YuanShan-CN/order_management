package models

import (
	"encoding/csv"
	"io"
	"strconv"
	"gorm.io/gorm"
)

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

// CSVHeaderWithID 带ID的订单CSV表头
var CSVHeaderWithID = []string{"ID", "日期", "地点", "内容", "店铺", "定金", "尾款", "已结算", "收入"}

// CSVHeaderWithoutID 不带ID的订单CSV表头
var CSVHeaderWithoutID = []string{"日期", "地点", "内容", "店铺", "定金", "尾款", "已结算", "收入"}

// ToCSVRowWithID 将订单转换为带ID的CSV行
func (o *Order) ToCSVRowWithID() []string {
	settledStr := "否"
	if o.Settled {
		settledStr = "是"
	}
	return []string{
		strconv.Itoa(int(o.ID)),
		ExtractDatePart(o.Date),
		o.Location,
		o.Content,
		o.Shop,
		strconv.FormatFloat(o.Deposit, 'f', 2, 64),
		strconv.FormatFloat(o.Balance, 'f', 2, 64),
		settledStr,
		strconv.FormatFloat(o.Income, 'f', 2, 64),
	}
}

// ToCSVRowWithoutID 将订单转换为不带ID的CSV行
func (o *Order) ToCSVRowWithoutID() []string {
	settledStr := "否"
	if o.Settled {
		settledStr = "是"
	}
	return []string{
		ExtractDatePart(o.Date),
		o.Location,
		o.Content,
		o.Shop,
		strconv.FormatFloat(o.Deposit, 'f', 2, 64),
		strconv.FormatFloat(o.Balance, 'f', 2, 64),
		settledStr,
		strconv.FormatFloat(o.Income, 'f', 2, 64),
	}
}

// WriteOrdersToCSV 将订单列表写入CSV写入器
// includeID: 是否包含ID列
func WriteOrdersToCSV(writer io.Writer, orders []Order, includeID bool) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// 写入表头
	var headers []string
	if includeID {
		headers = CSVHeaderWithID
	} else {
		headers = CSVHeaderWithoutID
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// 写入数据行
	for _, order := range orders {
		var row []string
		if includeID {
			row = order.ToCSVRowWithID()
		} else {
			row = order.ToCSVRowWithoutID()
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExtractDatePart 从日期字符串中提取日期部分
// 处理格式：2006-01-02T00:00:00+08:00 -> 2006-01-02
func ExtractDatePart(date string) string {
	if len(date) > 10 && (date[10] == 'T' || date[10] == ' ') {
		return date[:10]
	}
	return date
}
