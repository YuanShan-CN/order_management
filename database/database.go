package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/photographer/order_management/models"
)

var DB *gorm.DB

type sqlLogger struct{}

func (l sqlLogger) Print(v ...interface{}) {
	log.Println(v...)
}

func Connect() error {
	var err error
	DB, err = gorm.Open("mysql", "root:123456@tcp(localhost:3306)/order_management?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	DB.LogMode(true)
	DB.SetLogger(sqlLogger{})

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	DB.AutoMigrate(&models.Order{}, &models.Shop{})
	return nil
}

func Close() {
	DB.Close()
}