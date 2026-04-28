package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := flag.String("username", "", "用户名")
	password := flag.String("password", "", "新密码")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("请提供参数: -username <用户名> -password <新密码>")
	}

	normalizedUsername := *username
	if len(*username) >= 3 {
		normalizedUsername = fmt.Sprintf("%s", *username)
	}

	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer database.Close()

	var user models.User
	if err := database.DB.Where("username = ?", normalizedUsername).First(&user).Error; err != nil {
		log.Fatalf("用户 '%s' 不存在", normalizedUsername)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}

	if err := database.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		log.Fatal("更新密码失败:", err)
	}

	fmt.Printf("用户 '%s' 的密码已更新\n", normalizedUsername)
}