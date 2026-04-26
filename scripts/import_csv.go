package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint    `gorm:"primary_key"`
	Date      string  `gorm:"type:date;not null"`
	Location  string  `gorm:"type:varchar(200)"`
	Content   string  `gorm:"type:text"`
	Fee       float64 `gorm:"type:decimal(10,2)"`
	Settled   bool    `gorm:"default:false"`
	Shop      string  `gorm:"type:varchar(100)"`
}

func (Order) TableName() string {
	return "orders"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run import_csv.go <csv_file_path>")
		return
	}

	csvPath := os.Args[1]

	dsn := "root:123456@tcp(127.0.0.1:3307)/order_management?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect database: %v\n", err)
		return
	}

	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read CSV: %v\n", err)
		return
	}

	if len(records) < 2 {
		fmt.Println("No data in CSV file")
		return
	}

	var orders []Order
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 9 {
			continue
		}

		date := record[0]
		location := record[1]
		content := record[2]
		fee, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			continue
		}
		settled := record[4] == "是"
		shop := record[8]

		if date == "" || shop == "" {
			continue
		}

		orders = append(orders, Order{
			Date:     date,
			Location: location,
			Content:  content,
			Fee:      fee,
			Settled:  settled,
			Shop:     shop,
		})
	}

	if len(orders) == 0 {
		fmt.Println("No valid data to import")
		return
	}

	result := db.Create(&orders)
	if result.Error != nil {
		fmt.Printf("Failed to import data: %v\n", result.Error)
		return
	}

	fmt.Printf("Successfully imported %d orders\n", len(orders))
}
