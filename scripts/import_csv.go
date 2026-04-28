package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint      `gorm:"primary_key"`
	UserID    uint      `gorm:"not null"`
	Date      string    `gorm:"type:date;not null"`
	Location  string    `gorm:"type:varchar(200)"`
	Content   string    `gorm:"type:text"`
	Fee       float64   `gorm:"type:decimal(10,2)"`
	Settled   bool      `gorm:"default:false"`
	Shop      string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
}

func (Order) TableName() string {
	return "orders"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	var userID uint
	flag.UintVar(&userID, "user_id", 1, "User ID to assign to imported orders")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run import_csv.go -user_id <user_id> <csv_file_path>")
		return
	}

	csvPath := args[0]

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "order_management")

	if dbPassword == "" {
		fmt.Println("Error: DB_PASSWORD environment variable is not set")
		return
	}

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
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

		if len(record) < 6 {
			continue
		}

		date := record[0]
		shop := record[1]
		location := record[2]
		content := record[3]
		fee, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}
		settled := record[5] == "是"

		if date == "" || shop == "" {
			continue
		}

		if len(date) > 10 {
			date = date[:10]
		}

		createdAt, err := time.Parse("2006-01-02", date)
		if err != nil {
			createdAt = time.Now()
		}

		orders = append(orders, Order{
			UserID:    userID,
			Date:      date,
			Location:  location,
			Content:   content,
			Fee:       fee,
			Settled:   settled,
			Shop:      shop,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
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

	fmt.Printf("Successfully imported %d orders with user_id=%d\n", len(orders), userID)
}
