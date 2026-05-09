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

type StatsSummary struct {
	TotalOrders      int64   `json:"totalOrders"`
	TotalIncome      float64 `json:"totalIncome"`
	SettledOrders    int64   `json:"settledOrders"`
	SettledIncome    float64 `json:"settledIncome"`
	UnsettledOrders  int64   `json:"unsettledOrders"`
	UnsettledIncome  float64 `json:"unsettledIncome"`
	UnsettledDeposit float64 `json:"unsettledDeposit"`
}

type MonthlyStat struct {
	Month  int     `json:"month"`
	Income float64 `json:"income"`
}

type ShopStat struct {
	Name   string  `json:"name"`
	Count  int     `json:"count"`
	Income float64 `json:"income"`
}

type UnifiedStatsResponse struct {
	Summary StatsSummary  `json:"summary"`
	Monthly []MonthlyStat `json:"monthly"`
	Shops   []ShopStat    `json:"shops"`
}

type PaginatedOrdersResponse struct {
	Orders      []models.Order `json:"orders"`
	CurrentPage int            `json:"currentPage"`
	TotalPages  int            `json:"totalPages"`
	TotalOrders int64          `json:"totalOrders"`
	PageSize    int            `json:"pageSize"`
}

func GetUnifiedStats(c *gin.Context) {
	userID := getUserID(c)

	year := c.Query("year")
	month := c.Query("month")
	settled := c.Query("settled")
	shop := c.Query("shop")

	var response UnifiedStatsResponse

	db := buildStatsQuery(userID, year, month, "all", shop).Model(&models.Order{})
	if err := db.Select(`
		COUNT(*) as total_orders,
		COALESCE(SUM(deposit + balance), 0) as total_income,
		SUM(CASE WHEN settled = true THEN 1 ELSE 0 END) as settled_orders,
		COALESCE(SUM(CASE WHEN settled = true THEN income ELSE 0 END), 0) as settled_income,
		SUM(CASE WHEN settled = false THEN 1 ELSE 0 END) as unsettled_orders,
		COALESCE(SUM(CASE WHEN settled = false THEN balance ELSE 0 END), 0) as unsettled_income,
		COALESCE(SUM(CASE WHEN settled = false THEN deposit ELSE 0 END), 0) as unsettled_deposit
	`).Scan(&response.Summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询统计数据失败"})
		return
	}

	db = buildStatsQuery(userID, year, month, settled, shop).Model(&models.Order{})
	var monthlyStats []MonthlyStat
	if settled == "unsettled" {
		db.Select("MONTH(date) as month, COALESCE(SUM(balance), 0) as income").
			Group("MONTH(date)").Order("month").Scan(&monthlyStats)
	} else {
		db.Select("MONTH(date) as month, COALESCE(SUM(deposit + balance), 0) as income").
			Group("MONTH(date)").Order("month").Scan(&monthlyStats)
	}
	response.Monthly = make([]MonthlyStat, 12)
	for i := 0; i < 12; i++ {
		response.Monthly[i] = MonthlyStat{Month: i + 1, Income: 0}
	}
	for _, stat := range monthlyStats {
		if stat.Month >= 1 && stat.Month <= 12 {
			response.Monthly[stat.Month-1] = stat
		}
	}

	db = buildStatsQuery(userID, year, month, settled, shop).Model(&models.Order{})
	var shopStats []ShopStat
	if settled == "unsettled" {
		db.Select("shop as name, COUNT(*) as count, COALESCE(SUM(balance), 0) as income").
			Group("shop").Order("income DESC").Scan(&shopStats)
	} else {
		db.Select("shop as name, COUNT(*) as count, COALESCE(SUM(deposit + balance), 0) as income").
			Group("shop").Order("income DESC").Scan(&shopStats)
	}
	if shopStats == nil {
		response.Shops = []ShopStat{}
	} else {
		response.Shops = shopStats
	}

	c.JSON(http.StatusOK, response)
}

func buildStatsQuery(userID uint, year, month, settled, shop string) *gorm.DB {
	db := database.DB.Where("deleted_at IS NULL AND user_id = ?", userID)

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

	return db
}

func GetStatsSummary(c *gin.Context) {
	userID := getUserID(c)

	year := c.Query("year")
	month := c.Query("month")
	shop := c.Query("shop")

	db := buildStatsQuery(userID, year, month, "all", shop).Model(&models.Order{})

	var summary StatsSummary

	db.Select(`
		COUNT(*) as total_orders,
		COALESCE(SUM(deposit + balance), 0) as total_income,
		SUM(CASE WHEN settled = true THEN 1 ELSE 0 END) as settled_orders,
		COALESCE(SUM(CASE WHEN settled = true THEN income ELSE 0 END), 0) as settled_income,
		SUM(CASE WHEN settled = false THEN 1 ELSE 0 END) as unsettled_orders,
		COALESCE(SUM(CASE WHEN settled = false THEN balance ELSE 0 END), 0) as unsettled_income,
		COALESCE(SUM(CASE WHEN settled = false THEN deposit ELSE 0 END), 0) as unsettled_deposit
	`).Scan(&summary)

	c.JSON(http.StatusOK, summary)
}

func GetMonthlyStats(c *gin.Context) {
	userID := getUserID(c)

	year := c.Query("year")
	month := c.Query("month")
	settled := c.Query("settled")
	shop := c.Query("shop")

	db := buildStatsQuery(userID, year, month, settled, shop).Model(&models.Order{})

	var monthlyStats []MonthlyStat

	if settled == "unsettled" {
		db.Select("MONTH(date) as month, SUM(balance) as income").
			Group("MONTH(date)").
			Order("month").
			Scan(&monthlyStats)
	} else {
		db.Select("MONTH(date) as month, SUM(deposit + balance) as income").
			Group("MONTH(date)").
			Order("month").
			Scan(&monthlyStats)
	}

	result := make([]MonthlyStat, 12)
	for i := 0; i < 12; i++ {
		result[i] = MonthlyStat{Month: i + 1, Income: 0}
	}

	for _, stat := range monthlyStats {
		if stat.Month >= 1 && stat.Month <= 12 {
			result[stat.Month-1] = stat
		}
	}

	c.JSON(http.StatusOK, result)
}

func GetShopStats(c *gin.Context) {
	userID := getUserID(c)

	year := c.Query("year")
	month := c.Query("month")
	settled := c.Query("settled")
	shop := c.Query("shop")

	db := buildStatsQuery(userID, year, month, settled, shop).Model(&models.Order{})

	var shopStats []ShopStat

	if settled == "unsettled" {
		db.Select("shop as name, COUNT(*) as count, SUM(balance) as income").
			Group("shop").
			Order("income DESC").
			Scan(&shopStats)
	} else {
		db.Select("shop as name, COUNT(*) as count, SUM(deposit + balance) as income").
			Group("shop").
			Order("income DESC").
			Scan(&shopStats)
	}

	c.JSON(http.StatusOK, shopStats)
}

func GetShopOrderDetails(c *gin.Context) {
	userID := getUserID(c)

	year := c.Query("year")
	month := c.Query("month")
	settled := c.Query("settled")
	shopName := c.Query("shop")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	db := buildStatsQuery(userID, year, month, settled, shopName).Model(&models.Order{})

	var total int64
	db.Count(&total)

	offset := (page - 1) * pageSize
	var orders []models.Order
	db.Order("date DESC").Offset(offset).Limit(pageSize).Find(&orders)

	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, PaginatedOrdersResponse{
		Orders:      orders,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalOrders: total,
		PageSize:    pageSize,
	})
}

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
	case "Income":
		return "收入"
	case "Deposit":
		return "定金"
	case "Balance":
		return "尾款"
	case "Shop":
		return "店铺"
	default:
		return field
	}
}

// 提取日期的年月日部分
func extractDatePart(date string) string {
	// 如果是日期时间格式（如 2026-04-28T00:00:00+08:00），只保留前10个字符
	if len(date) > 10 && (date[10] == 'T' || date[10] == ' ') {
		return date[:10]
	}
	// 否则返回原日期
	return date
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
	case "income_desc":
		db = db.Order("income DESC")
	case "income_asc":
		db = db.Order("income ASC")
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
	order.CalculateIncome()
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
	order.CalculateIncome()
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

	headers := []string{"ID", "日期", "地点", "内容", "店铺", "定金", "尾款", "已结算", "收入"}
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
			extractDatePart(order.Date),
			order.Location,
			order.Content,
			order.Shop,
			strconv.FormatFloat(order.Deposit, 'f', 2, 64),
			strconv.FormatFloat(order.Balance, 'f', 2, 64),
			settledStr,
			strconv.FormatFloat(order.Income, 'f', 2, 64),
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

	headers := []string{"日期", "地点", "内容", "店铺", "定金", "尾款", "已结算", "收入"}
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
			extractDatePart(order.Date),
			order.Location,
			order.Content,
			order.Shop,
			strconv.FormatFloat(order.Deposit, 'f', 2, 64),
			strconv.FormatFloat(order.Balance, 'f', 2, 64),
			settledStr,
			strconv.FormatFloat(order.Income, 'f', 2, 64),
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无法解析CSV文件: %v\n请确保文件是标准的CSV格式，使用英文逗号分隔", err)})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV文件中没有数据\n请确保文件包含表头行和至少一行数据"})
		return
	}

	// 分析表头，确定列的位置
	header := records[0]
	dateIdx := -1
	shopIdx := -1
	locationIdx := -1
	contentIdx := -1
	incomeIdx := -1
	depositIdx := -1
	balanceIdx := -1
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
		case "收入":
			incomeIdx = i
		case "定金":
			depositIdx = i
		case "尾款":
			balanceIdx = i
		case "已结算":
			settledIdx = i
		}
	}

	// 检查必要的列是否存在
	if dateIdx == -1 || shopIdx == -1 || settledIdx == -1 {
		var missingColumns []string
		if dateIdx == -1 {
			missingColumns = append(missingColumns, "日期")
		}
		if shopIdx == -1 {
			missingColumns = append(missingColumns, "店铺")
		}
		if settledIdx == -1 {
			missingColumns = append(missingColumns, "已结算")
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("CSV格式不正确，缺少必要的列：%s\n期望的CSV格式示例：\n日期,店铺,地点,内容,定金,尾款,已结算\n2026-04-29,示例店,北京,购买商品,500.00,500.00,否",
				strings.Join(missingColumns, "、"))})
		return
	}

	var importedCount int
	var failedCount int
	var newShopCount int
	var errors []string

	// 使用map缓存已查询的shop，避免重复查询
	shopCache := make(map[string]bool)

	// 先查询该用户已有的所有shop，缓存到map中
	var existingShops []models.Shop
	database.DB.Where("user_id = ?", userID).Find(&existingShops)
	for _, s := range existingShops {
		shopCache[s.Name] = true
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) <= dateIdx || len(record) <= shopIdx || len(record) <= settledIdx {
			failedCount++
			errors = append(errors, fmt.Sprintf("第%d行: 数据列数不足，请确保包含所有必要的列（日期、店铺、已结算）", i+1))
			continue
		}

		date := strings.TrimSpace(record[dateIdx])
		shop := strings.TrimSpace(record[shopIdx])
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
			errors = append(errors, fmt.Sprintf("第%d行: 日期或店铺不能为空，请填写完整信息", i+1))
			continue
		}

		// 处理日期格式：支持 2026-04-28T00:00:00+08:00 和 2026-04-28
		if len(date) > 10 && (date[10] == 'T' || date[10] == ' ') {
			date = date[:10]
		}

		// 解析定金
		var deposit float64 = 0.0
		if depositIdx != -1 && len(record) > depositIdx {
			depositStr := strings.TrimSpace(record[depositIdx])
			if depositStr != "" {
				d, err := strconv.ParseFloat(depositStr, 64)
				if err == nil {
					deposit = d
				}
			}
		}

		// 解析尾款
		var balance float64 = 0.0
		if balanceIdx != -1 && len(record) > balanceIdx {
			balanceStr := strings.TrimSpace(record[balanceIdx])
			if balanceStr != "" {
				b, err := strconv.ParseFloat(balanceStr, 64)
				if err == nil {
					balance = b
				}
			}
		}

		settled := settledStr == "是" || settledStr == "true" || settledStr == "1"

		tempOrder := models.Order{
			Settled: settled,
			Deposit: deposit,
			Balance: balance,
		}
		tempOrder.CalculateIncome()
		income := tempOrder.Income

		// 如果CSV中指定了收入，且自动计算的收入为0，则使用CSV中的值
		if income == 0 && incomeIdx != -1 && len(record) > incomeIdx {
			incomeStr := strings.TrimSpace(record[incomeIdx])
			if incomeStr != "" {
				f, err := strconv.ParseFloat(incomeStr, 64)
				if err == nil {
					income = f
				}
			}
		}

		// 检查并自动创建shop（使用缓存map）
		if !shopCache[shop] {
			// 缓存中不存在，检查数据库
			var existingShop models.Shop
			result := database.DB.Where("name = ? AND user_id = ?", shop, userID).First(&existingShop)
			if result.Error == gorm.ErrRecordNotFound {
				// 数据库中也不存在，创建新shop
				newShop := models.Shop{
					UserID: userID,
					Name:   shop,
				}
				if err := database.DB.Create(&newShop).Error; err != nil {
					failedCount++
					errors = append(errors, fmt.Sprintf("第%d行: 创建店铺失败 - %v", i+1, err))
					continue
				}
				// 将新创建的shop加入缓存，并增加计数
				shopCache[shop] = true
				newShopCount++
			} else if result.Error != nil {
				failedCount++
				errors = append(errors, fmt.Sprintf("第%d行: 检查店铺失败 - %v", i+1, result.Error))
				continue
			} else {
				// 数据库中存在但缓存中没有，加入缓存
				shopCache[shop] = true
			}
		}

		order := models.Order{
			UserID:   userID,
			Date:     date,
			Location: location,
			Content:  content,
			Income:   income,
			Deposit:  deposit,
			Balance:  balance,
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

	message := fmt.Sprintf("成功导入 %d 条数据，失败 %d 条", importedCount, failedCount)
	if newShopCount > 0 {
		message = fmt.Sprintf("%s，新建 %d 个店铺", message, newShopCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"importedCount": importedCount,
		"failedCount":   failedCount,
		"newShopCount":  newShopCount,
		"errors":        errors,
		"message":       message,
	})
}
