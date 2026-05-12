package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
	"github.com/YuanShan-CN/order_management/src/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// shopValidate 全局验证器实例，用于店铺数据验证
var shopValidate = validator.New()

// validateShop 验证店铺数据
// 参数: shop - 待验证的店铺对象指针
// 返回: 验证错误信息，如果验证通过返回nil
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

// getShopFieldName 将字段名转换为中文显示名
// 参数: field - 字段名
// 返回: 中文显示名
func getShopFieldName(field string) string {
	switch field {
	case "Name":
		return "店铺名称"
	default:
		return field
	}
}

// getShopsCacheKey 生成用户店铺列表的缓存键
// 参数: userID - 用户ID
// 返回: 缓存键字符串
func getShopsCacheKey(userID uint) string {
	return "shops:user:" + strconv.FormatUint(uint64(userID), 10)
}

// GetShops 获取当前用户的店铺列表（带Redis缓存）
// 缓存策略：先从Redis读取，未命中则查询数据库并写入缓存，有效期1小时
// GET /api/shops
func GetShops(c *gin.Context) {
	userID := getUserID(c)
	cacheKey := getShopsCacheKey(userID)

	ctx := context.Background()
	var shops []models.Shop

	// 尝试从Redis获取缓存数据
	cached, err := redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &shops); err == nil {
			c.JSON(http.StatusOK, shops)
			return
		}
	}

	// 缓存未命中，查询数据库
	database.DB.Where("user_id = ?", userID).Order("id DESC").Find(&shops)

	// 将查询结果写入Redis缓存，有效期1小时
	data, _ := json.Marshal(shops)
	redis.Client.Set(ctx, cacheKey, data, time.Hour)

	c.JSON(http.StatusOK, shops)
}

// GetShop 获取单个店铺详情
// GET /api/shops/:id
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

// CreateShop 创建新店铺
// POST /api/shops
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

	// 清除店铺列表缓存，确保数据一致性
	redis.Client.Del(context.Background(), getShopsCacheKey(userID))

	c.JSON(http.StatusCreated, shop)
}

// UpdateShop 更新店铺信息
// PUT /api/shops/:id
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

	// 清除店铺列表缓存
	redis.Client.Del(context.Background(), getShopsCacheKey(userID))

	c.JSON(http.StatusOK, shop)
}

// DeleteShop 软删除店铺（移到回收站）
// DELETE /api/shops/:id
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

	// 清除店铺列表缓存
	redis.Client.Del(context.Background(), getShopsCacheKey(userID))

	c.JSON(http.StatusOK, gin.H{"message": "Shop moved to trash successfully"})
}

// GetDeletedShops 获取回收站中的店铺列表
// GET /api/shops/trash
func GetDeletedShops(c *gin.Context) {
	userID := getUserID(c)
	var shops []models.Shop
	database.DB.Unscoped().Where("deleted_at IS NOT NULL AND user_id = ?", userID).Find(&shops)
	c.JSON(http.StatusOK, shops)
}

// RestoreShop 从回收站恢复店铺
// PUT /api/shops/:id/restore
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

// ForceDeleteShop 永久删除店铺（不可恢复）
// DELETE /api/shops/:id/force
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
