package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/photographer/order_management/database"
	"github.com/photographer/order_management/models"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order
	db := database.DB.Where("deleted_at IS NULL")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "15"))
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "date_desc")
	settled := c.Query("settled")
	year := c.Query("year")
	month := c.Query("month")
	shop := c.Query("shop")

	if search != "" {
		db = db.Where("location LIKE ? OR content LIKE ? OR shop LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if settled == "settled" {
		db = db.Where("settled = ?", true)
	} else if settled == "unsettled" {
		db = db.Where("settled = ?", false)
	}

	if year != "" && year != "all" {
		db = db.Where("YEAR(date) = ?", year)
	}

	if month != "" && month != "all" {
		db = db.Where("MONTH(date) = ?", month)
	}

	if shop != "" && shop != "all" {
		db = db.Where("shop = ?", shop)
	}

	switch sort {
	case "date_asc":
		db = db.Order("date ASC")
	case "fee_desc":
		db = db.Order("fee DESC")
	case "fee_asc":
		db = db.Order("fee ASC")
	default:
		db = db.Order("date DESC")
	}

	var total int
	db.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Find(&orders)

	totalPages := (total + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"orders":       orders,
		"currentPage":  page,
		"totalPages":   totalPages,
		"totalOrders":  total,
		"pageSize":     pageSize,
	})
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&order)
	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Where("deleted_at IS NULL").First(&order, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&order).Update("deleted_at", gorm.Expr("NOW()"))
	c.JSON(http.StatusOK, gin.H{"message": "Order moved to trash successfully"})
}

func GetDeletedOrders(c *gin.Context) {
	var orders []models.Order
	db := database.DB.Unscoped().Where("deleted_at IS NOT NULL")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "15"))
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "deleted_at_desc")

	if search != "" {
		db = db.Where("location LIKE ? OR content LIKE ? OR shop LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	switch sort {
	case "deleted_at_asc":
		db = db.Order("deleted_at ASC")
	case "date_desc":
		db = db.Order("date DESC")
	case "date_asc":
		db = db.Order("date ASC")
	default:
		db = db.Order("deleted_at DESC")
	}

	var total int
	db.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Find(&orders)

	totalPages := (total + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"orders":       orders,
		"currentPage":  page,
		"totalPages":   totalPages,
		"totalOrders":  total,
		"pageSize":     pageSize,
	})
}

func RestoreOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").First(&order, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found in trash"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Unscoped().Model(&order).Update("deleted_at", gorm.Expr("NULL"))
	c.JSON(http.StatusOK, gin.H{"message": "Order restored successfully"})
}

func ForceDeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").First(&order, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found in trash"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Unscoped().Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "Order permanently deleted successfully"})
}