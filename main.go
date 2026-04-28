package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/router"
)

func main() {
	configPath := flag.String("c", "config.yaml", "Path to config file")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer database.Close()

	r := router.Setup()

	port := strconv.Itoa(config.AppConfig.Server.Port)
	log.Println("Server running on http://localhost:" + port)
	r.Run(":" + port)
}
