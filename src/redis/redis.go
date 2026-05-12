package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/YuanShan-CN/order_management/src/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
	cfg := config.AppConfig.Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	log.Printf("Redis connected: %s:%d/%d", cfg.Host, cfg.Port, cfg.DB)
	return nil
}

func Close() {
	if Client != nil {
		Client.Close()
	}
}