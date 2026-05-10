package database

import (
	"log"
	"os"
	"strconv"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	cfg := config.AppConfig.Database

	log.Printf("Using database timezone: UTC")
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.Database + "?charset=utf8mb4&parseTime=True&loc=UTC"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	err = DB.AutoMigrate(&models.Order{}, &models.Shop{}, &models.User{})
	if err != nil {
		return err
	}

	log.Printf("Database connected: %s:%d/%s", cfg.Host, cfg.Port, cfg.Database)
	return nil
}

func Close() {
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}
