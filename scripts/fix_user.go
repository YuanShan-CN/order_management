package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
)

func main() {
	oldUsername := flag.String("old", "", "旧用户名")
	newUsername := flag.String("new", "", "新用户名")
	flag.Parse()

	if *oldUsername == "" || *newUsername == "" {
		log.Fatal("请提供参数: -old <旧用户名> -new <新用户名>")
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

	result := database.DB.Model(&struct{}{}).Table("users").Where("username = ?", *oldUsername).Update("username", *newUsername)
	if result.Error != nil {
		log.Fatal("更新用户失败:", result.Error)
	}

	if result.RowsAffected == 0 {
		fmt.Printf("用户 '%s' 不存在\n", *oldUsername)
	} else {
		fmt.Printf("用户 '%s' 已更新为 '%s'\n", *oldUsername, *newUsername)
	}
}