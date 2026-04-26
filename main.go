package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/photographer/order_management/config"
	"github.com/photographer/order_management/controllers"
	"github.com/photographer/order_management/database"
	"github.com/photographer/order_management/models"
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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/orders", controllers.GetOrders)
	r.GET("/api/orders/:id", controllers.GetOrder)
	r.POST("/api/orders", controllers.CreateOrder)
	r.PUT("/api/orders/:id", controllers.UpdateOrder)
	r.DELETE("/api/orders/:id", controllers.DeleteOrder)
	r.GET("/api/orders/trash", controllers.GetDeletedOrders)
	r.POST("/api/orders/:id/restore", controllers.RestoreOrder)
	r.DELETE("/api/orders/:id/force-delete", controllers.ForceDeleteOrder)

	r.GET("/api/shops", controllers.GetShops)
	r.GET("/api/shops/:id", controllers.GetShop)
	r.POST("/api/shops", controllers.CreateShop)
	r.PUT("/api/shops/:id", controllers.UpdateShop)
	r.DELETE("/api/shops/:id", controllers.DeleteShop)
	r.GET("/api/shops/trash", controllers.GetDeletedShops)
	r.POST("/api/shops/:id/restore", controllers.RestoreShop)
	r.DELETE("/api/shops/:id/force-delete", controllers.ForceDeleteShop)

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/orders", func(c *gin.Context) {
		c.HTML(http.StatusOK, "orders.html", nil)
	})
	r.GET("/shops", func(c *gin.Context) {
		c.HTML(http.StatusOK, "shops.html", nil)
	})
	r.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", nil)
	})
	r.GET("/trash", func(c *gin.Context) {
		c.HTML(http.StatusOK, "trash.html", nil)
	})

	r.GET("/api/clean", func(c *gin.Context) {
		database.DB.Delete(&models.Order{})
		database.DB.Delete(&models.Shop{})
		c.JSON(http.StatusOK, gin.H{"message": "数据已清空"})
	})

	port := strconv.Itoa(config.AppConfig.Server.Port)
	log.Println("Server running on http://localhost:" + port)
	r.Run(":" + port)
}
