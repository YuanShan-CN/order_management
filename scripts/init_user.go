package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := flag.String("username", "", "用户名")
	password := flag.String("password", "", "密码")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("请提供用户名和密码参数: -username <用户名> -password <密码>")
	}

	normalizedUsername := strings.ToLower(strings.TrimSpace(*username))

	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer database.Close()

	var existingUser models.User
	if err := database.DB.Where("username = ?", normalizedUsername).First(&existingUser).Error; err == nil {
		log.Fatalf("用户 '%s' 已存在", normalizedUsername)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}

	user := models.User{
		Username: normalizedUsername,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Fatal("创建用户失败:", err)
	}

	fmt.Printf("用户 '%s' 创建成功，ID: %d\n", user.Username, user.ID)
}