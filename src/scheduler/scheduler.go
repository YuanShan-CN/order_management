// Package scheduler 提供定时任务功能，用于自动导出订单数据为CSV文件
// 主要功能：
// - 检测运行环境（Docker/本地）
// - 每小时检查一次，如果今天无导出文件则执行导出
// - 支持优雅关闭
package scheduler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/models"
)

// getLocation 获取配置的时区
func getLocation() *time.Location {
	timezone := config.AppConfig.Scheduler.Timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Warning: Failed to load timezone %s, using local time: %v", timezone, err)
		return time.Local
	}
	return loc
}

// globalScheduler 全局调度器实例
var globalScheduler *Scheduler

// Init 初始化并启动定时任务调度器
// 会自动检测运行环境，仅在Docker环境中启用定时任务
func Init() {
	if globalScheduler != nil {
		log.Println("Scheduler already initialized")
		return
	}

	globalScheduler = NewScheduler()
	globalScheduler.Start()
}

// Shutdown 关闭定时任务调度器，执行清理工作
func Shutdown() {
	if globalScheduler != nil {
		log.Println("Shutting down scheduler...")
		globalScheduler.Stop()
	}
}

// RunNow 立即执行一次导出（用于调试）
func RunNow() error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	globalScheduler.exportToCSV()
	return nil
}

// Scheduler 定时任务调度器
type Scheduler struct {
	isRunningInDocker bool          // 是否在Docker环境中运行
	csvPath           string        // CSV文件保存路径
	stopChan          chan struct{} // 停止信号通道
	exporting         bool          // 是否正在导出
	exportDone        chan struct{} // 导出完成信号
}

// NewScheduler 创建新的调度器实例
func NewScheduler() *Scheduler {
	isDocker := checkDockerEnvironment()

	// 根据运行环境设置CSV保存路径
	csvPath := "/app/data/csv"
	if !isDocker {
		csvPath = "./data/csv"
	}

	return &Scheduler{
		isRunningInDocker: isDocker,
		csvPath:           csvPath,
		stopChan:          make(chan struct{}),
		exportDone:        make(chan struct{}),
	}
}

// checkDockerEnvironment 检测当前是否在Docker容器中运行
// 返回：true表示在Docker环境，false表示本地环境
func checkDockerEnvironment() bool {
	// 检测/.dockerenv文件（Docker容器标志）
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// 检测环境变量DOCKER_CONTAINER
	if os.Getenv("DOCKER_CONTAINER") == "true" {
		return true
	}

	// macOS本地开发环境禁用
	if runtime.GOOS == "darwin" && os.Getenv("IS_DOCKER") != "true" {
		return false
	}

	// 没有数据库主机配置视为本地环境
	if os.Getenv("DB_HOST") == "" {
		return false
	}

	return true
}

// IsRunningInDocker 返回是否在Docker环境中运行
func (s *Scheduler) IsRunningInDocker() bool {
	return s.isRunningInDocker
}

// Start 启动定时任务调度器
// 仅在Docker环境中启用
func (s *Scheduler) Start() {
	log.Printf("isRunningInDocker: %v", s.isRunningInDocker)
	if !s.isRunningInDocker {
		log.Println("Running in local environment, scheduler disabled")
		return
	}

	// 确保CSV目录存在
	if err := s.ensureCSVDirectory(); err != nil {
		log.Printf("Failed to create CSV directory: %v", err)
		return
	}

	log.Printf("CSV export directory: %s", s.csvPath)

	// 启动日常导出任务
	go s.runDailyExport()

	log.Printf("Scheduler started, checking hourly for CSV export to %s", s.csvPath)
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopChan)

	// 如果正在导出，等待导出完成
	if s.exporting {
		log.Println("Waiting for current export to complete...")
		<-s.exportDone
	}

	log.Println("Scheduler stopped")
}

// ensureCSVDirectory 确保CSV保存目录存在
func (s *Scheduler) ensureCSVDirectory() error {
	return os.MkdirAll(s.csvPath, 0755)
}

// hasExportedToday 检查今天是否已经有导出文件
func (s *Scheduler) hasExportedToday(todayInUserZone time.Time) bool {
	dateStr := todayInUserZone.Format("2006-01-02")

	// 检查目录中的文件
	files, err := os.ReadDir(s.csvPath)
	if err != nil {
		log.Printf("Warning: Failed to read CSV directory: %v", err)
		return false
	}

	// 查找是否有今天日期的文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 文件名格式: {username}_{date}.csv
		if len(file.Name()) > len(dateStr)+1 && file.Name()[len(file.Name())-len(dateStr)-4:len(file.Name())-4] == dateStr {
			return true
		}
	}
	return false
}

// runDailyExport 执行每日CSV导出任务
// 每小时检查一次，如果今天没有导出文件则执行
func (s *Scheduler) runDailyExport() {
	loc := getLocation()

	log.Printf("Configured export timezone: %s", loc.String())
	log.Println("Starting hourly check for CSV export")

	for {
		now := time.Now().UTC()
		todayInUserZone := now.In(loc)

		// 检查今天是否已经导出过
		if !s.hasExportedToday(todayInUserZone) {
			log.Printf("No export found for %s, running export now", todayInUserZone.Format("2006-01-02"))
			s.exportToCSV()
		} else {
			log.Printf("Already exported for %s, skipping", todayInUserZone.Format("2006-01-02"))
		}

		// 等待到下一个整点小时（或者1小时后）
		nextCheck := now.Truncate(time.Hour).Add(time.Hour)
		waitDuration := nextCheck.Sub(now)
		log.Printf("Next check at (UTC): %s (in %v)",
			nextCheck.Format(time.RFC3339),
			waitDuration)

		select {
		case <-time.After(waitDuration):
			// 继续下一次检查
		case <-s.stopChan:
			// 收到停止信号，退出
			return
		}
	}
}

// exportToCSV 导出所有用户的订单数据到CSV文件
func (s *Scheduler) exportToCSV() {
	s.exporting = true
	defer func() {
		s.exporting = false
		select {
		case s.exportDone <- struct{}{}:
		default:
		}
	}()

	loc := getLocation()
	// 使用配置的时区作为用户可见的日期
	todayInUserZone := time.Now().In(loc)
	log.Printf("Starting scheduled CSV export... (date: %s)", todayInUserZone.Format("2006-01-02"))

	// 确保目录存在
	if err := s.ensureCSVDirectory(); err != nil {
		log.Printf("Failed to create CSV directory: %v", err)
		return
	}

	// 获取所有用户
	users, err := s.getAllUsers()
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return
	}

	// 为每个用户单独导出
	for _, user := range users {
		if err := s.exportUserOrders(user.ID, user.Username, todayInUserZone); err != nil {
			log.Printf("Failed to export orders for user %s: %v", user.Username, err)
			continue
		}
		log.Printf("Exported CSV for user: %s", user.Username)
	}

	log.Println("Scheduled CSV export completed")
}

// UserInfo 用户基本信息
type UserInfo struct {
	ID       uint   // 用户ID
	Username string // 用户名
}

// getAllUsers 获取所有用户信息
func (s *Scheduler) getAllUsers() ([]UserInfo, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]UserInfo, len(users))
	for i, user := range users {
		result[i] = UserInfo{
			ID:       user.ID,
			Username: user.Username,
		}
	}
	return result, nil
}

// exportUserOrders 导出指定用户的订单数据
// 参数：
//
//	userID - 用户ID
//	username - 用户名（用于文件名）
//	exportDate - 导出日期（用户可见的日期）
func (s *Scheduler) exportUserOrders(userID uint, username string, exportDate time.Time) error {
	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userID).
		Order("date DESC").Find(&orders).Error; err != nil {
		return err
	}

	// 如果没有订单则跳过
	if len(orders) == 0 {
		log.Printf("No orders found for user %s, skipping export", username)
		return nil
	}

	// 生成文件名：{用户名}_{日期}.csv
	filename := fmt.Sprintf("%s_%s.csv", username, exportDate.Format("2006-01-02"))
	filePath := filepath.Join(s.csvPath, filename)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 使用共享的CSV导出功能
	if err := models.WriteOrdersToCSV(file, orders, true); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	log.Printf("Exported %d orders to %s", len(orders), filePath)
	return nil
}

// RunOnce 手动执行一次CSV导出（仅在Docker环境中可用）
// 用于测试或临时导出
func (s *Scheduler) RunOnce() error {
	if !s.isRunningInDocker {
		log.Println("Running in local environment, cannot run scheduled export manually")
		return fmt.Errorf("scheduler only runs in Docker environment")
	}

	s.exportToCSV()
	return nil
}
