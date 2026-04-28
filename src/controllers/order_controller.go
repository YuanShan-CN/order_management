package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func validateOrder(order *models.Order) error {
	order.Location = strings.TrimSpace(order.Location)
	order.Content = strings.TrimSpace(order.Content)
	order.Shop = strings.TrimSpace(order.Shop)

	if err := validate.Struct(order); err != nil {
		var errMsg string
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := getFieldName(e.Field())
			switch e.Tag() {
			case "required":
				errMsg = fmt.Sprintf("%s不能为空", fieldName)
			case "datetime":
				errMsg = fmt.Sprintf("%s格式不正确，应为 YYYY-MM-DD", fieldName)
			case "numeric":
				errMsg = fmt.Sprintf("%s必须是数字", fieldName)
			default:
				errMsg = fmt.Sprintf("%s验证失败", fieldName)
			}
			return errors.New(errMsg)
		}
	}
	return nil
}

func getFieldName(field string) string {
	switch field {
	case "Date":
		return "日期"
	case "Location":
		return "地点"
	case "Content":
		return "内容"
	case "Fee":
		return "费用"
	case "Shop":
		return "店铺"
	default:
		return field
	}
}

func getUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func GetOrders(c *gin.Context) {
	userID := getUserID(c)
	var orders []models.Order
	db := database.DB.Where("deleted_at IS NULL AND user_id = ?", userID)

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

	var total int64
	db.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Find(&orders)

	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"currentPage": page,
		"totalPages":  totalPages,
		"totalOrders": total,
		"pageSize":    pageSize,
	})
}

func GetOrder(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	userID := getUserID(c)
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateOrder(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.UserID = userID
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err := validateOrder(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.UserID = userID
	database.DB.Save(&order)
	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Where("deleted_at IS NULL AND id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	userID := getUserID(c)
	var orders []models.Order
	db := database.DB.Unscoped().Where("deleted_at IS NOT NULL AND user_id = ?", userID)

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

	var total int64
	db.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Find(&orders)

	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"currentPage": page,
		"totalPages":  totalPages,
		"totalOrders": total,
		"pageSize":    pageSize,
	})
}

func RestoreOrder(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL AND id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	userID := getUserID(c)
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Unscoped().Where("deleted_at IS NOT NULL AND id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found in trash"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Unscoped().Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "Order permanently deleted successfully"})
}

func ExportOrdersCSV(c *gin.Context) {
	userID := getUserID(c)
	var orders []models.Order
	db := database.DB.Where("deleted_at IS NULL AND user_id = ?", userID)

	search := c.Query("search")
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

	db.Order("date DESC").Find(&orders)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=orders.csv")
	c.Header("Content-Transfer-Encoding", "binary")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"ID", "日期", "地点", "内容", "费用", "已结算", "店铺"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers"})
		return
	}

	for _, order := range orders {
		settledStr := "否"
		if order.Settled {
			settledStr = "是"
		}

		row := []string{
			strconv.Itoa(int(order.ID)),
			order.Date,
			order.Location,
			order.Content,
			strconv.FormatFloat(order.Fee, 'f', 2, 64),
			settledStr,
			order.Shop,
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
			return
		}
	}
}

func ExportStatsCSV(c *gin.Context) {
	userID := getUserID(c)
	var orders []models.Order
	db := database.DB.Where("deleted_at IS NULL AND user_id = ?", userID)

	year := c.Query("year")
	month := c.Query("month")
	settled := c.Query("settled")
	shop := c.Query("shop")

	if year != "" && year != "all" {
		db = db.Where("YEAR(date) = ?", year)
	}

	if month != "" && month != "all" {
		db = db.Where("MONTH(date) = ?", month)
	}

	if settled == "settled" {
		db = db.Where("settled = ?", true)
	} else if settled == "unsettled" {
		db = db.Where("settled = ?", false)
	}

	if shop != "" && shop != "all" {
		db = db.Where("shop = ?", shop)
	}

	db.Order("date DESC").Find(&orders)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=stats.csv")
	c.Header("Content-Transfer-Encoding", "binary")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"日期", "店铺", "地点", "内容", "费用", "已结算"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers"})
		return
	}

	for _, order := range orders {
		settledStr := "否"
		if order.Settled {
			settledStr = "是"
		}

		row := []string{
			order.Date,
			order.Shop,
			order.Location,
			order.Content,
			strconv.FormatFloat(order.Fee, 'f', 2, 64),
			settledStr,
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
			return
		}
	}
}

func ImportOrdersCSV(c *gin.Context) {
	userID := getUserID(c)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析CSV文件"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV文件中没有数据"})
		return
	}

	// 分析表头，确定列的位置
	header := records[0]
	dateIdx := -1
	shopIdx := -1
	locationIdx := -1
	contentIdx := -1
	feeIdx := -1
	settledIdx := -1

	for i, col := range header {
		col = strings.TrimSpace(col)
		switch col {
		case "日期":
			dateIdx = i
		case "店铺":
			shopIdx = i
		case "地点":
			locationIdx = i
		case "内容":
			contentIdx = i
		case "费用":
			feeIdx = i
		case "已结算":
			settledIdx = i
		}
	}

	// 检查必要的列是否存在
	if dateIdx == -1 || shopIdx == -1 || feeIdx == -1 || settledIdx == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV格式不正确，缺少必要的列（日期、店铺、费用、已结算）"})
		return
	}

	var importedCount int
	var failedCount int
	var errors []string

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) <= dateIdx || len(record) <= shopIdx || len(record) <= feeIdx || len(record) <= settledIdx {
			failedCount++
			errors = append(errors, fmt.Sprintf("第%d行: 数据列数不足", i+1))
			continue
		}

		date := strings.TrimSpace(record[dateIdx])
		shop := strings.TrimSpace(record[shopIdx])
		feeStr := strings.TrimSpace(record[feeIdx])
		settledStr := strings.TrimSpace(record[settledIdx])

		var location, content string
		if locationIdx != -1 && len(record) > locationIdx {
			location = strings.TrimSpace(record[locationIdx])
		}
		if contentIdx != -1 && len(record) > contentIdx {
			content = strings.TrimSpace(record[contentIdx])
		}

		if date == "" || shop == "" {
			failedCount++
			errors = append(errors, fmt.Sprintf("第%d行: 日期或店铺不能为空", i+1))
			continue
		}

		// 处理日期格式：支持 2026-04-28T00:00:00+08:00 和 2026-04-28
		if len(date) > 10 && (date[10] == 'T' || date[10] == ' ') {
			date = date[:10]
		}

		fee, err := strconv.ParseFloat(feeStr, 64)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("第%d行: 费用格式不正确", i+1))
			continue
		}

		settled := settledStr == "是" || settledStr == "true" || settledStr == "1"

		order := models.Order{
			UserID:   userID,
			Date:     date,
			Location: location,
			Content:  content,
			Fee:      fee,
			Settled:  settled,
			Shop:     shop,
		}

		if err := database.DB.Create(&order).Error; err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("第%d行: 导入失败 - %v", i+1, err))
			continue
		}

		importedCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"importedCount": importedCount,
		"failedCount":   failedCount,
		"errors":        errors,
		"message":       fmt.Sprintf("成功导入 %d 条数据，失败 %d 条", importedCount, failedCount),
	})
}
