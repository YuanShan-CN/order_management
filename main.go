package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/YuanShan-CN/order_management/src/database"
	"github.com/YuanShan-CN/order_management/src/redis"
	"github.com/YuanShan-CN/order_management/src/router"
	"github.com/YuanShan-CN/order_management/src/scheduler"
)

func main() {
	configPath := flag.String("c", "config.yaml", "Path to config file")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = redis.Connect()
	if err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}

	scheduler.Init()

	r := router.Setup()
	port := strconv.Itoa(config.AppConfig.Server.Port)
	addr := ":" + port

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Println("Server running on http://localhost:" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 关闭调度器
	scheduler.Shutdown()

	// 关闭服务器（等待最多 10 秒）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 关闭数据库
	database.Close()

	// 关闭 Redis
	redis.Close()

	log.Println("Server exited")
}
