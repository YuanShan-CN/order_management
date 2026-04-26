package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/photographer/order_management/database"
	"github.com/photographer/order_management/models"
)

func GetShops(c *gin.Context) {
	var shops []models.Shop
	database.DB.Where("deleted_at IS NULL").Find(&shops)
	c.JSON(http.StatusOK, shops)
}

func GetShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shop)
}

func CreateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shop)
}

func UpdateShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&shop)
	c.JSON(http.StatusOK, shop)
}

func DeleteShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Where("deleted_at IS NULL").First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&shop).Update("deleted_at", gorm.Expr("NOW()"))
	c.JSON(http.StatusOK, gin.H{"message": "Shop moved to trash successfully"})
}

func GetDeletedShops(c *gin.Context) {
	var shops []models.Shop
	database.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&shops)
	c.JSON(http.StatusOK, shops)
}

func RestoreShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found in trash"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Unscoped().Model(&shop).Update("deleted_at", gorm.Expr("NULL"))
	c.JSON(http.StatusOK, gin.H{"message": "Shop restored successfully"})
}

func ForceDeleteShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found in trash"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Unscoped().Delete(&shop)
	c.JSON(http.StatusOK, gin.H{"message": "Shop permanently deleted successfully"})
}