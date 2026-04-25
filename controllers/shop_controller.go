package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/photographer/order_management/database"
	"github.com/photographer/order_management/models"
)

func GetShops(c *gin.Context) {
	var shops []models.Shop
	database.DB.Find(&shops)
	c.JSON(http.StatusOK, shops)
}

func GetShop(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.First(&shop, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
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
		if gorm.IsRecordNotFoundError(err) {
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
	if err := database.DB.First(&shop, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Delete(&shop)
	c.JSON(http.StatusOK, gin.H{"message": "Shop deleted successfully"})
}