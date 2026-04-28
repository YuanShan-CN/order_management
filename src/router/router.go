package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/YuanShan-CN/order_management/src/controllers"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/middleware"
	"github.com/YuanShan-CN/order_management/src/models"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

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

	return r
}