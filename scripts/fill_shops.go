package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/YuanShan-CN/order_management/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	var userID uint
	flag.UintVar(&userID, "user_id", 1, "User ID to assign to created shops")
	flag.Parse()

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

	var shopNames []string
	err = db.Model(&models.Order{}).Where("user_id = ?", userID).Distinct("shop").Pluck("shop", &shopNames).Error
	if err != nil {
		fmt.Printf("Failed to get shop names: %v\n", err)
		return
	}

	fmt.Printf("Found %d unique shop names in orders table for user_id=%d\n", len(shopNames), userID)

	now := time.Now()
	createdCount := 0
	for _, name := range shopNames {
		if name == "" {
			continue
		}

		var existingShop models.Shop
		result := db.Where("name = ? AND user_id = ?", name, userID).First(&existingShop)

		if result.Error == gorm.ErrRecordNotFound {
			shop := models.Shop{
				UserID: userID,
				Name:   name,
			}
			shop.CreatedAt = now
			err := db.Create(&shop).Error
			if err != nil {
				fmt.Printf("Failed to create shop '%s': %v\n", name, err)
			} else {
				fmt.Printf("Created shop: %s (user_id=%d)\n", name, userID)
				createdCount++
			}
		} else if result.Error != nil {
			fmt.Printf("Failed to check shop '%s': %v\n", name, result.Error)
		} else {
			fmt.Printf("Shop already exists: %s (user_id=%d)\n", name, userID)
		}
	}

	fmt.Printf("\nDone! Created %d new shops for user_id=%d\n", createdCount, userID)
}