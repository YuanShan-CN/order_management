package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint    `gorm:"primary_key"`
	Date      string  `gorm:"type:date"`
	Location  string  `gorm:"type:varchar(200)"`
	Content   string  `gorm:"type:text"`
	Fee       float64 `gorm:"type:decimal(10,2)"`
	Shop      string  `gorm:"type:varchar(100)"`
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

	var countBefore int64
	db.Model(&Order{}).Count(&countBefore)
	fmt.Printf("Before: %d orders\n", countBefore)

	type DuplicateGroup struct {
		Date     string
		Location string
		Content  string
		Fee      float64
		Shop     string
		MinID    uint
	}

	var groups []DuplicateGroup
	err = db.Table("orders").
		Select("date, location, content, fee, shop, MIN(id) as min_id").
		Group("date, location, content, fee, shop").
		Having("COUNT(*) > 1").
		Scan(&groups).Error

	if err != nil {
		fmt.Printf("Failed to find duplicates: %v\n", err)
		return
	}

	fmt.Printf("Found %d duplicate groups\n", len(groups))

	deletedCount := 0
	for _, group := range groups {
		var toDelete []Order
		err := db.Where("date = ? AND location = ? AND content = ? AND fee = ? AND shop = ? AND id != ?",
			group.Date, group.Location, group.Content, group.Fee, group.Shop, group.MinID).
			Find(&toDelete).Error

		if err != nil {
			fmt.Printf("Failed to find duplicates for group: %v\n", err)
			continue
		}

		if len(toDelete) > 0 {
			err := db.Delete(&toDelete).Error
			if err != nil {
				fmt.Printf("Failed to delete duplicates: %v\n", err)
			} else {
				deletedCount += len(toDelete)
				fmt.Printf("Deleted %d duplicates for shop '%s'\n", len(toDelete), group.Shop)
			}
		}
	}

	var countAfter int64
	db.Model(&Order{}).Count(&countAfter)
	fmt.Printf("\nAfter: %d orders\n", countAfter)
	fmt.Printf("Total deleted: %d duplicates\n", deletedCount)
}
