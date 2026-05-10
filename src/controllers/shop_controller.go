package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var shopValidate = validator.New()

func validateShop(shop *models.Shop) error {
	shop.Name = strings.TrimSpace(shop.Name)

	if err := shopValidate.Struct(shop); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := getShopFieldName(e.Field())
			switch e.Tag() {
			case "required":
				return errors.New(fmt.Sprintf("%s不能为空", fieldName))
			default:
				return errors.New(fmt.Sprintf("%s验证失败", fieldName))
			}
		}
	}
	return nil
}

func getShopFieldName(field string) string {
	switch field {
	case "Name":
		return "店铺名称"
	default:
		return field
	}
}

func GetShops(c *gin.Context) {
	userID := getUserID(c)
	var shops []models.Shop
	database.DB.Where("user_id = ?", userID).Order("id DESC").Find(&shops)
	c.JSON(http.StatusOK, shops)
}

func GetShop(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&shop).Error; err != nil {
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
	userID := getUserID(c)
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateShop(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shop.UserID = userID
	if err := database.DB.Create(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shop)
}

func UpdateShop(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&shop).Error; err != nil {
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
	if err := validateShop(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shop.UserID = userID
	database.DB.Save(&shop)
	c.JSON(http.StatusOK, shop)
}

func DeleteShop(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&shop).Error; err != nil {
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
	userID := getUserID(c)
	var shops []models.Shop
	database.DB.Unscoped().Where("deleted_at IS NOT NULL AND user_id = ?", userID).Find(&shops)
	c.JSON(http.StatusOK, shops)
}

func RestoreShop(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL AND id = ? AND user_id = ?", id, userID).First(&shop).Error; err != nil {
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
	userID := getUserID(c)
	id := c.Param("id")
	var shop models.Shop
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL AND id = ? AND user_id = ?", id, userID).First(&shop).Error; err != nil {
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