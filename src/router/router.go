package router

import (
	"net/http"

	"github.com/YuanShan-CN/order_management/src/controllers"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/middleware"
	"github.com/YuanShan-CN/order_management/src/models"
	"github.com/YuanShan-CN/order_management/src/scheduler"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/logout", controllers.Logout)

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/auth/current", controllers.GetCurrentUser)
		authGroup.PUT("/auth/password", controllers.ChangePassword)
		authGroup.DELETE("/auth/account", controllers.DeleteAccount)

		authGroup.GET("/orders", controllers.GetOrders)
		authGroup.GET("/orders/export", controllers.ExportOrdersCSV)
		authGroup.GET("/orders/stats/export", controllers.ExportStatsCSV)
		authGroup.GET("/orders/stats", controllers.GetUnifiedStats)
		authGroup.GET("/orders/stats/summary", controllers.GetStatsSummary)
		authGroup.GET("/orders/stats/monthly", controllers.GetMonthlyStats)
		authGroup.GET("/orders/stats/shops", controllers.GetShopStats)
		authGroup.GET("/orders/stats/shop-details", controllers.GetShopOrderDetails)
		authGroup.POST("/orders/import", controllers.ImportOrdersCSV)
		authGroup.GET("/orders/:id", controllers.GetOrder)
		authGroup.POST("/orders", controllers.CreateOrder)
		authGroup.PUT("/orders/:id", controllers.UpdateOrder)
		authGroup.DELETE("/orders/:id", controllers.DeleteOrder)
		authGroup.GET("/orders/trash", controllers.GetDeletedOrders)
		authGroup.POST("/orders/:id/restore", controllers.RestoreOrder)
		authGroup.DELETE("/orders/:id/force-delete", controllers.ForceDeleteOrder)

		authGroup.GET("/shops", controllers.GetShops)
		authGroup.GET("/shops/:id", controllers.GetShop)
		authGroup.POST("/shops", controllers.CreateShop)
		authGroup.PUT("/shops/:id", controllers.UpdateShop)
		authGroup.DELETE("/shops/:id", controllers.DeleteShop)
		authGroup.GET("/shops/trash", controllers.GetDeletedShops)
		authGroup.POST("/shops/:id/restore", controllers.RestoreShop)
		authGroup.DELETE("/shops/:id/force-delete", controllers.ForceDeleteShop)
	}

	r.LoadHTMLGlob("templates/*")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
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

	// 调试用的定时任务接口
	r.GET("/api/scheduler/run-now", func(c *gin.Context) {
		err := scheduler.RunNow()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Export triggered"})
	})

	// 处理 Vite 开发服务器客户端请求（避免 404 日志）
	r.GET("/@vite/client", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	return r
}